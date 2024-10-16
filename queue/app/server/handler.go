package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) ReadAllQueue(c *gin.Context) {
	s.PrintGoroutine()
	level := c.DefaultQuery("level", "")
	c.JSON(http.StatusOK, s.Queue(level))
}
