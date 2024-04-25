package server

import (
	"compendium/models"
	"context"
	"fmt"
)

func (s *Server) GetTokenIdentity(token string) *models.Identity {
	i := s.db.Temp.IdentityRead(context.TODO(), token)
	if i.Token != "" {
		return &i
	}
	return nil
}

func (s *Server) GetCorpData(i *models.Identity, roleId string) models.CorpData {
	c := models.CorpData{}
	c.Members = []models.CorpMemberint{}

	if i.Guild.ID != "" {
		fmt.Printf("%+v\n", i)
		c.Roles = s.getRoles(i)
		cm := s.db.Temp.CorpMemberReadAllByGuildId(context.TODO(), i.Guild.ID)
		for _, member := range cm {
			if i.Guild.Icon == "tg" {
				if roleId == "" {
					c.Members = append(c.Members, member)
				} else {
					roles := s.db.Temp.ReadRoles(i.Guild.ID)
					for _, role := range roles {
						if role.Id == roleId {
							if s.db.Temp.ExistSubscribe(i.Guild.ID, role.Name, member.UserId) {
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
	return c
}
func (s *Server) getRoles(i *models.Identity) []models.CorpRole {
	if i.Guild.Icon == "tg" && i.User.Avatar == "tg" {
		everyone := []models.CorpRole{{
			Id:   "",
			Name: "Telegram",
		}}
		roles := s.db.Temp.ReadRoles(i.Guild.ID)
		if len(roles) > 0 {
			everyone = append(everyone, roles...)
		}
		return everyone
	} else {
		roles, err := GetRoles(i.Guild.ID)
		if err != nil {
			return nil
		}
		return roles
	}
}
