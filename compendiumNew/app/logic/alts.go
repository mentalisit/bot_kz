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
			return c.altsAdd(m, split[2])
		}
		if split[0] == "alts" && split[1] == "del" {
			return c.altsDel(m, split[2])
		}
	}
	return false
}

func (c *Hs) altsAdd(m models.IncomingMessage, altName string) bool {
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
			if alt == altName {
				c.sendChat(m, c.getText(m, "ALREADY_EXISTS"))
				return true
			}
		}
	}
	alts = append(alts, altName)
	if m.MAcc != nil {
		m.MAcc.Alts = alts
		m.MAcc, _ = c.db.V2.UpdateMultiAccountAlts(*m.MAcc)
	}

	if m.MultiAccount != nil {
		m.MultiAccount.Alts = alts
		_, err := c.db.Multi.UpdateMultiAccountAlts(*m.MultiAccount)
		if err != nil {
			c.log.ErrorErr(err)
		}
		_ = c.db.Multi.TechnologiesInsert(m.MultiAccount.UUID, altName, nil)
		c.log.Info(fmt.Sprintf("User %s alts new %+v", m.MultiAccount.Nickname, alts))
	} else {
		_ = c.tech.TechInsert(altName, m.NameId, nil)
		user.Alts = alts
		err := c.users.UsersUpdate(user)
		if err != nil {
			c.log.ErrorErr(err)
			return false
		}
		c.log.Info(fmt.Sprintf("User %s alts new %+v", user.Username, alts))
	}

	c.sendChat(m, fmt.Sprintf(c.getText(m, "ALTO_ADDED"), altName))
	_, _ = c.sendDM(m, fmt.Sprintf(c.getText(m, "LIST_ALTS"), alts))
	return true
}

func (c *Hs) altsDel(m models.IncomingMessage, altName string) bool {
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
			if alt != altName {
				alts = append(alts, alt)
			} else if alt == altName {
				c.sendChat(m, fmt.Sprintf(c.getText(m, "ALTO_REMOVED"), altName))
				if m.MultiAccount != nil {
					err = c.db.Multi.TechnologiesDelete(m.MultiAccount.UUID, altName)
					if err != nil {
						c.log.ErrorErr(err)
					}
				} else if u != nil {
					err = c.tech.TechDelete(altName, u.ID)
					if err != nil {
						c.log.ErrorErr(err)
					}
				}
			}
		}
		if m.MAcc != nil {
			m.MAcc.Alts = alts
			m.MAcc, _ = c.db.V2.UpdateMultiAccountAlts(*m.MAcc)
		}
		if m.MultiAccount != nil {
			m.MultiAccount.Alts = alts
			_, err = c.db.Multi.UpdateMultiAccountAlts(*m.MultiAccount)
			if err != nil {
				c.log.ErrorErr(err)
			}
			return true
		} else if u != nil {
			u.Alts = alts
			err = c.users.UsersUpdate(*u)
			if err != nil {
				c.log.ErrorErr(err)
				return false
			}
			err = c.corpMember.CorpMemberDeleteAlt(m.MGuild.GuildId(), m.NameId, altName)
			if err != nil {
				c.log.ErrorErr(err)
				return false
			}
			c.log.Info(fmt.Sprintf("User %s alts delete %+v", u.Username, altName))
			_, _ = c.sendDM(m, fmt.Sprintf(c.getText(m, "LIST_ALTS"), u.Alts))
			return true
		}
	}
	c.sendChat(m, fmt.Sprintf(c.getText(m, "NO_ALTOS_FOUND")))
	return true
}
