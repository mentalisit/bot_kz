package models

import (
	"fmt"
	"time"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

type Chat struct {
	ChatID   int64  `json:"chat_id"`
	ChatName string `json:"chat_name"`
}
type CompendiumCorpMember struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	AvatarURL   string `json:"avatar_url,omitempty"`
	TableSource string `json:"table_source"` // "my_compendium", "compendium", "hs_compendium"
}

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"user_name"`
	IsAdmin   bool   `json:"is_admin"`
	// Добавляем поле для хранения ролей пользователя
	Roles map[int64]bool `json:"roles,omitempty"`
}

func (u *User) GetUserName() string {
	if u.UserName != "" {
		return u.UserName
	}

	name := u.FirstName
	if u.LastName != "" {
		name += " " + u.LastName
	}

	return name
}

func (u *User) FormatMention() string {
	if u.UserName != "" {
		return "@" + u.UserName
	}
	return fmt.Sprintf("[%s](tg://user?id=%d)", EscapeMarkdownV2(u.GetUserName()), u.ID)
}
func (u *User) FormatMentionAll() string {
	return fmt.Sprintf("[%s](tg://user?id=%d)", EscapeMarkdownV2(u.GetUserName()), u.ID)
}

func (u *User) TgUser(user *tgbotapi.User) {
	u.ID = user.ID
	u.FirstName = user.FirstName
	u.LastName = user.LastName
	u.UserName = user.UserName
}

type Role struct {
	ID          int64     `json:"id"`
	ChatID      int64     `json:"chat_id"`
	Name        string    `json:"name"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	MemberCount int       `json:"member_count,omitempty"`
	IsMember    bool      `json:"is_member,omitempty"`
}

type UserRole struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
	ChatID int64 `json:"chat_id"`
}

type ChatPermission struct {
	ChatID  int64 `json:"chat_id"`
	UserID  int64 `json:"user_id"`
	IsAdmin bool  `json:"is_admin"`
}

type CreateRoleRequest struct {
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// ChatMember представляет участника чата в базе данных
type ChatMember struct {
	ChatID    int64  `json:"chat_id"`
	UserID    int64  `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"user_name"`
	IsAdmin   bool   `json:"is_admin"`
}

type ChatBackup struct {
	ChatID   int64     `json:"chat_id"`
	ChatName string    `json:"chat_name"`
	Members  []User    `json:"members"`
	Roles    []Role    `json:"roles"`
	BackupAt time.Time `json:"backup_at"`
}
