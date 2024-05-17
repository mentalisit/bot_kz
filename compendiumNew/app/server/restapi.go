package server

import (
	"bytes"
	"compendium/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) RunServerRestApi() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.POST("/compendium/inbox", s.InboxMessage)

	err := router.Run(":80")
	if err != nil {
		s.log.ErrorErr(err)
	}
}

func (s *Server) InboxMessage(c *gin.Context) {
	var data models.IncomingMessage
	if err := c.BindJSON(&data); err != nil {
		s.log.ErrorErr(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	s.In <- data
	c.JSON(http.StatusOK, "ok")
}

const apiname = "kz_bot"

func GetRoles(guildId string) ([]models.CorpRole, error) {
	data, err := json.Marshal(guildId)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post("http://"+apiname+"/discord/GetRoles", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		//resp, err = http.Post("http://192.168.100.155:802/data", "application/json", bytes.NewBuffer(data))
		return nil, err
	}

	var guildRole []models.CorpRole

	err = json.NewDecoder(resp.Body).Decode(&guildRole)
	if err != nil {
		return nil, err
	}

	return guildRole, nil
}
func CheckRoleDs(guildId, memderId, roleid string) bool {
	if roleid == "" {
		return true
	}
	m := models.CheckRoleStruct{
		GuildId:  guildId,
		MemberId: memderId,
		RoleId:   roleid,
	}
	data, err := json.Marshal(m)
	if err != nil {
		return false
	}

	resp, err := http.Post("http://"+apiname+"/discord/CheckRole", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		//resp, err = http.Post("http://192.168.100.155:802/data", "application/json", bytes.NewBuffer(data))
		return false
	}

	var ok bool

	err = json.NewDecoder(resp.Body).Decode(&ok)
	if err != nil {
		return false
	}
	return ok
}
