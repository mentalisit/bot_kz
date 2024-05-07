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
func (c *Hs) techImage() {
	mBytes, err := c.tech.TechGet(c.in.Name, c.in.NameId, c.in.GuildId)
	if err != nil {
		c.log.ErrorErr(err)
		c.sendChat("data not found")
		return
	}
	_, m := c.BytesToTechLevel(mBytes)
	userPic := imageGenerator.GenerateUser(
		c.in.Avatar,
		c.in.GuildAvatar,
		c.in.Name,
		c.in.GuildName,
		m)
	c.sendChatPic("", userPic)

}

func (c *Hs) techImageName() bool {
	after, _ := strings.CutPrefix(c.in.Text, "%")
	re := regexp.MustCompile(`^[тt] +<@(\d{17,20})> +[иi]$`)
	retg := regexp.MustCompile(`^[tт] +@(\w+) +[iи]$`)

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
			c.sendChat("данные не найдены")
			return false
		}
		if user == nil {
			c.sendChat("данные не найдены")
			return false
		}
		techBytes, err := c.tech.TechGet(user.Username, user.ID, c.in.GuildId)
		if err != nil {
			c.log.ErrorErr(err)
			c.sendChat("технические данные не найдены")
			return false
		}
		_, m := c.BytesToTechLevel(techBytes)
		userPic := imageGenerator.GenerateUser(
			user.AvatarURL,
			c.in.GuildAvatar,
			user.Username,
			c.in.GuildName,
			m)
		c.sendChatPic("", userPic)
		return true
	}
	return false
}
func (c *Hs) techImageNameAlt() bool {
	after, _ := strings.CutPrefix(c.in.Text, "%")
	re := regexp.MustCompile(`^[tт] +(\w+) +[iи]$`)

	matches := re.FindStringSubmatch(after)

	if len(matches) > 0 {
		userName := matches[1]
		techBytes, userID, err := c.tech.TechGetName(userName, c.in.GuildId)
		if err != nil && userID == "" {
			c.sendChat("данные не найдены")
			return false
		}

		user, err := c.users.UsersGetByUserId(userID)
		if err != nil {
			c.sendChat("данные не найдены")
			c.log.Info(err.Error())
			return false
		}
		if user == nil {
			c.sendChat("данные не найдены")
			return false
		}

		_, m := c.BytesToTechLevel(techBytes)
		userPic := imageGenerator.GenerateUser(
			user.AvatarURL,
			c.in.GuildAvatar,
			userName,
			c.in.GuildName,
			m)
		c.sendChatPic("", userPic)
		return true
	}
	return false
}
