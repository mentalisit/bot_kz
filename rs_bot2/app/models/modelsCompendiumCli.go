package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// CorpData represents corporation data structure
type CorpData struct {
	Members    []CorpMember        `json:"members"`
	Roles      []CorpRole          `json:"roles"`
	FilterID   string              `json:"filterId"`   // Current filter roleId
	FilterName string              `json:"filterName"` // Name of current filter roleId
	MGuild     MultiAccountGuildV2 `json:"MGuild"`
}

func (c *CorpData) Initialization() {
	c.Members = []CorpMember{}
	c.Roles = []CorpRole{{
		ID:   "",
		Name: "@everyone",
	}}
}

func (c *CorpData) AppendEveryone(typeRole string) {
	c.Roles = append(c.Roles, CorpRole{
		ID:       typeRole,
		Name:     fmt.Sprintf("@everyone (%s)", strings.ToUpper(typeRole)),
		TypeRole: typeRole,
	})
}

// CorpMember represents a member of a corporation.
type CorpMember struct {
	Name         string        `json:"name"`
	UserID       string        `json:"userId"`
	ClientUserID string        `json:"clientUserId"`
	Avatar       string        `json:"avatar"`
	Tech         TechLevels    `json:"tech"`
	AvatarURL    string        `json:"avatarUrl"`
	LocalTime    string        `json:"localTime"`   // localTime:"07:52 PM"
	LocalTime24  string        `json:"localTime24"` // localTime24:"19:52"
	TimeZone     string        `json:"timeZone"`    // timeZone:"UTC-5"
	ZoneOffset   int           `json:"zoneOffset"`  // zoneOffset:-300
	AfkFor       string        `json:"afkFor"`      // readable afk duration
	AfkWhen      int           `json:"afkWhen"`     // Unix Epoch when user returns
	TypeAccount  string        `json:"typeAccount,omitempty"`
	Multi        *MultiAccount `json:"multi,omitempty"`
}

func (v *CorpMember) GetType() string {
	if v.Multi == nil {
		return ""
	}
	if v.Multi.DiscordID != "" && v.Multi.TelegramID != "" && v.Multi.WhatsappID != "" {
		return "ma"
	}
	if v.Multi.TelegramID != "" {
		return "tg"
	}
	if v.Multi.DiscordID != "" {
		return "ds"
	}
	if v.Multi.WhatsappID != "" {
		return "wa"
	}
	return ""
}

type Technology struct {
	Uid  uuid.UUID  `db:"uid"`
	Tech TechLevels `db:"tech"`
	Name string     `db:"username"` // Маппим username из БД в Name структуры
}

// CorpRole represents a corporation role data structure
type CorpRole struct {
	ID       string `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	ChatID   int64  `db:"chat_id" json:"chat_id"`
	TypeRole string `db:"-" json:"typeRole,omitempty"`
}

func (r *CorpRole) GetRoleId() int64 {
	i, _ := strconv.ParseInt(r.ID, 10, 64)
	return i
}

type TechLevels map[int]TechLevel

// Value: превращает мапу в JSON для базы (Valuer)
func (t TechLevels) Value() (driver.Value, error) {
	if t == nil {
		return json.Marshal(make(TechLevels))
	}
	return json.Marshal(t)
}

// Scan: читает JSON из базы в мапу (Scanner)
func (t *TechLevels) Scan(src interface{}) error {
	if src == nil {
		*t = make(TechLevels)
		return nil
	}
	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("unsupported type for TechLevels: %T", src)
	}
	return json.Unmarshal(data, t)
}

type TechLevel struct {
	Level int   `json:"level"`
	Ts    int64 `json:"ts"`
}

type TechUser struct {
	Name string
	Tech TechLevels
}

type CompendiumTechReq struct {
	Uuid  string `json:"uuid"`
	Name  string `json:"name"`
	Id    string `json:"id"`
	Level int    `json:"level"`
}

type MultiAccountCorpMember struct {
	Uid        uuid.UUID `db:"uid"`
	GuildIds   UUIDArray `db:"guildids"`
	TimeZona   string    `db:"timezona"`
	ZonaOffset int       `db:"zonaoffset"`
	AfkFor     string    `db:"afkfor"`
}

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

// Remove удаляет элемент по значению и обновляет слайс напрямую
func (a *UUIDArray) Remove(gid uuid.UUID) {
	if a == nil || len(*a) == 0 {
		return
	}

	for i, v := range *a {
		if v == gid {
			// Перезаписываем дескриптор слайса по указателю
			*a = append((*a)[:i], (*a)[i+1:]...)
			return
		}
	}
}

type MultiAccountGuildV2 struct {
	GId       uuid.UUID     `db:"gid"`
	GuildName string        `db:"guildname"`
	Channels  GuildChannels `db:"channels"` // Наш новый тип
	AvatarUrl string        `db:"avatarurl"`
	Data      DataGuild     `db:"data" json:"data,omitempty"`
}
type DataGuild struct {
	PollChannels []Channel            `json:"poll_channels,omitempty"`
	Coordination map[string][]Channel `json:"coordination,omitempty"`
	Discussion   map[string][]Channel `json:"discussion,omitempty"`
	GameGuilds   []GameGuilds         `json:"game_guilds,omitempty"`
}
type GameGuilds struct {
	CorpName string `json:"corp_name"`
	CorpId   string `json:"corp_id"`
	General  bool   `json:"general"`
}

// Value implements driver.Valuer for DataGuild (JSONB)
func (d DataGuild) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements sql.Scanner for DataGuild (JSONB)
func (d *DataGuild) Scan(src interface{}) error {
	if src == nil {
		*d = DataGuild{}
		return nil
	}
	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("unsupported type for DataGuild: %T", src)
	}
	return json.Unmarshal(data, d)
}

type Channel struct {
	ChannelID     string `json:"channel_id"`
	ChannelName   string `json:"channel_name"`
	TypeMessenger string `json:"type_messenger"`
	Description   string `json:"description,omitempty"`
}

type GuildChannels map[string][]string

// Value: превращает мапу в JSON для базы (Valuer)
func (g GuildChannels) Value() (driver.Value, error) {
	if g == nil {
		return json.Marshal(make(GuildChannels))
	}
	return json.Marshal(g)
}

// Scan: читает JSON из базы в мапу (Scanner)
func (g *GuildChannels) Scan(src interface{}) error {
	if src == nil {
		*g = make(GuildChannels)
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
	return json.Unmarshal(data, g)
}

type DsMembersRoles struct {
	Userid  string
	RolesId []string
}
type Study struct {
	Uuid    uuid.UUID    `json:"uuid"`
	Name    string       `json:"name"`
	Studies StudiesArray `json:"studies"`
}

type StudiesArray []Studies

// Value: превращает слайс в JSON для базы (Valuer)
func (s StudiesArray) Value() (driver.Value, error) {
	if s == nil {
		return json.Marshal(make(StudiesArray, 0))
	}
	return json.Marshal(s)
}

// Scan: читает JSON из базы в слайс (Scanner)
func (s *StudiesArray) Scan(src interface{}) error {
	if src == nil {
		*s = make(StudiesArray, 0)
		return nil
	}
	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("unsupported type for StudiesArray: %T", src)
	}
	return json.Unmarshal(data, s)
}

type Studies struct {
	ModuleId string `json:"moduleId"`
	EndTime  int64  `json:"endTime"`
	Level    int    `json:"level"`
	Name     string `json:"name"`
}
