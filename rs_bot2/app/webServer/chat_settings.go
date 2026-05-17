package webServer

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getChatUserSettings(c *gin.Context) {
	uuidStr := c.Query("uuid")
	if uuidStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}

	gidStr := c.DefaultQuery("gid", "00000000-0000-0000-0000-000000000000")
	settings, err := s.db.GetChatUserSettings(uuidStr, gidStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   settings,
	})
}

func (s *Server) postChatUserSettings(c *gin.Context) {
	var req struct {
		UUID     string         `json:"uuid"`
		GID      string         `json:"gid"`
		Settings map[string]any `json:"settings"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if req.UUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}
	if req.GID == "" {
		req.GID = "00000000-0000-0000-0000-000000000000"
	}
	if req.Settings == nil {
		req.Settings = map[string]any{}
	}

	if err := s.db.SaveChatUserSettings(req.UUID, req.GID, req.Settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
