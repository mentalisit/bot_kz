package logic

import (
	"compendium/models"
	"fmt"
	"time"
)

func (c *Hs) logic(m models.IncomingMessage) {
	c.PrintGoroutine()

	if m.GuildId != "" {
		guild, _ := c.db.Multi.GuildGet(m.GuildId)
		if guild == nil {
			err := c.db.Multi.GuildInsert(models.MultiAccountGuild{
				GuildName: m.GuildName,
				Channels:  []string{m.GuildId},
				AvatarUrl: m.GuildAvatar,
			})
			if err != nil {
				c.log.ErrorErr(err)
			}
			guild, _ = c.db.Multi.GuildGet(m.GuildId)
		} else if guild.AvatarUrl != m.GuildAvatar {
			guild.AvatarUrl = m.GuildAvatar
			err := c.db.Multi.GuildUpdateAvatar(*guild)
			if err != nil {
				c.log.ErrorErr(err)
			}
		}
		m.MultiGuild = guild
	}

	multiAccount, _ := c.db.Multi.FindMultiAccountByUserId(m.NameId)
	if multiAccount != nil && multiAccount.TelegramID != "" && multiAccount.DiscordID != "" {
		if m.Avatar != "" {
			if multiAccount.AvatarURL != m.Avatar {
				multiAccount.AvatarURL = m.Avatar
				_, _ = c.db.Multi.UpdateMultiAccountAvatarUrl(*multiAccount)
			}
		}
		m.MultiAccount = multiAccount
	}

	fmt.Printf("logic: %+v %+v\n", time.Now().Format(time.RFC3339), m)
	if m.MultiAccount != nil {
		fmt.Printf("logic: %+v %+v\n", time.Now().Format(time.RFC3339), m.MultiAccount)
	}
	if c.connect(m) {
	} else if c.multiConnect(m) {
	} else if c.Help(m) {
	} else if c.techImage(m) {
	} else if c.techImageName(m) {
	} else if c.techImageNameAlt(m) {
	} else if c.logicRoles(m) {
	} else if c.createAlt(m) {
	} else if c.wskill(m) {
	} else if c.TzTime(m) {
	} else if c.setGameName(m) {
	} else if c.removeMember(m) {
	} else {
		c.log.Info(fmt.Sprintf("else %+v\n", m))
	}
}
