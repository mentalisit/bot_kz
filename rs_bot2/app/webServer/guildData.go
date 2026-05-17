package webServer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) GetGuildData(c *gin.Context) {
	gid := c.Query("gid")
	if gid == "" {
		c.JSON(http.StatusBadRequest, "No gid specified")
		return
	}
	GID, _ := uuid.Parse(gid)
	full, err := s.db.GuildGetFull(GID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, full)
}
