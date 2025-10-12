package models

type BridgeTempMemory struct {
	Timestamp int64
	RelayName string
	MessageDs []MessageIds
	MessageTg []MessageIds
	MessageWa []MessageIds
	Message   map[string]string
}
type MessageIds struct {
	MessageId string `json:"message_id"`
	ChatId    string `json:"chat_id"`
}
