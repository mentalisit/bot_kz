package webapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"telegram/telegram/roles"
	"telegram/telegram/types"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

type Handlers struct {
	bot       *tgbotapi.BotAPI
	auth      *AuthManager
	roles     *roles.Manager
	templates *template.Template
}

func (h *Handlers) handleWebApp(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("WebApp accessed r.URL.RawQuery %+v\n", r.URL.RawQuery)

	// Получаем ID чата из параметров
	chatIDStr := strings.TrimPrefix(r.URL.Query().Get("tgWebAppStartParam"), "chat")

	var chatID int64 = 0

	if chatIDStr != "" {
		if id, err := strconv.ParseInt(chatIDStr, 10, 64); err == nil {
			chatID = id
		}
	}

	log.Printf("Rendering WebApp for chat ID: %d", chatID)
	h.renderWebAppPage(w, chatID)
}

func (h *Handlers) renderWebAppPage(w http.ResponseWriter, chatID int64) {
	if h.roles == nil {
		h.renderErrorPage(w, "Система ролей не инициализирована")
		return
	}

	var rolesList []*types.Role
	var chatTitle string = "Общие роли"

	if chatID != 0 {
		// Роли для конкретного чата
		rolesList = h.roles.GetChatRoles(chatID)
		chatTitle = fmt.Sprintf("Роли чата (ID: %d)", chatID)
	} else {
		// Общие роли (без привязки к чату)
		rolesList = h.roles.GetAllRoles()
	}
	fmt.Printf("chatTitle:%s rolesList %+v\n", chatTitle, rolesList)

	data := struct {
		Roles     []*types.Role
		ChatID    int64
		ChatTitle string
	}{
		Roles:     rolesList,
		ChatID:    chatID,
		ChatTitle: chatTitle,
	}

	w.Header().Set("Content-Type", "text/html")

	if err := h.templates.Execute(w, data); err != nil {
		log.Printf("Template error: %v", err)
		w.Write([]byte(DefaultHTMLTemplate))
	}
}

func (h *Handlers) handleRolesAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("Roles API called")
	fmt.Printf("%+v\n", r.URL.RawQuery)

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	chatIDStr := r.URL.Query().Get("chat_id")

	var userID int64
	var chatID int64

	if userIDStr != "" {
		fmt.Sscanf(userIDStr, "%d", &userID)
	}
	if chatIDStr != "" {
		fmt.Sscanf(chatIDStr, "%d", &chatID)
	}
	fmt.Printf("UserID: %d, ChatID: %d\n", userID, chatID)
	if h.roles == nil {
		log.Println("Roles manager is nil!")
		h.sendJSONError(w, "Roles manager not initialized", http.StatusInternalServerError)
		return
	}

	var rolesList []*types.Role
	if chatID != 0 {
		// Роли для конкретного чата
		rolesList = h.roles.GetChatRoles(chatID)
	} else {
		// Все роли
		rolesList = h.roles.GetAllRoles()
	}

	log.Printf("Found %d roles for chat ID: %d", len(rolesList), chatID)

	type RoleWithSubscription struct {
		*types.Role
		Subscribed bool `json:"subscribed"`
	}

	result := make([]RoleWithSubscription, len(rolesList))
	for i, role := range rolesList {
		result[i] = RoleWithSubscription{
			Role:       role,
			Subscribed: h.roles.IsUserSubscribed(userID, role.ID),
		}
	}

	h.sendJSONSuccess(w, result)
}

func (h *Handlers) handleCreateRoleAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("Create role API called")
	fmt.Printf("%+v\n", r.URL.RawQuery)

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		CreatedBy   int64  `json:"created_by"`
		ChatID      int64  `json:"chat_id"`
		ChatTitle   string `json:"chat_title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		h.sendJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		h.sendJSONError(w, "Role name is required", http.StatusBadRequest)
		return
	}

	if h.roles == nil {
		h.sendJSONError(w, "Roles manager not initialized", http.StatusInternalServerError)
		return
	}

	log.Printf("Creating role: %s by user %d for chat %s", req.Name, req.CreatedBy, req.ChatTitle)

	role := h.roles.CreateRole(req.Name, req.Description, req.CreatedBy, req.ChatID, req.ChatTitle)

	// Уведомляем создателя
	go h.notifyRoleCreation(req.CreatedBy, role)

	log.Printf("Role created successfully: %s", role.ID)
	h.sendJSONSuccess(w, role)
}

func NewHandlers(bot *tgbotapi.BotAPI, auth *AuthManager, rolesManager *roles.Manager) *Handlers {
	// Создаем шаблон из встроенной строки
	tmpl := template.Must(template.New("index").Parse(DefaultHTMLTemplate))

	log.Println("WebApp handlers initialized with built-in template")

	return &Handlers{
		bot:       bot,
		auth:      auth,
		roles:     rolesManager,
		templates: tmpl,
	}
}

func (h *Handlers) Start() {
	http.HandleFunc("/", h.handleWebApp)
	http.HandleFunc("/api/roles", h.handleRolesAPI)
	http.HandleFunc("/api/subscribe", h.handleSubscribeAPI)
	http.HandleFunc("/api/unsubscribe", h.handleUnsubscribeAPI)
	http.HandleFunc("/api/create-role", h.handleCreateRoleAPI)
	http.HandleFunc("/api/delete-role", h.handleDeleteRoleAPI)
	http.HandleFunc("/api/user-info", h.handleUserInfo)

	log.Println("Web App server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (h *Handlers) handleSubscribeAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID int64  `json:"user_id"`
		RoleID string `json:"role_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if h.roles == nil {
		h.sendJSONError(w, "Roles manager not initialized", http.StatusInternalServerError)
		return
	}

	if h.roles.SubscribeToRole(req.UserID, req.RoleID) {
		go h.notifyUserAboutSubscription(req.UserID, req.RoleID, true)
		h.sendJSONSuccess(w, nil)
	} else {
		h.sendJSONError(w, "Role not found", http.StatusNotFound)
	}
}

func (h *Handlers) handleUnsubscribeAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID int64  `json:"user_id"`
		RoleID string `json:"role_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if h.roles == nil {
		h.sendJSONError(w, "Roles manager not initialized", http.StatusInternalServerError)
		return
	}

	if h.roles.UnsubscribeFromRole(req.UserID, req.RoleID) {
		go h.notifyUserAboutSubscription(req.UserID, req.RoleID, false)
		h.sendJSONSuccess(w, nil)
	} else {
		h.sendJSONError(w, "Role not found", http.StatusNotFound)
	}
}

func (h *Handlers) handleDeleteRoleAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		RoleID string `json:"role_id"`
		UserID int64  `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if h.roles == nil {
		h.sendJSONError(w, "Roles manager not initialized", http.StatusInternalServerError)
		return
	}

	if h.roles.DeleteRole(req.RoleID, req.UserID) {
		h.sendJSONSuccess(w, nil)
	} else {
		h.sendJSONError(w, "Role not found or access denied", http.StatusForbidden)
	}
}

func (h *Handlers) handleUserInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Читаем тело для отладки
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.sendJSONError(w, "Error reading body", http.StatusBadRequest)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	fmt.Printf("📨 Raw Body: %s\n", string(bodyBytes))

	// Исправленная структура
	var request struct {
		User     *types.TelegramUser `json:"user"`
		InitData string              `json:"initData"`
		Chat     interface{}         `json:"chat"`    // для групповых чатов
		ChatID   int64               `json:"chat_id"` // ← ДОБАВЬ ЭТО ПОЛЕ!
		ChatType string              `json:"chat_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.sendJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("👤 User: %+v\n", request.User)
	fmt.Printf("🔐 InitData: %s\n", request.InitData)
	fmt.Printf("💬 Chat: %+v\n", request.Chat)
	fmt.Printf("🆔 ChatID: %d\n", request.ChatID) // ← теперь будет правильное значение
	fmt.Printf("📱 ChatType: %s\n", request.ChatType)

	// Определяем реальный chat_id
	var realChatID int64
	if request.ChatID != 0 {
		realChatID = request.ChatID
	} else if request.Chat != nil {
		// Парсим chat объект если он есть
		if chatMap, ok := request.Chat.(map[string]interface{}); ok {
			if id, exists := chatMap["id"]; exists {
				if idFloat, ok := id.(float64); ok {
					realChatID = int64(idFloat)
				}
			}
		}
	}

	fmt.Printf("🎯 Real ChatID: %d\n", realChatID)

	if request.User != nil {
		h.auth.SaveUserSession(request.User)
		log.Printf("📱 WebApp opened by: @%s (ID: %d) in chat: %d",
			request.User.Username, request.User.ID, realChatID)
	}

	// Используй realChatID в дальнейшей логике
	h.sendJSONSuccess(w, nil)
}

func (h *Handlers) sendJSONSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(APIResponse{
		Status: "success",
		Data:   data,
	})
}

func (h *Handlers) sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(APIResponse{
		Status:  "error",
		Message: message,
	})
}

func (h *Handlers) notifyUserAboutSubscription(userID int64, roleID string, subscribed bool) {
	role := h.roles.GetRole(roleID)
	if role == nil {
		return
	}

	action := "отписался"
	if subscribed {
		action = "подписался"
	}

	message := fmt.Sprintf("🎭 Вы %s от роли \"%s\"", action, role.Name)
	msg := tgbotapi.NewMessage(userID, message)
	h.bot.Send(msg)
}

func (h *Handlers) notifyRoleCreation(userID int64, role *types.Role) {
	message := fmt.Sprintf("🎭 Роль \"%s\" успешно создана!", role.Name)
	msg := tgbotapi.NewMessage(userID, message)
	h.bot.Send(msg)
}

func (h *Handlers) renderErrorPage(w http.ResponseWriter, message string) {
	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>Ошибка</title>
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<style>
			body { font-family: Arial, sans-serif; padding: 20px; text-align: center; }
			.error { color: #dc3545; margin: 20px 0; }
		</style>
	</head>
	<body>
		<h1>⚠️ Ошибка</h1>
		<div class="error">%s</div>
		<p>Пожалуйста, попробуйте открыть Web App через Telegram бота.</p>
	</body>
	</html>`, message)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
