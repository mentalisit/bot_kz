package models

import (
	"database/sql/driver"
	"fmt"
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

func (m *MultiAccountCorpMember) Exist(gid uuid.UUID) bool {
	for _, id := range m.GuildIds {
		if id == gid {
			return true
		}
	}
	return false
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
