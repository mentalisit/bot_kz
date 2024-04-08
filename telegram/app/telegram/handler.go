package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"telegram/models"
	"telegram/telegram/restapi"
)

func (t *Telegram) callback(cb *tgbotapi.CallbackQuery) {
	callback := tgbotapi.NewCallback(cb.ID, cb.Data)
	if _, err := t.t.Request(callback); err != nil {
		t.log.ErrorErr(err)
	}
	ChatId := strconv.FormatInt(cb.Message.Chat.ID, 10) + fmt.Sprintf("/%d", cb.Message.MessageThreadID)
	ok, config := t.checkChannelConfigTG(ChatId)
	if ok {
		in := models.InMessage{
			Mtext:       cb.Data,
			Tip:         "tg",
			Name:        cb.From.String(),
			NameMention: "@" + t.nickName(cb.From, ChatId),
			Tg: struct {
				Mesid int
			}{
				Mesid: cb.Message.MessageID,
			},
			Config: config,
			Option: models.Option{
				Reaction: true},
		}

		err := restapi.SendRsBotApp(in)
		if err != nil {
			t.log.ErrorErr(err)
		}
	}
}
func (t *Telegram) ifPrivatMesage(m *tgbotapi.Message) {
	if m.Text == "/start" {
		_, err := t.t.Send(tgbotapi.NewMessage(m.From.ID,
			"Возможность писать сообщения в личку активирована /"+
				" The ability to write private messages is activated"))
		if err != nil {
			t.log.ErrorErr(err)
			return
		}
	} else {
		//нужно решить что тут делать
		t.log.Info(m.From.String() + ": " + m.Text)
	}

}

//func (t *Telegram) updatesComand(c *tgbotapi.Message) {
//	ChatId := strconv.FormatInt(c.Chat.ID, 10) + fmt.Sprintf("/%d", c.MessageThreadID)
//	if c.Command() == "chatid" {
//		t.SendChannelDelSecond(ChatId, ChatId, 20)
//	}
//	ok, config := t.checkChannelConfigTG(ChatId)
//	if ok {
//		MessageID := strconv.Itoa(c.MessageID)
//		switch c.Command() {
//		case "help":
//			t.help(config.TgChannel, MessageID)
//		case "helpqueue":
//			t.helpQueue(config.TgChannel, MessageID)
//		case "helpnotification":
//			t.helpNotification(config.TgChannel, MessageID)
//		case "helpevent":
//			t.helpEvent(config.TgChannel, MessageID)
//		case "helptop":
//			t.helpTop(config.TgChannel, MessageID)
//		case "helpicon":
//			t.helpIcon(config.TgChannel, MessageID)
//		}
//	} else {
//		switch c.Command() {
//		case "help":
//			t.SendChannelDelSecond(ChatId, "Активируйте бота командой \n.add", 60)
//		default:
//			t.SendChannelDelSecond(ChatId, "Вам не доступна данная команда \n /help", 60)
//		}
//	}
//}

func (t *Telegram) myChatMember(member *tgbotapi.ChatMemberUpdated) {
	ChatId := strconv.FormatInt(member.Chat.ID, 10) + "/0"
	if member.NewChatMember.Status == "member" {
		t.SendChannelDelSecond(ChatId, fmt.Sprintf("@%s мне нужны права админа для коректной работы", member.From.UserName), 60)
	} else if member.NewChatMember.Status == "administrator" {
		t.SendChannelDelSecond(ChatId, fmt.Sprintf("@%s спасибо ... я готов к работе \nАктивируй нужный режим бота,\n если сложности пиши мне @Mentalisit", member.From.UserName), 60)
	}
}

func (t *Telegram) chatMember(chMember *tgbotapi.ChatMemberUpdated) {
	if chMember.NewChatMember.IsMember {
		ChatId := strconv.FormatInt(chMember.Chat.ID, 10) + "/0"
		t.SendChannelDelSecond(ChatId,
			fmt.Sprintf("%s Добро пожаловать в наш чат ", chMember.NewChatMember.User.FirstName),
			60)
	}

}
func (t *Telegram) handlePoll(message *tgbotapi.Message) {
	if message.Poll != nil {
		text := "Запущен  ОПРОС  \n"
		text += message.Poll.Question
		text += "\nВарианты ответа:\n"
		for _, o := range message.Poll.Options {
			text += fmt.Sprintf(" %s\n", o.Text)
		}
		message.Text = text
	}
}
