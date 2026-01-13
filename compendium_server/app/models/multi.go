package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type MultiAccount struct {
	UUID             uuid.UUID   `db:"uuid"`
	Nickname         string      `db:"nickname"`
	TelegramID       string      `db:"telegram_id"`
	TelegramUsername string      `db:"telegram_username"`
	DiscordID        string      `db:"discord_id"`
	DiscordUsername  string      `db:"discord_username"`
	WhatsappID       string      `db:"whatsapp_id"`
	WhatsappUsername string      `db:"whatsapp_username"`
	CreatedAt        time.Time   `db:"created_at"`
	AvatarURL        string      `db:"avatarurl"`
	Alts             StringArray `db:"alts"`
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
func (m *MultiAccount) GetTelegramChatId() int64 {
	parseInt, _ := strconv.ParseInt(m.TelegramID, 10, 64)
	return parseInt
}

type MultiAccountGuildV2 struct {
	GId       uuid.UUID     `db:"gid"`
	GuildName string        `db:"guildname"`
	Channels  GuildChannels `db:"channels"` // Наш новый тип
	AvatarUrl string        `db:"avatarurl"`
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
		return nil
	}
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, m)
}

// UUIDArray позволяет sqlx автоматически работать с UUID[] в Postgres
type UUIDArray []uuid.UUID

func (a UUIDArray) Value() (driver.Value, error) {
	if a == nil {
		return "{}", nil
	}
	return pq.Array(a).Value()
}

func (a *UUIDArray) Scan(src interface{}) error {
	return pq.Array(a).Scan(src)
}

type MultiAccountCorpMember struct {
	Uid        uuid.UUID `db:"uid"`
	GuildIds   UUIDArray `db:"guildids"`
	TimeZona   string    `db:"timezona"`
	ZonaOffset int       `db:"zonaoffset"`
	AfkFor     string    `db:"afkfor"`
}
type Moving struct {
	MAcc       MultiAccount
	Tech       []Technology
	CorpMember MultiAccountCorpMember
}

type StringArray []string

// Метод Scan позволяет стандартному sql.Scan понимать этот тип
func (a *StringArray) Scan(src interface{}) error {
	return pq.Array((*[]string)(a)).Scan(src)
}

// Метод Value позволяет передавать этот тип в запросы без pq.Array()
func (a StringArray) Value() (driver.Value, error) {
	return pq.Array([]string(a)).Value()
}
