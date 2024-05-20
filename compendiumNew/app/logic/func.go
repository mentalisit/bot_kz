package logic

import (
	"compendium/logic/ds"
	"compendium/logic/tg"
	"compendium/models"
)

func (c *Hs) sendChat(m models.IncomingMessage, text string) {
	if m.Type == "ds" {
		err := ds.SendChannel(m.ChannelId, text)
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
	} else if m.Type == "tg" {
		err := tg.Send(m.ChannelId, text)
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
	}
}
func (c *Hs) sendChatTable(m models.IncomingMessage, text, table string) {
	if m.Type == "ds" {
		err := ds.SendChannel(m.ChannelId, text+"\n```"+table+"```")
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
	} else if m.Type == "tg" {
		err := tg.Send(m.ChannelId, text+"\n"+table)
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
	}
}

func (c *Hs) sendDM(m models.IncomingMessage, text string) error {
	if m.Type == "ds" {
		err := ds.SendChannel(m.DmChat, text)
		if err != nil {
			c.log.ErrorErr(err)
			return err
		}
	} else if m.Type == "tg" {
		err := tg.Send(m.DmChat, text)
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *Hs) sendChatPic(m models.IncomingMessage, text string, pic []byte) {
	if m.Type == "ds" {
		err := ds.SendChannelPic(m.ChannelId, text, pic)
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
	} else if m.Type == "tg" {
		err := tg.SendPic(m.ChannelId, text, pic)
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
	}

}
func (c *Hs) getText(m models.IncomingMessage, key string) string {
	return c.Dict.GetText(m.Language, key)
}
