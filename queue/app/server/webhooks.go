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

	var all []models.Webhook

	if eventType != "" {
		all, _ = s.kzbot.GetWebhooksByEventTypeAndAfterTs(eventType, afterTs)
	} else {
		all, _ = s.kzbot.GetWebhooksByAfter(afterTs)
	}

	c.JSON(http.StatusOK, all)
	return
}
