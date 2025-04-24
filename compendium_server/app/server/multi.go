package server

import (
	"compendium_s/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strings"
)

func (s *Server) GetCorpDataMulti(i *models.Identity, roleId string) *models.CorpData {
	c := models.CorpData{}
	c.Members = []models.CorpMember{}
	var members []models.CorpMember
	var guildDs, guildTg []string
	appendRolesByType := func(roles []models.CorpRole, nameMSG string) {
		if len(roles) > 0 {
			for _, role := range roles {
				if role.Name == "@everyone" {
					continue
				}
				c.Roles = append(c.Roles, models.CorpRole{
					Id:       role.Id,
					Name:     fmt.Sprintf("(%s) %s", strings.ToUpper(nameMSG), role.Name),
					TypeRole: nameMSG,
				})
			}
		}
	}
	loadMembersByType := func(channel, nameMSG string) {
		cm, _ := s.db.CorpMembersRead(channel)
		for _, m := range cm {
			m.Name = fmt.Sprintf("(%s) %s", strings.ToUpper(nameMSG), m.Name)
			m.TypeAccount = nameMSG
			members = append(members, m)
		}
	}
	CheckRoleDs := func(role models.CorpRole) {
		for _, gId := range guildDs {
			for _, member := range members {
				if member.TypeAccount == role.TypeRole {
					uid := member.UserId
					if strings.Contains(member.UserId, "/") {
						split := strings.Split(member.UserId, "/")
						uid = split[0]
					}
					if s.roles.CheckRoleDs(gId, uid, roleId) {
						c.Members = append(c.Members, member)
						return
					}
				}
			}
		}
	}
	CheckRoleTg := func(role models.CorpRole) {
		for _, gId := range guildTg {
			for _, member := range members {
				if s.db.GuildRolesExistSubscribe(gId, role.Name, member.UserId) {
					c.Members = append(c.Members, member)
				}
			}
		}
	}

	if i.Uid != nil && i.GId != nil {
		c.Roles = []models.CorpRole{{
			Id:   "",
			Name: "@everyone",
		}}
		memberMulti, _ := s.multi.CorpMembersRead(*i.GId)
		if len(memberMulti) != 0 {
			c.Roles = append(c.Roles, models.CorpRole{
				Id:   "ma",
				Name: "(MA)@everyone",
			})
			for _, member := range memberMulti {
				member.Name = fmt.Sprintf("(%s) %s", strings.ToUpper(member.TypeAccount), member.Name)
				members = append(members, member)
			}
		}
		guild, err := s.multi.GuildGet(i.GId)
		if err != nil {
			s.log.ErrorErr(err)
		}
		if guild != nil {
			for _, channel := range guild.Channels {
				if !strings.HasPrefix(channel, "-100") {
					guildDs = append(guildDs, channel)
				} else if strings.HasPrefix(channel, "-100") {
					guildTg = append(guildTg, channel)
				}
			}

			if len(guildDs) != 0 {
				c.Roles = append(c.Roles, models.CorpRole{
					Id:       "ds",
					Name:     "(DS)@everyone",
					TypeRole: "ds",
				})
			}
			if len(guildTg) != 0 {
				c.Roles = append(c.Roles, models.CorpRole{
					Id:   "tg",
					Name: "(TG)@everyone",
				})
			}
			for _, gds := range guildDs {
				s.roles.LoadGuild(gds)
				roles, err := s.roles.ds.GetRoles(gds)
				if err != nil {
					s.log.ErrorErr(err)
				}
				appendRolesByType(roles, "ds")
				loadMembersByType(gds, "ds")
			}
			for _, gtg := range guildTg {
				roles, err := s.db.GuildRolesRead(gtg)
				if err != nil {
					s.log.ErrorErr(err)
				}
				appendRolesByType(roles, "tg")
				loadMembersByType(gtg, "tg")
			}

			if roleId == "" {
				c.Members = members
			} else if roleId == "tg" || roleId == "ds" || roleId == "ma" {
				for _, member := range members {
					if member.TypeAccount == roleId {
						c.Members = append(c.Members, member)
					}
				}
			} else {
				var role models.CorpRole
				for _, roles := range c.Roles {
					if roleId == roles.Id {
						role = roles
					}
				}
				if role.TypeRole == "ds" {
					CheckRoleDs(role)
				} else if role.TypeRole == "tg" {
					CheckRoleTg(role)
				} else if role.TypeRole == "ma" {
					var cm []models.CorpMember
					for _, member := range members {
						if member.TypeAccount == role.TypeRole {
							cm = append(cm, member)
						}
					}
					for _, member := range cm {
						found := false
						for _, gId := range guildTg {
							if s.db.GuildRolesExistSubscribe(gId, role.Name, member.Multi.TelegramID) {
								found = true
								break
							}
						}

						for _, gId := range guildDs {
							uid := member.Multi.DiscordID
							if s.roles.CheckRoleDs(gId, uid, roleId) {
								found = true
								break
							}
						}
						if found {
							c.Members = append(c.Members, member)
						}
					}
				}
			}
		}
		sort.Slice(c.Members, func(i, j int) bool {
			// Проверка, чтобы индекс не выходил за пределы строки
			nameI := c.Members[i].Name
			nameJ := c.Members[j].Name

			// Игнорируем первые пять символов, если длина имени больше или равна пяти
			if len(nameI) >= 5 {
				nameI = nameI[5:]
			}
			if len(nameJ) >= 5 {
				nameJ = nameJ[5:]
			}

			return nameI < nameJ
		})
		return &c
	}
	return nil
}

func (s *Server) SyncTechMulti(c *gin.Context, i *models.Identity, mode, twin string) {
	userName := i.User.Username

	if twin != "" && twin != "default" {
		userName = twin
	}

	fmt.Printf("mode %s corporation %s Name %s\n", mode, i.Guild.Name, userName)

	if mode == "get" {
		sd := models.SyncData{
			TechLevels: models.TechLevels{},
			Ver:        2,
			InSync:     1,
		}
		techBytes, err := s.multi.TechnologiesGet(*i.Uid, userName)
		if err == nil && len(techBytes) > 0 {
			sd.TechLevels = sd.TechLevels.ConvertToTech(techBytes)
		}
		c.JSON(http.StatusOK, sd)
	} else if mode == "sync" {

		var data models.SyncData
		if err := c.BindJSON(&data); err != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		bytes, err := json.Marshal(data.TechLevels)
		if err != nil {
			s.log.ErrorErr(err)
		}
		err = s.multi.TechnologiesUpdate(*i.Uid, userName, bytes)
		if err != nil {
			s.log.ErrorErr(err)
		}

		// Используйте переменную data с полученными данными
		c.JSON(http.StatusOK, data)
	}
}
