package models

type Bridge2Config struct {
	Id                int                         `json:"id"`
	NameRelay         string                      `json:"nameRelay"`
	HostRelay         string                      `json:"hostRelay"`
	Role              []string                    `json:"role"`
	Channel           map[string][]Bridge2Configs `json:"channel"`
	ForbiddenPrefixes []string                    `json:"forbiddenPrefixes"`
}
type Bridge2Configs struct {
	ChannelId       string            `json:"channel_id"`
	GuildId         string            `json:"guild_id"`
	CorpChannelName string            `json:"corp_channel_name"`
	AliasName       string            `json:"alias_name"`
	MappingRoles    map[string]string `json:"mapping_roles"`
}
