package models

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type Identify struct {
	Id           int64
	MID          int
	SolarId      int64
	Author       string
	Count        int
	Participants string
	Points       int
	StartTime    string
}
type BattlesTop struct {
	Id       int
	CorpName string
	Name     string
	Level    int
	Count    int
}
type ScoreboardParams struct {
	Name              string
	ChannelWebhook    string
	ChannelScoreboard string
}
type InMessage struct {
	Mtext       string
	Tip         string
	NameNick    string
	Username    string
	UserId      string
	NameMention string
	RsTypeLevel string
	//Lvlkz       string
	//Timekz      string
	TimeRs int
	Ds     struct {
		Mesid   string
		Guildid string
		Avatar  string
	}
	Tg struct {
		Mesid int
	}
	Config CorporationConfig
	//Option //Option
	Opt Options
}

// TypeRedStar rs or drs or solo and level
func (i *InMessage) TypeRedStar() (DarkOrRed bool, level string) {
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
func (i *InMessage) IsDRS() bool {
	dark, _ := i.TypeRedStar()
	return dark
}
func (i *InMessage) SetLevelRsOrDrs(s string) {
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

	//i.Lvlkz = s
}
func (i *InMessage) SetTimeRs(s string) {
	timeRs, _ := strconv.Atoi(s)
	if timeRs == 0 {
		timeRs = 30 //default time
	}

	if timeRs > 180 {
		timeRs = 180
	}
	i.TimeRs = timeRs
	//i.Timekz = strconv.Itoa(timeRs)
}
func (i *InMessage) IfDiscord() bool {
	return i.Config.DsChannel != ""
}
func (i *InMessage) IfTelegram() bool {
	return i.Config.TgChannel != ""
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

type CorporationConfig struct {
	Type           int
	CorpName       string
	DsChannel      string
	TgChannel      string
	WaChannel      string
	Country        string
	DelMesComplite int
	MesidDsHelp    string
	MesidTgHelp    string
	Forward        bool
	Guildid        string
}

//type Option struct {
//	Reaction bool
//	InClient bool
//	Queue    bool
//	Pl30     bool
//	MinusMin bool
//	Edit     bool
//	Update   bool
//	Elsetrue bool
//}

const (
	OptionReaction        = "Reaction"
	OptionInClient        = "InClient"
	OptionQueue           = "Queue"
	OptionPl30            = "Pl30"
	OptionMinusMin        = "MinusMin"
	OptionMinusMinNext    = "MinusMinNext"
	OptionEdit            = "Edit"
	OptionUpdate          = "Update"
	OptionUpdateAutoHelp  = "UpdateAutoHelp"
	OptionMessageUpdateDS = "MessageUpdateDS"
	OptionMessageUpdateTG = "MessageUpdateTG"
	OptionElseTrue        = "ElseTrue"
	OptionQueueAll        = "QueueAll"
	OptionPlus            = "Plus"
)

type Users struct {
	User1 Sborkz
	User2 *Sborkz
	User3 *Sborkz
	User4 *Sborkz
}

func (u Users) GetAllUserId() (all []string, tg []string) {
	if u.User1.Tip == "tg" {
		tg = append(tg, u.User1.UserId)
	}
	all = append(all, u.User1.UserId)
	if u.User2 != nil {
		all = append(all, u.User2.UserId)
		if u.User2.Tip == "tg" {
			tg = append(tg, u.User2.UserId)
		}
	}
	if u.User3 != nil {
		all = append(all, u.User3.UserId)
		if u.User3.Tip == "tg" {
			tg = append(tg, u.User3.UserId)
		}
	}
	if u.User4 != nil {
		all = append(all, u.User4.UserId)
		if u.User4.Tip == "tg" {
			tg = append(tg, u.User4.UserId)
		}
	}
	return all, tg
}

type Sborkz struct {
	Id          int
	Corpname    string
	Name        string
	UserId      string
	Mention     string
	Tip         string
	Dsmesid     string
	Tgmesid     int
	Wamesid     string
	Time        string
	Date        string
	Lvlkz       string
	Numkzn      int
	Numberkz    int
	Numberevent int
	Eventpoints int
	Active      int
	Timedown    int
}

type Names struct {
	Name1 string
	Name2 string
	Name3 string
	Name4 string
}

type EmodjiUser struct {
	Id                                int
	Tip, Name, Em1, Em2, Em3, Em4     string
	Module1, Module2, Module3, Weapon string
}

type Timer struct {
	//Id       string `bson:"_id"`
	Dsmesid  string `bson:"dsmesid"`
	Dschatid string `bson:"dschatid"`
	Tgmesid  string `bson:"tgmesid"`
	Tgchatid string `bson:"tgchatid"`
	Timed    int    `bson:"timed"`
}

type Top struct {
	Name              string
	Numkz, Id, Points int
}

type QueueStruct struct {
	CorpName string
	Level    string
	Count    int
}

type RsEvent struct {
	Id          int
	CorpName    string
	NumEvent    int
	ActiveEvent int
	Number      int
}

type Events struct {
	ID     int64
	Number int
	Event  int
	Status bool
}

type CorporationHistory struct {
	CorpName  string
	ChannelDs string
}

type EntryScoreboard struct {
	DisplayName string
	RsLevel     int
	StarsCount  int
	Score       int
}
