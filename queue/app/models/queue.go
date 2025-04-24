package models

import "encoding/json"

type QueueStruct struct {
	CorpName string
	Level    string
	Count    int
}
type Tumcha struct {
	Name     string
	NameId   int64
	Level    int
	Vid      string
	Chatid   int
	Timedown int
}

// Webhook представляет структуру записи из таблицы webhooks
type Webhook struct {
	ID      int64           `json:"id"`
	TsUnix  int64           `json:"ts_unix"`
	Corp    string          `json:"corp"`
	Message json.RawMessage `json:"message"` // для хранения JSONB
}
