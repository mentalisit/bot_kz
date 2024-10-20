package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) EditMessage(c *gin.Context) {
	var m apiRs

	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("func EditMessage", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a := answer{Time: time.Now()}

	mID, _ := strconv.Atoi(m.MessageId)
	a.ArrError = s.tg.EditText(m.Channel, mID, m.Text, m.ParseMode)
	if a.ArrError != nil {
		s.log.ErrorErr(a.ArrError)
		s.log.InfoStruct("apiRs", m)
		c.JSON(http.StatusBadRequest, a)
		return
	}
	a.ArrBool = true

	c.JSON(http.StatusOK, a)
}

func (s *Server) EditMessageTextKey(c *gin.Context) {
	var m apiRs

	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("func EditMessageTextKey", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a := answer{Time: time.Now()}

	mID, _ := strconv.Atoi(m.MessageId)
	err := s.tg.EditMessageTextKey(m.Channel, mID, m.Text, m.LevelRs)
	if err != nil {
		fmt.Println(err)
	}
	a.ArrBool = true

	c.JSON(http.StatusOK, a)
}
