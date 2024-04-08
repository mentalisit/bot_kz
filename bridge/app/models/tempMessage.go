package models

import "sync"

type BridgeTempMemory struct {
	Timestamp int64
	RelayName string
	MessageDs []MessageIds
	MessageTg []MessageIds
	Wg        sync.WaitGroup
}
type MessageIds struct {
	MessageId string `json:"message_id"`
	ChatId    string `json:"chat_id"`
}
