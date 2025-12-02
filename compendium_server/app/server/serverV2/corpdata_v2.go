package serverV2

import (
	"compendium_s/models"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//
//import (
//	"compendium_s/models"
//	"fmt"
//
//	"github.com/google/uuid"
//)
//
//func (s *ServerV2) GetCorpData(i *models.IdentityV2, roleId string) *models.CorpData {
//	return s.GetCorpDataInternal(i, roleId, true)
//}

func (s *ServerV2) GetCorpDataInternal(i *models.IdentityV2, roleId string) *models.CorpDataV2 {
	c := models.CorpDataV2{}
	c.Initialization()
	var members []models.CorpMemberV2

	mGuild, err := s.db.GuildGet(i.GetGuildUUID())
	if err != nil {
		s.log.ErrorErr(err)
		return &c
	}
	members, _ = s.db.CorpMembersReadMulti(&mGuild.GId)

	guildDs, guildTg, _ := getCorpsTypeId(mGuild)

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
				if member.GetType() == role.TypeRole {
					uid := member.Multi.DiscordID
					if strings.Contains(member.Multi.DiscordID, "/") {
						split := strings.Split(member.Multi.DiscordID, "/")
						uid = split[0]
					}
					if s.roles.CheckRoleDs(gId, uid, roleId) {
						fmt.Printf("CheckRoleDs OK roleId %s member %+v\n", roleId, member)
						c.Members = append(c.Members, member)
					}
				}
			}
		}
	}
	CheckRoleTg := func(role models.CorpRole) {
		for _, member := range members {
			if member.GetType() == role.TypeRole {
				exist, _ := s.db.IsUserSubscribedToRole(member.Multi.GetTelegramChatId(), role.GetRoleId())
				if exist {
					c.Members = append(c.Members, member)
				}
			}
		}
	}

	if len(guildTg) != 0 {
		c.AppendEveryone("tg")
		for _, tgId := range guildTg {
			gTg, _ := strconv.ParseInt(tgId, 10, 64)
			roles, err := s.db.GetChatsRoles(gTg)
			if err != nil {
				s.log.ErrorErr(err)
			}
			if len(roles) != 0 {
				appendRolesByType(roles, "tg")
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

	if roleId == "" {
		c.Members = members
	} else if roleId == "tg" || roleId == "ds" || roleId == "ma" {
		for _, member := range members {
			if member.GetType() == roleId {
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
			s.log.Info(fmt.Sprintf("role type not exist %+v\n", role))
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
