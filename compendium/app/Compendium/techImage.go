package Compendium

import (
	"compendium/Compendium/imageGenerator"
	"context"
	"fmt"
	"regexp"
)

func (c *Compendium) techImage() {
	m := c.db.Temp.CorpMemberReadByUserIdByGuildId(context.TODO(), c.in.NameId, c.in.GuildId)
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
	re := regexp.MustCompile(`<@(\d{17,20})>`)
	retg := regexp.MustCompile(`@(\w+)`)

	matches := re.FindStringSubmatch(c.in.Text)
	matchestg := retg.FindStringSubmatch(c.in.Text)

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
