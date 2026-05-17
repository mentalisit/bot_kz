package models

import (
	"fmt"
	"slices"
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

type Options []string

func (o *Options) Contains(s string) bool {
	return slices.Contains(*o, s)
}
func (o *Options) Remove(s string) {
	for i, item := range *o {
		if item == s {
			*o = append((*o)[:i], (*o)[i+1:]...)
			return
		}
	}
}
func (o *Options) Add(s string) {
	if o.Contains(s) {
		panic(fmt.Sprintf("option %s already add to opt %+v\n", s, *o))
	}
	*o = append(*o, s)
}

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
	if i.ChannelName != "" {
		m[MChName] = i.ChannelName
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
	if i.UserAvatarUrl != "" {
		m[MUsAvatarUrl] = i.UserAvatarUrl
	}
	if i.Language != "" {
		m[MLang] = i.Language
	}
	m[MCreAt] = i.CreatedAt.Format(time.RFC3339)

	return m
}

const (
	MType        = "TypeMessenger"
	MMId         = "MessageId"
	MChId        = "ChannelId"
	MChName      = "ChannelName"
	MGuId        = "GuildId"
	MGuName      = "GuildName"
	MGuAvatarUrl = "GuildAvatarUrl"
	MAlias       = "Alias"
	MUsAvatarUrl = "UserAvatarUrl"
	MGC          = "GameCorporation"
	MGCId        = "GameCorporationId"
	MLang        = "Language"
	MConfPar     = "ConfigParamId"
	MCreAt       = "CreateAt"
)

const (
	OptionReaction = "Reaction"
	OptionInClient = "InClient"
)

type CorporationConfigV2 struct {
	Uid         string
	Channels    ChannelsMap
	Bonuses     []GameSettings
	HelpMessage HelpMessage
}
type HelpMessage map[string]*Info

type ChannelsMap map[string]*Info

type Info struct {
	TypeMessenger  string        `json:"TypeMessenger,omitempty"`
	MessageId      string        `json:"MessageId,omitempty"`
	ChannelId      string        `json:"ChannelId,omitempty"`
	ChannelName    string        `json:"ChannelName,omitempty"`
	GuildId        string        `json:"GuildId,omitempty"`
	GuildName      string        `json:"GuildName,omitempty"`
	GuildAvatarUrl string        `json:"GuildAvatarUrl,omitempty"`
	UserAvatarUrl  string        `json:"UserAvatarUrl,omitempty"`
	Language       string        `json:"Language,omitempty"`
	CreatedAt      time.Time     `json:"CreatedAt,omitempty"`
	Game           *GameSettings `json:"Game,omitempty"`
	Corp           *CorpSettings `json:"Corp,omitempty"`
}
type GameSettings struct {
	GameCorporation   string `json:"GameCorporation,omitempty"`
	GameCorporationId string `json:"GameCorporationId,omitempty"`
	Alias             string `json:"Alias,omitempty"`
	GameLevel         int    `json:"GameLevel,omitempty"`
	GameXP            int    `json:"GameXP,omitempty"`
}
type CorpSettings struct {
	AutoHelp            bool     `json:"AutoHelp,omitempty"`
	DeleteMessages      bool     `json:"DeleteMessages,omitempty"`
	DeleteMessagesDelay int      `json:"DeleteMessagesDelay,omitempty"`
	CustomText          bool     `json:"CustomText,omitempty"`
	HelpText            string   `json:"HelpText,omitempty"`
	DefaultNameRS       string   `json:"DefaultNameRS,omitempty"`
	DefaultNameDRS      string   `json:"DefaultNameDRS,omitempty"`
	MessengerInvites    []string `json:"MessengerInvites,omitempty"`
}
