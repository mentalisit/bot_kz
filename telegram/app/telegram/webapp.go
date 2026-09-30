package telegram

import (
	"context"
	"fmt"
	"telegram/models2"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

// UpdateChatMembersCache обновляет кэш участников чата
func (t *Telegram) UpdateChatMembersCache(chatID int64) error {
	chatConfig := tgbotapi.ChatAdministratorsConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: chatID}}
	members, err := t.t.GetChatAdministrators(chatConfig)
	if err != nil {
		return fmt.Errorf("failed to get chat members: %w", err)
	}

	users := make(map[int64]models2.User)
	for _, member := range members {
		user := models2.User{
			ID:        member.User.ID,
			FirstName: member.User.FirstName,
			LastName:  member.User.LastName,
			UserName:  member.User.UserName,
			IsAdmin:   member.IsAdministrator(),
		}
		users[user.ID] = user
	}

	// Обновляем в хранилище
	if err := t.Storage.Db.UpdateUserCache(chatID, users); err != nil {
		return fmt.Errorf("failed to update user cache: %w", err)
	}

	return nil
}

// updateAdminsForChat определяет администраторов чата
func (t *Telegram) updateAdminsForChat(ctx context.Context, chatID int64, users map[int64]models2.User) error {
	chatConfig := tgbotapi.ChatConfig{ChatID: chatID}
	chatConf := tgbotapi.ChatAdministratorsConfig{ChatConfig: chatConfig}
	admins, err := t.t.GetChatAdministrators(chatConf)
	if err != nil {
		return fmt.Errorf("failed to get chat admins: %w", err)
	}

	for _, admin := range admins {
		userID := admin.User.ID
		if user, exists := users[userID]; exists {
			user.IsAdmin = true
			users[userID] = user
		}
	}

	return nil
}
