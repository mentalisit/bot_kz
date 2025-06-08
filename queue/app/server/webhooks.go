package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"queue/models"
	"strconv"
)

func (s *Server) GetWebhooks(c *gin.Context) {
	s.PrintGoroutine()
	ts := c.DefaultQuery("ts", "")
	afterTs, _ := strconv.ParseInt(ts, 10, 64)
	eventType := c.DefaultQuery("eventType", "")
	corp := c.DefaultQuery("corp", "")

	if corp == "" {
		var all []models.Webhook
		if eventType != "" {
			rus, _ := s.kzbot.GetWebhooksByEventTypeAndCorpAfter("rus", eventType, afterTs)
			all = append(all, rus...)
			soyuz, _ := s.kzbot.GetWebhooksByEventTypeAndCorpAfter("soyuz", eventType, afterTs)
			all = append(all, soyuz...)
			c.JSON(http.StatusOK, all)
			return
		}
		rus, _ := s.kzbot.GetWebhooksByCorpAfter("rus", afterTs)
		all = append(all, rus...)
		soyuz, _ := s.kzbot.GetWebhooksByCorpAfter("soyuz", afterTs)
		all = append(all, soyuz...)
		c.JSON(http.StatusOK, all)
		return
	}

	after, err := s.kzbot.GetWebhooksByCorpAfter(corp, afterTs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, after)
}
