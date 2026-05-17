package models

type ChatMessage struct {
	ID            int                 `json:"id,omitempty"`
	GID           string              `json:"gid,omitempty"` // Guid context
	Type          string              `json:"type"`          // "text", "image", "voice", "history", "user_list", "edit_msg", "delete_msg", "read_msg", "req_status", "reaction", "reaction_update", "typing"
	Action        string              `json:"action,omitempty"`
	Room          string              `json:"room"`
	Sender        string              `json:"sender,omitempty"`
	SenderUUID    string              `json:"sender_uuid,omitempty"`
	Avatar        string              `json:"avatar,omitempty"`
	Text          string              `json:"text,omitempty"`           // content of message or base64 image/audio
	Emoji         string              `json:"emoji,omitempty"`          // for reaction type
	Timestamp     int64               `json:"timestamp,omitempty"`      // unix timestamp
	Users         []UserInfo          `json:"users,omitempty"`          // For user_list type
	History       []ChatMessage       `json:"history,omitempty"`        // For history type
	Readers       []ReaderInfo        `json:"readers,omitempty"`        // For req_status type response
	Reactions     map[string][]string `json:"reactions,omitempty"`      // emoji -> []user_uuid
	ReplyTo       *ReplyRef           `json:"reply_to,omitempty"`       // Reference to replied message
	VoiceDuration int                 `json:"voice_duration,omitempty"` // Duration in seconds for voice messages
	Edited        bool                `json:"edited,omitempty"`         // Message was edited
	Typing        bool                `json:"typing,omitempty"`         // For typing indicator
}

type ReplyRef struct {
	ID     int    `json:"id"`
	Sender string `json:"sender"`
	Text   string `json:"text"`
}

type ReaderInfo struct {
	UUID     string `json:"uuid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	ReadAt   int64  `json:"read_at"`
}

type UserInfo struct {
	UUID     string `json:"uuid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
