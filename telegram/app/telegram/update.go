package telegram

import (
	//tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"fmt"
	"strconv"
	"strings"
	"telegram/models2"
	"time"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/mentalisit/restapi/models"
)

func (t *Telegram) update() {
	ut := tgbotapi.NewUpdate(0)
	ut.Timeout = 60
	// Полный список возможных значений для AllowedUpdates:
	allowedUpdates := []string{
		"message",              // новые сообщения
		"edited_message",       // редактированные сообщения
		"channel_post",         // посты в каналах
		"edited_channel_post",  // редактированные посты в каналах
		"inline_query",         // инлайн-запросы
		"chosen_inline_result", // выбранные инлайн-результаты
		"callback_query",       // callback-кнопки
		"shipping_query",       // запросы доставки
		"pre_checkout_query",   // предварительные проверки оплаты
		"poll",                 // опросы
		"poll_answer",          // ответы на опросы
		"my_chat_member",       // изменения статуса бота
		"chat_member",          // изменения статуса участников
		"chat_join_request",    // запросы на вступление в чат
	}
	ut.AllowedUpdates = allowedUpdates
	updates := t.t.GetUpdatesChan(ut)

	for update := range updates {
		if update.InlineQuery != nil {
			t.handleInlineQuery(update.InlineQuery)
		} else if update.ChosenInlineResult != nil {
			t.handleChosenInlineResult(update.ChosenInlineResult)
		} else if update.CallbackQuery != nil {
			t.callback(update.CallbackQuery) //нажатия в чате
		} else if update.Message != nil {
			//fmt.Printf("update: %+v\n", update.Message)
			if update.Message.Poll != nil {
				if update.Message.From.String() == "mentalisit" {

					m := models.Request{
						Data:    make(map[string]string),
						Options: []string{},
					}
					m.Data["chatid"] = strconv.FormatInt(update.Message.Chat.ID, 10) + "/0"
					m.Data["question"] = update.Message.Poll.Question
					m.Data["url"] = "https://mentalisit.tsl.rocks/rs/settings/test/HadesTable.html"
					m.Data["author"] = update.Message.From.String()
					for _, option := range update.Message.Poll.Options {
						m.Options = append(m.Options, option.Text)
					}

					go t.SendPoll(m)

				}
				fmt.Printf("messagePoll: %+v\n", update.Message.Poll)
				//Chat:{ID:-1003311563636
				//Type:supergroup
				//Title:Gr1
				//IsForum:true
				//}
				//IsTopicMessage:true
				//ReplyToMessage:0x256dbc743208
				//Poll:0x256dbc60a60
				//
				//messagePoll: &{
				//	ID:5402602678322732086
				//	Question:ф
				//	QuestionEntities:[]
				//	Options:[
				//		{PersistentID:0 Text:1 TextEntities:[] VoterCount:0 AddedByUser: AddedByChat:<nil> AdditionDate:0}
				//		{PersistentID:1 Text:2 TextEntities:[] VoterCount:0 AddedByUser: AddedByChat:<nil> AdditionDate:0}]
				//	TotalVoterCount:0
				//	IsClosed:false
				//	IsAnonymous:false
				//	Type:regular
				//	AllowsRevoting:true
				//	}

			}
			t.updateMessage(update.Message)
		} else if update.EditedMessage != nil {
			if time.Now().Unix()-update.EditedMessage.Date < 300 {
				t.logicMix(update.EditedMessage, true)
			}
		} else if update.MyChatMember != nil {
			t.myChatMember(update.MyChatMember)
		} else if update.ChatMember != nil {
			t.chatMember(update.ChatMember)
		} else if update.ChatJoinRequest != nil {
			t.log.Info(fmt.Sprintf("ChatJoinRequest Chat %s From @%s\n",
				update.ChatJoinRequest.Chat.Title, update.ChatJoinRequest.From.String()))
		} else if update.PollAnswer != nil {
			fmt.Printf("pollAnswer %+v\n", update.PollAnswer)
			//todo need update in db
			//tg, bridgeConfig := t.bridgeCheckChannelConfigTg(ChatId)
			//if tg {
			//	fmt.Println("cb.Data " + cb.Data)
			//	if strings.HasPrefix(cb.Data, "17") {
			//		mes := models.ToBridgeMessage{
			//			ChatId:        ChatId,
			//			Config:        &bridgeConfig,
			//			Text:          ".poll " + cb.Data,
			//			Tip:           "tg",
			//			MesId:         strconv.Itoa(cb.Message.MessageID),
			//			GuildId:       strconv.FormatInt(cb.GetInaccessibleMessage().Chat.ID, 10),
			//			TimestampUnix: cb.Message.Time().Unix(),
			//			Sender:        ReplaceCyrillicToLatin(cb.From.String()),
			//			SenderId:      strconv.FormatInt(cb.From.ID, 10),
			//		}
			//
			//		if mes.Text != "" {
			//			t.api.SendBridgeAppRecover(mes)
			//		}
			//	}
			//}
		} else {
			fmt.Printf("else %+v \n", update)
		}
	}
}
func (t *Telegram) updateMessage(m *tgbotapi.Message) {
	switch m.Text {
	case "/start":
		t.handleStartCommand(m)
	case "/webapp", "/roles":
		//t.webApp.RemoveReplyKeyboard(m.Chat.ID)
		t.SendWebAppButtonSmart(m.Chat.ID)
	case "/chatroles":
		// Специальная команда для управления ролями в текущем чате
		t.SendWebAppButtonSmart(m.Chat.ID)
	}

	if m.IsCommand() {
		t.ifCommand(m)
	} else if m.Chat.IsPrivate() { //если пишут боту в личку
		t.ifPrivatMesage(m)
	} else if m.LeftChatMember != nil {
		_ = t.Storage.Db.RemoveUserFromChat(m.Chat.ID, m.LeftChatMember.ID)
	} else { //остальные сообщения
		t.logicMix(m, false)
	}
}

// saveMessageToStorage сохраняет сообщение Telegram в БД
func (t *Telegram) saveMessageToStorage(m *tgbotapi.Message) {
	// Получаем communityID (UUID) для чата
	communityID, err := t.Storage.Db.GetGildUUIDMyCompendium(m.Chat.ID)
	if err != nil {
		mg := models2.MultiAccountGuildV2{
			GId:       uuid.New(),
			GuildName: m.Chat.Title,
			Channels:  make(models2.GuildChannels),
		}
		mg.Channels["tg"] = append(mg.Channels["tg"], fmt.Sprint(m.Chat.ID))
		save, _ := t.Storage.Db.GuildSave(mg)
		if save.GId != uuid.Nil {
			communityID = &save.GId
		} else {
			// Если не нашли guild, используем нулевой UUID или логируем
			t.log.Info(fmt.Sprintf("Не удалось получить communityID для чата %d: %v", m.Chat.ID, err))
			return
		}

	}

	// Сохраняем сообщение
	if err := t.Storage.Db.SaveTelegramMessage(*communityID, m); err != nil {
		// Собираем информацию о вложениях для детализации ошибки
		var attachmentInfo string
		if m.Photo != nil {
			attachmentInfo = fmt.Sprintf("Photo: %d sizes", len(m.Photo))
		} else if m.Document != nil {
			attachmentInfo = fmt.Sprintf("Document: %s (file_id: %s)", m.Document.FileName, m.Document.FileID)
		} else if m.Video != nil {
			attachmentInfo = fmt.Sprintf("Video: file_id=%s", m.Video.FileID)
		} else if m.Audio != nil {
			attachmentInfo = fmt.Sprintf("Audio: file_id=%s", m.Audio.FileID)
		} else if m.Voice != nil {
			attachmentInfo = fmt.Sprintf("Voice: file_id=%s", m.Voice.FileID)
		} else if m.Sticker != nil {
			attachmentInfo = fmt.Sprintf("Sticker: file_id=%s", m.Sticker.FileID)
		}

		t.log.Error(fmt.Sprintf("Ошибка сохранения сообщения %d в чате %d (communityID: %s): %v | Вложения: %s | Текст: %q",
			m.MessageID, m.Chat.ID, *communityID, err, attachmentInfo, m.Text))
	}
}

func (t *Telegram) SendWebAppButtonSmart(chatID int64) {
	//t.webApp.SendWebAppButtonSmart(chatID)
	fmt.Println("SendWebAppButtonSmart")

}

func (t *Telegram) handleStartCommand(message *tgbotapi.Message) {
	args := message.CommandArguments()
	fmt.Printf("Start command with args: '%s'\n", args)

	// Обрабатываем глубокие ссылки в формате: startapp=chat123456789
	if strings.HasPrefix(args, "chat") {
		// Извлекаем chat_id из аргументов: "chat-123456789"
		chatIDStr := strings.TrimPrefix(args, "chat")
		var chatID int64
		if id, err := strconv.ParseInt(chatIDStr, 10, 64); err == nil {
			chatID = id
			fmt.Printf("Processing deep link for chat ID: %d\n", chatID)
			t.openWebAppForGroup(message.Chat.ID, chatID)
			return
		} else {
			fmt.Printf("Error parsing chat ID from '%s': %v\n", chatIDStr, err)
		}
	}

}

// Открывает Web App для группы через глубокую ссылку
func (t *Telegram) openWebAppForGroup(userChatID int64, groupChatID int64) {
	// Получаем информацию о группе
	chat, err := t.t.GetChat(tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{ChatID: groupChatID},
	})

	var chatTitle string
	if err == nil && chat.Title != "" {
		chatTitle = chat.Title
	} else {
		chatTitle = fmt.Sprintf("Группа ID: %d", groupChatID)
	}
	fmt.Printf("chatTitle %s ID %+v\n", chatTitle, groupChatID)
	webAppURL := fmt.Sprintf("https://webapp.mentalisit.myds.me/?chat_id=%d", groupChatID)

	msg := tgbotapi.NewMessage(userChatID,
		fmt.Sprintf("🎭 *Управление ролями для \"%s\"*\n\nОткрываю панель управления...", chatTitle))
	msg.ParseMode = "Markdown"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonWebApp(
				"📋 Открыть управление ролями",
				tgbotapi.WebAppInfo{URL: webAppURL},
			),
		),
	)
	msg.ReplyMarkup = keyboard

	if _, err := t.t.Send(msg); err != nil {
		fmt.Printf("Error sending WebApp for group: %v\n", err)
		// Fallback - отправляем обычную ссылку
		fallbackMsg := tgbotapi.NewMessage(userChatID,
			fmt.Sprintf("🎭 *Управление ролями для \"%s\"*\n\n[Открыть в браузере](%s)",
				chatTitle, webAppURL))
		fallbackMsg.ParseMode = "Markdown"
		t.t.Send(fallbackMsg)
	} else {
		fmt.Printf("✅ Opened WebApp for group '%s' (ID: %d)\n", chatTitle, groupChatID)
	}
}
