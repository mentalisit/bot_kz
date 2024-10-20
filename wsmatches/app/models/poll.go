package models

type PollStruct struct {
	Author      string            `json:"author"`
	Question    string            `json:"question"`
	Options     []string          `json:"options"`
	CreateTime  int64             `json:"createTime"`
	UrlPoll     string            `json:"urlPoll"`
	Config      BridgeConfig      `json:"config"`
	Votes       []Votes           `json:"votes"`
	PollMessage map[string]string `json:"pollMessage"`
}
type Votes struct {
	Type     string `json:"type"`
	Channel  string `json:"channel"`
	UserName string `json:"userName"`
	Answer   string `json:"answer"`
}
type BridgeConfig struct {
	Id                int              `json:"id"`
	NameRelay         string           `json:"nameRelay"`
	HostRelay         string           `json:"hostRelay"`
	Role              []string         `json:"role"`
	ChannelDs         []BridgeConfigDs `json:"channelDs"`
	ChannelTg         []BridgeConfigTg `json:"channelTg"`
	ForbiddenPrefixes []string         `json:"forbiddenPrefixes"`
}
type BridgeConfigDs struct {
	ChannelId       string            `json:"channel_id"`
	GuildId         string            `json:"guild_id"`
	CorpChannelName string            `json:"corp_channel_name"`
	AliasName       string            `json:"alias_name"`
	MappingRoles    map[string]string `json:"mapping_roles"`
}
type BridgeConfigTg struct {
	ChannelId       string            `json:"channel_id"`
	CorpChannelName string            `json:"corp_channel_name"`
	AliasName       string            `json:"alias_name"`
	MappingRoles    map[string]string `json:"mapping_roles"`
}
