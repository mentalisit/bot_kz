package models

import (
	"time"

	"github.com/google/uuid"
)

// DiscordTokenResponse содержит ответ OAuth от Discord
type DiscordTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// DiscordUserResponse содержит информацию о пользователе Discord
type DiscordUserResponse struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Email         string `json:"email"`
	GlobalName    string `json:"global_name"`
}

// OAuthConfig конфигурация OAuth
type OAuthConfig struct {
	DiscordClientID     string
	DiscordClientSecret string
	DiscordRedirectURI  string
	DiscordAuthURL      string
	DiscordTokenURL     string
	DiscordUserURL      string
}

// MultiAccount представляет мульти-аккаунт пользователя
type MultiAccount struct {
	UUID             uuid.UUID `json:"uuid"`
	Nickname         string    `json:"nickname"`
	TelegramID       string    `json:"telegram_id"`
	TelegramUsername string    `json:"telegram_username"`
	DiscordID        string    `json:"discord_id"`
	DiscordUsername  string    `json:"discord_username"`
	WhatsappID       string    `json:"whatsapp_id"`
	WhatsappUsername string    `json:"whatsapp_username"`
	CreatedAt        time.Time `json:"created_at"`
	AvatarURL        string    `json:"avatar_url"`
	Alts             []string  `json:"alts"`
}

//type Button struct {
//	Text string
//	Data string
//}

//type Timer struct {
//	//Id       string `bson:"_id"`
//	Dsmesid  string `bson:"dsmesid"`
//	Dschatid string `bson:"dschatid"`
//	Tgmesid  string `bson:"tgmesid"`
//	Tgchatid string `bson:"tgchatid"`
//	Timed    int    `bson:"timed"`
//}

type Timer struct {
	//Id       string `bson:"_id"`
	Tip    string `bson:"tip"`
	ChatId string `bson:"chatId"`
	MesId  string `bson:"mesId"`
	Timed  int    `bson:"timed"`
}
