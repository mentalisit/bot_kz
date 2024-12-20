package models

import (
	"fmt"
	"slices"
)

type InMessage struct {
	Mtext         string
	Tip           string
	NameNick      string
	Username      string
	UserId        string
	NameMention   string
	Lvlkz, Timekz string
	Ds            struct {
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

func (u Users) GetAllUserId() []string {
	var all []string
	all = append(all, u.User1.UserId)
	if u.User2 != nil {
		all = append(all, u.User2.UserId)
	}
	if u.User3 != nil {
		all = append(all, u.User3.UserId)
	}
	if u.User4 != nil {
		all = append(all, u.User4.UserId)
	}
	return all
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
