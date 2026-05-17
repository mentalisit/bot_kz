package models

import "time"

type Chat struct {
	ChatID   int64  `json:"chat_id"`
	ChatName string `json:"chat_name"`
}

// ChatUser represents a user in a chat with their roles
type ChatUser struct {
	ID        int64          `json:"id"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	UserName  string         `json:"user_name"`
	IsAdmin   bool           `json:"is_admin"`
	Roles     map[int64]bool `json:"roles,omitempty"`
}

// Role represents a role in a chat
type Role struct {
	ID          int64     `json:"id"`
	ChatID      int64     `json:"chat_id"`
	Name        string    `json:"name"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	MemberCount int       `json:"member_count,omitempty"`
	IsMember    bool      `json:"is_member,omitempty"`
}

// CreateRoleRequest is the request body for creating a role
type CreateRoleRequest struct {
	Name string `json:"name"`
}

// CompendiumCorpMember represents a corp member from compendium tables
type CompendiumCorpMember struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	AvatarURL   string `json:"avatar_url,omitempty"`
	TableSource string `json:"table_source"` // "my_compendium", "hs_compendium"
}

type ChatAccess struct {
	ChatID   int64  `json:"chat_id"`
	ChatName string `json:"chat_name"`
	UserID   int64  `json:"user_id"`
	Status   string `json:"status"`
}
