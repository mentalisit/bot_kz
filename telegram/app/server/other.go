package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"telegram/models"
	"time"
)

func (s *Server) CheckAdmin(c *gin.Context) {
	var m apiRs

	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("func CheckAdmin", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a := answer{Time: time.Now()}

	a.ArrBool = s.tg.CheckAdminTg(m.Channel, m.UserName)

	c.JSON(http.StatusOK, a)
}

func (s *Server) GetAvatarUrl(c *gin.Context) {
	var m apiRs

	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("func GetAvatarUrl", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var a string

	a = s.tg.GetAvatarUrl(m.UserName)

	c.JSON(http.StatusOK, a)
}

func (s *Server) telegramSendPoll(c *gin.Context) {
	var m models.Request
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("telegramSendPoll", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, s.tg.SendPoll(m))
}
