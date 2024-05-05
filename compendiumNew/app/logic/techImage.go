package logic

import (
	"compendium/logic/imageGenerator"
	"compendium/models"
	"encoding/json"
	"regexp"
	"strings"
)

func (c *Hs) techImage() {
	mBytes, err := c.tech.TechGet(c.in.Name, c.in.NameId, c.in.GuildId)
	if err != nil {
		c.log.ErrorErr(err)
	}
	m := make(map[int][2]int)
	err = json.Unmarshal(mBytes, &m)
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
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

	if len(matches) > 0 && len(matchestg) > 0 {
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
			c.log.ErrorErr(err)
			c.sendChat("данные не найдены")
			return false
		}
		if user == nil {
			return false
		}
		techBytes, err := c.tech.TechGet(user.Username, user.ID, c.in.GuildId)
		if err != nil {
			c.log.ErrorErr(err)
			c.sendChat("технические данные не найдены")
			return false
		}
		tech := make(map[int][2]int)
		err = json.Unmarshal(techBytes, &tech)
		if err != nil {
			c.log.ErrorErr(err)
			return false
		}
		userPic := imageGenerator.GenerateUser(
			user.AvatarURL,
			c.in.GuildAvatar,
			user.Username,
			c.in.GuildName,
			tech)
		c.sendChatPic("", userPic)
		return true
	}
	return false
}
