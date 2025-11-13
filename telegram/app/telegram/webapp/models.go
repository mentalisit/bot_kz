package webapp

import (
	"telegram/telegram/types"
	"time"
)

type WebAppData struct {
	User     *types.TelegramUser `json:"user"`
	ChatID   int64               `json:"chat_id"`
	ChatType string              `json:"chat_type"`
	AuthDate int64               `json:"auth_date"`
	QueryID  string              `json:"query_id"`
	Hash     string              `json:"hash"`
}

type UserSession struct {
	UserID    int64
	Username  string
	FirstName string
	LastName  string
	ChatType  string
	LastSeen  time.Time
}

type APIResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}
