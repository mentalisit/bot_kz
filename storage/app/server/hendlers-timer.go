package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"storage/models"
)

func (s *Server) InsertTimer(c *gin.Context) {
	var br models.Timer
	if err := c.ShouldBindJSON(&br); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.db.TimerInsert(br)
	c.JSON(http.StatusOK, gin.H{"status": "Insert done"})
}
func (s *Server) DeleteMessageTimer(c *gin.Context) {
	mm := s.db.TimerDeleteMessage()
	c.JSON(http.StatusOK, mm)
}
