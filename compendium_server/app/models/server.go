package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type Identity struct {
	User     User                 `json:"user"`
	Guild    Guild                `json:"guild"`
	Token    string               `json:"token"`
	MAccount *MultiAccount        `json:"mAccount"`
	MGuild   *MultiAccountGuildV2 `json:"mGuild"`
}

func (i Identity) Value() (driver.Value, error) {
	return json.Marshal(i)
}

// Scan реализует интерфейс sql.Scanner для чтения из БД
func (i *Identity) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, i)
}

type Code struct {
	Code      string   `db:"code"`
	Timestamp int64    `db:"time"`
	Identity  Identity `db:"identity"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	//Discriminator string   `json:"discriminator"`
	//Avatar        string   `json:"avatar"`
	AvatarURL string   `json:"avatarUrl"`
	Alts      []string `json:"alts"`
	GameName  string   `json:"gameName"`
}

type Guild struct {
	URL  string `json:"url"`
	ID   string `json:"id"`
	Name string `json:"name"`
	//Icon string `json:"icon"`
	Type string `json:"type"`
}

type TechLevel struct {
	Ts    int64 `json:"ts"`
	Level int   `json:"level"`
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

func (t *TechLevels) ConvertToTech(b []byte) TechLevels {
	m := make(TechLevels)
	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil
	}
	return m
}

type Technology struct {
	Uid  uuid.UUID  `db:"uid"`
	Tech TechLevels `db:"tech"`
	Name string     `db:"username"` // Маппим username из БД в Name структуры
}

type SyncData struct {
	Ver        int        `json:"ver"`        // Версия данных
	InSync     int        `json:"inSync"`     // Флаг синхронизации
	TechLevels TechLevels `json:"techLevels"` // Коллекция уровней технологии
}
type CorpData struct {
	Members    []CorpMember `json:"members"`
	Roles      []CorpRole   `json:"roles"`
	FilterId   string       `json:"filterId"`   // Current filter roleId
	FilterName string       `json:"filterName"` // Name of current filter roleId
}
type CorpRole struct {
	Id       string `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	ChatId   int64  `db:"chat_id" json:"chat_id"`
	TypeRole string `db:"-" json:"type_role,omitempty"`
}

func (r *CorpRole) GetRoleId() int64 {
	i, _ := strconv.ParseInt(r.Id, 10, 64)
	return i
}

func (c *CorpData) Initialization() {
	c.Members = []CorpMember{}
	c.Roles = []CorpRole{{
		Id:   "",
		Name: "@everyone",
	}}
}
func (c *CorpData) AppendEveryone(typeRole string) {
	c.Roles = append(c.Roles, CorpRole{
		Id:       typeRole,
		Name:     fmt.Sprintf("(%s)@everyone", strings.ToUpper(typeRole)),
		TypeRole: typeRole,
	})
}

type CorpMember struct {
	Name        string     `json:"name"`
	UserId      string     `json:"userId"`
	GuildId     string     `json:"guildId"`
	Avatar      string     `json:"avatar"`
	Tech        TechLevels `json:"tech"`
	AvatarUrl   string     `json:"avatarUrl"`
	LocalTime   string     `json:"localTime"`   //localTime:"07:52 PM"
	LocalTime24 string     `json:"localTime24"` //localTime24:"19:52"
	TimeZone    string     `json:"timeZone"`    //timeZone:"UTC-5"
	ZoneOffset  int        `json:"zoneOffset"`  //zoneOffset:-300
	AfkFor      string     `json:"afkFor"`      // readable afk duration
	AfkWhen     int        `json:"afkWhen"`     // Unix Epoch when user returns
	TypeAccount string
	Multi       *MultiAccount
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
