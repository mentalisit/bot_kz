package models

import (
	"encoding/json"
	"time"
)

type Identity struct {
	User  User   `json:"user"`
	Guild Guild  `json:"guild"`
	Token string `json:"token"`
	//Type  string `json:"type"`
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
}

type Guild struct {
	URL  string `json:"url"`
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
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
	Name       string         `json:"name"`
	UserId     string         `json:"userId"`
	GuildId    string         `json:"guildId"`
	Avatar     string         `json:"avatar"`
	Tech       TechLevelArray `json:"tech"`
	AvatarUrl  string         `json:"avatarUrl"`
	TimeZone   string         `json:"timeZone"`
	LocalTime  string         `json:"localTime"`
	ZoneOffset int            `json:"zoneOffset"` // TZ offset in minutes
	AfkFor     string         `json:"afkFor"`     // readable afk duration
	AfkWhen    int            `json:"afkWhen"`    // Unix Epoch when user returns
}

type CorpRole struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type Corporation struct {
	Name string `json:"Name"`
	Id   string `json:"Id"`
}
type FilterCorps struct {
	Corp []Corporation `json:"Corp"`
}
type Corp struct {
	MaxPage int           `json:"MaxPage"`
	Matches []Corporation `json:"matches"`
}

type CorporationsData struct {
	Corporation1Name  string    `json:"Corporation1Name"`
	Corporation1Id    string    `json:"Corporation1Id"`
	Corporation2Name  string    `json:"Corporation2Name"`
	Corporation2Id    string    `json:"Corporation2Id"`
	Corporation1Score int       `json:"Corporation1Score"`
	Corporation2Score int       `json:"Corporation2Score"`
	DateEnded         time.Time `json:"DateEnded"`
	MatchId           string    `json:"MatchId"`
}
type Match struct {
	Corporation1Name  string    `json:"Corporation1Name"`
	Corporation1Id    string    `json:"Corporation1Id"`
	Corporation2Name  string    `json:"Corporation2Name"`
	Corporation2Id    string    `json:"Corporation2Id"`
	Corporation1Score int       `json:"Corporation1Score"`
	Corporation2Score int       `json:"Corporation2Score"`
	DateEnded         time.Time `json:"DateEnded"`
}
type Ws struct {
	MaxPage int                `json:"MaxPage"`
	Matches []CorporationsData `json:"matches"`
}
