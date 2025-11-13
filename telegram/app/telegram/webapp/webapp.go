package webapp

import (
	"fmt"
	"telegram/telegram/roles"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

type WebApp struct {
	bot      *tgbotapi.BotAPI
	auth     *AuthManager
	handlers *Handlers
}

func NewWebApp(bot *tgbotapi.BotAPI, rolesManager *roles.Manager) *WebApp {
	auth := NewAuthManager(bot.Token)
	handlers := NewHandlers(bot, auth, rolesManager)

	return &WebApp{
		bot:      bot,
		auth:     auth,
		handlers: handlers,
	}
}

func (w *WebApp) Start() {
	fmt.Println("Web App server starting on :8080")
	w.handlers.Start()
}

// Умная отправка Web App кнопки
func (w *WebApp) SendWebAppButtonSmart(chatID int64) {
	fmt.Println("sendWebAppButtonSmart")

	webAppURL := "https://webapp.mentalisit.myds.me/"
	botUsername := w.bot.Self.UserName

	// Для ЛИЧНЫХ чатов используем INLINE кнопку с WebApp
	if chatID > 0 {
		w.sendWebAppButtonPrivate(chatID, webAppURL)
	} else {
		// Для ГРУПП используем специальную кнопку
		w.sendWebAppButtonGroup(chatID, webAppURL, botUsername)
	}
}

// Для личных чатов - INLINE кнопка с WebApp
func (w *WebApp) sendWebAppButtonPrivate(chatID int64, webAppURL string) {
	msg := tgbotapi.NewMessage(chatID,
		`🎭 *Управление ролями*

Используйте кнопку ниже чтобы открыть Web App для создания ролей и управления подписками.`)
	msg.ParseMode = "Markdown"

	// Добавляем chat_id в URL Web App
	webAppURLWithChat := fmt.Sprintf("%s?chat_id=%d", webAppURL, chatID)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonWebApp(
				"📋 Открыть управление ролями",
				tgbotapi.WebAppInfo{URL: webAppURLWithChat},
			),
		),
	)

	msg.ReplyMarkup = keyboard

	if _, err := w.bot.Send(msg); err != nil {
		fmt.Println("Error sending WebApp inline button:", err)
	} else {
		fmt.Println("✅ Sent WebApp inline button successfully")
	}
}

// Для групп - используем глубокие ссылки для мгновенного открытия Web App
func (w *WebApp) sendWebAppButtonGroup(chatID int64, webAppURL, botUsername string) {
	// Получаем информацию о чате
	chat, err := w.bot.GetChat(tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{ChatID: chatID},
	})

	var chatTitle string
	if err == nil && chat.Title != "" {
		chatTitle = chat.Title
	} else {
		chatTitle = fmt.Sprintf("Чат ID: %d", chatID)
	}

	text := fmt.Sprintf(`🎭 *Управление ролями в "%s"*

Создавайте и управляйте ролями specifically для этого чата.

_Выберите способ открытия 👇_`, chatTitle)

	// Создаем глубокую ссылку которая откроет Web App сразу
	deepLink := fmt.Sprintf("https://t.me/%s?startapp=chat%d", botUsername, chatID)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(
				"🚀 Открыть управление ролями",
				deepLink,
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(
				"🚀 Открыть управление ролями2",
				deepLink+"chat",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(
				"🚀 Открыть управление ролями3",
				deepLink+"&user=vasya",
			),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard

	if _, err := w.bot.Send(msg); err != nil {
		fmt.Println("Error sending group WebApp message:", err)
		w.sendGroupFallback(chatID, botUsername)
	} else {
		fmt.Printf("✅ Sent group WebApp deep link for chat '%s' (ID: %d)\n", chatTitle, chatID)
	}
}

// Fallback для групп если Web App не работает
func (w *WebApp) sendGroupFallback(chatID int64, botUsername string) {
	text := `🎭 *Управление ролями в группе*

Чтобы управлять ролями, откройте бота в личных сообщениях.

Нажмите кнопку ниже чтобы перейти к боту 👇`

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(
				"📱 Открыть бота",
				fmt.Sprintf("https://t.me/%s", botUsername),
			),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard

	if _, err := w.bot.Send(msg); err != nil {
		fmt.Println("Error sending group fallback message:", err)
	} else {
		fmt.Println("✅ Sent group fallback button")
	}
}

// Резервная кнопка для личных чатов
func (w *WebApp) sendFallbackButton(chatID int64, webAppURL string) {
	msg := tgbotapi.NewMessage(chatID,
		"🎭 Управление ролями\n\n"+
			"Используйте команды:\n"+
			"/createrole - создать роль\n"+
			"/listroles - список ролей\n"+
			"/myroles - мои подписки")

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(
				"🌐 Открыть Web App в браузере",
				webAppURL,
			),
		),
	)

	msg.ReplyMarkup = keyboard
	w.bot.Send(msg)
}

// Удаляет Reply Keyboard (если она была показана ранее)
func (w *WebApp) RemoveReplyKeyboard(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "⌨️ Убираю клавиатуру...")

	// Создаем пустую клавиатуру с флагом удаления
	removeKeyboard := tgbotapi.NewRemoveKeyboard(true)

	msg.ReplyMarkup = removeKeyboard

	// Отправляем и сразу удаляем сообщение
	if sentMsg, err := w.bot.Send(msg); err == nil {
		// Удаляем сообщение через секунду
		go func() {
			// time.Sleep(1 * time.Second)
			deleteMsg := tgbotapi.NewDeleteMessage(chatID, sentMsg.MessageID)
			w.bot.Send(deleteMsg)
		}()
	}
}
