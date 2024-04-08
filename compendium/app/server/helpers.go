package server

import (
	"compendium/models"
	"context"
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
		c.Roles = s.getRoles(i)
		cm := s.db.Temp.CorpMemberReadAllByGuildId(context.TODO(), i.Guild.ID)
		for _, member := range cm {
			if i.Guild.Icon == "tg" {
				c.Members = append(c.Members, member)
			} else if CheckRole(i.Guild.ID, member.UserId, roleId) {
				c.Members = append(c.Members, member)
			}
		}
	}
	return c
}
func (s *Server) getRoles(i *models.Identity) []models.CorpRole {
	if i.Guild.Icon == "tg" && i.User.Avatar == "tg" {
		return []models.CorpRole{{
			Id:   "tg",
			Name: "Telegram",
		}}
	} else {
		roles, err := GetRoles(i.Guild.ID)
		if err != nil {
			return nil
		}
		return roles
	}
}
