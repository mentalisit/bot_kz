package Compendium

import (
	"compendium/models"
	"context"
	"fmt"
	"strings"
)

func (c *Compendium) logic(m models.IncomingMessage) {
	fmt.Printf("logic %+v\n", m)
	cutPrefix, found := strings.CutPrefix(m.Text, "%")

	if found {
		c.in = m
		switch cutPrefix {
		case "help", "помощь", "справка":
			c.help()
		case "подключить", "connect":
			c.connect()
		case "т и", "t i":
			c.techImage()
		default:

			c.regular(cutPrefix)
		}
	}
}
func (c *Compendium) regular(text string) {
	if c.techImageName() {
	} else if c.logicRoles() {
	} else if c.createAlt() {
	} else {
		c.log.Info(text)
	}
}
func (c *Compendium) help() {
	c.sendChat("на текуший момент доступны команды:\n" +
		"'%help' или '%помощь', '%справка' для получения текущей справки \n" +
		"'%connect' или '%коннект' для подключения приложения\n" +
		"'%t i' или '%т и' для получения изображения с вашими модулями\n" +
		"'%t @Name i' или '%т @имя и' для получения изображения с модулями другого игрока\n")
}

func (c *Compendium) createAlt() bool {
	after, _ := strings.CutPrefix(c.in.Text, "%")
	split := strings.Split(after, " ")
	if len(split) == 3 {
		if split[0] == "alts" && split[1] == "add" {
			u, err := c.db.Temp.UserReadByUserIdByUsername(context.Background(), c.in.NameId, c.in.Name)
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
			c.db.Temp.UserUpdateAlts(context.Background(), u.Username, u.ID, alts)
			c.log.Info(fmt.Sprintf("User %s alts new %+v", u.Username, alts))
			c.sendChat("alto added " + split[2])
			_ = c.sendDM(fmt.Sprintf("List of your alts %+v", alts))
			return true
		}
	}
	return false
}
