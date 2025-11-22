package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"telegram/models"
	"telegram/storage"
	"time"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/gorilla/mux"
)

type WebAppHandler struct {
	storage *storage.Storage
	bot     *tgbotapi.BotAPI
}

func NewWebAppHandler(storage *storage.Storage, bot *tgbotapi.BotAPI) *WebAppHandler {
	return &WebAppHandler{
		storage: storage,
		bot:     bot,
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
		adminDB := make(map[int64]models.User)
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

	var req models.CreateRoleRequest
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

	role := &models.Role{
		ChatID:    chatID,
		Name:      req.Name,
		CreatedBy: userID,
	}

	if err := h.storage.Db.CreateRole(r.Context(), role); err != nil {
		log.Printf("Error creating role: %v", err)
		h.sendError(w, "failed to create role", http.StatusInternalServerError)
		return
	}

	h.sendJSON(w, models.SuccessResponse{
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

	var req models.CreateRoleRequest
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

	h.sendJSON(w, models.SuccessResponse{
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

	h.sendJSON(w, models.SuccessResponse{
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

	h.sendJSON(w, models.SuccessResponse{
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

	h.sendJSON(w, models.SuccessResponse{
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
	json.NewEncoder(w).Encode(models.ErrorResponse{
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
	var roleUsers []models.User
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
