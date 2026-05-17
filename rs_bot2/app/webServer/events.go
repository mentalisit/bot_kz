package webServer

import (
	"net/http"
	"rs/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) getEvents(c *gin.Context) {
	scheduleAll := s.db.ReadEventScheduleAll()

	if scheduleAll == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "database read error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   scheduleAll,
	})
}

func (s *Server) getEventId(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID parameter is required",
		})
		return
	}
	season, _ := strconv.Atoi(id)

	corp := c.Query("corp")

	var all []models.PlayerStats

	if corp != "" {
		var corps []string
		scoreboardReadByUid := s.db.ScoreboardReadByUid(corp)
		if scoreboardReadByUid != nil && scoreboardReadByUid.Game != nil && len(scoreboardReadByUid.Game) != 0 {
			for _, g := range scoreboardReadByUid.Game {
				corps = append(corps, g.GameCorporation)
			}
		}
		all = s.MergeAndSumStats(season, corps)
	}

	if len(all) == 0 {
		all, _ = s.db.BattlesGetAllId(season)
	}

	if len(all) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "database read error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   all,
	})
}

func (s *Server) getEventUser(c *gin.Context) {
	playerName := c.Query("name")
	season := c.Query("season")

	if playerName == "" || season == "" {
		c.JSON(400, gin.H{"status": "error", "message": "Параметры обязательны"})
		return
	}
	eventId, _ := strconv.Atoi(season)

	eventGames, err := s.db.BattlesGetForEvent(playerName, eventId)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   eventGames,
	})
}

func (s *Server) MergeAndSumStats(season int, corporations []string) []models.PlayerStats {
	// Используем map для группировки по имени игрока
	mergedMap := make(map[string]models.PlayerStats)

	for _, corp := range corporations {
		stats, _ := s.db.BattlesGetAll(corp, season)
		for _, p := range stats {
			if existing, ok := mergedMap[p.Player]; ok {
				// Если игрок уже есть, суммируем данные
				existing.Points += p.Points
				existing.Runs += p.Runs
				// Level обычно берем максимальный
				if p.Level > existing.Level {
					existing.Level = p.Level
				}
				mergedMap[p.Player] = existing
			} else {
				mergedMap[p.Player] = p
			}
		}
	}

	// Превращаем map обратно в слайс
	result := make([]models.PlayerStats, 0, len(mergedMap))
	for _, v := range mergedMap {
		result = append(result, v)
	}
	return result
}
