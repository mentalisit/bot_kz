package server

import (
	"github.com/gin-gonic/gin"
	"kz_bot/models"
	"net/http"
)

func (s *Server) discordSendBridge(c *gin.Context) {
	var m models.BridgeSendToMessenger
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("discordSendBridge", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	messageDs := s.cl.Ds.SendBridgeFuncRest(m)
	c.JSON(http.StatusOK, messageDs)
}

func (s *Server) discordDel(c *gin.Context) {
	var m models.DeleteMessageStruct
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("discordDel", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.cl.Ds.DeleteMessage(m.Channel, m.MessageId)
	c.JSON(http.StatusOK, gin.H{"status": "Message sent to Discord successfully"})
}

func (s *Server) discordSendDelSecond(c *gin.Context) {
	var m models.SendTextDeleteSeconds
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("discordSendDelSec", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.cl.Ds.SendChannelDelSecond(m.Channel, m.Text, m.Seconds)
	c.JSON(http.StatusOK, gin.H{"status": "Message sent to Discord successfully"})
}

func (s *Server) discordSendText(c *gin.Context) {
	var m models.SendText
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("discordSendText", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mid := s.cl.Ds.Send(m.Channel, m.Text)
	c.JSON(http.StatusOK, mid)
}

func (s *Server) discordEditMessage(c *gin.Context) {
	var m models.EditText
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("discordEditMessage", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.cl.Ds.EditMessage(m.Channel, m.MessageId, m.Text)
	c.JSON(http.StatusOK, "ok")
}

func (s *Server) discordGetRoles(c *gin.Context) {
	var guildId string
	if err := c.ShouldBindJSON(&guildId); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("discordGetRoles", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roles := s.cl.Ds.GetRoles(guildId)
	c.JSON(http.StatusOK, roles)
}

func (s *Server) discordCheckRole(c *gin.Context) {
	var m models.CheckRoleStruct
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("discordCheckRole", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok := s.cl.Ds.CheckRole(m.GuildId, m.MemberId, m.RoleId)
	c.JSON(http.StatusOK, ok)
}
func (s *Server) discordSendPic(c *gin.Context) {
	var m models.SendPic
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("discordSendPic", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := s.cl.Ds.SendPic(m.Channel, m.Text, m.Pic)
	if err != nil {
		s.log.ErrorErr(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Message sent to Discord successfully"})
}
func (s *Server) discordGetAvatarUrl(c *gin.Context) {
	userid := c.Query("userid")
	if userid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid must not be empty"})
		return
	}
	urlAvatar := s.cl.Ds.GetAvatarUrl(userid)
	c.JSON(http.StatusOK, urlAvatar)
}
