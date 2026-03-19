package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"telegram/models2"
	"telegram/storage"
	"time"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mentalisit/logger"
	"github.com/mentalisit/restapi/models"
)

type WebAppHandler struct {
	storage *storage.Storage
	bot     *tgbotapi.BotAPI
	log     *logger.Logger
}

func NewWebAppHandler(storage *storage.Storage, bot *tgbotapi.BotAPI, log *logger.Logger) *WebAppHandler {
	h := &WebAppHandler{
		storage: storage,
		bot:     bot,
		log:     log,
	}
	h.loadConfig()
	return h
}

var DiscordOAuthConfig = OAuthConfig{}

func (h *WebAppHandler) loadConfig() {
	// Discord OAuth конфигурация
	DiscordOAuthConfig = OAuthConfig{
		DiscordClientID:     h.storage.Conf.DiscordClientID,
		DiscordClientSecret: h.storage.Conf.DiscordClientSecret,
		DiscordRedirectURI:  "https://webapp.mentalisit.myds.me/auth/callback/discord",
		DiscordAuthURL:      "https://discord.com/api/oauth2/authorize",
		DiscordTokenURL:     "https://discord.com/api/oauth2/token",
		DiscordUserURL:      "https://discord.com/api/users/@me",
	}
}

// HealthCheck проверяет доступность API
func (h *WebAppHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.sendJSON(w, map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0",
	})
}

// GetUserChats возвращает список чатов пользователя
func (h *WebAppHandler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		h.sendError(w, "user_id is required", http.StatusBadRequest)
		return
	}
	//fmt.Printf("GetUserChats %s\n", userIDStr)

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	chats, err := h.storage.Db.GetUserChats(r.Context(), userID)
	if err != nil {
		log.Printf("Error getting user chats: %v", err)
		h.sendError(w, "failed to get user chats", http.StatusInternalServerError)
		return
	}
	for i, chat := range chats {
		if chat.ChatName == "" {
			getChat, _ := h.bot.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: chat.ChatID}})
			if getChat.Title != "" {
				chats[i].ChatName = getChat.Title
				h.storage.Db.UpdateChatTitle(context.Background(), chat.ChatID, getChat.Title)
			}
		}
		chatConfig := tgbotapi.ChatConfig{ChatID: chat.ChatID}
		chatConf := tgbotapi.ChatAdministratorsConfig{ChatConfig: chatConfig}
		admins, _ := h.bot.GetChatAdministrators(chatConf)
		chatAdmins, _ := h.storage.Db.GetChatAdmins(context.Background(), chat.ChatID)
		adminDB := make(map[int64]models2.User)
		for _, admin := range chatAdmins {
			adminDB[admin.ID] = admin
		}

		for _, admin := range admins {
			if admin.CanChangeInfo || admin.Status == "creator" || admin.CanRestrictMembers {
				_ = h.storage.Db.SetChatAdmin(context.Background(), chat.ChatID, admin.User.ID, true)
				delete(adminDB, admin.User.ID)
			}
		}

		if len(adminDB) != 0 {
			for _, user := range adminDB {
				_ = h.storage.Db.RemoveChatAdmin(context.Background(), chat.ChatID, user.ID)
			}
		}

	}

	h.sendJSON(w, chats)
}

// GetChatRoles возвращает роли чата
func (h *WebAppHandler) GetChatRoles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatIDStr := vars["chatId"]
	userIDStr := r.URL.Query().Get("user_id")
	//fmt.Printf("GetChatRoles %s %s\n", chatIDStr, userIDStr)
	if userIDStr == "" {
		h.sendError(w, "user_id is required", http.StatusBadRequest)
		return
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid chat_id", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	roles, err := h.storage.Db.GetChatRoles(r.Context(), chatID, userID)
	if err != nil {
		log.Printf("Error getting chat roles: %v", err)
		h.sendError(w, "failed to get chat roles", http.StatusInternalServerError)
		return
	}

	h.sendJSON(w, roles)
}

// CreateRole создает новую роль
func (h *WebAppHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatIDStr := vars["chatId"]
	userIDStr := r.URL.Query().Get("user_id")
	//fmt.Printf("CreateRole %s %s\n", chatIDStr, userIDStr)

	if userIDStr == "" {
		h.sendError(w, "user_id is required", http.StatusBadRequest)
		return
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid chat_id", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	// Проверяем права администратора
	isAdmin, err := h.storage.Db.IsChatAdmin(r.Context(), chatID, userID)
	if err != nil {
		log.Printf("Error checking admin status: %v", err)
		h.sendError(w, "failed to check permissions", http.StatusInternalServerError)
		return
	}

	if !isAdmin {
		h.sendError(w, "admin rights required", http.StatusForbidden)
		return
	}

	var req models2.CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		h.sendError(w, "role name is required", http.StatusBadRequest)
		return
	}

	if len(req.Name) > 100 {
		h.sendError(w, "role name too long", http.StatusBadRequest)
		return
	}

	role := &models2.Role{
		ChatID:    chatID,
		Name:      req.Name,
		CreatedBy: userID,
	}

	if err := h.storage.Db.CreateRole(r.Context(), role); err != nil {
		log.Printf("Error creating role: %v", err)
		h.sendError(w, "failed to create role", http.StatusInternalServerError)
		return
	}

	h.sendJSON(w, models2.SuccessResponse{
		Message: "Role created successfully",
		Success: true,
	})
}

// UpdateRole обновляет имя роли
func (h *WebAppHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatIDStr := vars["chatId"]
	roleIDStr := vars["roleId"]
	userIDStr := r.URL.Query().Get("user_id")

	if userIDStr == "" {
		h.sendError(w, "user_id is required", http.StatusBadRequest)
		return
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid chat_id", http.StatusBadRequest)
		return
	}

	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid role_id", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	isAdmin, err := h.storage.Db.IsChatAdmin(r.Context(), chatID, userID)
	if err != nil {
		log.Printf("Error checking admin status: %v", err)
		h.sendError(w, "failed to check permissions", http.StatusInternalServerError)
		return
	}

	if !isAdmin {
		h.sendError(w, "admin rights required", http.StatusForbidden)
		return
	}

	var req models2.CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	newName := req.Name
	if newName == "" {
		h.sendError(w, "role name is required", http.StatusBadRequest)
		return
	}

	if len(newName) > 100 {
		h.sendError(w, "role name too long", http.StatusBadRequest)
		return
	}

	if err := h.storage.Db.UpdateRoleName(r.Context(), roleID, chatID, newName); err != nil {
		log.Printf("Error updating role: %v", err)
		h.sendError(w, "failed to update role", http.StatusInternalServerError)
		return
	}

	h.sendJSON(w, models2.SuccessResponse{
		Message: "Role updated successfully",
		Success: true,
	})
}

// DeleteRole удаляет роль
func (h *WebAppHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatIDStr := vars["chatId"]
	roleIDStr := vars["roleId"]
	userIDStr := r.URL.Query().Get("user_id")
	//fmt.Printf("DeleteRole %s %s %s\n", chatIDStr, userIDStr, roleIDStr)

	if userIDStr == "" {
		h.sendError(w, "user_id is required", http.StatusBadRequest)
		return
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid chat_id", http.StatusBadRequest)
		return
	}

	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid role_id", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	// Проверяем права администратора
	isAdmin, err := h.storage.Db.IsChatAdmin(r.Context(), chatID, userID)
	if err != nil {
		log.Printf("Error checking admin status: %v", err)
		h.sendError(w, "failed to check permissions", http.StatusInternalServerError)
		return
	}

	if !isAdmin {
		h.sendError(w, "admin rights required", http.StatusForbidden)
		return
	}

	if err := h.storage.Db.DeleteRole(r.Context(), roleID, chatID); err != nil {
		log.Printf("Error deleting role: %v", err)
		h.sendError(w, "failed to delete role", http.StatusInternalServerError)
		return
	}

	h.sendJSON(w, models2.SuccessResponse{
		Message: "Role deleted successfully",
		Success: true,
	})
}

// JoinRole добавляет пользователя в роль
func (h *WebAppHandler) JoinRole(w http.ResponseWriter, r *http.Request) {
	h.handleRoleMembership(w, r, true)
}

// LeaveRole удаляет пользователя из роли
func (h *WebAppHandler) LeaveRole(w http.ResponseWriter, r *http.Request) {
	h.handleRoleMembership(w, r, false)
}

func (h *WebAppHandler) handleRoleMembership(w http.ResponseWriter, r *http.Request, join bool) {
	vars := mux.Vars(r)
	chatIDStr := vars["chatId"]
	roleIDStr := vars["roleId"]
	userIDStr := r.URL.Query().Get("user_id")
	//fmt.Printf("handleRoleMembership %s %s %s\n", chatIDStr, userIDStr, roleIDStr)

	if userIDStr == "" {
		h.sendError(w, "user_id is required", http.StatusBadRequest)
		return
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid chat_id", http.StatusBadRequest)
		return
	}

	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid role_id", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	if join {
		if err := h.storage.Db.JoinRole(r.Context(), userID, roleID, chatID); err != nil {
			log.Printf("Error joining role: %v", err)
			h.sendError(w, "failed to join role", http.StatusInternalServerError)
			return
		}
	} else {
		if err := h.storage.Db.LeaveRole(r.Context(), userID, roleID); err != nil {
			log.Printf("Error leaving role: %v", err)
			h.sendError(w, "failed to leave role", http.StatusInternalServerError)
			return
		}
	}

	action := "joined"
	if !join {
		action = "left"
	}

	h.sendJSON(w, models2.SuccessResponse{
		Message: fmt.Sprintf("Successfully %s role", action),
		Success: true,
	})
}

// GetChatUsers возвращает пользователей чата
func (h *WebAppHandler) GetChatUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatIDStr := vars["chatId"]
	//fmt.Printf("GetChatUsers %s\n", chatIDStr)

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid chat_id", http.StatusBadRequest)
		return
	}

	users, err := h.storage.Db.GetChatUsers(r.Context(), chatID)
	if err != nil {
		log.Printf("Error getting chat users: %v", err)
		h.sendError(w, "failed to get chat users", http.StatusInternalServerError)
		return
	}

	h.sendJSON(w, users)
}

// SetUserRole назначает/снимает роль пользователя (для админов)
func (h *WebAppHandler) SetUserRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatIDStr := vars["chatId"]
	userIDStr := vars["userId"]
	roleIDStr := vars["roleId"]
	adminIDStr := r.URL.Query().Get("admin_id")
	//fmt.Printf("SetUserRole %s %s %s %+v\n", chatIDStr, userIDStr, roleIDStr, adminIDStr)

	if adminIDStr == "" {
		h.sendError(w, "admin_id is required", http.StatusBadRequest)
		return
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid chat_id", http.StatusBadRequest)
		return
	}

	targetUserID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid role_id", http.StatusBadRequest)
		return
	}

	adminID, err := strconv.ParseInt(adminIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid admin_id", http.StatusBadRequest)
		return
	}

	// Проверяем права администратора
	isAdmin, err := h.storage.Db.IsChatAdmin(r.Context(), chatID, adminID)
	if err != nil {
		log.Printf("Error checking admin status: %v", err)
		h.sendError(w, "failed to check permissions", http.StatusInternalServerError)
		return
	}

	if !isAdmin {
		h.sendError(w, "admin rights required", http.StatusForbidden)
		return
	}

	assign := r.Method == "POST"

	if err := h.storage.Db.SetUserRole(r.Context(), targetUserID, roleID, chatID, assign); err != nil {
		log.Printf("Error setting user role: %v", err)
		h.sendError(w, "failed to set user role", http.StatusInternalServerError)
		return
	}

	action := "assigned"
	if !assign {
		action = "removed"
	}

	h.sendJSON(w, models2.SuccessResponse{
		Message: fmt.Sprintf("Role %s successfully", action),
		Success: true,
	})
}

// GetUserPermissions возвращает права пользователя в чате
func (h *WebAppHandler) GetUserPermissions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatIDStr := vars["chatId"]
	userIDStr := r.URL.Query().Get("user_id")
	//fmt.Printf("GetUserPermissions %s %s\n", chatIDStr, userIDStr)

	if userIDStr == "" {
		h.sendError(w, "user_id is required", http.StatusBadRequest)
		return
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid chat_id", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	isAdmin, err := h.storage.Db.IsChatAdmin(r.Context(), chatID, userID)
	if err != nil {
		log.Printf("Error checking permissions: %v", err)
		h.sendError(w, "failed to check permissions", http.StatusInternalServerError)
		return
	}

	h.sendJSON(w, map[string]interface{}{
		"is_admin": isAdmin,
		"chat_id":  chatID,
		"user_id":  userID,
	})
}

// Вспомогательные методы
func (h *WebAppHandler) sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (h *WebAppHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models2.ErrorResponse{
		Error: message,
	})
}

// GetRoleMembers возвращает участников конкретной роли
func (h *WebAppHandler) GetRoleMembers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatIDStr := vars["chatId"]
	roleIDStr := vars["roleId"]

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid chat_id", http.StatusBadRequest)
		return
	}

	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid role_id", http.StatusBadRequest)
		return
	}

	// Получаем всех пользователей чата
	users, err := h.storage.Db.GetChatUsers(r.Context(), chatID)
	if err != nil {
		log.Printf("Error getting chat users: %v", err)
		h.sendError(w, "failed to get chat users", http.StatusInternalServerError)
		return
	}

	// Получаем информацию о роли
	var roleName string
	err = h.storage.Db.GetRoleName(r.Context(), roleID, &roleName)
	if err != nil {
		log.Printf("Error getting role name: %v", err)
		// Продолжаем выполнение, даже если не получили имя роли
	}

	// Фильтруем пользователей по роли
	var roleUsers []models2.User
	if roleName == "all" {
		// Для роли "all" возвращаем всех пользователей
		roleUsers = users
	} else {
		// Для обычных ролей фильтруем по наличию роли
		for _, user := range users {
			if user.Roles[roleID] {
				roleUsers = append(roleUsers, user)
			}
		}
	}

	h.sendJSON(w, roleUsers)
}

// authDiscord перенаправляет на Discord OAuth
func (h *WebAppHandler) AuthDiscord(w http.ResponseWriter, r *http.Request) {
	// Получаем Telegram initData из параметров запроса
	initData := r.URL.Query().Get("init_data")
	if initData == "" {
		initData = r.Header.Get("X-Telegram-Init-Data")
	}

	// Создаем state с включенными Telegram данными
	stateData := map[string]string{
		"id":       uuid.New().String(),
		"initData": initData,
	}
	stateJSON, _ := json.Marshal(stateData)
	state := base64.StdEncoding.EncodeToString(stateJSON)

	discordOAuthURL := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&response_type=code&scope=identify&state=%s",
		DiscordOAuthConfig.DiscordAuthURL,
		DiscordOAuthConfig.DiscordClientID,
		url.QueryEscape(DiscordOAuthConfig.DiscordRedirectURI),
		state,
	)

	http.Redirect(w, r, discordOAuthURL, http.StatusFound)
}

// authDiscordCallback обрабатывает callback от Discord OAuth
func (h *WebAppHandler) AuthDiscordCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		log.Printf("No authorization code received from Discord")
		http.Error(w, "Не получен authorization code", http.StatusBadRequest)
		return
	}

	// Извлекаем Telegram данные из state
	var telegramInitData string
	if state != "" {
		stateData, err := base64.StdEncoding.DecodeString(state)
		if err == nil {
			var stateMap map[string]string
			if json.Unmarshal(stateData, &stateMap) == nil {
				telegramInitData = stateMap["initData"]
			}
		}
	}

	// Обмениваем code на access token
	tokenResp, err := h.exchangeDiscordCode(code)
	if err != nil {
		log.Printf("Failed to exchange Discord code: %v", err)
		http.Error(w, "Ошибка авторизации Discord", http.StatusInternalServerError)
		return
	}

	// Получаем информацию о пользователе
	userResp, err := h.getDiscordUser(tokenResp.AccessToken)
	if err != nil {
		log.Printf("Failed to get Discord user info: %v", err)
		http.Error(w, "Не удалось получить информацию о пользователе", http.StatusInternalServerError)
		return
	}

	// Создаем токен для фронтенда
	token := fmt.Sprintf("discord_%s", userResp.ID)

	// Формируем URL для редиректа с данными Discord и Telegram
	params := url.Values{}
	params.Set("discord_token", token)
	params.Set("discord_id", userResp.ID)
	params.Set("discord_name", userResp.Username)

	if telegramInitData != "" {
		params.Set("telegram_init_data", telegramInitData)
	}

	redirectURL := "/auth/link?" + params.Encode()

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// exchangeDiscordCode обменивает authorization code на access token
func (h *WebAppHandler) exchangeDiscordCode(code string) (*DiscordTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", DiscordOAuthConfig.DiscordRedirectURI)
	data.Set("client_id", DiscordOAuthConfig.DiscordClientID)
	data.Set("client_secret", DiscordOAuthConfig.DiscordClientSecret)

	req, err := http.NewRequest("POST", DiscordOAuthConfig.DiscordTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("не удалось обменять code на token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ошибка Discord API [%d]: %s", resp.StatusCode, string(body))
	}

	var tokenResp DiscordTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("не удалось декодировать ответ: %w", err)
	}

	return &tokenResp, nil
}

// getDiscordUser получает информацию о пользователе Discord
func (h *WebAppHandler) getDiscordUser(accessToken string) (*DiscordUserResponse, error) {
	req, err := http.NewRequest("GET", DiscordOAuthConfig.DiscordUserURL, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить информацию о пользователе: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ошибка Discord User API [%d]: %s", resp.StatusCode, string(body))
	}

	var userResp DiscordUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, fmt.Errorf("не удалось декодировать ответ: %w", err)
	}

	return &userResp, nil
}

// AuthLinkPage показывает страницу с информацией об авторизации Discord
func (h *WebAppHandler) AuthLinkPage(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры из URL
	discordToken := r.URL.Query().Get("discord_token")
	discordID := r.URL.Query().Get("discord_id")
	discordName := r.URL.Query().Get("discord_name")
	telegramInitData := r.URL.Query().Get("telegram_init_data")

	// Парсим Telegram initData для получения информации о пользователе
	var telegramUser map[string]interface{}
	var telegramUserID int64
	var telegramUsername string
	var telegramFirstName string
	var telegramLastName string

	if telegramInitData != "" {
		// initData имеет формат: user={json}&chat_instance=...&...
		// Нужно правильно декодировать URL и извлечь JSON пользователя

		// Сначала декодируем URL
		decodedInitData, err := url.QueryUnescape(telegramInitData)
		if err != nil {
			log.Printf("Failed to decode telegram initData: %v", err)
		} else {
			// Ищем user= в декодированных данных
			if strings.Contains(decodedInitData, "user=") {
				// Разбираем параметры
				params := strings.Split(decodedInitData, "&")
				for _, param := range params {
					if strings.HasPrefix(param, "user=") {
						userJSON := strings.TrimPrefix(param, "user=")

						// Декодируем JSON пользователя (может быть двойное кодирование)
						userJSON, err = url.QueryUnescape(userJSON)
						if err != nil {
							log.Printf("Failed to decode user JSON: %v", err)
							continue
						}

						// Парсим JSON
						if json.Unmarshal([]byte(userJSON), &telegramUser) == nil {
							if id, ok := telegramUser["id"].(float64); ok {
								telegramUserID = int64(id)
							}
							if username, ok := telegramUser["username"].(string); ok {
								telegramUsername = username
							}
							if firstName, ok := telegramUser["first_name"].(string); ok {
								telegramFirstName = firstName
							}
							if lastName, ok := telegramUser["last_name"].(string); ok {
								telegramLastName = lastName
							}
						} else {
							log.Printf("Failed to parse user JSON: %s", userJSON)
						}
						break
					}
				}
			}
		}
	}

	// Формируем HTML форму для ввода основного никнейма
	html := `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Discord авторизация</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #5865F2, #7289DA);
            color: white;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            margin: 0;
            padding: 20px;
        }
        .container {
            background: rgba(255, 255, 255, 0.1);
            backdrop-filter: blur(10px);
            border-radius: 20px;
            padding: 30px;
            max-width: 600px;
            width: 100%;
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
        }
        h1 {
            text-align: center;
            margin-bottom: 30px;
            font-size: 28px;
        }
        .data-display {
            background: rgba(255, 255, 255, 0.1);
            border-radius: 10px;
            padding: 20px;
            margin-bottom: 30px;
            font-family: monospace;
            font-size: 14px;
            line-height: 1.5;
        }
        .form-group {
            margin-bottom: 20px;
        }
        label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
            color: #00ff88;
        }
        input[type="text"] {
            width: 100%;
            padding: 12px;
            border: 2px solid rgba(255, 255, 255, 0.3);
            border-radius: 8px;
            background: rgba(255, 255, 255, 0.1);
            color: white;
            font-size: 16px;
            box-sizing: border-box;
        }
        input[type="text"]::placeholder {
            color: rgba(255, 255, 255, 0.6);
        }
        input[type="text"]:focus {
            outline: none;
            border-color: #00ff88;
            background: rgba(255, 255, 255, 0.2);
        }
        .btn {
            background: #00ff88;
            color: #5865F2;
            border: none;
            padding: 15px 30px;
            border-radius: 10px;
            font-size: 16px;
            font-weight: bold;
            cursor: pointer;
            width: 100%;
            transition: all 0.3s ease;
        }
        .btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 20px rgba(0, 0, 0, 0.2);
        }
        .btn:disabled {
            opacity: 0.6;
            cursor: not-allowed;
            transform: none;
        }
        .status {
            text-align: center;
            margin-top: 20px;
            padding: 10px;
            border-radius: 8px;
            display: none;
        }
        .status.success {
            background: rgba(0, 255, 136, 0.2);
            color: #00ff88;
        }
        .status.error {
            background: rgba(255, 0, 0, 0.2);
            color: #ff6b6b;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎯 Финальный шаг авторизации</h1>

        <div class="data-display">`

	if discordID != "" {
		html += fmt.Sprintf("discord_id=%s<br>", discordID)
	}
	if discordName != "" {
		html += fmt.Sprintf("discord_name=%s<br>", discordName)
	}

	if telegramUserID != 0 {
		html += "<br>"
		html += fmt.Sprintf("telegram_id=%d<br>", telegramUserID)
		if telegramFirstName != "" {
			html += fmt.Sprintf("first_name=%s<br>", telegramFirstName)
		}
		if telegramLastName != "" {
			html += fmt.Sprintf("last_name=%s<br>", telegramLastName)
		}
		if telegramUsername != "" {
			html += fmt.Sprintf("username=%s<br>", telegramUsername)
		}
	}

	html += `</div>

        <form id="nicknameForm">
            <div class="form-group">
                <label for="mainNickname">Основной никнейм:</label>
                <input type="text" id="mainNickname" name="main_nickname"
                       placeholder="Введите ваш основной никнейм" required>
            </div>

            <button type="submit" class="btn" id="submitBtn">Отправить данные</button>
        </form>

        <div id="statusMessage" class="status"></div>
    </div>

    <script>
        // Инициализация Telegram WebApp
        if (window.Telegram && window.Telegram.WebApp) {
            window.Telegram.WebApp.ready();
            window.Telegram.WebApp.expand();
            console.log('Telegram WebApp initialized');
        } else {
            console.log('Not in Telegram WebApp environment');
        }

        const form = document.getElementById('nicknameForm');
        const submitBtn = document.getElementById('submitBtn');
        const statusMessage = document.getElementById('statusMessage');

        form.addEventListener('submit', async function(e) {
            e.preventDefault();

            const mainNickname = document.getElementById('mainNickname').value.trim();
            if (!mainNickname) {
                showStatus('Пожалуйста, введите основной никнейм', 'error');
                return;
            }

            submitBtn.disabled = true;
            submitBtn.textContent = 'Отправка...';

            try {
                const response = await fetch('/api/auth/link', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        discord_id: '` + discordID + `',
                        discord_name: '` + discordName + `',
                        discord_token: '` + discordToken + `',
                        telegram_id: ` + fmt.Sprintf("%d", telegramUserID) + `,
                        first_name: '` + telegramFirstName + `',
                        last_name: '` + telegramLastName + `',
                        username: '` + telegramUsername + `',
                        main_nickname: mainNickname
                    })
                });

                if (response.ok) {
                    showStatus('✅ Данные успешно отправлены! Перенаправление в редактор ролей...', 'success');
                    console.log('Data sent successfully, will redirect to main page');


                    // Перенаправление на главную страницу через 2 секунды
                    setTimeout(() => {
                        returnToMainPage();
                    }, 2000);

                    // Запасной вариант - перенаправление на главную через 5 секунд
                    setTimeout(() => {
                        console.log('Fallback: redirecting to main page');
                        window.location.href = 'https://webapp.mentalisit.myds.me/';
                    }, 5000);
                } else {
                    const error = await response.text();
                    showStatus('❌ Ошибка: ' + error, 'error');
                }
            } catch (error) {
                showStatus('❌ Ошибка сети: ' + error.message, 'error');
            } finally {
                submitBtn.disabled = false;
                submitBtn.textContent = 'Отправить данные';
            }
        });

        function showStatus(message, type) {
            statusMessage.textContent = message;
            statusMessage.className = 'status ' + type;
            statusMessage.style.display = 'block';
        }

        function returnToMainPage() {
            console.log('Returning to main page...');

            if (window.Telegram && window.Telegram.WebApp) {
                try {
                    // Сначала попробуем отправить данные боту
                    if (typeof window.Telegram.WebApp.sendData === 'function') {
                        window.Telegram.WebApp.sendData(JSON.stringify({
                            action: 'auth_complete',
                            discord_id: '` + discordID + `',
                            telegram_id: ` + fmt.Sprintf("%d", telegramUserID) + `,
                            main_nickname: 'submitted'
                        }));
                        console.log('Data sent to bot');
                    }

                    // Через небольшую задержку перенаправляем на главную страницу
                    setTimeout(() => {
                        window.location.href = 'https://webapp.mentalisit.myds.me/';
                    }, 500);

                } catch (error) {
                    console.error('Error sending data:', error);
                    // Если не получилось отправить данные, просто перенаправляем
                    window.location.href = '/';
                }
            } else {
                // Для тестирования вне Telegram - перенаправление на главную
                console.log('Not in Telegram WebApp, redirecting to main page');
                window.location.href = 'https://webapp.mentalisit.myds.me/';
            }
        }

        function closeWebApp() {
            returnToMainPage();
        }
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func (h *WebAppHandler) SubmitAuthData(w http.ResponseWriter, r *http.Request) {
	var data struct {
		DiscordID    string `json:"discord_id"`
		DiscordName  string `json:"discord_name"`
		DiscordToken string `json:"discord_token"`
		TelegramID   int64  `json:"telegram_id"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Username     string `json:"username"`
		MainNickname string `json:"main_nickname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Failed to decode auth data: %v", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	telegramIDStr := fmt.Sprintf("%d", data.TelegramID)

	// Проверяем, существует ли уже аккаунт для этого Telegram пользователя
	existingAccount, err := h.storage.Db.FindMultiAccountByTelegramID(telegramIDStr)
	if err != nil {
		log.Printf("Error checking existing account: %v", err)
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	var account *models.MultiAccount

	if existingAccount != nil {
		// Обновляем существующий аккаунт Discord данными
		existingAccount.DiscordID = data.DiscordID
		existingAccount.DiscordUsername = data.DiscordName
		existingAccount.Nickname = data.MainNickname

		discordAcc, _ := h.storage.Db.FindMultiAccountByDiscordID(data.DiscordID)
		if discordAcc != nil {
			// Объединяем аккаунты
			account, err = h.mergeAccounts(existingAccount, discordAcc)
			if err != nil {
				http.Error(w, "Ошибка объединения аккаунтов", http.StatusInternalServerError)
				return
			}
		} else {
			h.log.InfoStruct("*existingAccount ", *existingAccount)

			updatedAccount, err := h.storage.Db.UpdateMultiAccount(*existingAccount)
			if err != nil {
				log.Printf("Error updating account: %v", err)
				http.Error(w, "Ошибка обновления аккаунта", http.StatusInternalServerError)
				return
			}
			account = updatedAccount
		}

	} else {
		// Создаем новый аккаунт
		newAccount := models.MultiAccount{
			UUID:             uuid.New(),
			Nickname:         data.MainNickname,
			TelegramID:       telegramIDStr,
			TelegramUsername: data.Username,
			DiscordID:        data.DiscordID,
			DiscordUsername:  data.DiscordName,
			AvatarURL:        "",
			Alts:             []string{},
		}

		createdAccount, err := h.storage.Db.CreateMultiAccount(newAccount)
		h.log.InfoStruct("CreateMultiAccount ", createdAccount)
		if err != nil {
			log.Printf("Error creating account: %v", err)
			http.Error(w, "Ошибка создания аккаунта", http.StatusInternalServerError)
			return
		}
		account = createdAccount
	}

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Данные успешно сохранены. Возврат в редактор ролей...",
		"uuid":    account.UUID.String(),
	})
}

// CheckDiscordData проверяет, есть ли у Telegram пользователя Discord данные
func (h *WebAppHandler) CheckDiscordData(w http.ResponseWriter, r *http.Request) {
	telegramIDStr := r.URL.Query().Get("telegram_id")
	if telegramIDStr == "" {
		http.Error(w, "Не указан telegram_id", http.StatusBadRequest)
		return
	}

	// Ищем аккаунт в базе данных по Telegram ID
	account, err := h.storage.Db.FindMultiAccountByTelegramID(telegramIDStr)
	if err != nil {
		log.Printf("Error finding multi account by telegram ID %s: %v", telegramIDStr, err)
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	var discordData map[string]interface{}

	if account != nil && account.DiscordID != "" {
		// Discord аккаунт найден
		discordData = map[string]interface{}{
			"has_discord":  true,
			"discord_id":   account.DiscordID,
			"discord_name": account.DiscordUsername,
			"nickname":     account.Nickname,
			"uuid":         account.UUID.String(),
		}
		// Добавляем аватар только если есть ссылка
		if account.AvatarURL != "" {
			discordData["avatar"] = account.AvatarURL
		}
	} else {
		// Discord аккаунт не найден
		discordData = map[string]interface{}{
			"has_discord": false,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(discordData)
}

// GetCorpMembers возвращает участников корпорации из всех источников
func (h *WebAppHandler) GetCorpMembers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatIDStr := vars["chatId"]

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid chat_id", http.StatusBadRequest)
		return
	}

	// Получаем участников из всех источников
	var allMembers []models2.CompendiumCorpMember

	// My Compendium
	if members, err := h.storage.Db.GetCorpMembersMyCompendium(r.Context(), chatID); err == nil {
		allMembers = append(allMembers, members...)
	} else {
		log.Printf("Error getting corp members from my_compendium: %v", err)
	}

	// HS Compendium
	if members, err := h.storage.Db.GetCorpMembersHSCompendium(r.Context(), chatID); err == nil {
		allMembers = append(allMembers, members...)
	} else {
		log.Printf("Error getting corp members from hs_compendium: %v", err)
	}

	h.sendJSON(w, allMembers)
}

// RemoveCorpMember удаляет участника из корпорации по указанному источнику
func (h *WebAppHandler) RemoveCorpMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatIDStr := vars["chatId"]
	userIDStr := vars["userId"]

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "invalid chat_id", http.StatusBadRequest)
		return
	}

	// Получаем источник из query параметра
	tableSource := r.URL.Query().Get("tableSource")
	if tableSource == "" {
		h.sendError(w, "tableSource parameter is required", http.StatusBadRequest)
		return
	}
	fmt.Printf("RemoveCorpMember %s %d %s\n", tableSource, chatID, userIDStr)
	// Выбираем функцию удаления в зависимости от источника
	var removeErr error
	switch tableSource {
	case "my_compendium":
		removeErr = h.storage.Db.RemoveCorpMemberMyCompendium(r.Context(), chatID, userIDStr)
	case "hs_compendium":
		removeErr = h.storage.Db.RemoveCorpMemberHSCompendium(r.Context(), chatID, userIDStr)
	default:
		h.sendError(w, "invalid tableSource parameter", http.StatusBadRequest)
		return
	}

	if removeErr != nil {
		log.Printf("Error removing corp member from %s: %v", tableSource, removeErr)
		h.sendError(w, "failed to remove corp member", http.StatusInternalServerError)
		return
	}

	h.sendJSON(w, map[string]string{"status": "success"})
}
