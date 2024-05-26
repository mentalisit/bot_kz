package logic

import (
	"compendium/logic/generate"
	"compendium/models"
	"fmt"
	"strings"
	"time"
)

func (c *Hs) connect(m models.IncomingMessage) {
	_, err := c.sendDM(m, fmt.Sprintf(c.getText(m, "CODE_FOR_CONNECT"), m.GuildName))
	if err != nil && err.Error() == "forbidden" {
		c.sendChat(m, fmt.Sprintf(c.getText(m, "ERROR_SEND"), m.MentionName))
		return
	} else if err != nil {
		c.log.ErrorErr(err)
	}
	c.sendChat(m, fmt.Sprintf(c.getText(m, "INSTRUCTIONS_SEND"), m.MentionName))

	newIdentify, cm := generate.GenerateIdentity(m)

	tokenOld, _ := c.listUser.ListUserGetToken(m.NameId, m.GuildId)
	if tokenOld != "" {
		newIdentify.Token = tokenOld
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
	mid1, errs1 := c.sendDM(m, code)
	if errs1 != nil {
		c.log.ErrorErr(errs1)
		return
	}
	urlLink := "https://mentalisit.github.io/HadesSpace/"
	urlLinkAdd := "compendiumTech?c=" + code + "&lang=" + m.Language + "&client=1"
	mid2, errs2 := c.sendDM(m, fmt.Sprintf(c.getText(m, "PLEASE_PASTE_CODE"), urlLink, urlLink+urlLinkAdd))
	if errs2 != nil {
		c.log.ErrorErr(errs2)
		return
	}
	go func() {
		time.Sleep(10 * time.Minute)
		err = c.editMessage(m, m.DmChat, mid1, c.getText(m, "CODE_OUTDATED"))
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
		err = c.deleteMessage(m, m.DmChat, mid2)
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
	}()
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
