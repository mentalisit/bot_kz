package models

import "github.com/google/uuid"

type BridgeConfigV2 struct {
	Uid      uuid.UUID                     `json:"uid"`
	LocalGid string                        `json:"local_gid"`
	Channel  map[string][]BridgeChannelsV2 `json:"channel"`
	Conf     BridgeConfV2                  `json:"conf,omitempty"`
}
type BridgeChannelsV2 struct {
	ChannelId       string            `json:"channel_id"`
	GuildId         string            `json:"guild_id"`
	CorpChannelName string            `json:"corp_channel_name,omitempty"`
	AliasName       string            `json:"alias_name,omitempty"`
	MappingRoles    map[string]string `json:"mapping_roles,omitempty"`
}
type BridgeConfV2 struct {
}
