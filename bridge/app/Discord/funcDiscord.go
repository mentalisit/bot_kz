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
	err := d.MarshalDataDiscordDel(s)
	if err != nil {
		d.log.ErrorErr(err)
		d.log.InfoStruct("DeleteMessageDs s: ", s)
	}
}
func (d *Discord) SendChannelDelSecondDs(chatId, text string, second int) {
	s := models.SendTextDeleteSeconds{
		Text:    text,
		Channel: chatId,
		Seconds: second,
	}
	err := d.MarshalDataDiscordSendDel(s)
	if err != nil {
		d.log.ErrorErr(err)
		d.log.InfoStruct("SendChannelDelSecondDs s:", s)
	}
}
