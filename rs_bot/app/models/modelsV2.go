package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type InMessageV2 struct {
	Text        string
	Tip         string
	NameNick    string
	Username    string
	UserId      string
	NameMention string
	Messenger   Info
	Config      CorporationConfigV2
	Options     Options
}

func (i *InMessageV2) GetNameMention() string {
	if i.NameMention == "@" {
		return fmt.Sprintf("[%s](tg://user?id=%s)", i.Username, i.UserId)
	}
	return i.NameMention
}

type Rs struct {
	RsTypeLevel string
	TimeRs      string
	AltName     string
	Money       bool
}

func (i *Rs) GetTimeRs() int64 {
	timeRs, err := strconv.ParseInt(i.TimeRs, 10, 64)
	if err != nil {
		return 0
	}

	return timeRs
}

func (i *Rs) SetLevelRsOrDrs(s string) {
	if strings.HasPrefix(s, "solo") {
		i.RsTypeLevel = s
		//i.Lvlkz = s
		return
	}

	lvl, _ := strconv.Atoi(s)
	if lvl == 0 {
		i.RsTypeLevel = s
	} else {
		if lvl >= 7 {
			i.RsTypeLevel = "drs" + s
		} else {
			i.RsTypeLevel = "rs" + s
		}
	}
}

// TypeRedStar rs or drs or solo and level
func (i *Rs) TypeRedStar() (DarkOrRed bool, level string) {
	after, found := strings.CutPrefix(i.RsTypeLevel, "rs")
	if found {
		return false, after
	}
	after, found = strings.CutPrefix(i.RsTypeLevel, "drs")
	if found {
		return true, after
	}
	after, found = strings.CutPrefix(i.RsTypeLevel, "d")
	if found {
		return true, after
	}
	after, found = strings.CutPrefix(i.RsTypeLevel, "solo")
	if found {
		return true, after
	}

	return false, i.RsTypeLevel
}

type CorporationConfigV2 struct {
	Uid         string
	Channels    ChannelsMap
	HelpMessage HelpMessage
}
type HelpMessage map[string]*Info

type ChannelsMap map[string]*Info

type Info struct {
	TypeMessenger     string          `json:"TypeMessenger,omitempty"`
	MessageId         string          `json:"MessageId,omitempty"`
	ChannelId         string          `json:"ChannelId,omitempty"`
	ChannelName       string          `json:"ChannelName,omitempty"`
	GuildId           string          `json:"GuildId,omitempty"`
	GuildName         string          `json:"GuildName,omitempty"`
	GuildAvatarUrl    string          `json:"GuildAvatarUrl,omitempty"`
	Alias             string          `json:"Alias,omitempty"`
	UserAvatarUrl     string          `json:"UserAvatarUrl,omitempty"`
	GameCorporation   string          `json:"GameCorporation,omitempty"`
	GameCorporationId string          `json:"GameCorporationId,omitempty"`
	Language          string          `json:"Language,omitempty"`
	ConfigParamId     int             `json:"ConfigParamId,omitempty"`
	CreatedAt         time.Time       `json:"CreatedAt,omitempty"`
	Options           map[string]bool `json:"Options,omitempty"`
}

const (
	MType        = "TypeMessenger"
	MMId         = "MessageId"
	MChId        = "ChannelId"
	MGuId        = "GuildId"
	MGuName      = "GuildName"
	MGuAvatarUrl = "GuildAvatarUrl"
	MAlias       = "Alias"
	MUsAvatarUrl = "UserAvatarUrl"
)

func (i Info) ToMap() map[string]string {
	m := make(map[string]string)

	if i.TypeMessenger != "" {
		m[MType] = i.TypeMessenger
	}
	if i.MessageId != "" {
		m[MMId] = i.MessageId
	}
	if i.ChannelId != "" {
		m[MChId] = i.ChannelId
	}
	if i.GuildId != "" {
		m[MGuId] = i.GuildId
	}
	if i.GuildName != "" {
		m[MGuName] = i.GuildName
	}
	if i.GuildAvatarUrl != "" {
		m[MGuAvatarUrl] = i.GuildAvatarUrl
	}
	if i.Alias != "" {
		m[MAlias] = i.Alias
	}
	if i.UserAvatarUrl != "" {
		m[MUsAvatarUrl] = i.UserAvatarUrl
	}

	return m
}

type CorpInfo struct {
	ID         int64     `json:"id" db:"id"`
	ConfigName string    `json:"config_name" db:"config_name"`
	CorpName   string    `json:"corp_name" db:"corp_name"`
	CorpID     string    `json:"corp_id" db:"corp_id"`
	Level      int       `json:"level" db:"level"`
	Percent    int       `json:"percent" db:"percent"`
	XP         int       `json:"xp" db:"xp"`
	Webhook    bool      `json:"webhook" db:"webhook"`
	DateEnded  time.Time `json:"date_ended" db:"date_ended"`
	LastUpdate time.Time `json:"last_update" db:"last_update"`
}
