package server

import (
	"github.com/gin-gonic/gin"
	"kz_bot/models"
	"net/http"
	"strconv"
)

func (s *Server) telegramSendBridge(c *gin.Context) {
	var m models.BridgeSendToMessenger
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("telegramSendBridge", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	messageTg := s.cl.Tg.SendBridgeFuncRest(m)
	c.JSON(http.StatusOK, messageTg)
}

func (s *Server) telegramSendText(c *gin.Context) {
	var m models.SendText
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("telegramSendText", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := s.cl.Tg.Send(m.Channel, m.Text)
	if err != nil {
		s.log.ErrorErr(err)
		if err.Error() == "Forbidden: bot can't initiate conversation with a user" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "Message sent to Telegram successfully"})
}
func (s *Server) telegramSendPic(c *gin.Context) {
	var m models.SendPic
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("telegramSendPic", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := s.cl.Tg.SendPic(m.Channel, m.Text, m.Pic)
	if err != nil {
		s.log.ErrorErr(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Message sent to Telegram successfully"})
}
func (s *Server) telegramDel(c *gin.Context) {
	var m models.DeleteMessageStruct
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("telegramDel", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := strconv.Atoi(m.MessageId)
	if err != nil {
		s.log.ErrorErr(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.cl.Tg.DelMessage(m.Channel, id)
	c.JSON(http.StatusOK, gin.H{"status": "Message sent to telegram successfully"})
}
func (s *Server) telegramSendDelSecond(c *gin.Context) {
	var m models.SendTextDeleteSeconds
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("telegramSendDelSec", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.cl.Tg.SendChannelDelSecond(m.Channel, m.Text, m.Seconds)
	c.JSON(http.StatusOK, gin.H{"status": "Message sent to Discord successfully"})
}
