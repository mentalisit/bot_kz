package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Identity struct {
	User   User                 `json:"user"`
	Guild  Guild                `json:"guild"`
	Token  string               `json:"token"`
	MAcc   *MultiAccount        `json:"mAcc"`
	MGuild *MultiAccountGuildV2 `json:"mGuild"`
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

//type IdentityGET struct {
//	User  User    `json:"user"`
//	Guild []Guild `json:"guilds"`
//	Token string  `json:"token"`
//}

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

func (t *TechLevels) ConvertToTech(b []byte) map[int]TechLevel {
	//var m map[int]TechLevel
	m := make(map[int]TechLevel)
	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil
	}
	return m
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
type CorpMember struct {
	Name        string     `db:"username" json:"name"` // Колонка username в базу -> Name в Go
	UserId      string     `db:"uid" json:"userId"`    // UUID в базу -> string в Go (sqlx сконвертирует)
	GuildId     string     `json:"guildId"`            // Это поле обычно заполняется отдельно
	Tech        TechLevels `db:"tech" json:"tech"`     // Наш умный тип JSONB
	AvatarUrl   string     `db:"avatarurl" json:"avatarUrl"`
	LocalTime   string     `json:"localTime"`   //localTime:"07:52 PM"
	LocalTime24 string     `json:"localTime24"` //localTime24:"19:52"
	TimeZone    string     `json:"timeZone"`    //timeZone:"UTC-5"
	ZoneOffset  int        `json:"zoneOffset"`  //zoneOffset:-300
	AfkFor      string     `json:"afkFor"`      // readable afk duration
	AfkWhen     int        `json:"afkWhen"`     // Unix Epoch when user returns
	MAcc        *MultiAccount
	MGuild      *MultiAccountGuildV2
}

type WsKill struct {
	Id           int64  `db:"id" json:"id"` // Добавим Id, так как он есть в БД
	GuildId      string `db:"guildid"`
	ChatId       string `db:"chatid"`
	UserName     string `db:"username"`
	Mention      string `db:"mention"`
	ShipName     string `db:"shipname"`
	TimestampEnd int64  `db:"timestampend"`
	Language     string `db:"language"`
}
type TechTable struct {
	Id      int64
	Name    string
	NameId  string
	GuildId string
	Tech    []byte
}
