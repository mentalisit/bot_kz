package server

import (
	"compendium_s/models"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) GetTokenIdentity(token string) *models.Identity {
	var i models.Identity
	if strings.HasPrefix(token, "my_compendium_") {
		i.Token = token
		// удаляем префикс "my_compendium_"
		uid, gid, err := GetTokenData(token[14:])
		if err == nil {
			if uid.String() == "00000000-0000-0000-0000-000000000000" {
				s.log.Error("да ну пиздец заебал ")
				return nil
			}
			i.MAccount, err = s.dbV2.FindMultiAccountUUID(uid)
			if err != nil {
				fmt.Println(err.Error())
				return nil
			}
			i.MultiAccount = i.MAccount
			i.User = multiToUser(i.MAccount)

			gGet, err := s.dbV2.GuildGet(gid)
			if err != nil {
				s.log.ErrorErr(err)
			}
			i.MGuild = gGet
			i.Guild = multiToGuild(gGet)

			return &i
		}
		s.log.ErrorErr(err)
	}
	return s.GetTokenIdentityByOldToken(token)
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
	if strings.HasPrefix(token, "Multi_") {
		i.Token = token
		// удаляем префикс "Multi_"
		uid, gid, err := GetTokenData(token[6:])
		if err == nil {
			multiAccount, _ := s.multi.FindMultiAccountUUID(uid)
			if multiAccount != nil && multiAccount.UUID == uid {
				i.MultiAccount = multiAccount
				i.MAccount = multiAccount
			} else {
				mAcc, _ := s.dbV2.FindMultiAccountUUID(uid)
				if mAcc != nil && mAcc.UUID == uid {
					i.MAccount = mAcc
					i.MultiAccount = mAcc
				}
			}
			i.User = multiToUser(i.MAccount)

			i.MGuild, _ = s.dbV2.GuildGet(gid)
			if i.MGuild != nil {
				i.Guild = multiToGuild(i.MGuild)
			}
			i.Token, _ = JWTGenerateTokenV2(i.MAccount.UUID, i.MGuild.GId)
			go s.searchAndMove(i)

			return &i
		}
		s.log.ErrorErr(err)

	}
	if strings.HasPrefix(token, "identity_") {
		i.Token = token
		userId, GID, err := GetTokenUserData(token[9:])
		if err == nil {
			mAcc, _ := s.dbV2.FindMultiAccountByUserId(userId)
			if mAcc != nil &&
				(mAcc.DiscordID == userId || mAcc.TelegramID == userId || mAcc.WhatsappID == userId) {
				i.MAccount = mAcc
				i.MultiAccount = mAcc
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
				i.MultiAccount = i.MAccount
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
		if err != nil {
			token2 := s.db.ListUserGetByMatch(token)
			userid, guildid, err = s.db.ListUserGetUserIdAndGuildId(token2)
			if err != nil {
				s.log.Info("get user by token: " + token + " " + err.Error())
				return nil
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
			i.MultiAccount = ma
			i.MAccount = ma
		} else {
			user, err := s.db.UsersGetByUserId(userid)
			if err != nil {
				s.log.ErrorErr(err)
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
			i.MultiAccount = i.MAccount
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
	if i.MGuild == nil || i.MAccount == nil {
		s.log.InfoStruct("what ", i)
		return
	}
	m1 := s.multi.SearchOldData(i)
	m2 := s.db.SearchOldData(i)
	var guilds map[uuid.UUID]struct{}
	guilds = make(map[uuid.UUID]struct{})
	for _, id := range m1.CorpMember.GuildIds {
		guilds[id] = struct{}{}
	}
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
	m.CorpMember.TimeZona = m1.CorpMember.TimeZona
	if m.CorpMember.TimeZona == "" {
		m.CorpMember.TimeZona = m2.CorpMember.TimeZona
	}
	m.CorpMember.ZonaOffset = m1.CorpMember.ZonaOffset
	if m.CorpMember.ZonaOffset == 0 {
		m.CorpMember.ZonaOffset = m2.CorpMember.ZonaOffset
	}
	m.CorpMember.AfkFor = m1.CorpMember.AfkFor
	if m.CorpMember.AfkFor == "" {
		m.CorpMember.AfkFor = m2.CorpMember.AfkFor
	}
	techMap := make(map[string]models.TechLevels)
	for _, member := range m1.Tech {
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
		err = s.dbV2.CorpMemberInsert(m.CorpMember)
		if err != nil {
			s.log.ErrorErr(err)
		}
		for _, tech := range m.Tech {
			err = s.dbV2.TechnologiesUpdate(m.MAcc.UUID, tech.Name, tech.Tech)
			if err != nil {
				s.log.ErrorErr(err)
			}
		}
	}
	fmt.Printf("searchAndMove complete %+v\n", m)
	s.multi.DeleteOldClient(m.MAcc.UUID)
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
