package webapp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"sync"
	"telegram/telegram/types"
	"time"
)

type AuthManager struct {
	userSessions map[int64]*UserSession
	mux          sync.RWMutex
	botToken     string
}

func NewAuthManager(botToken string) *AuthManager {
	return &AuthManager{
		userSessions: make(map[int64]*UserSession),
		botToken:     botToken,
	}
}

// Парсинг и валидация данных Web App
func (a *AuthManager) ParseAndValidateWebAppData(initData string) (*WebAppData, error) {
	if initData == "" {
		return nil, fmt.Errorf("empty init data")
	}

	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("invalid init data format")
	}

	// Валидируем подпись
	if err := a.validateInitData(values); err != nil {
		return nil, fmt.Errorf("signature validation failed: %v", err)
	}

	return a.parseWebAppData(values)
}

func (a *AuthManager) parseWebAppData(values url.Values) (*WebAppData, error) {
	var webAppData WebAppData

	// Парсим пользователя
	if userStr := values.Get("user"); userStr != "" {
		if err := json.Unmarshal([]byte(userStr), &webAppData.User); err != nil {
			return nil, fmt.Errorf("invalid user data")
		}
	}

	if webAppData.User == nil {
		return nil, fmt.Errorf("user data not found")
	}

	// Парсим остальные поля
	webAppData.QueryID = values.Get("query_id")
	webAppData.ChatType = values.Get("chat_type")
	webAppData.Hash = values.Get("hash")

	if chatInstance := values.Get("chat_instance"); chatInstance != "" {
		fmt.Sscanf(chatInstance, "%d", &webAppData.ChatID)
	}
	if authDate := values.Get("auth_date"); authDate != "" {
		fmt.Sscanf(authDate, "%d", &webAppData.AuthDate)
	}

	return &webAppData, nil
}

// Валидация данных инициализации
func (a *AuthManager) validateInitData(values url.Values) error {
	receivedHash := values.Get("hash")
	if receivedHash == "" {
		return fmt.Errorf("hash not found")
	}

	// Удаляем хеш из значений для проверки
	values.Del("hash")

	// Сортируем ключи
	var keys []string
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Создаем data-check-string
	var dataCheckStrings []string
	for _, key := range keys {
		value := values.Get(key)
		dataCheckStrings = append(dataCheckStrings, fmt.Sprintf("%s=%s", key, value))
	}
	dataCheckString := strings.Join(dataCheckStrings, "\n")

	// Вычисляем секретный ключ
	secretKey := hmac.New(sha256.New, []byte("WebAppData"))
	secretKey.Write([]byte(a.botToken))

	// Вычисляем хеш
	computedHash := hex.EncodeToString(
		hmac.New(sha256.New, secretKey.Sum(nil)).
			Sum([]byte(dataCheckString)),
	)

	// Сравниваем хеши
	if computedHash != receivedHash {
		return fmt.Errorf("hash mismatch")
	}

	return nil
}

// Управление сессиями
func (a *AuthManager) SaveUserSession(user *types.TelegramUser) {
	a.mux.Lock()
	defer a.mux.Unlock()

	a.userSessions[user.ID] = &UserSession{
		UserID:    user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		LastSeen:  time.Now(),
	}
}

func (a *AuthManager) GetUserSession(userID int64) *UserSession {
	a.mux.RLock()
	defer a.mux.RUnlock()

	return a.userSessions[userID]
}
