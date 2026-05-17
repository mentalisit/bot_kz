package webServer

import (
	"rs/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type universalWs struct {
	Action string `json:"action"`
	Data   struct {
		WsName        string `json:"wsName"`
		ChannelType   string `json:"channelType"`
		ChannelId     string `json:"channelId"`
		TypeMessenger string `json:"typeMessenger"`
		Gid           string `json:"gid"`
	} `json:"data"`
}

func (s *Server) WsUniversal(c *gin.Context) {
	var ws universalWs
	if err := c.ShouldBindJSON(&ws); err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Некорректный формат данных"})
		return
	}
	uid, err := uuid.Parse(ws.Data.Gid)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Некорректный формат uuid"})
		return
	}
	full, err := s.db.GuildGetFull(uid)
	if err != nil || full == nil {
		c.JSON(400, gin.H{"status": "error", "message": "ошибка поиска гильдии"})
		return
	}

	switch ws.Action {
	case "createWs":
		{
			// Инициализируем мапы если они nil
			if full.Data.Discussion == nil {
				full.Data.Discussion = make(map[string][]models.Channel)
			}
			if full.Data.Coordination == nil {
				full.Data.Coordination = make(map[string][]models.Channel)
			}

			if len(full.Data.Discussion[ws.Data.WsName]) == 0 {
				full.Data.Discussion[ws.Data.WsName] = []models.Channel{}
			}

			if len(full.Data.Coordination[ws.Data.WsName]) == 0 {
				full.Data.Coordination[ws.Data.WsName] = []models.Channel{}
			}
		}
	case "deletePollChannel":
		{
			var channels []models.Channel
			for _, channel := range full.Data.PollChannels {
				if channel.ChannelID != ws.Data.ChannelId {
					channels = append(channels, channel)
				}
			}
			full.Data.PollChannels = channels // Всегда устанавливаем, даже если пустой
		}
	case "deleteChannel":
		{
			// Инициализируем мапы если они nil
			if full.Data.Coordination == nil {
				full.Data.Coordination = make(map[string][]models.Channel)
			}
			if full.Data.Discussion == nil {
				full.Data.Discussion = make(map[string][]models.Channel)
			}

			if ws.Data.ChannelType == "coordination" {
				for k, channel := range full.Data.Coordination {
					if k == ws.Data.WsName {
						var newChannels []models.Channel
						for _, ch := range channel {
							if ch.ChannelID != ws.Data.ChannelId {
								newChannels = append(newChannels, ch)
							}
						}
						full.Data.Coordination[k] = newChannels // Всегда устанавливаем, даже если пустой
					}
				}
			}
			if ws.Data.ChannelType == "discussion" {
				for k, channel := range full.Data.Discussion {
					if k == ws.Data.WsName {
						var newChannels []models.Channel
						for _, ch := range channel {
							if ch.ChannelID != ws.Data.ChannelId {
								newChannels = append(newChannels, ch)
							}
						}
						full.Data.Discussion[k] = newChannels // Всегда устанавливаем, даже если пустой
					}
				}
			}
		}

	}

	err = s.db.UpdateGuildData(uid, full.Data)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "ошибка обновления гильдии " + err.Error()})
		return
	}
	c.JSON(200, full)
}
