package webServer

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /api/chart/corp?corp=CorpName&period=month&level=5
func (s *Server) getChartCorp(c *gin.Context) {
	corpName := c.Query("corp")
	period := c.DefaultQuery("period", "month")
	level := c.Query("level") // опционально

	if corpName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "corp обязателен"})
		return
	}

	data, err := s.db.GetChartDataByCorpAndPeriod(corpName, period, level)
	if err != nil {
		s.log.ErrorErr(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка получения данных"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"chart":  data,
	})
}

// GET /api/chart/user?name=PlayerName&period=month&level=5
func (s *Server) getChartUser(c *gin.Context) {
	name := c.Query("name")
	period := c.DefaultQuery("period", "month")
	level := c.Query("level")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "name обязателен"})
		return
	}

	data, err := s.db.GetChartDataByUser(name, period, level)
	if err != nil {
		s.log.ErrorErr(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка получения данных"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"chart":  data,
	})
}

// GET /api/chart/levels?corp=CorpName
func (s *Server) getChartLevels(c *gin.Context) {
	corpName := c.Query("corp")
	if corpName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "corp обязателен"})
		return
	}

	levels, err := s.db.GetAvailableLevels(corpName)
	if err != nil {
		s.log.ErrorErr(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка получения уровней"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"levels": levels,
	})
}

// GET /api/chart/corps
func (s *Server) getChartCorps(c *gin.Context) {
	corps, err := s.db.GetAvailableCorps()
	if err != nil {
		s.log.ErrorErr(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка получения списка корпораций"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"corps":  corps,
	})
}

// GET /api/chart/popularity?corp=CorpName&period=month
func (s *Server) getChartPopularity(c *gin.Context) {
	corpName := c.Query("corp")
	period := c.DefaultQuery("period", "month")

	if corpName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "corp обязателен"})
		return
	}

	data, err := s.db.GetChartPopularityByCorp(corpName, period)
	if err != nil {
		s.log.ErrorErr(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка получения данных популярности"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"chart":  data, // массив моделей ChartPoint
	})
}
