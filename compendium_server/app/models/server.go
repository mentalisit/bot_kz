package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Identity struct {
	User         User                 `json:"user"`
	Guild        Guild                `json:"guild"`
	Token        string               `json:"token"`
	MultiAccount *MultiAccount        `json:"multiAccount"`
	MAccount     *MultiAccount        `json:"mAccount"`
	MGuild       *MultiAccountGuildV2 `json:"mGuild"`
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

//type TechLevelArray map[int][2]int

type Technology struct {
	Tech TechLevels
	Name string
}

//func (a *TechLevelArray) ConvertToTech(b []byte) TechLevelArray {
//	//var m map[int]TechLevel
//	m := make(map[int]TechLevel)
//	err := json.Unmarshal(b, &m)
//	if err != nil {
//		return TechLevelArray{}
//	}
//	var mi = make(TechLevelArray)
//	for i, le := range m {
//		mi[i] = [2]int{le.Level}
//	}
//	return mi
//}

func (l *TechLevels) ConvertToTech(b []byte) TechLevels {
	m := make(TechLevels)
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

//type WsKill struct {
//	GuildId      string
//	ChatId       string
//	UserName     string
//	Mention      string
//	ShipName     string
//	TimestampEnd int64
//	Language     string
//}
