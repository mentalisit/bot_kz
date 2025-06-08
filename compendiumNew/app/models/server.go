package models

import (
	"encoding/json"
)

type Identity struct {
	User         User               `json:"user"`
	Guild        Guild              `json:"guild"`
	Token        string             `json:"token"`
	MultiAccount *MultiAccount      `json:"multiAccount"`
	MultiGuild   *MultiAccountGuild `json:"multiGuild"`
}
type Code struct {
	Code      string
	Timestamp int64
	Identity  Identity
}

//type IdentityGET struct {
//	User  User    `json:"user"`
//	Guild []Guild `json:"guilds"`
//	Token string  `json:"token"`
//}

type User struct {
	ID            string   `json:"id"`
	Username      string   `json:"username"`
	Discriminator string   `json:"discriminator"`
	Avatar        string   `json:"avatar"`
	AvatarURL     string   `json:"avatarUrl"`
	Alts          []string `json:"alts"`
	GameName      string   `json:"gameName"`
}

type Guild struct {
	URL  string `json:"url"`
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Type string `json:"type"`
}

type TechLevel struct {
	Ts    int64 `json:"ts"`
	Level int   `json:"level"`
}
type TechLevels map[int]TechLevel
type TechLevelArray map[int][2]int

func (a *TechLevelArray) ConvertToTech(b []byte) TechLevelArray {
	//var m map[int]TechLevel
	m := make(map[int]TechLevel)
	err := json.Unmarshal(b, &m)
	if err != nil {
		return TechLevelArray{}
	}
	var mi = make(TechLevelArray)
	for i, le := range m {
		mi[i] = [2]int{le.Level}
	}
	return mi
}
func (a *TechLevelArray) ConvertToByte() []byte {
	techBytes, _ := json.Marshal(a)
	return techBytes
}
func (l *TechLevels) ConvertToTech(b []byte) map[int]TechLevel {
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
	Name         string         `json:"name"`
	UserId       string         `json:"userId"`
	GuildId      string         `json:"guildId"`
	Avatar       string         `json:"avatar"`
	Tech         TechLevelArray `json:"tech"`
	AvatarUrl    string         `json:"avatarUrl"`
	LocalTime    string         `json:"localTime"`   //localTime:"07:52 PM"
	LocalTime24  string         `json:"localTime24"` //localTime24:"19:52"
	TimeZone     string         `json:"timeZone"`    //timeZone:"UTC-5"
	ZoneOffset   int            `json:"zoneOffset"`  //zoneOffset:-300
	AfkFor       string         `json:"afkFor"`      // readable afk duration
	AfkWhen      int            `json:"afkWhen"`     // Unix Epoch when user returns
	MultiAccount *MultiAccount
	MultiGuild   *MultiAccountGuild
}

type CorpRole struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type WsKill struct {
	GuildId      string
	ChatId       string
	UserName     string
	Mention      string
	ShipName     string
	TimestampEnd int64
	Language     string
}
type TechTable struct {
	Id      int64
	Name    string
	NameId  string
	GuildId string
	Tech    []byte
}
