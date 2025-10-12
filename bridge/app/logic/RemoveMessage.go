package logic

import (
	"strings"
)

func (b *Bridge) RemoveMessage() {
	m := make(map[string]string)
	m[b.in.ChatId] = b.in.MesId
	linkedID, err := b.storage.GetMapByLinkedID(m)
	if err != nil {
		b.log.ErrorErr(err)
		return
	}
	if linkedID != nil && len(linkedID) > 0 {
		for chatId, mID := range linkedID {
			if strings.HasPrefix(chatId, "-100") { //tg
				go b.telegram.DeleteMessage(chatId, mID)
			} else if strings.HasSuffix(chatId, "@g.us") { //wa
				go b.whatsapp.DeleteMessage(chatId, mID)
			} else if IsPurelyNumeric(chatId) { //ds
				go b.discord.DeleteMessageDs(chatId, mID)
			}
		}
	}
}

func (b *Bridge) EditMessage() {
	if b.in.Tip == "dse" {
		b.in.Tip = "ds"
	} else if b.in.Tip == "tge" {
		b.in.Tip = "tg"
	} else if b.in.Tip == "wae" {
		b.in.Tip = "wa"
	}
	b.in.Text = "EDIT: " + b.in.Text
}
