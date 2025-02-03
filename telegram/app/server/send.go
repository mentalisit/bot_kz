package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) Send(c *gin.Context) {
	s.log.Info("using server http")
	var m apiRs

	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("func Send", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a := answer{Time: time.Now()}

	a.ArrString, a.ArrError = s.tg.SendChannel(m.Channel, m.Text, m.ParseMode)
	if a.ArrError != nil {
		s.log.ErrorErr(a.ArrError)
		s.log.InfoStruct("apiRs", m)
		c.JSON(http.StatusBadRequest, a)
		return
	}

	c.JSON(http.StatusOK, a)
}

func (s *Server) SendDel(c *gin.Context) {
	s.log.Info("using server http")
	var m apiRs

	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("func SendDel", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a := answer{Time: time.Now()}

	a.ArrBool, a.ArrError = s.tg.SendChannelDelSecond(m.Channel, m.Text, m.Second)
	if a.ArrError != nil {
		s.log.ErrorErr(a.ArrError)
		s.log.InfoStruct("apiRs", m)
		c.JSON(http.StatusBadRequest, a)
		return
	}

	c.JSON(http.StatusOK, a)
}

func (s *Server) SendHelp(c *gin.Context) {
	s.log.Info("using server http")
	var m apiRs

	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("func SendHelp", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ifUser := false
	if m.Second == 1 {
		ifUser = true
	}

	a := answer{Time: time.Now()}

	a.ArrString, a.ArrError = s.tg.SendHelp(m.Channel, m.Text, m.MessageId, ifUser)
	if a.ArrError != nil {
		s.log.ErrorErr(a.ArrError)
		s.log.InfoStruct("apiRs", m)
		c.JSON(http.StatusBadRequest, a)
		return
	}

	c.JSON(http.StatusOK, a)
}

func (s *Server) SendEmbed(c *gin.Context) {
	s.log.Info("using server http")
	var m apiRs

	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("func SendEmbed", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a := answer{Time: time.Now()}

	a.ArrInt, a.ArrError = s.tg.SendEmbed(m.LevelRs, m.Channel, m.Text)
	if a.ArrError != nil {
		s.log.ErrorErr(a.ArrError)
		s.log.InfoStruct("apiRs", m)
		c.JSON(http.StatusBadRequest, a)
		return
	}
	a.ArrString = strconv.Itoa(a.ArrInt)

	c.JSON(http.StatusOK, a)
}

func (s *Server) SendEmbedTime(c *gin.Context) {
	s.log.Info("using server http")
	var m apiRs

	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("func SendEmbedTime", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a := answer{Time: time.Now()}

	a.ArrInt, a.ArrError = s.tg.SendEmbedTime(m.Channel, m.Text)
	if a.ArrError != nil {
		s.log.ErrorErr(a.ArrError)
		s.log.InfoStruct("apiRs", m)
		c.JSON(http.StatusBadRequest, a)
		return
	}
	a.ArrString = strconv.Itoa(a.ArrInt)

	c.JSON(http.StatusOK, a)
}

func (s *Server) SendChatTyping(c *gin.Context) {
	s.log.Info("using server http")
	var m apiRs

	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("func ChatTyping", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a := answer{Time: time.Now()}

	a.ArrError = s.tg.ChatTyping(m.Channel)
	a.ArrBool = true

	c.JSON(http.StatusOK, a)
}
