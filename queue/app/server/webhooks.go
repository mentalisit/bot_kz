package server

import (
	"net/http"
	"queue/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetWebhooks(c *gin.Context) {
	s.PrintGoroutine()
	ts := c.DefaultQuery("ts", "")
	afterTs, _ := strconv.ParseInt(ts, 10, 64)
	eventType := c.DefaultQuery("eventType", "")

	var all []models.Webhook

	if eventType != "" {
		all, _ = s.kzbot.GetWebhooksByEventTypeAndAfterTs(eventType, afterTs)
	} else {
		all, _ = s.kzbot.GetWebhooksByAfter(afterTs)
	}

	c.JSON(http.StatusOK, all)
	return
}

func (s *Server) GetBattlesAll(c *gin.Context) {
	s.PrintGoroutine()

	all, err := s.kzbot.BattlesGetSeasonAll(48)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, all)
	return
}
