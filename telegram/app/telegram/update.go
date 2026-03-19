package telegram

import (
	//tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
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
		_ = t.Storage.Db.RemoveUserFromChat(context.Background(), m.Chat.ID, m.LeftChatMember.ID)
	} else { //остальные сообщения
		t.logicMix(m, false)
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
