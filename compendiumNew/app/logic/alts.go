package logic

import (
	"compendium/models"
	"encoding/json"
	"fmt"
	"strings"
)

func (c *Hs) createAlt(m models.IncomingMessage) bool {
	after, _ := strings.CutPrefix(m.Text, "%")
	split := strings.Split(after, " ")
	if len(split) == 3 {
		if split[0] == "alts" && split[1] == "add" {
			u, err := c.users.UsersGetByUserId(m.NameId)
			if err != nil {
				c.log.ErrorErr(err)
				return false
			}
			alts := u.Alts
			if len(alts) > 0 {
				for _, alt := range alts {
					if alt == split[2] {
						c.sendChat(m, c.getText(m, "ALREADY_EXISTS"))
						return true
					}
				}
			}
			alts = append(alts, split[2])
			tech := make(map[int]models.TechLevel)
			tech[701] = models.TechLevel{
				Ts:    0,
				Level: 0,
			}
			techBytes, _ := json.Marshal(tech)

			_ = c.tech.TechInsert(split[2], m.NameId, m.GuildId, techBytes)
			u.Alts = alts
			err = c.users.UsersUpdate(*u)
			if err != nil {
				c.log.ErrorErr(err)
				return false
			}

			c.log.Info(fmt.Sprintf("User %s alts new %+v", u.Username, alts))
			c.sendChat(m, fmt.Sprintf(c.getText(m, "ALTO_ADDED"), split[2]))
			_ = c.sendDM(m, fmt.Sprintf(c.getText(m, "LIST_ALTS"), alts))
			return true
		}
		if split[0] == "alts" && split[1] == "del" {
			u, err := c.users.UsersGetByUserId(m.NameId)
			if err != nil {
				c.log.ErrorErr(err)
				return false
			}
			if len(u.Alts) > 0 {
				var alts []string
				for _, alt := range u.Alts {
					if alt != split[2] {
						alts = append(alts, alt)
					} else if alt == split[2] {
						c.sendChat(m, fmt.Sprintf(c.getText(m, "ALTO_REMOVED"), split[2]))
						err = c.tech.TechDelete(split[2], u.ID, m.GuildId)
						if err != nil {
							c.log.ErrorErr(err)
						}
					}
				}
				u.Alts = alts
				err = c.users.UsersUpdate(*u)
				if err != nil {
					c.log.ErrorErr(err)
					return false
				}
			} else {
				c.sendChat(m, fmt.Sprintf("NO_ALTOS_FOUND"))
			}

			c.log.Info(fmt.Sprintf("User %s alts delete %+v", u.Username, split[2]))
			_ = c.sendDM(m, fmt.Sprintf(c.getText(m, "LIST_ALTS"), u.Alts))
			return true
		}
	}
	return false
}
