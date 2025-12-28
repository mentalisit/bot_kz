package handlers

import (
	"net/http"
	"telegram/storage"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/gorilla/mux"
	"github.com/mentalisit/logger"
)

func SetupRouter(storage *storage.Storage, bot *tgbotapi.BotAPI, log *logger.Logger) http.Handler {
	handler := NewWebAppHandler(storage, bot, log)

	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api").Subrouter()

	// Discord OAuth routes
	router.HandleFunc("/auth/discord", handler.AuthDiscord).Methods("GET")
	router.HandleFunc("/auth/callback/discord", handler.AuthDiscordCallback).Methods("GET")
	router.HandleFunc("/auth/link", handler.AuthLinkPage).Methods("GET")

	// API для отправки данных авторизации
	api.HandleFunc("/auth/link", handler.SubmitAuthData).Methods("POST")

	// API для проверки Discord данных пользователя
	api.HandleFunc("/auth/check-discord", handler.CheckDiscordData).Methods("GET")
	// В разделе Chat routes добавьте:
	api.HandleFunc("/chat/{chatId}/roles/{roleId}/members", handler.GetRoleMembers).Methods("GET")

	// User routes
	api.HandleFunc("/user/chats", handler.GetUserChats).Methods("GET")

	router.HandleFunc("/api/health", handler.HealthCheck).Methods("GET")

	// Chat routes
	api.HandleFunc("/chat/{chatId}/roles", handler.GetChatRoles).Methods("GET")
	api.HandleFunc("/chat/{chatId}/roles", handler.CreateRole).Methods("POST")
	api.HandleFunc("/chat/{chatId}/roles/{roleId}", handler.UpdateRole).Methods("PUT")
	api.HandleFunc("/chat/{chatId}/roles/{roleId}", handler.DeleteRole).Methods("DELETE")
	api.HandleFunc("/chat/{chatId}/roles/{roleId}/join", handler.JoinRole).Methods("POST")
	api.HandleFunc("/chat/{chatId}/roles/{roleId}/leave", handler.LeaveRole).Methods("POST")
	api.HandleFunc("/chat/{chatId}/users", handler.GetChatUsers).Methods("GET")
	api.HandleFunc("/chat/{chatId}/users/{userId}/roles/{roleId}", handler.SetUserRole).Methods("POST", "DELETE")
	api.HandleFunc("/chat/{chatId}/permissions", handler.GetUserPermissions).Methods("GET")

	// Serve static files (frontend)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./docker/templates/")))

	return router
}
