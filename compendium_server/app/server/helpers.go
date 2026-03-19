package server

import (
	"compendium_s/models"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) GetTokenIdentity(token string) *models.Identity {
	var i models.Identity
	i.Token = token

	if strings.HasPrefix(token, "Multi_") {
		token = token[6:] // удаляем префикс "Multi_"
	} else if strings.HasPrefix(token, "my_compendium_") {
		token = token[14:] // удаляем префикс "my_compendium_"
	} else {
		return s.GetTokenIdentityByOldToken(token)
	}

	uid, gid, err := GetTokenData(token)
	if err == nil {
		if uid.String() == "00000000-0000-0000-0000-000000000000" {
			s.log.Error("да ну пиздец заебал ")
			return nil
		}
		i.MAccount, err = s.dbV2.FindMultiAccountUUID(uid)
		if err != nil || i.MAccount == nil {
			fmt.Println(err.Error())
			return nil
		}
		i.User = multiToUser(i.MAccount)

		gGet, err := s.dbV2.GuildGet(gid)
		if err != nil {
			s.log.ErrorErr(err)
		}
		i.MGuild = gGet
		i.Guild = multiToGuild(gGet)

		return &i
	}
	return nil
}

func (s *Server) CheckCode(code string) models.Identity {
	var i models.Identity
	//coder, err := s.db.CodeGet(code)
	coder, err := s.dbV2.CodeGet(code)
	if err != nil {
		fmt.Println("CheckCode " + err.Error())
	}

	if coder != nil && coder.Code == code {
		if time.Now().Unix() < coder.Timestamp+600 {
			i = coder.Identity
		} else {
			go s.CleanOldCodes()
		}
	}

	return i
}
func (s *Server) CleanOldCodes() {
	all := s.db.CodeAllGet()
	names := make(map[string]models.Code)
	for _, m := range all {
		value, exists := names[m.Identity.User.Username]
		if !exists {
			names[m.Identity.User.Username] = m
		} else {
			if value.Timestamp < m.Timestamp {
				names[m.Identity.User.Username] = m
				s.dbV2.CodeDelete(value.Code)
				fmt.Printf("Delete Code %+v\n", value)
			} else {
				s.dbV2.CodeDelete(m.Code)
				fmt.Printf("Delete Code %+v\n", m)
			}
		}
	}
}

func extractAndValidateCheckHeaders(c *gin.Context) (token, roleId, mGuild string, err error) {
	token = c.GetHeader("authorization")
	if token == "" {
		return "", "", "", errors.New("missing token")
	}
	roleId = c.Query("roleId")
	if roleId == "" {
		roleId = c.GetHeader("X-Role-ID")
	}
	mGuild = c.Query("corpId")
	if mGuild == "" {
		mGuild = c.GetHeader("X-Corp-ID")
	}
	return token, roleId, mGuild, nil

}

func (s *Server) GetTokenIdentityByOldToken(token string) *models.Identity {
	var i models.Identity
	if strings.HasPrefix(token, "identity_") {
		i.Token = token
		userId, GID, err := GetTokenUserData(token[9:])
		if err == nil {
			mAcc, _ := s.dbV2.FindMultiAccountByUserId(userId)
			if mAcc != nil &&
				(mAcc.DiscordID == userId || mAcc.TelegramID == userId || mAcc.WhatsappID == userId) {
				i.MAccount = mAcc
				i.User = multiToUser(mAcc)
			} else {
				user, err := s.db.UsersGetByUserId(userId)
				if err != nil {
					s.log.ErrorErr(err)
				}
				i.User = *user
				ma := models.MultiAccount{
					Nickname:  user.GameName,
					AvatarURL: user.AvatarURL,
					Alts:      user.Alts,
				}
				if strings.Contains(user.ID, "@") {
					ma.WhatsappUsername = user.Username
					ma.WhatsappID = user.ID
				} else if len(userId) < 13 {
					ma.TelegramUsername = user.Username
					ma.TelegramID = user.ID
				} else if len(userId) > 16 {
					ma.DiscordUsername = user.Username
					ma.DiscordID = user.ID
				}

				if ma.Nickname == "" {
					ma.Nickname = user.Username
				}

				i.MAccount, err = s.dbV2.CreateMultiAccountFull(ma)
				if err != nil {
					s.log.ErrorErr(err)
				}
			}

			i.MGuild, _ = s.dbV2.GuildGet(GID)
			if i.MGuild != nil {
				i.Guild = multiToGuild(i.MGuild)
			}

			if i.MAccount != nil && i.MAccount.Nickname != "" {
				i.Token, _ = JWTGenerateTokenV2(i.MAccount.UUID, i.MGuild.GId)
			}

			go s.searchAndMove(i)
			return &i

		}
	}

	if token != "" { //search in old database
		userid, guildid, err := s.db.ListUserGetUserIdAndGuildId(token)
		if strings.HasPrefix(userid, "1862488373") {
			userid = "1862488373"
		}
		if err != nil {
			userid, guildid, err = ParseComplexToken(token)
			if err == nil && userid != "" && guildid != "" {
				//найдено
			} else {
				token2 := s.db.ListUserGetByMatch(token)
				userid, guildid, err = s.db.ListUserGetUserIdAndGuildId(token2)
				if err != nil {
					s.log.Info("get user by token: " + token + " " + err.Error())
					return nil
				}
			}
		}

		i.MGuild, _ = s.dbV2.GuildGetById(guildid)
		if i.MGuild != nil {
			i.Guild = multiToGuild(i.MGuild)
		} else {
			s.log.Info("guildId " + guildid + " not found")
		}

		ma, _ := s.dbV2.FindMultiAccountByUserId(userid)
		if ma != nil && ma.Nickname != "" {
			i.MAccount = ma
		} else {
			user, err := s.db.UsersGetByUserId(userid)
			if err != nil {
				s.log.ErrorErr(err)
				s.log.Info(fmt.Sprintf("get userID: %s\n", userid))
				fmt.Printf("get userID: %s token %s\n", userid, token)
				return nil
			}
			ma = &models.MultiAccount{
				Nickname:  user.GameName,
				AvatarURL: user.AvatarURL,
				Alts:      user.Alts,
			}
			if strings.Contains(user.ID, "@") {
				ma.WhatsappUsername = user.Username
				ma.WhatsappID = user.ID
			} else if len(userid) < 13 {
				ma.TelegramUsername = user.Username
				ma.TelegramID = user.ID
			} else if len(userid) > 16 {
				ma.DiscordUsername = user.Username
				ma.DiscordID = user.ID
			}

			if ma.Nickname == "" {
				ma.Nickname = user.Username
			}

			i.MAccount, err = s.dbV2.CreateMultiAccountFull(*ma)
			if err != nil {
				s.log.ErrorErr(err)
			}
		}

		if i.MGuild != nil && i.MAccount != nil {
			i.Token, _ = JWTGenerateTokenV2(i.MAccount.UUID, i.MGuild.GId)
			i.User = multiToUser(i.MAccount)
			go s.searchAndMove(i)
			return &i
		}
	}
	return nil
}
func multiToUser(m *models.MultiAccount) models.User {
	return models.User{
		ID:        m.UUID.String(),
		Username:  m.Nickname,
		AvatarURL: m.AvatarURL,
		Alts:      m.Alts,
		GameName:  m.Nickname,
	}
}
func multiToGuild(g *models.MultiAccountGuildV2) models.Guild {
	return models.Guild{
		URL:  g.AvatarUrl,
		ID:   g.GId.String(),
		Name: g.GuildName,
		Type: "mg",
	}
}
func (s *Server) searchAndMove(i models.Identity) {
	if i.MAccount == nil && i.User.ID != "" {
		id, err := s.dbV2.FindMultiAccountByUserId(i.User.ID)
		if err != nil {
			s.log.ErrorErr(err)
		}
		i.MAccount = id

	}
	if i.MGuild == nil || i.MAccount == nil {
		s.log.InfoStruct("what ", i)
		return
	}
	m2 := s.db.SearchOldData(i)
	var guilds map[uuid.UUID]struct{}
	guilds = make(map[uuid.UUID]struct{})
	for _, id := range m2.CorpMember.GuildIds {
		guilds[id] = struct{}{}
	}
	m := models.Moving{
		MAcc: *i.MAccount,
	}
	for u, _ := range guilds {
		m.CorpMember.GuildIds = append(m.CorpMember.GuildIds, u)
	}
	m.CorpMember.Uid = i.MAccount.UUID
	if m.CorpMember.TimeZona == "" {
		m.CorpMember.TimeZona = m2.CorpMember.TimeZona
	}
	if m.CorpMember.AfkFor == "" {
		m.CorpMember.AfkFor = m2.CorpMember.AfkFor
	}
	techMap := make(map[string]models.TechLevels)
	for _, member := range m2.Tech {
		if techMap[member.Name] == nil {
			techMap[member.Name] = member.Tech
		} else {
			for module, data := range member.Tech {
				if techMap[member.Name][module].Level == 0 || techMap[member.Name][module].Ts < data.Ts {
					techMap[member.Name][module] = data
				}
			}
		}
	}
	for name, levels := range techMap {
		m.Tech = append(m.Tech, models.Technology{
			Tech: levels,
			Name: name,
		})
	}
	//read complete

	accountUUID, err := s.dbV2.FindMultiAccountUUID(m.MAcc.UUID)
	if err == nil && accountUUID != nil {
		if len(m.CorpMember.GuildIds) != 0 {
			fmt.Printf("searchAndMove CorpMember %+v\n", m.CorpMember)
			err = s.dbV2.CorpMemberInsert(m.CorpMember)
			if err != nil {
				s.log.ErrorErr(err)
			}
		}
		if len(m.Tech) != 0 {
			fmt.Printf("searchAndMove Technology %+v\n", m.Tech)
			for _, tech := range m.Tech {
				err = s.dbV2.TechnologiesUpdate(m.MAcc.UUID, tech.Name, tech.Tech)
				if err != nil {
					s.log.ErrorErr(err)
				}
			}
		}
	}

	if m.MAcc.TelegramID != "" {
		s.db.DeleteOldClient(m.MAcc.TelegramID)
	}
	if m.MAcc.DiscordID != "" {
		s.db.DeleteOldClient(m.MAcc.DiscordID)
	}
	if m.MAcc.WhatsappID != "" {
		s.db.DeleteOldClient(m.MAcc.WhatsappID)
	}
	return
}

func ParseComplexToken(raw string) (userid, guildId string, err error) {
	if strings.HasPrefix(raw, "ds") || strings.HasPrefix(raw, "tg") {
		// 1. Определяем платформу
		remainder := ""
		if strings.HasPrefix(raw, "tg") {
			remainder = strings.TrimPrefix(raw, "tg")
		} else if strings.HasPrefix(raw, "ds") {
			remainder = strings.TrimPrefix(raw, "ds")
		} else {
			return "", "", fmt.Errorf("unknown prefix")
		}

		// 2. Делим по первой точке: [ChatID].[Все остальное]
		parts := strings.SplitN(remainder, ".", 2)
		if len(parts) < 2 {
			return "", "", fmt.Errorf("invalid format")
		}
		guildId = parts[0]

		// 3. Используем регулярное выражение, чтобы отделить цифры (UserID) от букв (Secret)
		// ^(\d+) ищет только цифры в начале строки
		re := regexp.MustCompile(`^(\d+)(.*)`)
		matches := re.FindStringSubmatch(parts[1])

		if len(matches) > 2 {
			userid = matches[1] // Только цифры
		}

		return userid, guildId, nil
	}
	return "", "", fmt.Errorf("invalid token")

}
