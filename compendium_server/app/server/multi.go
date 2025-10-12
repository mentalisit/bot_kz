package server

import (
	"compendium_s/models"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

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
		techBytes, err := s.multi.TechnologiesGet(i.MultiAccount.UUID, userName)
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
		err = s.multi.TechnologiesUpdate(i.MultiAccount.UUID, userName, bytes)
		if err != nil {
			s.log.ErrorErr(err)
		}

		// Используйте переменную data с полученными данными
		c.JSON(http.StatusOK, data)
	}
}
func (s *Server) SyncTechMultiGuild(c *gin.Context, i *models.Identity, mode, twin string) {
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
		techBytes, err := s.db.TechGet(userName, i.User.ID, i.MultiGuild.GId.String())
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
		err = s.db.TechUpdate(userName, i.User.ID, i.MultiGuild.GId.String(), bytes)
		if err != nil {
			s.log.ErrorErr(err)
		}

		// Используйте переменную data с полученными данными
		c.JSON(http.StatusOK, data)
	}
}
func getCorpsTypeId(mg *models.MultiAccountGuild) (guildDs, guildTg, guildWa []string) {
	for _, channel := range mg.Channels {
		if channel == "DM" {
			//fmt.Println(channel)
		} else if strings.Contains(channel, "@") {
			guildWa = append(guildWa, channel)
		} else if !strings.HasPrefix(channel, "-100") {
			guildDs = append(guildDs, channel)
		} else if strings.HasPrefix(channel, "-100") {
			guildTg = append(guildTg, channel)
		}
	}
	return guildDs, guildTg, guildWa
}
func corpMemberDetectType(cm []models.CorpMember) []models.CorpMember {
	var members []models.CorpMember
	if len(cm) == 0 {
		return members
	}
	for _, m := range cm {
		m.TypeAccount = ""
		uId := m.UserId
		if strings.Contains(m.UserId, "/") {
			split := strings.Split(m.UserId, "/")
			uId = split[0]
		}
		if len(uId) < 12 {
			m.TypeAccount = "tg"
		} else if strings.Contains(uId, "@") {
			m.TypeAccount = "wa"
		} else if len(uId) > 12 && len(uId) < 24 {
			m.TypeAccount = "ds"
		}
		if m.TypeAccount == "" {
			fmt.Printf("TypeAccount len(%d) userid %s \n", len(m.UserId), m.UserId)
		}
		if m.TypeAccount != "" {
			m.Name = fmt.Sprintf("(%s) %s", strings.ToUpper(m.TypeAccount), m.Name)
			members = append(members, m)
		}
	}
	return members
}
func (s *Server) GetCorpDataMultiGuild(i *models.Identity, roleId string) *models.CorpData {
	c := models.CorpData{}
	c.Initialization()
	var members []models.CorpMember

	memberMulti, _ := s.multi.CorpMembersRead(i.MultiGuild.GId)
	fmt.Printf("memberMulti gid %+v %+v\n", i.MultiGuild.GId, memberMulti)
	if len(memberMulti) != 0 {
		c.AppendEveryone("ma")
		for _, member := range memberMulti {
			member.Name = fmt.Sprintf("(%s) %s", strings.ToUpper(member.TypeAccount), member.Name)
			members = append(members, member)
		}
	}

	cm, _ := s.db.CorpMembersRead(i.MultiGuild.GId.String())
	fmt.Printf("cm %s %+v\n", i.MultiGuild.GId.String(), cm)
	cm = corpMemberDetectType(cm)
	members = append(members, cm...)

	guildDs, guildTg, _ := getCorpsTypeId(i.MultiGuild)
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
	CheckRoleDs := func(role models.CorpRole) {
		for _, gId := range guildDs {
			fmt.Printf("members total Len %d\nCheckRoleDs gId %+v\n", len(members), gId)
			for _, member := range members {
				if member.TypeAccount == role.TypeRole {
					uid := member.UserId
					if strings.Contains(member.UserId, "/") {
						split := strings.Split(member.UserId, "/")
						uid = split[0]
					}
					if s.roles.CheckRoleDs(gId, uid, roleId) {
						fmt.Printf("CheckRoleDs OK roleId %s member %+v\n", roleId, member)
						c.Members = append(c.Members, member)
					}
				} else if member.TypeAccount == "ma" {
					uid := member.Multi.DiscordID
					if strings.Contains(member.UserId, "/") {
						split := strings.Split(member.UserId, "/")
						uid = split[0]
					}
					if s.roles.CheckRoleDs(gId, uid, roleId) {
						fmt.Printf("CheckRoleMa OK roleId %s member %+v\n", roleId, member)
						c.Members = append(c.Members, member)
					}
				}
			}
		}
	}
	CheckRoleTg := func(role models.CorpRole) {
		gId := i.MultiGuild.GId.String()
		for _, member := range members {
			if member.TypeAccount == role.TypeRole {
				exist := s.db.GuildRolesExistSubscribe(gId, role.Name[5:], member.UserId)
				if exist {
					c.Members = append(c.Members, member)
				}
			} else if member.TypeAccount == "ma" {
				if s.db.GuildRolesExistSubscribe(gId, role.Name[5:], member.Multi.TelegramID) {
					c.Members = append(c.Members, member)
				}
			}

		}
	}

	if len(guildDs) != 0 {
		c.AppendEveryone("ds")
		for _, gds := range guildDs {
			s.roles.LoadGuild(gds)
			roles, err := s.roles.ds.GetRoles(gds)
			if err != nil {
				s.log.ErrorErr(err)
			}
			appendRolesByType(roles, "ds")
		}
	}
	if len(guildTg) != 0 {
		c.AppendEveryone("tg")
		roles, err := s.db.GuildRolesRead(i.MultiGuild.GId.String())
		if err != nil {
			s.log.ErrorErr(err)
		}
		if len(roles) != 0 {
			appendRolesByType(roles, "tg")
		}
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
		if role.TypeRole == "" {
			fmt.Printf("role type not exist %+v\n", role)
		}

		if role.TypeRole == "ds" {
			CheckRoleDs(role)
		} else if role.TypeRole == "tg" {
			CheckRoleTg(role)
		} else {
			s.log.InfoStruct("role %+v\n", role)
		}
	}
	if len(c.Members) != 0 {
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
	}

	return &c
}
