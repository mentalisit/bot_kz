package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

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
	MCreAt       = "CreatedAt"
	MOptions     = "Options"
	MAutoHelp    = "AutoHelp"
	MCleanChat   = "CleanChat"
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
	MAcc        *MultiAccount
}

func (i *InMessageV2) GetNameMention() string {
	if i.NameMention == "@" {
		return fmt.Sprintf("[%s](tg://user?id=%s)", i.Username, i.UserId)
	}
	return i.NameMention
}

type Rs struct {
	RsTypeLevel            string
	TimeRs                 string
	AltName                string
	Money                  bool
	CountQueue             int
	NumberName             int
	NumberLevel            int
	Ch                     string
	Info                   *Info
	LevelRsPingDs          string
	NameRsLevel            string
	QueueMessages          map[string]QueueMessages
	TextQueueCompliteBonus string
	//MessagesMutex sync.Mutex
	U              []QueueActive
	MessageIdsChan chan map[string]QueueMessages
}

func (i *Rs) GetTimeRs() int64 {
	timeRs, err := strconv.ParseInt(i.TimeRs, 10, 64)
	if err != nil {
		return 0
	}

	return timeRs
}

func NewRs() *Rs {
	return &Rs{
		MessageIdsChan: make(chan map[string]QueueMessages, 50),
	}
}

func (i *Rs) SetLevelRsOrDrs(in *InMessageV2, s string) {
	lvl, _ := strconv.Atoi(s)
	if lvl == 0 {
		i.RsTypeLevel = s
		return
	}

	if in.Config.Channels[in.Messenger.ChannelId].Corp != nil {
		corp := in.Config.Channels[in.Messenger.ChannelId].Corp
		if corp != nil {
			if corp.DefaultNameDRS != "" && lvl >= 7 {
				i.RsTypeLevel = corp.DefaultNameDRS + s
				return
			} else if corp.DefaultNameRS != "" && lvl <= 6 {
				i.RsTypeLevel = corp.DefaultNameRS + s
				return
			}
		}
	}

	if lvl >= 7 {
		i.RsTypeLevel = "drs" + s
	} else {
		i.RsTypeLevel = "rs" + s
	}
}

func (i *Rs) getCorpTypeStar() (DarkOrRed bool, level string) {
	if i.Info != nil && i.Info.Corp != nil {
		c := i.Info.Corp

		if c.DefaultNameRS != "" {
			after, found := strings.CutPrefix(i.RsTypeLevel, c.DefaultNameRS)
			if found {
				return false, after
			}
		}
		if c.DefaultNameDRS != "" {
			after, found := strings.CutPrefix(i.RsTypeLevel, c.DefaultNameDRS)
			if found {
				return true, after
			}
		}
	}

	after, found := strings.CutPrefix(i.RsTypeLevel, "drs")
	if found {
		return true, after
	}
	after, found = strings.CutPrefix(i.RsTypeLevel, "rs")
	if found {
		return false, after
	}
	return false, ""
}

// TypeRedStar rs or drs or solo and level
func (i *Rs) TypeRedStar() (DarkOrRed bool, level string) {
	if i.Info != nil && i.Info.Corp != nil {
		return i.getCorpTypeStar()
	}
	after, found := strings.CutPrefix(i.RsTypeLevel, "rs")
	if found {
		return false, after
	}
	after, found = strings.CutPrefix(i.RsTypeLevel, "drs")
	if found {
		return true, after
	}

	return false, i.RsTypeLevel
}

func (i *Rs) GetTitle(getText func(lang string, key string) string) string {
	return fmt.Sprintf("%s %s", getText(i.Info.Language, "queue"), i.GetRsNameLevel(getText))
}

func (i *Rs) GetRsNameLevel(getText func(lang string, key string) string) string {
	darkOrRed, level := i.TypeRedStar()
	if darkOrRed {
		if i.Info.Corp != nil && i.Info.Corp.DefaultNameDRS != "" {
			return i.Info.Corp.DefaultNameDRS + level
		}
		return getText(i.Info.Language, "drs") + level
	}
	if i.Info.Corp != nil && i.Info.Corp.DefaultNameRS != "" {
		return i.Info.Corp.DefaultNameRS + level
	}
	return getText(i.Info.Language, "rs") + level
}

// GetTypeRs true=Dark , false=Red
func (i *Rs) GetTypeRs() bool {
	darkOrRed, _ := i.getCorpTypeStar()
	return darkOrRed
}

// GetLevelRs return level star
func (i *Rs) GetLevelRs() string {
	_, level := i.getCorpTypeStar()
	return level
}

type CorporationConfigV2 struct {
	Uid         string
	Channels    ChannelsMap
	Bonuses     []GameSettings
	HelpMessage HelpMessage
}

func (c *CorporationConfigV2) GetGameLevelInCorpInfo(f func(c string) (*CorpInfo, error)) {
	for s := range c.Channels {
		if c.Channels[s].Game != nil && c.Channels[s].Game.GameCorporationId != "" {
			i, err := f(c.Channels[s].Game.GameCorporationId)
			if err == nil && i != nil {
				if c.Channels[s].Game == nil {
					c.Channels[s].Game = &GameSettings{}
				}
				c.Channels[s].Game.GameXP = i.XP
				c.Channels[s].Game.GameLevel = i.GetLevelByXP()
			}

		}
	}
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
	if !i.CreatedAt.IsZero() {
		m[MCreAt] = i.CreatedAt.Format(time.RFC3339)
	}

	return m
}

type Other struct {
	Uuid     string
	DataType string
	Data     Info
	Read     bool
}

func (q *Other) JsonMarshalDataWeb() []byte {
	messagesJSON, _ := json.Marshal(q.Data)
	return messagesJSON
}
