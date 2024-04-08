package models

type ToBridgeMessage struct {
	Text          string              `json:"text"`
	Sender        string              `json:"sender"`
	Tip           string              `json:"tip"`
	ChatId        string              `json:"chatId"`
	MesId         string              `json:"mesId"`
	GuildId       string              `json:"guildId"`
	TimestampUnix int64               `json:"timestampUnix"`
	Extra         []FileInfo          `json:"extra"`
	Avatar        string              `json:"avatar"`
	Reply         *BridgeMessageReply `json:"reply"`
	Config        *BridgeConfig       `json:"config"`
}
type FileInfo struct {
	Name   string `json:"name"`
	Data   []byte `json:"data"`
	URL    string `json:"URL"`
	Size   int64  `json:"size"`
	FileID string `json:"fileID"`
}
type BridgeMessageReply struct {
	TimeMessage int64  `json:"time_message"`
	Text        string `json:"text"`
	Avatar      string `json:"avatar"`
	UserName    string `json:"userName"`
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

type BridgeSendToMessenger struct {
	Text      string              `json:"text"`
	Sender    string              `json:"sender"`
	ChannelId []string            `json:"channelId"`
	Avatar    string              `json:"avatar"`
	Extra     []FileInfo          `json:"extra"`
	Reply     *BridgeMessageReply `json:"reply"`
}
type MessageIds struct {
	MessageId string
	ChatId    string
}
