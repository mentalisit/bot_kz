package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// MultiAccountGuildV2 - версия для V2 с Channels как map[string][]string
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
		return nil
	}
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, m)
}

func (m *MultiAccountGuildV2) ChannelsBytes() []byte {
	channelsBytes, _ := json.Marshal(m.Channels)
	return channelsBytes
}
func (m *MultiAccountGuildV2) GuildId() string {
	if m == nil {
		return ""
	}
	return m.GId.String()
}

type CorpMemberV2 struct {
	Name        string     `json:"name"`
	UserUUID    string     `json:"userUuid"`
	GuildUUID   string     `json:"guildUuid"`
	Avatar      string     `json:"avatar"`
	Tech        TechLevels `json:"tech"`
	AvatarUrl   string     `json:"avatarUrl"`
	LocalTime   string     `json:"localTime"`   //localTime:"07:52 PM"
	LocalTime24 string     `json:"localTime24"` //localTime24:"19:52"
	TimeZone    string     `json:"timeZone"`    //timeZone:"UTC-5"
	ZoneOffset  int        `json:"zoneOffset"`  //zoneOffset:-300
	AfkFor      string     `json:"afkFor"`      // readable afk duration
	AfkWhen     int        `json:"afkWhen"`     // Unix Epoch when user returns
	MAcc        *MultiAccount
}

func (v *CorpMemberV2) GetType() string {
	if v.MAcc == nil {
		return ""
	}
	if v.MAcc.DiscordID != "" && v.MAcc.TelegramID != "" {
		return "ma"
	}
	if v.MAcc.TelegramID != "" {
		return "tg"
	}
	if v.MAcc.DiscordID != "" {
		return "ds"
	}
	if v.MAcc.WhatsappID != "" {
		return "wa"
	}
	return ""
}

type CorpDataV2 struct {
	Members    []CorpMemberV2 `json:"members"`
	Roles      []CorpRole     `json:"roles"`
	FilterId   string         `json:"filterId"`   // Current filter roleId
	FilterName string         `json:"filterName"` // Name of current filter roleId
}

func (c *CorpDataV2) Initialization() {
	c.Members = []CorpMemberV2{}
	c.Roles = []CorpRole{{
		Id:   "",
		Name: "@everyone",
	}}
}
func (c *CorpDataV2) AppendEveryone(typeRole string) {
	c.Roles = append(c.Roles, CorpRole{
		Id:       typeRole,
		Name:     fmt.Sprintf("(%s)@everyone", strings.ToUpper(typeRole)),
		TypeRole: typeRole,
	})
}

type CorpRole struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	TypeRole string
	ChatId   int64
}

func (r *CorpRole) GetRoleId() int64 {
	i, _ := strconv.ParseInt(r.Id, 10, 64)
	return i
}
