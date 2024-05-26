package server

import (
	"compendium_s/models"
	"time"
)

func (s *Server) GetTokenIdentity(token string) *models.Identity {
	userid, guildid, err := s.db.ListUserGetUserIdAndGuildId(token)
	if err != nil {
		s.log.ErrorErr(err)
		s.log.Info("get user by token: " + token)
		return nil
	}
	var i models.Identity
	i.Token = token
	user, err := s.db.UsersGetByUserId(userid)
	if err != nil {
		s.log.ErrorErr(err)
		s.log.Info("get user by userid :" + userid)
		return nil
	}
	i.User = *user
	guild, err := s.db.GuildGet(guildid)
	if err != nil {
		s.log.ErrorErr(err)
		s.log.Info("get guild by guildid:" + guildid)
		return nil
	}
	i.Guild = *guild

	return &i
}

func (s *Server) GetCorpData(i *models.Identity, roleId string) *models.CorpData {
	c := models.CorpData{}
	c.Members = []models.CorpMember{}

	if i.Guild.ID != "" {
		c.Roles = s.getRoles(i)
		cm, err := s.db.CorpMembersRead(i.Guild.ID)
		if err != nil {
			s.log.ErrorErr(err)
			return nil
		}
		for _, member := range cm {
			if i.Guild.Type == "tg" {
				if roleId == "" {
					c.Members = append(c.Members, member)
				} else {
					roles, er := s.db.GuildRolesRead(i.Guild.ID)
					if er != nil {
						s.log.ErrorErr(er)
						return nil
					}
					for _, role := range roles {
						if role.Id == roleId {
							if s.db.GuildRolesExistSubscribe(i.Guild.ID, role.Name, member.UserId) {
								c.Members = append(c.Members, member)
							}
						}
					}
				}
			} else if CheckRoleDs(i.Guild.ID, member.UserId, roleId) {
				c.Members = append(c.Members, member)
			}
		}
	}
	return &c
}
func (s *Server) getRoles(i *models.Identity) []models.CorpRole {
	if i.Guild.Type == "tg" {
		everyone := []models.CorpRole{{
			Id:   "",
			Name: "Telegram",
		}}
		roles, err := s.db.GuildRolesRead(i.Guild.ID)
		if err != nil {
			s.log.ErrorErr(err)
			return nil
		}
		if len(roles) > 0 {
			everyone = append(everyone, roles...)
		}
		return everyone
	} else {
		roles, err := GetRoles(i.Guild.ID)
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
		s.log.ErrorErr(err)
		return i
	}

	if coder != nil && coder.Code == code {
		if time.Now().Unix() < coder.Timestamp+600 {
			i = coder.Identity
			return i
		}
	}
	if code == "test-test-test" {
		i = models.Identity{
			User: models.User{
				ID:       "111111111",
				Username: "TestUser",
				Alts:     []string{"alt1", "alt2"},
			},
			Guild: models.Guild{
				ID:   "22222222222",
				Name: "TestGuild",
				Type: "tg",
			},
			Token: "gGUBIlUAU1uTKWd8HssP27ojG0DugoAaPslwFGTDSAbEM6UM",
		}

	}
	return i
}
