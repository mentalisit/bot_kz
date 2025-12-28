package logic

import (
	"compendium/logic/generate"
	"compendium/models"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const urlLink = "https://mentalisit.github.io/HadesSpace/"

func (c *Hs) connect(m models.IncomingMessage) (conn bool) {
	if helperCommand(m.Text, "connect") || helperCommand(m.Text, "подключить") || helperCommand(m.Text, "підключити") {
		conn = true
	}
	if !conn {
		return false
	}

	text := fmt.Sprintf(c.getText(m, "CODE_FOR_CONNECT"), m.MGuild.GuildName)
	mid1, err := c.sendDM(m, text)
	if err != nil && err.Error() == "forbidden" {
		if c.checkMoron(m) {
			c.log.InfoStruct("moron", m)
		} else {
			c.sendChat(m, fmt.Sprintf(c.getText(m, "ERROR_SEND"), m.MentionName))
			return
		}

	} else if err != nil {
		c.log.ErrorErr(err)
		c.log.InfoStruct("connect error ", err)
		return
	}

	c.sendChat(m, fmt.Sprintf(c.getText(m, "INSTRUCTIONS_SEND"), m.MentionName))

	newIdentify := generate.GenerateIdentity(m)

	code := c.generateCodeAndSave(newIdentify)
	cmNew := false
	oldCM, _ := c.db.V2.CorpMemberByUId(newIdentify.MAccount.UUID)
	if oldCM != nil {
		if !oldCM.Exist(newIdentify.MGuild.GId) {
			oldCM.GuildIds = append(oldCM.GuildIds, newIdentify.MGuild.GId)
			err = c.db.V2.CorpMemberUpdate(*oldCM)
			if err != nil {
				c.log.ErrorErr(err)
			}
			cmNew = true
		}
	} else {
		err = c.db.V2.CorpMemberInsert(models.MultiAccountCorpMember{
			Uid:      m.MAcc.UUID,
			GuildIds: []uuid.UUID{m.MGuild.GId},
		})
		if err != nil {
			c.log.ErrorErr(err)
		}
		cmNew = true
	}

	mid2, errs1 := c.sendDM(m, code)
	if errs1 != nil {
		c.log.ErrorErr(errs1)
		return
	}
	if cmNew {
		urlLinkAdd := "compendiumTech?c=" + code + "&lang=" + m.Language + "&client=1"
		mid3, errs2 := c.sendDM(m, fmt.Sprintf(c.getText(m, "PLEASE_PASTE_CODE"), urlLink, urlLink+urlLinkAdd))
		if errs2 != nil {
			c.log.ErrorErr(errs2)
			return
		}

		go c.timerEditMessage(m, mid1, mid2, mid3, newIdentify.Token)
	} else {
		sendDM, _ := c.sendDM(m, "Подготавливаю секретную ссылку")
		go func() {
			time.Sleep(10 * time.Second)
			_ = c.deleteMessage(m, m.DmChat, mid1)
			_ = c.deleteMessage(m, m.DmChat, mid2)

			links := "https://mentalisit.github.io/HadesSpace/compendiumTech?secretToken=" + newIdentify.Token
			text = fmt.Sprintf(c.getText(m, "SECRET_LINK"), links, m.MGuild.GuildName)
			err = c.editMessage(m, m.DmChat, sendDM, text, "MarkdownV2")
			if err != nil {
				c.log.ErrorErr(err)
				return
			}
		}()
	}
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
		err := c.db.V2.CodeInsert(m)
		if err != nil {
			c.log.ErrorErr(err)
			c.log.InfoStruct("CodeInsert", m)
		}
	}()

	return m.Code
}

func (c *Hs) timerEditMessage(m models.IncomingMessage, mid1, mid2, mid3, token string) {
	time.Sleep(3 * time.Minute)
	links := "https://mentalisit.github.io/HadesSpace/compendiumTech?secretToken=" + token
	text := fmt.Sprintf(c.getText(m, "SECRET_LINK"), links, m.MGuild.GuildName)
	if m.Type == "ds" || m.Type == "tg" {
		err := c.editMessage(m, m.DmChat, mid1, text, "MarkdownV2")
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
	} else if m.Type == "wa" {
		c.sendDM(m, text)
	}
	err := c.deleteMessage(m, m.DmChat, mid2)
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

func (c *Hs) checkMoron(in models.IncomingMessage) bool {
	if len(c.moron) > 0 {
		if c.moron[in] != 0 {
			return true
		}
	}
	c.moron[in] += 1
	return false
}
