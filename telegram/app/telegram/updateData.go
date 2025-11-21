package telegram

import (
	"context"
	"telegram/models"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func (t *Telegram) SaveMember(c *tgbotapi.Chat, u *tgbotapi.User) {
	chat := models.Chat{
		ChatID:   c.ID,
		ChatName: c.Title,
	}
	if u.IsBot {
		return
	}

	user := models.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		UserName:  u.UserName,
		IsAdmin:   false,
	}
	t.UpdateUserFromMessage(chat, user)
}

// UpdateUserFromMessage обновляет данные пользователя из сообщения
func (t *Telegram) UpdateUserFromMessage(chat models.Chat, user models.User) {
	// Сначала убедимся, что чат существует в бд
	_, err := t.Storage.Db.GetChat(chat.ChatID)
	if err != nil {
		// Если чат не существует, создаем его
		go func() {
			if chat.ChatName != "" {
				_ = t.Storage.Db.CreateOrUpdateChat(context.Background(), chat.ChatID, chat.ChatName)
			}
		}()
	}

	// Обновляем/добавляем пользователя в чат
	dbUser := models.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UserName:  user.UserName,
		IsAdmin:   false, // По умолчанию не админ, будет обновлено отдельно
	}
	_ = t.Storage.Db.AddUserToChat(context.Background(), chat.ChatID, dbUser)
}
