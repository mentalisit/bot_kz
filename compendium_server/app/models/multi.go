package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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
