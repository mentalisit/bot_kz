package telegram

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"telegram/models"
	"telegram/telegram/handlers"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

// StartWebApp запускает веб-сервер для веб-приложения
func (t *Telegram) StartWebApp(port string) {
	router := handlers.SetupRouter(t.Storage, t.t, t.log)

	t.server = &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Starting web app server on port %s", port)

	go func() {
		if err := t.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start web server: %v", err)
		}
	}()
}

// StopWebApp останавливает веб-сервер
func (t *Telegram) StopWebApp() {
	if t.server != nil {
		t.server.Shutdown(context.Background())
	}
}

// UpdateChatMembersCache обновляет кэш участников чата
func (t *Telegram) UpdateChatMembersCache(chatID int64) error {
	chatConfig := tgbotapi.ChatAdministratorsConfig{tgbotapi.ChatConfig{ChatID: chatID}}
	members, err := t.t.GetChatAdministrators(chatConfig)
	if err != nil {
		return fmt.Errorf("failed to get chat members: %w", err)
	}

	users := make(map[int64]models.User)
	for _, member := range members {
		user := models.User{
			ID:        member.User.ID,
			FirstName: member.User.FirstName,
			LastName:  member.User.LastName,
			UserName:  member.User.UserName,
			IsAdmin:   member.IsAdministrator(),
		}
		users[user.ID] = user
	}

	// Обновляем в хранилище
	ctx := context.Background()
	if err := t.Storage.Db.UpdateUserCache(ctx, chatID, users); err != nil {
		return fmt.Errorf("failed to update user cache: %w", err)
	}

	return nil
}

// updateAdminsForChat определяет администраторов чата
func (t *Telegram) updateAdminsForChat(ctx context.Context, chatID int64, users map[int64]models.User) error {
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
