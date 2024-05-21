package logic

import (
	"compendium/logic/ds"
	"compendium/logic/tg"
	"compendium/models"
	"fmt"
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
func (c *Hs) sendFormatedText(m models.IncomingMessage, Text string, data [][]string) {
	// Определяем максимальную длину для каждого столбца
	colWidths := make([]int, len(data[0]))
	for _, row := range data {
		for i, col := range row {
			if len(col) > colWidths[i] {
				colWidths[i] = len(col)
			}
		}
	}

	// Формируем формат для печати строк
	format := ""
	for _, width := range colWidths {
		format += fmt.Sprintf("%%-%ds  ", width)
	}
	format = format[:len(format)-2] // Убираем последний лишний пробел

	// Печатаем строки с выравниванием
	text := ""
	for _, row := range data {
		text += fmt.Sprintf(format+"\n", row[0], row[1], row[2])
	}
	c.sendChatTable(m, Text, text)
}
