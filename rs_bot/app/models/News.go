package models

type News struct {
	Id            int    `json:"id"`
	ChatId        string `json:"chat_id"`
	Language      string `json:"language"`
	MessengerType string `json:"messenger_type"`
}
