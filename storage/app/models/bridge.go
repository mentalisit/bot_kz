package models

type BridgeConfig struct {
	Id                int              `json:"id"`
	NameRelay         string           `json:"nameRelay"`
	HostRelay         string           `json:"hostRelay"`
	Role              []string         `json:"role"`
	ChannelDs         []BridgeConfigDs `json:"channelDs"`
	ChannelTg         []BridgeConfigTg `json:"channelTg"`
	ForbiddenPrefixes []string         `json:"forbiddenPrefixes"`
	Prefix            string           `json:"prefix"`
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
