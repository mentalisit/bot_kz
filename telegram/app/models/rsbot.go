package models

type InMessage struct {
	Mtext         string
	Tip           string
	Name          string
	NameMention   string
	Lvlkz, Timekz string
	Ds            struct {
		Mesid   string
		Nameid  string
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
