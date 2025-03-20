package server

import (
	"compendium_s/models"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"time"
)

func (s *Server) GetTokenIdentity(token string) *models.Identity {
	userid, guildid, err := s.db.ListUserGetUserIdAndGuildId(token)
	if err != nil {
		token2 := s.db.ListUserGetByMatch(token)
		userid, guildid, err = s.db.ListUserGetUserIdAndGuildId(token2)
		if err != nil {
			s.log.Info("get user by token: " + token + " " + err.Error())
			return nil
		}
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
	together := s.GetCorpDataIfTogether(i, roleId)
	if together != nil {
		return together
	}
	if i.Guild.Type == "ds" {
		s.roles.LoadGuild(i.Guild.ID)
	}
	c := models.CorpData{}
	c.Members = []models.CorpMember{}

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
func (s *Server) getRoles(i models.Guild) []models.CorpRole {
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
func GenerateToken() string {
	// Вычисляем необходимый размер байт для указанной длины токена
	tokenBytes := make([]byte, 174)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return ""
	}

	// Кодируем байты в строку base64
	token := base64.URLEncoding.EncodeToString(tokenBytes)
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

func (s *Server) GetCorpDataIfTogether(i *models.Identity, roleId string) *models.CorpData {
	compatible, err := s.listOfCompatible(&i.Guild)
	if err != nil {
		s.log.ErrorErr(err)
		return nil
	}
	if compatible != nil {
		var d models.CorpData
		var Ids, Itg models.Guild
		if i.Guild.Type == "ds" {
			Ids = i.Guild
			Itg = *compatible
		} else {
			Itg = i.Guild
			Ids = *compatible
		}

		s.roles.LoadGuild(Ids.ID)

		d.Members = []models.CorpMember{}
		var memberDs []models.CorpMember
		var memberTg []models.CorpMember

		d.Roles = append(d.Roles, models.CorpRole{
			Id:   "",
			Name: "@everyone",
		})

		rolesTg := s.getRoles(Itg)
		for _, roles := range rolesTg {
			if roles.Name == "@everyone" {
				d.Roles = append(d.Roles, models.CorpRole{
					Id:   "tg",
					Name: "(TG) " + roles.Name,
				})
				continue
			}
			d.Roles = append(d.Roles, models.CorpRole{
				Id:   roles.Id,
				Name: "(TG) " + roles.Name,
			})
		}
		rolesDs := s.getRoles(Ids)
		for _, roles := range rolesDs {
			if roles.Name == "@everyone" {
				d.Roles = append(d.Roles, models.CorpRole{
					Id:   "ds",
					Name: "(DS) " + roles.Name,
				})
				continue
			}
			d.Roles = append(d.Roles, models.CorpRole{
				Id:   roles.Id,
				Name: "(DS) " + roles.Name,
			})
		}

		cmDs, _ := s.db.CorpMembersRead(Ids.ID)
		for _, m := range cmDs {
			m.Name = "(DS) " + m.Name
			memberDs = append(memberDs, m)
		}
		cmTg, _ := s.db.CorpMembersRead(Itg.ID)
		for _, m := range cmTg {
			m.Name = "(TG) " + m.Name
			memberTg = append(memberTg, m)
		}

		if roleId == "" {
			d.Members = append(d.Members, memberDs...)
			d.Members = append(d.Members, memberTg...)
		} else if roleId == "tg" {
			d.Members = append(d.Members, memberTg...)
		} else if roleId == "ds" {
			d.Members = append(d.Members, memberDs...)
		} else {
			var role models.CorpRole
			var tip string
			for _, roles := range rolesDs {
				if roleId == roles.Id {
					role = roles
					tip = "ds"
				}
			}
			for _, roles := range rolesTg {
				if roleId == roles.Id {
					role = roles
					tip = "tg"
				}
			}

			if tip == "ds" {
				for _, member := range memberDs {
					uid := member.UserId
					if strings.Contains(member.UserId, "/") {
						split := strings.Split(member.UserId, "/")
						uid = split[0]
					}
					if s.roles.CheckRoleDs(Ids.ID, uid, roleId) {
						d.Members = append(d.Members, member)
					}
				}
			} else if tip == "tg" {
				for _, member := range memberTg {
					if s.db.GuildRolesExistSubscribe(Itg.ID, role.Name, member.UserId) {
						d.Members = append(d.Members, member)
					}
				}
			}
		}
		sort.Slice(d.Members, func(i, j int) bool {
			// Проверка, чтобы индекс не выходил за пределы строки
			nameI := d.Members[i].Name
			nameJ := d.Members[j].Name

			// Игнорируем первые пять символов, если длина имени больше или равна пяти
			if len(nameI) >= 5 {
				nameI = nameI[5:]
			}
			if len(nameJ) >= 5 {
				nameJ = nameJ[5:]
			}

			return nameI < nameJ
		})
		return &d
	}
	return nil
}
func (s *Server) listOfCompatible(g *models.Guild) (*models.Guild, error) {
	if g.Type == "tg" {
		if g.Name == "Свободный Флот" && g.ID == "-1002125982067" {
			return s.db.GuildGet("1062696191575457892")
		} else if g.Name == "HS СССР" && g.ID == "-1002467616555" {
			return s.db.GuildGet("632245873769971732")
		} else if g.Name == "ТКЗ" && g.ID == "-1001697997137" { //Корпорация русь
			return s.db.GuildGet("716771579278917702")
		} else if g.Name == "IX Легион" && g.ID == "-1002298028181" {
			return s.db.GuildGet("398761209022644224")
		} else if g.Name == "DV NEBULA" && g.ID == "-1002014251679" {
			return s.db.GuildGet("656495834195558402")
		} else if g.ID == "-1002421683868" { //Hades Star: Eden.
			return s.db.GuildGet("1347552010383261797")
		}

	} else if g.Type == "ds" {
		if g.Name == "Свободный Флот" && g.ID == "1062696191575457892" {
			return s.db.GuildGet("-1002125982067")
		} else if g.Name == "СССР  (HS)" && g.ID == "632245873769971732" {
			return s.db.GuildGet("-1002467616555")
		} else if g.Name == "Корпорация  \"РУСЬ\"" && g.ID == "716771579278917702" {
			return s.db.GuildGet("-1001697997137")
		} else if g.Name == "IX Легион" && g.ID == "398761209022644224" {
			return s.db.GuildGet("-1002298028181")
		} else if g.Name == "ГОРИЗОНТ" && g.ID == "656495834195558402" {
			return s.db.GuildGet("-1002014251679")
		} else if g.ID == "1347552010383261797" {
			return s.db.GuildGet("-1002421683868")
		}

	}
	return nil, nil
}
