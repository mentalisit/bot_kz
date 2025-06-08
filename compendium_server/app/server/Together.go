package server

import (
	"compendium_s/models"
	"sort"
	"strings"
)

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
