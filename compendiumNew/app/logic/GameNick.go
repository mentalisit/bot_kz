package logic

import (
	"compendium/models"
	"fmt"
	"regexp"
	"strings"
)

// multi OK
func (c *Hs) setGameName(m models.IncomingMessage) bool {
	after, found := strings.CutPrefix(m.Text, "%")
	if found {
		re := regexp.MustCompile(`^(nick|ник|нік) ("([^"]+)"|(\S+))$`)
		matches := re.FindStringSubmatch(after)
		if len(matches) > 0 {
			gameName := ""
			if matches[3] != "" {
				gameName = matches[3]
			} else if matches[3] == "" && matches[2] != "" {
				gameName = matches[2]
			}
			if m.MultiAccount != nil {
				_ = c.db.Multi.TechnologiesUpdateUsername(m.MultiAccount.UUID, m.MultiAccount.Nickname, gameName)
				m.MultiAccount.Nickname = gameName
				_, err := c.db.Multi.UpdateMultiAccountNickname(*m.MultiAccount)
				if err != nil {
					c.log.ErrorErr(err)
				}
			} else {
				user, err := c.users.UsersGetByUserId(m.NameId)
				if err != nil {
					c.log.ErrorErr(err)
					text := fmt.Sprintf(c.getText(m, "YOU_ARE_NOT_FOUND"), m.MentionName)
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
			}
			text := fmt.Sprintf(c.getText(m, "GAME_NAME_SET"), m.MentionName, gameName)
			c.sendChat(m, text)
			return true
		} else {
			split := strings.Split(after, " ")
			if split[0] == "ник" || split[0] == "nick" || split[0] == "нік" {
				c.sendChat(m, fmt.Sprintf(c.getText(m, "HELP_NICKNAME"), m.MentionName))
				return true

			}
		}
	}
	return false
}
