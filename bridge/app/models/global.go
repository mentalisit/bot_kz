package models

type SendText struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
}
type DeleteMessageStruct struct {
	MessageId string `json:"message_id"`
	Channel   string `json:"channel"`
}
type SendTextDeleteSeconds struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
	Seconds int    `json:"seconds"`
}
