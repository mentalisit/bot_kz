package logic

import (
	"compendium/models"
	"fmt"
	"strings"
	"time"
)

func (c *Hs) logic(m models.IncomingMessage) {
	c.PrintGoroutine()

	if strings.Contains(m.Type, "DM") && !strings.HasPrefix(m.Text, "%") {
		return
	}
	if m.MultiGuild == nil {
		c.log.InfoStruct("m.MultiGuild==nil ", m)
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
		fmt.Printf("logic MultiAccount: %+v %+v\n", time.Now().Format(time.RFC3339), m.MultiAccount)
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
		c.sendChat(m, c.getText(m, "ErrorRequest"))
		fmt.Printf("else Corp:%s %+v\n", m.MultiGuild.GuildName, m)
	}
}
