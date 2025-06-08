package server

import (
	"compendium_s/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func (s *Server) CheckIdentityHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Authorization")
	code := c.GetHeader("authorization")

	// Проверка наличия кода в запросе и его длины
	if code == "" || len(code) != 14 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
		return
	}

	identity := s.CheckCode(code)

	if identity.Token != "" {
		c.JSON(http.StatusOK, identity)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
}

func (s *Server) CheckConnectHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Authorization, content-type")
	token := c.GetHeader("authorization")

	// Локальная проверка
	i := s.GetTokenIdentity(token)
	if i != nil {
		if i.User.GameName != "" {
			i.User.Username = i.User.GameName
		}
		c.JSON(http.StatusOK, i)
		return
	}
	c.JSON(http.StatusForbidden, nil)
}

func (s *Server) CheckCorpDataHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Authorization")

	token := c.GetHeader("authorization")
	roleId := c.Query("roleId")
	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "missing token"})
		return
	}

	cacheKey := token + ":" + roleId

	// Проверяем кэш
	s.cacheMutex.Lock()
	entry, exists := s.cacheReq[cacheKey]
	if exists && time.Since(entry.timestamp) <= 2*time.Second {
		// Есть свежий кэш, возвращаем
		s.cacheMutex.Unlock()
		c.JSON(http.StatusOK, entry.data)
		return
	}
	s.cacheMutex.Unlock()

	// Кэш отсутствует или устарел — получаем заново
	i := s.GetTokenIdentity(token)
	if i != nil {
		result := s.GetCorpData(i, roleId)

		// Сохраняем в кэш
		s.cacheMutex.Lock()
		s.cacheReq[cacheKey] = cacheEntry{
			data:      result,
			timestamp: time.Now(),
		}
		s.cacheMutex.Unlock()

		c.JSON(http.StatusOK, result)
		return
	}

	c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
}

func (s *Server) CheckRefreshHandler(c *gin.Context) {
	s.PrintGoroutine()
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Authorization, content-type")

	token := c.GetHeader("authorization")
	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "missing token"})
		return
	}

	token = s.refreshToken(token)

	i := s.GetTokenIdentity(token)
	if i != nil {
		if i.User.GameName != "" {
			i.User.Username = i.User.GameName
		}
		c.JSON(http.StatusOK, i)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
}

func (s *Server) CheckSyncTechHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Authorization, content-type")

	mode := c.Param("mode")
	twin := c.DefaultQuery("twin", "")
	token := c.GetHeader("authorization")

	i := s.GetTokenIdentity(token)

	if i == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}
	if i.MultiAccount != nil && i.MultiGuild != nil {
		s.SyncTechMulti(c, i, mode, twin)
		return
	}
	userId := i.User.ID
	userName := i.User.Username
	guildId := i.MultiGuild.GId.String()
	if twin != "" && twin != "default" {
		userName = twin
	}
	if userName == i.User.GameName {
		userName = i.User.Username
	}

	fmt.Printf("mode %s corporation %s Name %s\n", mode, i.Guild.Name, userName)

	if mode == "get" {
		sd := models.SyncData{
			TechLevels: models.TechLevels{},
			Ver:        1,
			InSync:     1,
		}
		techBytes, err := s.db.TechGet(userName, userId, guildId)
		if err == nil && len(techBytes) > 0 {
			sd.TechLevels = sd.TechLevels.ConvertToTech(techBytes)
		}
		c.JSON(http.StatusOK, sd)
	} else if mode == "sync" {

		var data models.SyncData
		if err := c.BindJSON(&data); err != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		bytes, err := json.Marshal(data.TechLevels)
		if err != nil {
			s.log.ErrorErr(err)
		}
		err = s.db.TechUpdate(userName, userId, i.MultiGuild.GId.String(), bytes)
		if err != nil {
			s.log.ErrorErr(err)
		}

		// Используйте переменную data с полученными данными
		c.JSON(http.StatusOK, data)
	}
}

func (s *Server) links(c *gin.Context) {
	htmlContent := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Links</title>
    <style>
        html, body {
            height: 100%;
            margin: 0;
            padding: 0;
            display: flex;
            justify-content: center;
            align-items: center;
            font-family: Arial, sans-serif;
            background-color: #f4f4f4; // Светлый фон всей страницы
        }
        .centered-content {
            width: 100%; // Ширина контента равна ширине страницы
            max-width: 600px; // Максимальная ширина контента
            text-align: center;
            box-shadow: 0 4px 8px rgba(0,0,0,0.1);
            background-color: white;
            padding: 20px;
            border-radius: 8px;
        }
        ul {
            list-style-type: none;
            padding: 0;
        }
        li {
            margin: 10px 0;
        }
        a {
            text-decoration: none;
            color: #3366cc;
        }
        a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="centered-content">
        <h1>Select link:</h1>
        <ul>
            <li><a href="https://discord.com/oauth2/authorize?client_id=909526127305953290&scope=bot+applications.commands&permissions=141533113424" target="_blank">Invite Discord Bot</a></li>
            <li><a href="https://t.me/gote1st_bot" target="_blank">Invite Telegram Bot</a></li>
            <li><a href="https://discord.com/users/582882137842122773" target="_blank">Message Discord Bot Author</a></li>
            <li><a href="https://t.me/mentalisit" target="_blank">Message Telegram Bot Author</a></li>
        </ul>
    </div>
</body>
</html>
`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
}

func (s *Server) api(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "empty token"})
		return
	}
	i := s.GetTokenIdentity(token)
	if i == nil || i.Token == "" || i.Guild.ID != "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
		return
	}
	userid := c.Query("userid")
	if userid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid must not be empty"})
		return
	}

	read, _ := s.db.CorpMembersRead(i.Guild.ID)
	if len(read) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "guildid empty members"})
		return
	}
	var memb []models.CorpMember
	for _, member := range read {
		if strings.Contains(member.UserId, userid) {
			memb = append(memb, member)
		}
	}
	if len(memb) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "member not found"})
		return
	}
	c.JSON(http.StatusOK, memb)
}
