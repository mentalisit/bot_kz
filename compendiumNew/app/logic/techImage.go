package logic

import (
	"compendium/logic/imageGenerator"
	"compendium/models"
	"encoding/json"
	"regexp"
	"strings"
)

func (c *Hs) BytesToTechLevel(b []byte) (map[int]models.TechLevel, models.TechLevelArray) {
	var m map[int]models.TechLevel
	m = make(map[int]models.TechLevel)
	err := json.Unmarshal(b, &m)
	if err != nil {
		c.log.ErrorErr(err)
		return nil, nil
	}
	var mi = make(models.TechLevelArray)
	for i, le := range m {
		mi[i] = [2]int{le.Level}
	}
	return m, mi
}
func (c *Hs) techImage(m models.IncomingMessage) {
	mBytes, err := c.tech.TechGet(m.Name, m.NameId, m.GuildId)
	if err != nil {
		c.log.ErrorErr(err)
		c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
		return
	}
	_, mt := c.BytesToTechLevel(mBytes)
	userPic := imageGenerator.GenerateUser(
		m.Avatar,
		m.GuildAvatar,
		m.Name,
		m.GuildName,
		mt)
	c.sendChatPic(m, "", userPic)

}

func (c *Hs) techImageName(m models.IncomingMessage) bool {
	after, _ := strings.CutPrefix(m.Text, "%")
	re := regexp.MustCompile(`^[тt] +<@(\d{17,20})> +[иiі]$`)
	retg := regexp.MustCompile(`^[tт] +@(\w+) +[iиі]$`)

	matches := re.FindStringSubmatch(after)
	matchestg := retg.FindStringSubmatch(after)

	if len(matches) > 0 || len(matchestg) > 0 {
		var user *models.User
		var err error
		if len(matches) > 0 {
			userID := matches[1]
			user, err = c.users.UsersGetByUserId(userID)
		} else if len(matchestg) > 0 {
			name := matchestg[1]
			user, err = c.users.UsersGetByUserName(name)
		}
		if err != nil {
			c.log.Info(err.Error())
			c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
			return false
		}
		if user == nil {
			c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
			return false
		}
		techBytes, err := c.tech.TechGet(user.Username, user.ID, m.GuildId)
		if err != nil {
			c.log.ErrorErr(err)
			c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
			return false
		}
		_, mt := c.BytesToTechLevel(techBytes)
		userPic := imageGenerator.GenerateUser(
			user.AvatarURL,
			m.GuildAvatar,
			user.Username,
			m.GuildName,
			mt)
		c.sendChatPic(m, "", userPic)
		return true
	}
	return false
}
func (c *Hs) techImageNameAlt(m models.IncomingMessage) bool {
	after, _ := strings.CutPrefix(m.Text, "%")
	re := regexp.MustCompile(`^[tт] +(\w+) +[iиі]$`)

	matches := re.FindStringSubmatch(after)

	if len(matches) > 0 {
		userName := matches[1]
		techBytes, userID, err := c.tech.TechGetName(userName, m.GuildId)
		if err != nil && userID == "" {
			c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
			return false
		}

		user, err := c.users.UsersGetByUserId(userID)
		if err != nil {
			c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
			c.log.Info(err.Error())
			return false
		}
		if user == nil {
			c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
			return false
		}

		_, mt := c.BytesToTechLevel(techBytes)
		userPic := imageGenerator.GenerateUser(
			user.AvatarURL,
			m.GuildAvatar,
			userName,
			m.GuildName,
			mt)
		c.sendChatPic(m, "", userPic)
		return true
	}
	return false
}
