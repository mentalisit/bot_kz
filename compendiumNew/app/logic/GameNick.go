package logic

import (
	"compendium/models"
	"fmt"
	"regexp"
	"strings"
)

func (c *Hs) setGameName(m models.IncomingMessage) bool {
	after, found := strings.CutPrefix(m.Text, "%")
	if found {
		re := regexp.MustCompile(`^(nick|ник) ("([^"]+)"|(\S+))$`)
		matches := re.FindStringSubmatch(after)
		if len(matches) > 0 {
			gameName := ""
			if matches[3] != "" {
				gameName = matches[3]
			} else if matches[3] == "" && matches[2] != "" {
				gameName = matches[2]
			}
			user, err := c.users.UsersGetByUserId(m.NameId)
			if err != nil {
				c.log.ErrorErr(err)
				text := fmt.Sprintf("%s you are not found in the database, please send %%connect", m.MentionName)
				c.sendChat(m, text)
				return true
			}
			user.GameName = gameName
			err = c.users.UsersUpdate(*user)
			if err != nil {
				c.log.ErrorErr(err)
				c.log.InfoStruct("UsersUpdate", user)
				return false
			}
			text := fmt.Sprintf("%s, game name set to '%s'", m.MentionName, gameName)
			c.sendChat(m, text)
			return true
		} else {
			split := strings.Split(after, " ")
			if split[0] == "ник" || split[0] == "nick" {
				helpText := "The %nick command is used to set the name,\n" +
					"`%nick name` if the name does not contain spaces\n" +
					"or\n" +
					"`%nick \"my name\"` if the name consists of several words\n" +
					"example\n" +
					"`%nick Vasya`\n" +
					"`%nick \"Vasya Ivanov\"`"
				c.sendChat(m, fmt.Sprintf("%s, %s", m.MentionName, helpText))
				return true

			}
		}
	}
	return false
}
