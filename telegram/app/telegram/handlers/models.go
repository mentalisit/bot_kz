package handlers

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
