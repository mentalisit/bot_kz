package original

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
)

var urlOriginalServer = "https://bot.hs-compendium.com/compendium"

func IdentityHandler(c *gin.Context, code string) {
	// Проксируем запрос на оригинальный сервер
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlOriginalServer+"/applink/identities?ver=2&code=1", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Передаём заголовок Authorization
	req.Header.Set("Authorization", code)

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to original server"})
		return
	}
	defer resp.Body.Close()

	// Передаём ответ обратно клиенту
	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
	return
}

func ConnectHandler(c *gin.Context, token, GuildID string) {
	// Проксирование на оригинальный сервер
	client := &http.Client{}
	body, _ := json.Marshal(map[string]string{"guild_id": GuildID})

	req, err := http.NewRequest("POST", urlOriginalServer+"/applink/connect", bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to original server"})
		return
	}
	defer resp.Body.Close()

	// Читаем ответ от оригинального сервера
	respBody, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", respBody)
}

func CorpDataHandler(c *gin.Context, token, roleId string) {
	// Если не найдено - проксируем на оригинальный сервер
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf(urlOriginalServer+"/cmd/corpdata?roleId=%s", url.QueryEscape(roleId)), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	// Передаём токен
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to original server"})
		return
	}
	defer resp.Body.Close()

	// Читаем и возвращаем ответ от оригинального сервера
	respBody, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", respBody)
}

func RefreshHandler(c *gin.Context, token string) {
	// Токен не найден → отправляем запрос на оригинальный сервер
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlOriginalServer+"/applink/refresh", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to original server"})
		return
	}
	defer resp.Body.Close()

	// Читаем и возвращаем ответ от оригинального сервера
	respBody, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", respBody)
}

func SyncTechHandler(c *gin.Context, token, mode string) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", urlOriginalServer+"/cmd/syncTech/"+mode, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to original server"})
		return
	}
	defer resp.Body.Close()

	// Читаем и возвращаем ответ от оригинального сервера
	respBody, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", respBody)
}
