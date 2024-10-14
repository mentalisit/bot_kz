package server

import (
	"github.com/gin-gonic/gin"
	"kz_bot/models"
	"kz_bot/pkg/utils"
	"net/http"
)

func (s *Server) discordSendBridge(c *gin.Context) {
	ch := utils.WaitForMessage("discordSendBridge")
	defer close(ch)
	var m models.BridgeSendToMessenger
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("discordSendBridge", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	messageDs := s.cl.DS.SendBridgeFuncRest(m)
	if len(messageDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "len(messageDs)==0"})
		s.log.InfoStruct("discordSendBridge", c.Request.Body)
		return
	}
	c.JSON(http.StatusOK, messageDs)
}

func (s *Server) discordDel(c *gin.Context) {
	ch := utils.WaitForMessage("discordDel")
	defer close(ch)
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
	ch := utils.WaitForMessage("discordSendDelSecond")
	defer close(ch)
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
	ch := utils.WaitForMessage("discordSendText")
	defer close(ch)
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
	ch := utils.WaitForMessage("discordEditMessage")
	defer close(ch)
	var m models.EditText
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("discordEditMessage", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.cl.DS.EditMessage(m.Channel, m.MessageId, m.Text)
	c.JSON(http.StatusOK, "ok")
}

func (s *Server) discordGetRoles(c *gin.Context) {
	ch := utils.WaitForMessage("discordGetRoles")
	defer close(ch)
	var guildId string
	if err := c.ShouldBindJSON(&guildId); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("discordGetRoles", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roles := s.cl.DS.GetRoles(guildId)
	c.JSON(http.StatusOK, roles)
}

func (s *Server) discordCheckRole(c *gin.Context) {
	ch := utils.WaitForMessage("discordCheckRole")
	defer close(ch)
	var m models.CheckRoleStruct
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("discordCheckRole", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok := s.cl.DS.CheckRole(m.GuildId, m.MemberId, m.RoleId)
	c.JSON(http.StatusOK, ok)
}

func (s *Server) discordGetMembersRoles(c *gin.Context) {
	ch := utils.WaitForMessage("discordGetMembersRoles")
	defer close(ch)
	var guildId string
	if err := c.ShouldBindJSON(&guildId); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("discordGetMembersRoles", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roles := s.cl.DS.GetMembersRoles(guildId)
	c.JSON(http.StatusOK, roles)
}

func (s *Server) discordSendPic(c *gin.Context) {
	ch := utils.WaitForMessage("discordSendPic")
	defer close(ch)
	var m models.SendPic
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("discordSendPic", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := s.cl.DS.SendPic(m.Channel, m.Text, m.Pic)
	if err != nil {
		s.log.ErrorErr(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Message sent to Discord successfully"})
}
func (s *Server) discordGetAvatarUrl(c *gin.Context) {
	ch := utils.WaitForMessage("discordGetAvatarUrl")
	defer close(ch)
	userid := c.Query("userid")
	if userid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid must not be empty"})
		return
	}
	urlAvatar := s.cl.DS.GetAvatarUrl(userid)
	c.JSON(http.StatusOK, urlAvatar)
}
