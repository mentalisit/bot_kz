package webServer

import (
	"net/http"
	"rs/models"

	"github.com/gin-gonic/gin"
)

type MobileSettingDelete struct {
	UUID      string `json:"uuid" binding:"required"`
	Action    string `json:"action" binding:"required"`
	Mode      string `json:"mode,omitempty"`
	ChannelId string `json:"channelId,omitempty"`
}

func (s *Server) deleteMobileSetting(c *gin.Context) {
	var req MobileSettingDelete

	// Парсим JSON из тела запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Invalid JSON format: " + err.Error(),
		})
		return
	}

	// Обрабатываем действие
	switch req.Action {
	case "removeRsChannel":
		{
			conf := s.db.ReadConfigV2Uid(req.UUID)
			if conf != nil && conf.Channels[req.ChannelId] != nil {
				if len(conf.Channels) == 1 {
					s.db.DeleteConfigV2(*conf)
					c.JSON(http.StatusOK, gin.H{
						"status": "success",
						"uuid":   req.UUID,
					})
				} else {
					delete(conf.Channels, req.ChannelId)
					s.db.UpdateConfigV2Channels(*conf)
					c.JSON(http.StatusOK, gin.H{
						"status": "success",
						"uuid":   req.UUID,
					})
				}
			}
		}
	case "removeBridgeChannel":
		{
			existing, found := s.db.ReadBridgeConfigByNameRelay(req.UUID)
			if found {
				newChannel := make(map[string][]models.Bridge2Configs)
				for s2, configs := range existing.Channel {
					for _, config := range configs {
						if config.ChannelId != req.ChannelId {
							if newChannel[s2] == nil {
								newChannel[s2] = []models.Bridge2Configs{}
								newChannel[s2] = append(newChannel[s2], config)
							}
						}
					}
				}
				if len(newChannel) == 0 {
					s.db.DeleteBridge2Chat(existing)
					c.JSON(http.StatusOK, gin.H{
						"status": "success",
					})
				} else {
					existing.Channel = newChannel
					err := s.db.UpdateBridgeConfig(existing)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"status": "error",
							"error":  err.Error(),
						})
						return
					}
					c.JSON(http.StatusOK, gin.H{
						"status": "success",
					})
				}
			}
		}
	case "removeScoreboardChannel":
		{
			params := s.db.ScoreboardReadByUid(req.UUID)
			if params == nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status": "error",
				})
				return
			}
			if params.Uid != "" {
				var newChannel []models.ChannelsInfo
				for _, ch := range params.Channels {
					if ch.ChannelId != req.ChannelId {
						newChannel = append(newChannel, ch)
					}
				}

				if len(newChannel) == 0 {
					s.db.ScoreboardDeleteByUid(params.Uid)
				} else {
					params.Channels = newChannel
					s.db.ScoreboardUpdateParamChannels(*params)
				}

				c.JSON(http.StatusOK, gin.H{
					"status": "success",
				})
			}

		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Unknown action: " + req.Action,
		})
	}
}
