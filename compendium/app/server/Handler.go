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
	fmt.Println("Last parameter:", mode)

	token := c.GetHeader("authorization")

	i := s.GetTokenIdentity(token)
	if i == nil {
		fmt.Println("i==nil")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}
	if mode == "get" {
		sd := models.SyncData{
			TechLevels: models.TechLevels{},
			Ver:        1,
			InSync:     1,
		}
		cm := s.db.Temp.CorpMemberReadByUserId(context.TODO(), i.User.ID)
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
		go s.db.Temp.CorpMemberTechUpdate(context.TODO(), i.User.ID, data.TechLevels)

		// Используйте переменную data с полученными данными
		c.JSON(http.StatusOK, data)
	}
}

//requestedMethod := c.GetHeader("Access-Control-Request-Method")
//requestedHeaders := c.GetHeader("Access-Control-Request-Headers")
//fmt.Println("Requested method:", requestedMethod)
//fmt.Println("Requested headers:", requestedHeaders)
