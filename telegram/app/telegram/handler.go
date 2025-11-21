package telegram

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"telegram/models"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func (t *Telegram) callback(cb *tgbotapi.CallbackQuery) {
	callback := tgbotapi.NewCallback(cb.ID, cb.Data)
	if _, err := t.t.Request(callback); err != nil {
		t.log.ErrorErr(err)
	}

	t.SaveMember(&cb.Message.Chat, cb.From)

	ChatId := strconv.FormatInt(cb.Message.Chat.ID, 10) + fmt.Sprintf("/%d", cb.Message.MessageThreadID)
	ok, config := t.checkChannelConfigTG(ChatId)
	if ok {
		in := models.InMessage{
			Mtext:       cb.Data,
			Tip:         "tg",
			Username:    cb.From.String(),
			UserId:      strconv.FormatInt(cb.From.ID, 10),
			NameMention: "@" + cb.From.UserName,
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
	if m.Text == ".паника" {
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
			NameMention: "@" + m.From.UserName,
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

func (t *Telegram) ifCommand(m *tgbotapi.Message) {
	switch m.Text {
	case "start":
		if m.CommandArguments() == "open_roles" {
			t.SendWebAppButtonSmart(m.Chat.ID)
		} else {
			_, err := t.t.Send(tgbotapi.NewMessage(m.From.ID,
				"Возможность писать сообщения в личку активирована, дальнейшее взаимодействие только с чата корпорации.  %помощь/"+
					" The ability to write private messages is activated, further interaction only through the corporate chat. %help"))
			if err != nil {
				t.log.ErrorErr(err)
				return
			}
		}
	}
}

func (t *Telegram) myChatMember(member *tgbotapi.ChatMemberUpdated) {
	t.SaveMember(&member.Chat, member.NewChatMember.User)

	ChatId := strconv.FormatInt(member.Chat.ID, 10) + "/0"
	if member.NewChatMember.Status == "member" {
		t.SendChannelDelSecond(ChatId, fmt.Sprintf("@%s мне нужны права админа для коректной работы", member.From.UserName), "", 60)
	} else if member.NewChatMember.Status == "administrator" {
		t.SendChannelDelSecond(ChatId, fmt.Sprintf("@%s спасибо ... я готов к работе \nАктивируй нужный режим бота,\n если сложности пиши мне @Mentalisit", member.From.UserName), "", 60)
	}
}

func (t *Telegram) chatMember(chMember *tgbotapi.ChatMemberUpdated) {
	// Если участник вышел или был удален
	if chMember.NewChatMember.Status == "left" || chMember.NewChatMember.Status == "kicked" {
		_ = t.Storage.Db.RemoveUserFromChat(context.Background(), chMember.Chat.ID, chMember.NewChatMember.User.ID)
	} else {
		// Обновляем информацию
		t.SaveMember(&chMember.Chat, chMember.NewChatMember.User)
	}

	if chMember.NewChatMember.IsMember {
		ChatId := strconv.FormatInt(chMember.Chat.ID, 10) + "/0"
		t.SendChannelDelSecond(ChatId,
			fmt.Sprintf("%s Добро пожаловать в наш чат ", chMember.NewChatMember.User.FirstName),
			"", 60)
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

func (t *Telegram) handleInlineQuery(inlineQuery *tgbotapi.InlineQuery) {
	fmt.Printf("Inline Query от %s: %s\n",
		inlineQuery.From.UserName,
		inlineQuery.Query)

	// Query - это текст после @username_бота
	userQuery := inlineQuery.Query
	fmt.Printf("Пользователь ищет: '%s'\n", userQuery)

	// Создаем результаты для инлайн-режима
	var results []interface{}
	if userQuery == "" {
		article := tgbotapi.NewInlineQueryResultArticle(
			"1",
			"Создать роль",
			"Вы искали: "+userQuery,
		)
		article.Description = "ну тут понятно"
		results = append(results, article)

		article = tgbotapi.NewInlineQueryResultArticle(
			"2",
			"Удалить роль",
			"Вы искали: "+userQuery,
		)
		article.Description = "и тут тоже"
		results = append(results, article)

		article = tgbotapi.NewInlineQueryResultArticle(
			"3",
			"Подписаться на роль",
			"Вы искали: "+userQuery,
		)
		article.Description = "для пинга "
		results = append(results, article)
		article = tgbotapi.NewInlineQueryResultArticle(
			"4",
			"Отписаться от роли",
			"Вы искали: "+userQuery,
		)
		article.Description = "ну и иди нахрен"
		results = append(results, article)
		article = tgbotapi.NewInlineQueryResultArticle(
			"5",
			"Список ролей ",
			"Вы искали: "+userQuery,
		)
		article.Description = "ну если надо"
		results = append(results, article)
	}

	if userQuery != "" {
		// Пример: создаем текстовый результат
		article := tgbotapi.NewInlineQueryResultArticle(
			"1",
			"Результат для: "+userQuery,
			"Вы искали: "+userQuery,
		)
		article.Description = "Описание результата"
		results = append(results, article)
	}

	// Отправляем результаты
	inlineConfig := tgbotapi.InlineConfig{
		InlineQueryID: inlineQuery.ID,
		Results:       results,
		CacheTime:     0,
	}

	_, err := t.t.Request(inlineConfig)
	if err != nil {
		log.Println("Ошибка отправки инлайн-результатов:", err)
	}
}

func (t *Telegram) handleChosenInlineResult(result *tgbotapi.ChosenInlineResult) {
	fmt.Printf("Пользователь выбрал результат:\n")
	fmt.Printf("ID результата: %s\n", result.ResultID)
	fmt.Printf("Запрос: %s\n", result.Query)
	fmt.Printf("Пользователь: %s (ID: %d)\n",
		result.From.UserName, result.From.ID)

	// Обрабатываем выбор в зависимости от ID результата
	switch result.ResultID {
	case "weather_moscow":
		fmt.Println("✅ Пользователь выбрал погоду в Москве")
		// Здесь можно отправить уведомление, сохранить в БД и т.д.

	case "weather_spb":
		fmt.Println("✅ Пользователь выбрал погоду в СПб")

	case "weather_sochi":
		fmt.Println("✅ Пользователь выбрал погоду в Сочи")

	default:
		fmt.Printf("❌ Неизвестный выбор: %s\n", result.ResultID)
	}

	// Дополнительная информация
	fmt.Printf("Inline Message ID: %s\n", result.InlineMessageID)
}
