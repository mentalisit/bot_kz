package logic

import (
	"fmt"
	"strings"
)

func (c *Hs) createAlt() bool {
	after, _ := strings.CutPrefix(c.in.Text, "%")
	split := strings.Split(after, " ")
	if len(split) == 3 {
		if split[0] == "alts" && split[1] == "add" {
			u, err := c.users.UsersGetByUserId(c.in.NameId)
			if err != nil {
				c.log.ErrorErr(err)
				return false
			}
			alts := u.Alts
			if len(alts) > 0 {
				for _, alt := range alts {
					if alt == split[2] {
						c.sendChat("already exists")
						return true
					}
				}
			}
			alts = append(alts, split[2])
			u.Alts = alts
			err = c.users.UsersUpdate(*u)
			if err != nil {
				c.log.ErrorErr(err)
				return false
			}

			c.log.Info(fmt.Sprintf("User %s alts new %+v", u.Username, alts))
			c.sendChat("alto added " + split[2])
			_ = c.sendDM(fmt.Sprintf("List of your alts %+v", alts))
			return true
		}
		if split[0] == "alts" && split[1] == "del" {
			u, err := c.users.UsersGetByUserId(c.in.NameId)
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
						c.sendChat("alto removed " + split[2])
						err = c.tech.TechDelete(u.Username, u.ID, c.in.GuildId)
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
				c.sendChat("no altos found")
			}

			c.log.Info(fmt.Sprintf("User %s alts delete %+v", u.Username, split[2]))
			_ = c.sendDM(fmt.Sprintf("List of your alts %+v", u.Alts))
			return true
		}
	}
	return false
}
