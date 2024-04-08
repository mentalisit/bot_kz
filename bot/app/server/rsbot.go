package server

import (
	"github.com/gin-gonic/gin"
	"kz_bot/models"
	"net/http"
)

func (s *Server) inboxRsBot(c *gin.Context) {
	var m models.InMessage
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("inboxRsBot", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.cl.Tg.ChanRsMessage <- m
	c.JSON(http.StatusOK, "OK")
}
