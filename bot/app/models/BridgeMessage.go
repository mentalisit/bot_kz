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

type Request struct {
	Data    map[string]string `json:"data"`
	Options []string          `json:"options"`
}
