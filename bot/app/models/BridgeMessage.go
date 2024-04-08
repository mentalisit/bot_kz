package models

type BridgeMessage struct {
	Text          string
	Sender        string
	Tip           string
	ChatId        string
	MesId         string
	GuildId       string
	TimestampUnix int64
	FileUrl       []string
	Extra         map[string][]interface{}
	Avatar        string
	Reply         *BridgeMessageReply
	Config        *BridgeConfig
}
