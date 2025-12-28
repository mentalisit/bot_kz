package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
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

func (m *MultiAccountGuild) GuildId() string {
	if m == nil {
		return ""
	}
	return m.GId.String()
}

func (m *MultiAccountGuild) GetMapChannel() map[string][]string {
	channels := make(map[string][]string)
	for _, channel := range m.Channels {
		if strings.Contains(channel, "@") {
			channels["wa"] = append(channels["wa"], channel)
		} else if strings.HasPrefix(channel, "-100") {
			channels["tg"] = append(channels["tg"], channel)
		} else {
			channels["ds"] = append(channels["ds"], channel)
		}
	}
	return channels
}

type MultiAccountCorpMember struct {
	Uid        uuid.UUID
	GuildIds   []uuid.UUID
	TimeZona   string
	ZonaOffset int
	AfkFor     string
}

func (m *MultiAccountCorpMember) Exist(gid uuid.UUID) bool {
	for _, id := range m.GuildIds {
		if id == gid {
			return true
		}
	}
	return false
}
