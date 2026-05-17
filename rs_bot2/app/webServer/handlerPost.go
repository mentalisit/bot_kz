package webServer

import (
	"fmt"
	"net/http"
	"rs/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MobileSettingRequest структура для приема данных настроек
type MobileSettingRequest struct {
	UUID            string                    `json:"uuid" binding:"required"`
	Action          string                    `json:"action" binding:"required"`
	Mode            string                    `json:"mode,omitempty"`
	Channels        models.ChannelsMap        `json:"channels,omitempty"`
	Bonuses         []models.GameSettings     `json:"bonuses,omitempty"`
	BridgeConfig    models.Bridge2Config      `json:"bridge_config,omitempty"`
	Subscribed      bool                      `json:"subscribed,omitempty"`
	Language        string                    `json:"Language,omitempty"`
	ScoreboardParam models.ScoreboardParamsV2 `json:"scoreboard_param,omitempty"`
}

func (s *Server) postMobileSetting(c *gin.Context) {
	var req MobileSettingRequest

	// Парсим JSON из тела запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Invalid JSON format: " + err.Error(),
		})
		return
	}
	fmt.Printf("req: %+v\n", req)

	// Валидируем UUID
	uid, err := uuid.Parse(req.UUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Invalid UUID format",
			"uuid":   req.UUID,
		})
		return
	}

	// Обрабатываем действие
	switch req.Action {
	case "setMode":
		s.handleSetMode(c, uid, req)

	case "saveRSSettings":
		s.handleSaveRSSettings(c, uid, req)

	case "saveBridge2Settings":
		s.handleSaveBridge2Settings(c, uid, req)

	case "saveNewsSettings":
		s.handleSaveNewsSettings(c, uid, req)

	case "saveScoreboardSettings":
		s.handleSaveScoreboardSettings(c, uid, req)

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Unknown action: " + req.Action,
		})
	}
}

// handleSetMode обрабатывает установку режима
func (s *Server) handleSetMode(c *gin.Context, uid uuid.UUID, req MobileSettingRequest) {
	if req.Mode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Mode is required for setMode action",
		})
		return
	}

	// TODO: Сохранить режим в базу данных
	fmt.Printf("postMobileSetting: uuid=%s, action=%s, mode=%s", uid.String(), req.Action, req.Mode)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Mode updated successfully",
		"uuid":    uid.String(),
		"mode":    req.Mode,
	})
}

// handleSaveRSSettings обрабатывает сохранение настроек RS
func (s *Server) handleSaveRSSettings(c *gin.Context, uid uuid.UUID, req MobileSettingRequest) {
	if len(req.Channels) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Channels are required for saveRSSettings action",
		})
		return
	}

	fmt.Printf(" handleSaveRSSettings %s\n %+v\n", uid.String(), req)

	other, err := s.db.GetOtherByUUID(uid)
	if err != nil || other == nil {
		s.log.ErrorErr(err)
		return
	}
	for s2, _ := range req.Channels {
		if req.Channels[s2].MessageId != "" {
			req.Channels[s2].MessageId = ""
		}
		if req.Channels[s2].Game != nil && req.Channels[s2].Game.GameCorporationId != "" {
			g := req.Channels[s2].Game
			corpInfo, _ := s.db.ReadCorpInfoByCorpID(g.GameCorporationId)
			if corpInfo != nil && corpInfo.CorpID != "" {
				if g.GameXP != corpInfo.XP {
					corpInfo.XP = g.GameXP
					corpInfo.Level = corpInfo.GetLevelByXP()
				}
				err = s.db.UpdateCorpInfo(*corpInfo)
				if err != nil {
					s.log.ErrorErr(err)
				}
			} else {
				corpInfo = &models.CorpInfo{
					CorpName:  g.GameCorporation,
					CorpID:    g.GameCorporationId,
					Level:     g.GameLevel,
					XP:        g.GameXP,
					DateEnded: time.Now().UTC(),
				}
				corpInfo.Level = corpInfo.GetLevelByXP()

				_, err = s.db.CreateCorpInfo(*corpInfo)
				if err != nil {
					s.log.ErrorErr(err)
				}
			}
		}
		if req.Channels[s2].Game != nil && req.Channels[s2].Game.GameCorporationId != "" {
			corpInfo, _ := s.db.ReadCorpInfoByCorpID(req.Channels[s2].Game.GameCorporationId)
			if corpInfo != nil && corpInfo.CorpID != "" {

			}
		}
	}

	exist, conf2 := s.checkConfigRs(other.Data.ChannelId)
	if exist {
		if len(req.Bonuses) != 0 {
			conf2.Bonuses = req.Bonuses
			s.db.UpdateConfigV2Bonuses(conf2)
		}
		for channelId, channelInfo := range req.Channels {
			conf2.Channels[channelId] = channelInfo
		}
		s.db.UpdateConfigV2Channels(conf2)
	} else {
		s.db.InsertConfigV2(models.CorporationConfigV2{
			Uid:         uid.String(),
			Channels:    req.Channels,
			Bonuses:     req.Bonuses,
			HelpMessage: make(models.HelpMessage),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         "success",
		"message":        "RS settings saved successfully",
		"uuid":           uid.String(),
		"channels_saved": len(req.Channels),
	})
}

// handleSaveBridge2Settings обрабатывает сохранение настроек Bridge2
func (s *Server) handleSaveBridge2Settings(c *gin.Context, uid uuid.UUID, req MobileSettingRequest) {
	if req.BridgeConfig.NameRelay == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "bridge_config is required for saveBridge2Settings action",
		})
		return
	}

	fmt.Printf("handleSaveBridge2Settings %s\n %+v\n", uid.String(), req.BridgeConfig)

	// Проверяем существует ли запись по name_relay
	existing, found := s.db.ReadBridgeConfigByNameRelay(req.BridgeConfig.NameRelay)
	if found {
		// Обновляем существующую запись
		req.BridgeConfig.Id = existing.Id
		err := s.db.UpdateBridgeConfig(req.BridgeConfig)
		if err != nil {
			s.log.ErrorErr(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  "Failed to update bridge config: " + err.Error(),
			})
			return
		}
	} else {
		// Создаем новую запись
		err := s.db.InsertBridgeConfig(req.BridgeConfig)
		if err != nil {
			s.log.ErrorErr(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  "Failed to insert bridge config: " + err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "Bridge2 settings saved successfully",
		"uuid":      uid.String(),
		"nameRelay": req.BridgeConfig.NameRelay,
	})
}

// handleSaveNewsSettings обрабатывает сохранение настроек новостей
func (s *Server) handleSaveNewsSettings(c *gin.Context, uid uuid.UUID, req MobileSettingRequest) {
	fmt.Printf("handleSaveNewsSettings %s\n subscribed=%v, language=%s\n", uid.String(), req.Subscribed, req.Language)

	other, err := s.db.GetOtherByUUID(uid)
	if err != nil || other == nil {
		s.log.ErrorErr(err)
		return
	}
	if req.Subscribed {
		s.db.InsertNews(other.Data.ChannelId, req.Language, other.Data.TypeMessenger)
	} else {
		s.db.DeleteNews(other.Data.ChannelId)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"message":    "News settings saved successfully",
		"uuid":       uid.String(),
		"subscribed": req.Subscribed,
		"language":   req.Language,
	})
}

// handleSaveScoreboardSettings обрабатывает сохранение настроек scoreboard
func (s *Server) handleSaveScoreboardSettings(c *gin.Context, uid uuid.UUID, req MobileSettingRequest) {

	fmt.Printf("handleSaveScoreboardSettings %s\n %+v\n", uid.String(), req.ScoreboardParam)

	sc := req.ScoreboardParam
	sc.Uid = uid.String()

	s.db.ScoreboardInsertParam(sc)

	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"message":       "Scoreboard settings saved successfully",
		"uuid":          uid.String(),
		"channelsCount": len(req.ScoreboardParam.Channels),
	})
}

func (s *Server) postChatChannel(c *gin.Context) {
	var req models.ChatChannel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ID == "" || req.Name == "" || req.CreatorUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id, name and creator_uuid are required"})
		return
	}

	if req.GID == "" {
		req.GID = "00000000-0000-0000-0000-000000000000"
	}

	if err := s.db.SaveChatChannel(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save channel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
