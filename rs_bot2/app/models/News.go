package models

type News struct {
	Uid           string `json:"uid"`
	ChatId        string `json:"chat_id"`
	Language      string `json:"language"`
	MessengerType string `json:"messenger_type"`
}
