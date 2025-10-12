package models

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
	Option Option
}

type Option struct {
	Reaction bool
	InClient bool
	Queue    bool
	Pl30     bool
	MinusMin bool
	Edit     bool
	Update   bool
	Elsetrue bool
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
