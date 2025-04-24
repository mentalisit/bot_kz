package telegram

import (
	"fmt"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"strconv"
	"strings"
	"telegram/models"
)

func (t *Telegram) callback(cb *tgbotapi.CallbackQuery) {
	callback := tgbotapi.NewCallback(cb.ID, cb.Data)
	if _, err := t.t.Request(callback); err != nil {
		t.log.ErrorErr(err)
	}
	ChatId := strconv.FormatInt(cb.Message.Chat.ID, 10) + fmt.Sprintf("/%d", cb.Message.MessageThreadID)
	ok, config := t.checkChannelConfigTG(ChatId)
	if ok {
		name := t.nickName(cb.From, ChatId)
		in := models.InMessage{
			Mtext:       cb.Data,
			Tip:         "tg",
			Username:    name,
			UserId:      strconv.FormatInt(cb.From.ID, 10),
			NameMention: "@" + name,
			Tg: struct {
				Mesid int
			}{
				Mesid: cb.Message.MessageID,
			},
			Config: config,
			Option: models.Option{
				Reaction: true},
		}

		t.api.SendRsBotAppRecover(in)
	}
	tg, bridgeConfig := t.bridgeCheckChannelConfigTg(ChatId)
	if tg {
		fmt.Println("cb.Data " + cb.Data)
		if strings.HasPrefix(cb.Data, "17") {
			mes := models.ToBridgeMessage{
				ChatId:        ChatId,
				Config:        &bridgeConfig,
				Text:          ".poll " + cb.Data,
				Tip:           "tg",
				MesId:         strconv.Itoa(cb.Message.MessageID),
				GuildId:       strconv.FormatInt(cb.GetInaccessibleMessage().Chat.ID, 10),
				TimestampUnix: cb.Message.Time().Unix(),
				Sender:        ReplaceCyrillicToLatin(cb.From.String()),
			}

			if mes.Text != "" {
				t.api.SendBridgeAppRecover(mes)
			}
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
	} else if m.Text == ".паника" {
		t.log.Panic(".паника " + m.From.String())
	} else if strings.HasPrefix(m.Text, "%") {
		i := models.IncomingMessage{
			Text:        m.Text,
			DmChat:      strconv.FormatInt(m.From.ID, 10),
			Name:        m.From.String(),
			MentionName: "@" + m.From.String(),
			NameId:      strconv.FormatInt(m.From.ID, 10),
			NickName:    "", //нет способа извлечь ник кроме member.CustomTitle
			Type:        "tgDM",
		}

		if m.From != nil && m.From.LanguageCode != "" {
			i.Language = m.From.LanguageCode
		}
		t.api.SendCompendiumAppRecover(i)
	} else {
		in := models.InMessage{
			Mtext:       m.Text,
			Tip:         "tgDM",
			Username:    m.From.String(),
			UserId:      strconv.FormatInt(m.From.ID, 10),
			NameMention: "@" + m.From.String(),
			Tg: struct {
				Mesid int
			}{
				Mesid: m.MessageID,
			},
			Config: models.CorporationConfig{
				TgChannel: strconv.FormatInt(m.Chat.ID, 10),
			},
		}

		t.api.SendRsBotAppRecover(in)
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
