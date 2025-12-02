package NewLogic

import (
	"compendium/models"
	"strings"
)

func helperCommand(text, command string) bool {
	removePrefix, pr := strings.CutPrefix(text, "%")
	if pr {
		toLower := strings.ToLower(removePrefix)
		af := strings.Split(toLower, " ")
		for _, s := range af {
			if s == command {
				return true
			}
		}
	}

	return false
}

func (c *HsLogic) getText(m models.IncomingMessage, key string) string {
	return c.Dict.GetText(m.Language, key)
}

func (c *HsLogic) sendDM(m models.IncomingMessage, text string) (string, error) {
	if m.Type == "ds" {
		mid, err := c.ds.SendChannel(m.DmChat, text)
		if err != nil {
			c.log.ErrorErr(err)
			return "", err
		}
		return mid, nil
	} else if m.Type == "tg" {
		mid, err := c.tg.Send(m.DmChat, text, "")
		if err != nil {
			return "", err
		}
		return mid, nil
	} else if m.Type == "wa" {
		mid, err := c.wa.Send(m.DmChat, text)
		if err != nil {
			return "", err
		}
		return mid, nil
	}
	return "", nil
}

func (c *HsLogic) sendChat(m models.IncomingMessage, text string) {
	if m.Type == "ds" {
		_, err := c.ds.SendChannel(m.ChannelId, text)
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
	} else if m.Type == "tg" {
		_, err := c.tg.Send(m.ChannelId, text, "")
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
	} else if m.Type == "wa" {
		_, err := c.wa.Send(m.ChannelId, text)
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
	}
}

func (c *HsLogic) deleteMessage(m models.IncomingMessage, chat, mid string) error {
	if m.Type == "ds" {
		err := c.ds.DeleteMessage(chat, mid)
		if err != nil {
			return err
		}
	} else if m.Type == "tg" {
		err := c.tg.DeleteMessage(chat, mid)
		if err != nil {
			return err
		}
	} else if m.Type == "wa" {
		err := c.wa.DeleteMessage(chat, mid)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *HsLogic) editMessage(m models.IncomingMessage, chat, mid, text, ParseMode string) error {
	if m.Type == "ds" {
		err := c.ds.EditMessage(chat, mid, text)
		if err != nil {
			return err
		}
	} else if m.Type == "tg" {
		err := c.tg.EditMessage(chat, mid, text, ParseMode)
		if err != nil {
			return err
		}
	} else if m.Type == "wa" {
		err := c.wa.EditMessage(chat, mid, text, ParseMode)
		if err != nil {
			return err
		}
	}
	return nil
}
