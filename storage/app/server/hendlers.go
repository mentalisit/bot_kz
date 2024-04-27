package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"storage/models"
	"storage/rsbotbd"
)

func (s *Server) DBReadBridgeConfig(c *gin.Context) {
	config, err := s.db.DBReadBridgeConfig()
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, config)
}
func (s *Server) UpdateBridgeChat(c *gin.Context) {
	s.log.Info("UpdateBridgeChat")
	var br models.BridgeConfig

	if err := c.ShouldBindJSON(&br); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.db.UpdateBridgeChat(br)
	c.JSON(http.StatusOK, gin.H{"status": "Update done"})
}
func (s *Server) InsertBridgeChat(c *gin.Context) {
	s.log.Info("InsertBridgeChat")
	var br models.BridgeConfig
	if err := c.ShouldBindJSON(&br); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.db.UpdateBridgeChat(br)
	c.JSON(http.StatusOK, gin.H{"status": "Insert done"})
}
func (s *Server) DBReadRsBotMySQL(c *gin.Context) {
	c.JSON(http.StatusOK, rsbotbd.GetQueue())
}

func (s *Server) DBReadRsConfig(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.ReadConfigRs())
}
func (s *Server) InsertRsConfig(c *gin.Context) {
	s.log.Info("InsertRsConfig")
	var Rs models.CorporationConfig
	if err := c.ShouldBindJSON(&Rs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.db.InsertConfigRs(Rs)
	c.JSON(http.StatusOK, gin.H{"status": "Insert done"})
}
func (s *Server) UpdateRsConfig(c *gin.Context) {
	s.log.Info("UpdateRsConfig")
	var br models.CorporationConfig

	if err := c.ShouldBindJSON(&br); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.db.UpdateRsConfig(br)
	c.JSON(http.StatusOK, gin.H{"status": "Update done"})
}
func (s *Server) DeleteRsConfig(c *gin.Context) {
	s.log.Info("DeleteRsConfig")
	var br models.CorporationConfig
	if err := c.ShouldBindJSON(&br); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	s.db.DeleteConfigRs(br)
	c.JSON(http.StatusOK, gin.H{"status": "Delete done"})
}
