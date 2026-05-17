package models2

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/google/uuid"
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

// EscapeMarkdownV2 экранирует специальные символы для MarkdownV2
func EscapeMarkdownV2(text string) string {
	// Список специальных символов в MarkdownV2
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}

	escaped := text
	for _, char := range specialChars {
		escaped = strings.ReplaceAll(escaped, char, "\\"+char)
	}
	return escaped
}

type ChatAccess struct {
	ChatID   int64  `json:"chat_id"`
	ChatName string `json:"chat_name"`
	UserID   int64  `json:"user_id"`
	Status   string `json:"status"`
}

type MultiAccountGuildV2 struct {
	GId       uuid.UUID
	GuildName string
	Channels  GuildChannels `db:"channels"` // Наш новый тип
	AvatarUrl string
}
type GuildChannels map[string][]string

// Value преобразует map в JSON для базы данных
func (m GuildChannels) Value() (driver.Value, error) {
	if m == nil {
		return json.Marshal(map[string][]string{})
	}
	return json.Marshal(m)
}

// Scan преобразует JSON из базы данных в map
func (m *GuildChannels) Scan(src interface{}) error {
	if src == nil {
		*m = make(GuildChannels)
		return nil
	}
	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("unsupported type for GuildChannels: %T", src)
	}
	return json.Unmarshal(data, m)
}

// UUIDArray позволяет sqlx автоматически работать с UUID[] в Postgres
type UUIDArray []uuid.UUID

// Value: превращает слайс UUID в PostgreSQL array literal для базы (Valuer)
func (u UUIDArray) Value() (driver.Value, error) {
	if u == nil || len(u) == 0 {
		return "{}", nil
	}
	strs := make([]string, len(u))
	for i, id := range u {
		strs[i] = id.String()
	}
	return "{" + strings.Join(strs, ",") + "}", nil
}

// Scan: читает PostgreSQL UUID[] массив из базы в слайс UUID (Scanner)
func (u *UUIDArray) Scan(src interface{}) error {
	if src == nil {
		*u = make(UUIDArray, 0)
		return nil
	}
	var source string
	switch v := src.(type) {
	case []byte:
		source = string(v)
	case string:
		source = v
	default:
		return fmt.Errorf("unsupported type for UUIDArray: %T", src)
	}
	// PostgreSQL array format: {uuid1,uuid2,...}
	s := strings.Trim(source, "{}")
	if s == "" {
		*u = make(UUIDArray, 0)
		return nil
	}
	parts := strings.Split(s, ",")
	result := make(UUIDArray, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		id, err := uuid.Parse(p)
		if err != nil {
			return fmt.Errorf("failed to parse UUID %q: %w", p, err)
		}
		result = append(result, id)
	}
	*u = result
	return nil
}
