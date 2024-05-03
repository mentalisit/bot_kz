package server

import (
	"compendium/Compendium/generate"
	"compendium/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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

	identity := generate.CheckCode(code)

	// Проверка на наличие токена в полученной идентификации
	if identity.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Outdated or invalid code"})
		return
	}
	c.JSON(http.StatusOK, identity)

	// Запуск асинхронной операции вставки идентификации в базу данных
	go s.db.Temp.IdentityInsert(context.TODO(), identity)
}

func (s *Server) CheckConnectHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Authorization, content-type")
	token := c.GetHeader("authorization")
	i := s.GetTokenIdentity(token)
	c.JSON(http.StatusOK, i)
}

func (s *Server) CheckCorpDataHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Authorization")
	token := c.GetHeader("authorization")
	roleId := c.Query("roleId")

	i := s.GetTokenIdentity(token)
	if i == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid code"})
		return
	}

	c.JSON(http.StatusOK, s.GetCorpData(i, roleId))

}
func (s *Server) CheckRefreshHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Authorization, content-type")

	token := c.GetHeader("authorization")
	i := s.GetTokenIdentity(token)
	c.JSON(http.StatusOK, i)
}
func (s *Server) CheckSyncTechHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Authorization, content-type")

	mode := c.Param("mode")
	twin := c.DefaultQuery("twin", "")

	token := c.GetHeader("authorization")

	i := s.GetTokenIdentity(token)
	if i.User.Username == "" {
		s.log.Error("no Name")
	}

	if i == nil {
		fmt.Println("i==nil")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}
	userId := i.User.ID
	userName := i.User.Username
	guildId := i.Guild.ID
	if twin != "" && twin != "default" {
		userName = twin
	}

	fmt.Printf("mode %s TwinOrName %s\n", mode, userName)

	if mode == "get" {
		sd := models.SyncData{
			TechLevels: models.TechLevels{},
			Ver:        1,
			InSync:     1,
		}
		cm := s.db.Temp.CorpMemberReadByUserIdByName(context.Background(), userId, guildId, userName)
		if len(cm.Tech) > 0 {
			sd.TechLevels = cm.Tech
		}

		c.JSON(http.StatusOK, sd)
	} else if mode == "sync" {

		var data models.SyncData
		if err := c.BindJSON(&data); err != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		go s.db.Temp.CorpMemberTechUpdate(context.TODO(), userId, guildId, userName, data.TechLevels)

		// Используйте переменную data с полученными данными
		c.JSON(http.StatusOK, data)
	}
}

//requestedMethod := c.GetHeader("Access-Control-Request-Method")
//requestedHeaders := c.GetHeader("Access-Control-Request-Headers")
//fmt.Println("Requested method:", requestedMethod)
//fmt.Println("Requested headers:", requestedHeaders)

func (s *Server) getWsMatches(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Authorization, content-type")
	limit := c.DefaultQuery("limit", "")
	page := c.DefaultQuery("page", "1")
	filter := c.DefaultQuery("filter", "")

	result := s.getMatchesAll(limit, page, filter)

	c.JSON(http.StatusOK, result)
}

//	func (s *Server) getWsMatchesAll(c *gin.Context) {
//		c.Header("Access-Control-Allow-Origin", "*")
//		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
//		c.Header("Access-Control-Allow-Headers", "Authorization, content-type")
//		limit := c.DefaultQuery("limit", "")
//		page := c.DefaultQuery("page", "1")
//		filter := c.DefaultQuery("filter", "")
//
//		result := s.getMatchesAll(limit, page, filter)
//
//		c.JSON(http.StatusOK, result)
//	}
func (s *Server) getWsCorps(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Authorization, content-type")
	limit := c.DefaultQuery("limit", "")
	page := c.DefaultQuery("page", "1")

	result := s.getCorps(limit, page)

	c.JSON(http.StatusOK, result)
}
func (s *Server) docs(c *gin.Context) {
	htmlContent := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Docs</title>
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
        <h1>Select endpoint:</h1>
        <ul>
            <li><a href="/ws/corps">list corporations</a></li>
            <li><a href="/ws/corps?limit=20">list corporations limit 20</a></li>
            <li><a href="/ws/matches">list matches</a></li>
            <li><a href="/ws/matches?limit=20">list matches limit 20</a></li>
        </ul>
    </div>
</body>
</html>
`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
}
