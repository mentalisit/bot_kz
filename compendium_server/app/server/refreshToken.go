package server

import (
	"compendium_s/config"
	"compendium_s/models"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (s *Server) refreshToken2(token string) *models.Identity {
	i := s.GetTokenIdentity(token)
	if i == nil {
		s.log.Info("refreshToken2: token not found")
		return nil
	}

	if i.MAccount == nil {
		if i.MultiAccount != nil {
			i.MAccount = i.MultiAccount
		} else {
			i.MAccount = &models.MultiAccount{
				AvatarURL: i.User.AvatarURL,
				Alts:      i.User.Alts,
			}
			if i.User.GameName != "" {
				i.MAccount.Nickname = i.User.GameName
			} else {
				i.MAccount.Nickname = i.User.Username
			}
			if strings.Contains(i.User.ID, "@") {
				// 1. WhatsApp — определяем по спецсимволу
				i.MAccount.WhatsappID = i.User.ID
				i.MAccount.WhatsappUsername = i.User.Username
			} else {
				// Если это число, пробуем понять какое
				idLen := len(i.User.ID)

				if idLen >= 17 {
					// 2. Discord — всё, что длинное и числовое
					i.MAccount.DiscordID = i.User.ID
					i.MAccount.DiscordUsername = i.User.Username
				} else {
					// 3. Telegram — всё остальное короткое
					// (включая будущие 12, 13... значные ID)
					i.MAccount.TelegramID = i.User.ID
					i.MAccount.TelegramUsername = i.User.Username
				}
			}
		}
	}

	full, err := s.dbV2.CreateMultiAccountFull(*i.MAccount)
	if err != nil {
		s.log.ErrorErr(err)
	}
	i.MAccount = full

	if i.MultiAccount != nil {
		member, _ := s.multi.CorpMemberByUId(i.MultiAccount.UUID)
		if member != nil {
			err = s.dbV2.CorpMemberInsert(*member)
			if err != nil {
				s.log.ErrorErr(err)
			}
			technologies := s.multi.TechnologiesGetMember(member.Uid)
			if technologies != nil {
				for _, technology := range technologies {
					err := s.dbV2.TechnologiesUpdate(member.Uid, technology.Name, technology.Tech)
					if err != nil {
						s.log.ErrorErr(err)
					}
				}
			}
		}
	} else {
		read, err := s.db.CorpMemberRead(i.User.ID)
		if err != nil {
			s.log.ErrorErr(err)
		}
		member := models.MultiAccountCorpMember{
			Uid: i.MultiAccount.UUID,
		}
		tech := make(map[string]models.TechLevels)

		for _, m := range read {
			if m.TimeZone != "" {
				member.TimeZona = m.TimeZone
			}
			if m.ZoneOffset != 0 {
				member.ZonaOffset = m.ZoneOffset
			}
			if m.AfkFor != "" {
				member.AfkFor = m.AfkFor
			}
			ggid, _ := uuid.Parse(m.GuildId)
			if ggid.String() != "" {
				member.GuildIds = append(member.GuildIds, ggid)
			}
			if tech[m.Name] == nil {
				tech[m.Name] = make(models.TechLevels)
			}
			for module, data := range m.Tech {
				if tech[m.Name][module].Level == 0 || tech[m.Name][module].Ts < data.Ts {
					tech[m.Name][module] = data
				}
			}
		}
		for s2, levels := range tech {
			err := s.dbV2.TechnologiesUpdate(member.Uid, s2, levels)
			if err != nil {
				s.log.ErrorErr(err)
			}
		}
	}
	i.Token, _ = JWTGenerateTokenV2(i.MAccount.UUID, i.MGuild.GId)

	return i
}

func JWTGenerateTokenV2(uuid, gid uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"uuid": uuid,
		"gid":  gid,
		"exp":  time.Now().AddDate(1, 0, 0).Unix(), // токен на год
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.Instance.Postgress.Password))
	if err != nil {
		return "", err
	}

	// добавляем префикс
	return "my_compendium_" + signedToken, nil
}
