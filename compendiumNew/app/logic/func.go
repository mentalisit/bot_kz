package logic

import (
	"compendium/logic/ds"
	"compendium/logic/tg"
)

func (c *Hs) sendChat(text string) {
	if c.in.Type == "ds" {
		err := ds.SendChannel(c.in.ChannelId, text)
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
	} else if c.in.Type == "tg" {
		err := tg.Send(c.in.ChannelId, text)
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
	}
}

func (c *Hs) sendDM(text string) error {
	if c.in.Type == "ds" {
		err := ds.SendChannel(c.in.DmChat, text)
		if err != nil {
			c.log.ErrorErr(err)
			return err
		}
	} else if c.in.Type == "tg" {
		err := tg.Send(c.in.DmChat, text)
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *Hs) sendChatPic(text string, pic []byte) {
	if c.in.Type == "ds" {
		err := ds.SendChannelPic(c.in.ChannelId, text, pic)
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
	} else if c.in.Type == "tg" {
		err := tg.SendPic(c.in.ChannelId, text, pic)
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
	}

}
