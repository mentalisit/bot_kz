package server

import (
	"compendium_s/models"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (s *Server) GetTokenIdentity(token string) *models.Identity {
	var i models.Identity
	if strings.HasPrefix(token, "Multi_") {
		i.Token = token
		// удаляем префикс "Multi_"
		uid, gid, err := GetTokenData(token[6:])
		if err == nil {
			multiAccount, _ := s.multi.FindMultiAccountUUID(uid)
			if multiAccount != nil {
				i.User = models.User{
					ID:        multiAccount.UUID.String(),
					Username:  multiAccount.Nickname,
					AvatarURL: multiAccount.AvatarURL,
					Alts:      multiAccount.Alts,
					GameName:  multiAccount.Nickname,
				}
			}
			i.MultiAccount = multiAccount

			multiAccountGuild, _ := s.multi.GuildGet(&gid)
			if multiAccountGuild != nil {
				i.Guild = models.Guild{
					URL:  multiAccountGuild.AvatarUrl,
					ID:   gid.String(),
					Name: multiAccountGuild.GuildName,
					Type: "mg",
				}
			}
			i.MultiGuild = multiAccountGuild

			return &i
		} else {
			s.log.ErrorErr(err)
		}
	} else if strings.HasPrefix(token, "identity_") {
		i.Token = token
		userId, GID, err := GetTokenUserData(token[9:])
		if err == nil {

			user, err := s.db.UsersGetByUserId(userId)
			if err != nil {
				s.log.ErrorErr(err)
			}
			i.User = *user

			multiAccountGuild, err := s.multi.GuildGet(&GID)
			if err != nil {
				s.log.ErrorErr(err)
			}
			if multiAccountGuild != nil {
				i.Guild = models.Guild{
					URL:  multiAccountGuild.AvatarUrl,
					ID:   GID.String(),
					Name: multiAccountGuild.GuildName,
					Type: "mg",
				}
			}
			i.MultiGuild = multiAccountGuild

			return &i
		}
	}

	fmt.Println("oldToken ", token)

	userid, guildid, err := s.db.ListUserGetUserIdAndGuildId(token)
	if err != nil {
		token2 := s.db.ListUserGetByMatch(token)
		userid, guildid, err = s.db.ListUserGetUserIdAndGuildId(token2)
		if err != nil {
			s.log.Info("get user by token: " + token + " " + err.Error())
			return nil
		}
	}

	i.MultiGuild, _ = s.multi.GuildGetById(guildid)
	if i.MultiGuild != nil {
		i.Guild = models.Guild{
			URL:  i.MultiGuild.AvatarUrl,
			ID:   i.MultiGuild.GId.String(),
			Name: i.MultiGuild.GuildName,
			Type: "mg",
		}
	}

	ma, _ := s.multi.FindMultiAccountByUserId(userid)
	if ma != nil {
		i.MultiAccount = ma
		var gg *models.MultiAccountGuild
		corpMember, _ := s.multi.CorpMemberByUId(ma.UUID)
		if corpMember != nil {
			for _, id := range corpMember.GuildIds {
				g, _ := s.multi.GuildGet(&id)
				for _, channel := range g.Channels {
					if channel == guildid {
						i.MultiGuild = g
						gg = g
						break
					}
				}
			}
		}
		if i.MultiGuild != nil && i.MultiAccount != nil {
			i.Token, _ = JWTGenerateToken(i.MultiAccount.UUID, i.MultiGuild.GId, ma.Nickname)
			i.User = models.User{
				Username:  ma.Nickname,
				AvatarURL: ma.AvatarURL,
				Alts:      ma.Alts,
			}
			if gg != nil {
				i.Guild = models.Guild{
					URL:  gg.AvatarUrl,
					ID:   gg.GId.String(),
					Name: gg.GuildName,
					Type: "ma",
				}
			}
			return &i
		}

	}

	i.Token = token
	user, err := s.db.UsersGetByUserId(userid)
	if err != nil {
		s.log.ErrorErr(err)
		s.log.Info("get user by userid :" + userid)
		return nil
	}
	i.User = *user
	//guild, err := s.db.GuildGet(guildid)
	//if err != nil {
	//	parse, _ := uuid.Parse(guildid)
	//	if parse.String() == guildid {
	//		get, err := s.multi.GuildGet(&parse)
	//		if err == nil {
	//			i.MultiGuild = get
	//			i.Guild = models.Guild{
	//				URL:  get.AvatarUrl,
	//				ID:   get.GId.String(),
	//				Name: get.GuildName,
	//				Type: "mg",
	//			}
	//			return &i
	//		}
	//	}
	//}
	//if err != nil {
	//	s.log.ErrorErr(err)
	//	s.log.Info("get guild by guildid:" + guildid)
	//	return nil
	//}
	//i.Guild = *guild

	return &i
}

func (s *Server) getRoles(i models.Guild) []models.CorpRole {
	if len(i.ID) > 24 {
		parse, _ := uuid.Parse(i.ID)
		GUILD, err := s.multi.GuildGet(&parse)
		if err != nil {
			return nil
		}
		var rolesAll []models.CorpRole
		for _, channel := range GUILD.Channels {
			if strings.HasPrefix(channel, "-100") {
				everyone := []models.CorpRole{{
					Id:   "",
					Name: "@everyone",
				}}
				roles, err := s.db.GuildRolesRead(channel)
				if err != nil {
					s.log.ErrorErr(err)
				}
				if len(roles) > 0 {
					everyone = append(everyone, roles...)
				}
				rolesAll = append(rolesAll, everyone...)
			} else {
				roles, err := s.roles.ds.GetRoles(channel)
				if err != nil {
					s.log.ErrorErr(err)
					continue
				}
				rolesAll = append(rolesAll, roles...)
			}
		}
		return rolesAll
	}

	if i.Type == "tg" {
		everyone := []models.CorpRole{{
			Id:   "",
			Name: "@everyone",
		}}
		roles, err := s.db.GuildRolesRead(i.ID)
		if err != nil {
			s.log.ErrorErr(err)
		}
		if len(roles) > 0 {
			everyone = append(everyone, roles...)
		}
		return everyone
	} else {
		roles, err := s.roles.ds.GetRoles(i.ID)
		if err != nil {
			s.log.ErrorErr(err)
			return nil
		}
		return roles
	}
}

func (s *Server) CheckCode(code string) models.Identity {
	var i models.Identity
	coder, err := s.db.CodeGet(code)
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
				s.db.CodeDelete(value.Code)
				fmt.Printf("Delete Code %+v\n", value)
			} else {
				s.db.CodeDelete(m.Code)
				fmt.Printf("Delete Code %+v\n", m)
			}
		}
	}
}
func (s *Server) refreshToken(token string) string {
	if strings.HasPrefix(token, "Multi_") || strings.HasPrefix(token, "identity_") {
		return token
	}

	if len(token) < 60 {
		newToken := s.checkPrefixToken(token)
		err := s.db.ListUserUpdateToken(token, newToken)
		if err != nil {
			return token
		}
		return newToken
	}
	return token
}

func (s *Server) checkPrefixToken(token string) string {
	userid, guildid, err := s.db.ListUserGetUserIdAndGuildId(token)
	if err != nil || userid == "" || guildid == "" {
		s.log.Info("err || userid == nil || guildid == nil")
		return token
	}
	guildGet, err := s.db.GuildGet(guildid)
	if err != nil || guildGet.Name == "" {
		s.log.Info("err || guildGet.Name == nil")
		return token
	}

	newToken := guildGet.Type + guildid + "." + userid + GenerateToken()
	return newToken
}
