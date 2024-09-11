package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) QueueApi(c *gin.Context) {
	var m []DataCaprican
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("QueueApi", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	merged1 = sortData(m)

	c.JSON(http.StatusOK, "OK")
}
