package logic

import (
	"compendium/models"
	"fmt"
	"strings"
)

func (c *Hs) createAlt(m models.IncomingMessage) bool {
	after, _ := strings.CutPrefix(m.Text, "%")
	split := strings.Split(after, " ")
	if len(split) == 3 {
		if split[0] == "alts" && split[1] == "add" {
			var alts []string
			var user models.User
			if m.MultiAccount != nil {
				alts = m.MultiAccount.Alts
			} else {
				u, err := c.users.UsersGetByUserId(m.NameId)
				if err != nil {
					c.log.ErrorErr(err)
					return false
				}
				alts = u.Alts
				user = *u
			}

			if len(alts) > 0 {
				for _, alt := range alts {
					if alt == split[2] {
						c.sendChat(m, c.getText(m, "ALREADY_EXISTS"))
						return true
					}
				}
			}
			alts = append(alts, split[2])

			if m.MultiAccount != nil {
				m.MultiAccount.Alts = alts
				_, err := c.db.Multi.UpdateMultiAccountAlts(*m.MultiAccount)
				if err != nil {
					c.log.ErrorErr(err)
				}
				_ = c.db.Multi.TechnologiesInsert(m.MultiAccount.UUID, split[2], nil)
				c.log.Info(fmt.Sprintf("User %s alts new %+v", m.MultiAccount.Nickname, alts))
			} else {
				_ = c.tech.TechInsert(split[2], m.NameId, m.GuildId, nil)
				user.Alts = alts
				err := c.users.UsersUpdate(user)
				if err != nil {
					c.log.ErrorErr(err)
					return false
				}
				c.log.Info(fmt.Sprintf("User %s alts new %+v", user.Username, alts))
			}

			c.sendChat(m, fmt.Sprintf(c.getText(m, "ALTO_ADDED"), split[2]))
			_, _ = c.sendDM(m, fmt.Sprintf(c.getText(m, "LIST_ALTS"), alts))
			return true
		}
		if split[0] == "alts" && split[1] == "del" {
			var u *models.User
			var err error
			var uAlts []string
			if m.MultiAccount != nil {
				uAlts = m.MultiAccount.Alts
			} else {
				u, err = c.users.UsersGetByUserId(m.NameId)
				if err != nil {
					c.log.ErrorErr(err)
					return false
				}
				uAlts = u.Alts
			}

			if len(uAlts) > 0 {
				var alts []string
				for _, alt := range uAlts {
					if alt != split[2] {
						alts = append(alts, alt)
					} else if alt == split[2] {
						c.sendChat(m, fmt.Sprintf(c.getText(m, "ALTO_REMOVED"), split[2]))
						if m.MultiAccount != nil {
							err = c.db.Multi.TechnologiesDelete(m.MultiAccount.UUID, split[2])
							if err != nil {
								c.log.ErrorErr(err)
							}
						} else {
							err = c.tech.TechDelete(split[2], u.ID, m.GuildId)
							if err != nil {
								c.log.ErrorErr(err)
							}
						}

					}
				}
				if m.MultiAccount != nil {
					m.MultiAccount.Alts = alts
					_, err = c.db.Multi.UpdateMultiAccountAlts(*m.MultiAccount)
					if err != nil {
						c.log.ErrorErr(err)
					}
					//todo удалить корп мембера безопасно
				} else {
					u.Alts = alts
					if u == nil {
						c.log.Error("u == nil")
						return false
					}
					err = c.users.UsersUpdate(*u)
					if err != nil {
						c.log.ErrorErr(err)
						return false
					}
					err = c.corpMember.CorpMemberDeleteAlt(m.GuildId, m.NameId, split[2])
					if err != nil {
						c.log.ErrorErr(err)
						return false
					}
					c.log.Info(fmt.Sprintf("User %s alts delete %+v", u.Username, split[2]))
					_, _ = c.sendDM(m, fmt.Sprintf(c.getText(m, "LIST_ALTS"), u.Alts))
					return true
				}
			} else {
				c.sendChat(m, fmt.Sprintf(c.getText(m, "NO_ALTOS_FOUND")))
				return true
			}
		}
	}
	return false
}
