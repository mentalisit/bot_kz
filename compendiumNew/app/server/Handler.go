package server

import (
	"compendium/models"
	"encoding/json"
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

	identity := s.CheckCode(code)

	// Проверка на наличие токена в полученной идентификации
	if identity.Token == "" {
		fmt.Println(code, identity)
		c.JSON(http.StatusForbidden, gin.H{"error": "Outdated or invalid code"})
		return
	}
	c.JSON(http.StatusOK, identity)
}

func (s *Server) CheckConnectHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Authorization, content-type")
	token := c.GetHeader("authorization")
	i := s.GetTokenIdentity(token)
	if i == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
	}
	if i != nil && i.User.GameName != "" {
		i.User.Username = i.User.GameName
	}
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
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid code"})
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
	if i == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid code"})
		return
	}
	if i.User.GameName != "" {
		i.User.Username = i.User.GameName
	}
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

	if i == nil || i.User.Username == "" || i.Guild.Name == "" {
		fmt.Println("i==nil")
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
		return
	}
	userId := i.User.ID
	userName := i.User.Username
	guildId := i.Guild.ID
	if twin != "" && twin != "default" {
		userName = twin
	}
	if userName == i.User.GameName {
		userName = i.User.Username
	}

	fmt.Printf("mode %s TwinOrName %s\n", mode, userName)

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
		err = s.db.TechUpdate(userName, userId, guildId, bytes)
		if err != nil {
			s.log.ErrorErr(err)
			return
		}

		// Используйте переменную data с полученными данными
		c.JSON(http.StatusOK, data)
	}
}

//requestedMethod := c.GetHeader("Access-Control-Request-Method")
//requestedHeaders := c.GetHeader("Access-Control-Request-Headers")
//fmt.Println("Requested method:", requestedMethod)
//fmt.Println("Requested headers:", requestedHeaders)
