package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) DeleteMessage(c *gin.Context) {
	s.log.Info("using server http")
	var m apiRs

	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("func DeleteMessage", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a := answer{Time: time.Now()}

	mID, _ := strconv.Atoi(m.MessageId)
	err := s.tg.DelMessage(m.Channel, mID)
	if err != nil {
		fmt.Printf("DeleteMessage %+v %+v\n", m, err)
	}
	a.ArrBool = true

	c.JSON(http.StatusOK, a)
}

func (s *Server) DeleteMessageSecond(c *gin.Context) {
	s.log.Info("using server http")
	var m apiRs

	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("func DeleteMessage", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a := answer{Time: time.Now()}

	a.ArrError = s.tg.DelMessageSecond(m.Channel, m.MessageId, m.Second)
	if a.ArrError != nil {
		s.log.ErrorErr(a.ArrError)
		s.log.InfoStruct("apiRs", m)
		c.JSON(http.StatusBadRequest, a)
		return
	}
	a.ArrBool = true

	c.JSON(http.StatusOK, a)
}
