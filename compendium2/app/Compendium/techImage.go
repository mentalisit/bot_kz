package Compendium

import (
	"compendium/Compendium/imageGenerator"
	"context"
	"fmt"
	"regexp"
	"strings"
)

func (c *Compendium) techImage() {
	m := c.db.Temp.CorpMemberReadByUserIdByGuildIdByName(context.TODO(), c.in.NameId, c.in.GuildId, c.in.Name)
	if m.Name == "" {
		c.sendChat("данные не найдены")
	} else {
		userPic := imageGenerator.GenerateUser(
			c.in.Avatar,
			c.in.GuildAvatar,
			m.Name,
			c.in.GuildName,
			m.Tech)
		c.sendChatPic("", userPic)
	}
}

func (c *Compendium) techImageName() bool {
	after, _ := strings.CutPrefix(c.in.Text, "%")
	re := regexp.MustCompile(`^[тt] +<@(\d{17,20})> +[иi]$`)
	retg := regexp.MustCompile(`^[tт] +@(\w+) +[iи]$`)

	matches := re.FindStringSubmatch(after)
	matchestg := retg.FindStringSubmatch(after)

	if len(matches) > 0 {
		userID := matches[1]
		fmt.Println("ID пользователя Discord:", userID)
		m := c.db.Temp.CorpMemberReadByUserIdByGuildId(context.TODO(), userID, c.in.GuildId)
		if m.Name == "" {
			c.sendChat("данные не найдены")
		} else {
			userPic := imageGenerator.GenerateUser(
				m.AvatarUrl,
				c.in.GuildAvatar,
				m.Name,
				c.in.GuildName,
				m.Tech)
			c.sendChatPic("", userPic)
		}
		return true
	} else if len(matchestg) > 0 {
		name := matchestg[1]
		fmt.Println("name пользователя :", name)
		m := c.db.Temp.CorpMemberReadByNameByGuildId(context.TODO(), name, c.in.GuildId)
		if m.Name == "" {
			c.sendChat("данные не найдены")
		} else {
			userPic := imageGenerator.GenerateUser(
				m.AvatarUrl,
				c.in.GuildAvatar,
				m.Name,
				c.in.GuildName,
				m.Tech)
			c.sendChatPic("", userPic)
		}
		return true
	}
	return false
}
