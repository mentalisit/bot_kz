package webServer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"rs/config"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const DiscordRedirectURI = "https://mentalisit.tsl.rocks/rs/settings/settings.html"

type linkDs struct {
	UUID     uuid.UUID `json:"UUID,omitempty"`
	Code     string    `json:"Code"`
	Provider string    `json:"Provider"`
	Method   string    `json:"Method"`
}

type DiscordTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type TelegramAuthRequest struct {
	UUID      uuid.UUID `json:"UUID"`
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	AuthDate  int64     `json:"auth_date"`
	Hash      string    `json:"hash"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	PhotoURL  string    `json:"photo_url"`
}

type DiscordUserResponse struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	GlobalName    string `json:"global_name"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Locale        string `json:"locale"`
}

func (s *Server) postSettingDsLink(c *gin.Context) {
	var req struct {
		UUID     uuid.UUID `json:"UUID"`
		Code     string    `json:"Code"`
		Provider string    `json:"Provider"`
		Method   string    `json:"Method"` // oauth или manual_code
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid JSON"})
		return
	}

	var userID, username string

	// Приводим названия провайдеров к короткому виду бота (discord -> ds, telegram -> tg)
	pvdNormalized := req.Provider
	if pvdNormalized == "discord" {
		pvdNormalized = "ds"
	}
	if pvdNormalized == "telegram" {
		pvdNormalized = "tg"
	}

	if req.Method == "manual_code" {
		// --- ПРОВЕРКА В ПАМЯТИ ---
		val, ok := s.LinkCodes.Load(req.Code)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "Код не найден или устарел"})
			return
		}

		data := val.(LinkCodeData)
		if time.Now().After(data.Expires) {
			s.LinkCodes.Delete(req.Code)
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "Срок действия кода истек"})
			return
		}

		if data.Provider != pvdNormalized {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "Код предназначен для другого мессенджера (" + data.Provider + " != " + pvdNormalized + ")"})
			return
		}

		userID = data.UserID
		username = data.Username
		s.LinkCodes.Delete(req.Code) // Удаляем после использования

	} else {
		// --- OAUTH ЛОГИКА (Для Discord) ---
		tokenResp, err := s.exchangeDiscordCode(req.Code)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "Failed to exchange code: " + err.Error()})
			return
		}

		userResp, err := s.getDiscordUser(tokenResp.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to fetch user: " + err.Error()})
			return
		}
		userID = userResp.ID
		username = userResp.Username
	}

	// Сохраняем в БД (Общее для обоих методов)
	fmt.Printf("update %s %s %s %s \n", req.UUID, req.Provider, userID, username)
	err := s.db.UpdateMultiAccountSocial(req.UUID, req.Provider, userID, username)
	if err != nil {
		s.log.ErrorErr(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to save: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "username": username})
}

// AddLinkCode - Эту функцию будет вызывать бот при команде /link
func (s *Server) AddLinkCode(code, userID, username, provider string) {
	s.LinkCodes.Store(code, LinkCodeData{
		UserID:   userID,
		Username: username,
		Provider: provider,
		Expires:  time.Now().Add(10 * time.Minute), // Код живет 10 минут
	})
}

func (s *Server) deleteSettingLink(c *gin.Context) {
	var req struct {
		UUID     uuid.UUID `json:"UUID"`
		Provider string    `json:"Provider"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid JSON"})
		return
	}

	if req.UUID == uuid.Nil || req.Provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "UUID and Provider are required"})
		return
	}

	// Обнуляем данные в БД
	err := s.db.UpdateMultiAccountSocial(req.UUID, req.Provider, "", "")
	if err != nil {
		s.log.ErrorErr(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to unlink account: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (s *Server) exchangeDiscordCode(code string) (*DiscordTokenResponse, error) {
	cfg := config.Instance
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", DiscordRedirectURI)
	data.Set("client_id", cfg.DiscordClientID)
	data.Set("client_secret", cfg.DiscordClientSecret)

	req, err := http.NewRequest("POST", "https://discord.com/api/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("discord API error [%d]: %s", resp.StatusCode, string(body))
	}

	var tokenResp DiscordTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func (s *Server) getDiscordUser(accessToken string) (*DiscordUserResponse, error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("discord API error [%d]: %s", resp.StatusCode, string(body))
	}

	var userResp DiscordUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, err
	}

	return &userResp, nil
}

func (s *Server) postSettingTgLink(c *gin.Context) {
	var req TelegramAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid Data"})
		return
	}

	// 1. Проверяем валидность данных от Telegram
	cfg := config.Instance
	if !s.verifyTelegramAuth(req, cfg.Token.TokenTelegram) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "Invalid Telegram Hash"})
		return
	}

	// 2. Проверяем, не устарел ли запрос (более 24 часов)
	now := time.Now().Unix()
	if now-req.AuthDate > 86400 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "Auth data expired"})
		return
	}

	// 3. Сохраняем в БД (ID и Username)
	err := s.db.UpdateMultiAccountSocial(req.UUID, "telegram", fmt.Sprintf("%d", req.ID), req.Username)
	if err != nil {
		s.log.ErrorErr(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "username": req.Username})
}

func (s *Server) verifyTelegramAuth(req TelegramAuthRequest, botToken string) bool {
	// Telegram Hash Verification Algorithm
	// 1. Collect all data into k=v pairs (only non-empty fields except hash)
	data := map[string]string{
		"id":         fmt.Sprintf("%d", req.ID),
		"auth_date":  fmt.Sprintf("%d", req.AuthDate),
		"username":   req.Username,
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"photo_url":  req.PhotoURL,
	}

	var keys []string
	for k, v := range data {
		if v != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var dataCheckArr []string
	for _, k := range keys {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("%s=%s", k, data[k]))
	}
	dataCheckString := strings.Join(dataCheckArr, "\n")

	// 2. Calculate secret key (SHA256 of bot token)
	sha := sha256.New()
	sha.Write([]byte(botToken))
	secretKey := sha.Sum(nil)

	// 3. HMAC-SHA256
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	return calculatedHash == req.Hash
}
