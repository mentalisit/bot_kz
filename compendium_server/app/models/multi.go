package models

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type MultiAccount struct {
	UUID             uuid.UUID
	Nickname         string
	TelegramID       string
	TelegramUsername string
	DiscordID        string
	DiscordUsername  string
	WhatsappID       string
	WhatsappUsername string
	CreatedAt        time.Time
	AvatarURL        string
	Alts             []string
}

func (m *MultiAccount) GetTextUsername() string {
	text := fmt.Sprintf("Твой текущий НикНейм %s\n", m.Nickname)
	if m.DiscordUsername != "" {
		text = text + "Discord UserName: " + m.DiscordUsername + "\n"
	}
	if m.TelegramUsername != "" {
		text = text + "Telegram UserName: " + m.TelegramUsername + "\n"
	}
	if m.WhatsappUsername != "" {
		text = text + "Whatsapp UserName: " + m.WhatsappUsername + "\n"
	}
	return text
}

type AccountLinkCode struct {
	Code      string
	UUID      uuid.UUID
	ExpiresAt time.Time
	CreatedAt time.Time
}

type MultiAccountGuild struct {
	GId       uuid.UUID
	GuildName string
	Channels  []string
	AvatarUrl string
}

type MultiAccountCorpMember struct {
	Uid        uuid.UUID
	GuildIds   []uuid.UUID
	TimeZona   string
	ZonaOffset int
	AfkFor     string
}
