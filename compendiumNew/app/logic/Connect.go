package logic

import (
	"compendium/logic/generate"
	"compendium/models"
	"fmt"
	"strings"
	"time"
)

const urlLink = "https://mentalisit.github.io/HadesSpace/"

func (c *Hs) connect(m models.IncomingMessage) (conn bool) {
	if helperCommand(m.Text, "connect") ||
		helperCommand(m.Text, "подключить") ||
		helperCommand(m.Text, "підключити") {
		conn = true
	}
	if !conn {
		return false
	}
	text := fmt.Sprintf(c.getText(m, "CODE_FOR_CONNECT"), m.GuildName)
	mid1, err := c.sendDM(m, text)
	if err != nil && err.Error() == "forbidden" {
		c.sendChat(m, fmt.Sprintf(c.getText(m, "ERROR_SEND"), m.MentionName))
		return
	} else if err != nil {
		c.log.ErrorErr(err)
		c.log.InfoStruct("connect error ", err)
		return
	}
	c.sendChat(m, fmt.Sprintf(c.getText(m, "INSTRUCTIONS_SEND"), m.MentionName))

	newIdentify, cm := generate.GenerateIdentity(m)

	tokenOld, _ := c.listUser.ListUserGetToken(m.NameId, m.GuildId)
	if tokenOld != "" {
		prefixToken := m.Type + m.GuildId + "." + m.NameId
		if len(tokenOld) < 60 {
			newToken := prefixToken + generate.GenerateToken()
			err = c.listUser.ListUserUpdateToken(tokenOld, newToken)
			if err == nil {
				newIdentify.Token = newToken
			}
		} else {
			if strings.Contains(tokenOld, prefixToken) {
				newIdentify.Token = tokenOld
			} else {
				newIdentify.Token = prefixToken + tokenOld
			}
		}
	}

	code := c.generateCodeAndSave(newIdentify)

	err = c.guilds.GuildInsert(newIdentify.Guild)
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
	err = c.users.UsersInsert(newIdentify.User)
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
	err = c.corpMember.CorpMemberInsert(cm)
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
	err = c.listUser.ListUserInsert(newIdentify.Token, newIdentify.User.ID, newIdentify.Guild.ID)
	if err != nil {
		c.log.ErrorErr(err)
		return
	}

	mid2, errs1 := c.sendDM(m, code)
	if errs1 != nil {
		c.log.ErrorErr(errs1)
		return
	}

	urlLinkAdd := "compendiumTech?c=" + code + "&lang=" + m.Language + "&client=1"
	mid3, errs2 := c.sendDM(m, fmt.Sprintf(c.getText(m, "PLEASE_PASTE_CODE"), urlLink, urlLink+urlLinkAdd))
	if errs2 != nil {
		c.log.ErrorErr(errs2)
		return
	}

	go c.timerEditMessage(m, mid1, mid2, mid3)

	return
}

func (c *Hs) generateCodeAndSave(Identity models.Identity) string {
	segments := []string{generate.RandString(4), generate.RandString(4), generate.RandString(4)}

	m := models.Code{
		Code:      strings.Join(segments, "-"),
		Timestamp: time.Now().Unix(),
		Identity:  Identity,
	}

	go func() {
		err := c.db.DB.CodeInsert(m)
		if err != nil {
			c.log.ErrorErr(err)
			c.log.InfoStruct("CodeInsert", m)
		}
	}()

	return m.Code
}
func (c *Hs) timerEditMessage(m models.IncomingMessage, mid1, mid2, mid3 string) {
	time.Sleep(10 * time.Minute)

	token, _ := c.listUser.ListUserGetToken(m.NameId, m.GuildId)
	links := "https://mentalisit.github.io/HadesSpace/compendiumTech?secretToken=" + token
	text := fmt.Sprintf(c.getText(m, "SECRET_LINK"), links, m.GuildName)
	err := c.editMessage(m, m.DmChat, mid1, text, "MarkdownV2")
	if err != nil {
		c.log.ErrorErr(err)
		return
	}

	err = c.deleteMessage(m, m.DmChat, mid2)
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
	err = c.deleteMessage(m, m.DmChat, mid3)
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
}
