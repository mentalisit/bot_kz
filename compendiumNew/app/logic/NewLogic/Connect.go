package NewLogic

import (
	"compendium/models"
	"fmt"
	"time"
)

func (c *HsLogic) connect(m models.IncomingMessage) (conn bool) {
	if helperCommand(m.Text, "connect") ||
		helperCommand(m.Text, "подключить") ||
		helperCommand(m.Text, "підключити") {
		conn = true
	}
	if !conn {
		return false
	}

	text := fmt.Sprintf("выполняется подготовка для %s\n", m.MultiAccount.Nickname)
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

	newIdentify := c.GenerateIdentity(m)

	code := c.generateCodeAndSave(newIdentify)
	sendDMCode, _ := c.sendDM(m, code)

	sendDM, _ := c.sendDM(m, "Подготавливаю секретную ссылку")

	go func() {
		time.Sleep(5 * time.Second)
		_ = c.deleteMessage(m, m.DmChat, mid1)

		links := "https://mentalisit.github.io/HadesSpace/compendiumTech?secretToken=" + newIdentify.Token
		text = fmt.Sprintf(c.getText(m, "SECRET_LINK"), links, "")
		err = c.editMessage(m, m.DmChat, sendDM, text, "MarkdownV2")
		if err != nil {
			c.log.ErrorErr(err)
		}
		time.Sleep(1 * time.Minute)
		_ = c.deleteMessage(m, m.DmChat, sendDMCode)
	}()

	return
}

func (c *HsLogic) checkMoron(in models.IncomingMessage) bool {
	if len(c.moron) > 0 {
		if c.moron[in] != 0 {
			return true
		}
	}
	c.moron[in] += 1
	return false
}
