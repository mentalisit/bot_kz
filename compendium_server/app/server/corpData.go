package server

import (
	"compendium_s/models"
	"fmt"
	"strings"
)

func (s *Server) GetCorpData(i *models.Identity, roleId string) *models.CorpData {
	c := models.CorpData{}
	c.Members = []models.CorpMember{}

	if i.MultiGuild != nil {
		return s.GetCorpDataMultiGuild(i, roleId)
	}

	fmt.Println("GetCorpData use old")

	together := s.GetCorpDataIfTogether(i, roleId)
	if together != nil {
		return together
	}
	if i.Guild.Type == "ds" {
		s.roles.LoadGuild(i.Guild.ID)
	}

	if i.Guild.ID != "" {
		c.Roles = s.getRoles(i.Guild)
		cm, err := s.db.CorpMembersRead(i.Guild.ID)
		if err != nil {
			s.log.ErrorErr(err)
			//return nil
		}
		var roles []models.CorpRole
		if i.Guild.Type == "tg" && roleId != "" {
			roles, err = s.db.GuildRolesRead(i.Guild.ID)
			if err != nil {
				s.log.ErrorErr(err)
			}
		}

		for _, member := range cm {
			if i.Guild.Type == "tg" {
				if roleId == "" || roleId == "tg" {
					c.Members = append(c.Members, member)
				} else {
					for _, role := range roles {
						if role.Id == roleId {
							if s.db.GuildRolesExistSubscribe(i.Guild.ID, role.Name, member.UserId) {
								c.Members = append(c.Members, member)
							}
						}
					}
				}
			} else if i.Guild.Type == "ds" {
				uid := member.UserId
				if strings.Contains(member.UserId, "/") {
					split := strings.Split(member.UserId, "/")
					uid = split[0]
				}
				if s.roles.ds.CheckRoleDs(i.Guild.ID, uid, roleId) {
					c.Members = append(c.Members, member)
				}
			}
		}
	}

	return &c
}
