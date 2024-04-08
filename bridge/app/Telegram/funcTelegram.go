package tg

import (
	"bridge/models"
	"github.com/mentalisit/logger"
	"strconv"
)

type Telegram struct {
	log *logger.Logger
}

func NewTelegram(log *logger.Logger) *Telegram {
	return &Telegram{log: log}
}

func (t *Telegram) DeleteMessage(ChatId string, MesId int) {
	s := models.DeleteMessageStruct{
		MessageId: strconv.Itoa(MesId),
		Channel:   ChatId,
	}
	t.MarshalDelTelegram(s)
}
func (t *Telegram) SendChannelDelSecond(chatId, text string, second int) {
	s := models.SendTextDeleteSeconds{
		Text:    text,
		Channel: chatId,
		Seconds: second,
	}
	t.MarshalSendDelTelegram(s)
}
