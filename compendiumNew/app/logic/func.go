package logic

import (
	"compendium/models"
	"fmt"
	"runtime"
	"strings"
	"time"
)

func (c *Hs) sendChat(m models.IncomingMessage, text string) {
	if m.Type == "ds" {
		_, err := c.ds.SendChannel(m.ChannelId, text)
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
	} else if m.Type == "tg" {
		_, err := c.tg.Send(m.ChannelId, text)
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
	}
}
func (c *Hs) sendChatTable(m models.IncomingMessage, text, table string) {
	if m.Type == "ds" {
		_, err := c.ds.SendChannel(m.ChannelId, text+"\n```"+table+"```")
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
	} else if m.Type == "tg" {
		_, err := c.tg.Send(m.ChannelId, text+"\n"+table)
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
	}
}

func (c *Hs) sendDM(m models.IncomingMessage, text string) (string, error) {
	if m.Type == "ds" {
		mid, err := c.ds.SendChannel(m.DmChat, text)
		if err != nil {
			c.log.ErrorErr(err)
			return "", err
		}
		return mid, nil
	} else if m.Type == "tg" {
		mid, err := c.tg.Send(m.DmChat, text)
		if err != nil {
			return "", err
		}
		return mid, nil
	}
	return "", nil
}
func (c *Hs) sendChatPic(m models.IncomingMessage, text string, pic []byte) {
	if m.Type == "ds" {
		err := c.ds.SendChannelPic(m.ChannelId, text, pic)
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
	} else if m.Type == "tg" {
		err := c.tg.SendPic(m.ChannelId, text, pic)
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

func (c *Hs) deleteMessage(m models.IncomingMessage, chat, mid string) error {
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
	}
	return nil
}

func (c *Hs) editMessage(m models.IncomingMessage, chat, mid, text, ParseMode string) error {
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
	}
	return nil
}
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

func (c *Hs) PrintGoroutine() {
	goroutine := runtime.NumGoroutine()
	tm := time.Now()
	mdate := (tm.Format("2006-01-02"))
	mtime := (tm.Format("15:04"))
	text := fmt.Sprintf(" %s %s Горутин  %d\n", mdate, mtime, goroutine)
	if goroutine > 120 {
		c.log.Info(text)
		c.log.Panic(text)
	} else if goroutine > 50 && goroutine%10 == 0 {
		c.log.Info(text)
	}

	fmt.Println(text)
}
