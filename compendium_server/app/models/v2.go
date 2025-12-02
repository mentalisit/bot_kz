package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// MultiAccountGuildV2 - версия для V2 с Channels как map[string]string
type MultiAccountGuildV2 struct {
	GId       uuid.UUID
	GuildName string
	Channels  map[string]string
	AvatarUrl string
}

// IdentityV2 - версия для V2 с MultiAccountGuildV2
type IdentityV2 struct {
	Token        string        `json:"token"`
	MultiAccount *MultiAccount `json:"multiAccount"`
	GuildId      string        `json:"guild_id"`
}

func (i *IdentityV2) GetGuildUUID() *uuid.UUID {
	gid, _ := uuid.Parse(i.GuildId)
	return &gid
}

type CorpMemberV2 struct {
	Name        string         `json:"name"`
	UserUUID    string         `json:"userUuid"`
	GuildUUID   string         `json:"guildUuid"`
	Avatar      string         `json:"avatar"`
	Tech        TechLevelArray `json:"tech"`
	AvatarUrl   string         `json:"avatarUrl"`
	LocalTime   string         `json:"localTime"`   //localTime:"07:52 PM"
	LocalTime24 string         `json:"localTime24"` //localTime24:"19:52"
	TimeZone    string         `json:"timeZone"`    //timeZone:"UTC-5"
	ZoneOffset  int            `json:"zoneOffset"`  //zoneOffset:-300
	AfkFor      string         `json:"afkFor"`      // readable afk duration
	AfkWhen     int            `json:"afkWhen"`     // Unix Epoch when user returns
	Multi       *MultiAccount
}

func (v *CorpMemberV2) GetType() string {
	if v.Multi == nil {
		return ""
	}
	if v.Multi.DiscordID != "" && v.Multi.TelegramID != "" {
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
