package ds

import (
	"bridge/models"
	"github.com/mentalisit/logger"
)

type Discord struct {
	log *logger.Logger
}

func NewDiscord(log *logger.Logger) *Discord {
	return &Discord{log: log}
}

func (d *Discord) DeleteMessageDs(ChatId, MesId string) {
	s := models.DeleteMessageStruct{
		MessageId: MesId,
		Channel:   ChatId,
	}
	d.MarshalDataDiscordDel(s)
}
func (d *Discord) SendChannelDelSecondDs(chatId, text string, second int) {
	s := models.SendTextDeleteSeconds{
		Text:    text,
		Channel: chatId,
		Seconds: second,
	}
	d.MarshalDataDiscordSendDel(s)
}
