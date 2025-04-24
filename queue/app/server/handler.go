package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"queue/rsbotbd/GetDataSborkz"
)

func (s *Server) ReadAllQueue(c *gin.Context) {
	s.PrintGoroutine()
	level := c.DefaultQuery("level", "")
	c.JSON(http.StatusOK, s.Queue(level))
}

func (s *Server) ReadAllQueueActive0(c *gin.Context) {
	s.log.InfoStruct("ReadAllQueueActive0", c.ClientIP())
	s.PrintGoroutine()
	type SborkzMinimal struct {
		Corpname string
		Name     string
		UserId   string
		Lvlkz    string
		Timedown int
	}
	var active []SborkzMinimal
	all := s.kzbot.SelectSborkzActive()
	if len(all) > 0 {
		for _, sb := range all {
			if sb.Tip == "tg" {

				active = append(active, SborkzMinimal{
					Corpname: sb.Corpname,
					Name:     sb.Name,
					UserId:   sb.UserId,
					Lvlkz:    sb.Lvlkz,
					Timedown: sb.Timedown,
				})
			}
		}
	}
	c.JSON(http.StatusOK, active)
}

func (s *Server) Left(c *gin.Context) {
	s.PrintGoroutine()
	var rawData json.RawMessage

	// Читаем тело запроса как необработанные JSON-данные
	if err := c.ShouldBindJSON(&rawData); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("Left", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Пробуем десериализовать как []DataCaprican
	var m []int64
	err := json.Unmarshal(rawData, &m)
	if err == nil {
		c.JSON(http.StatusOK, "OK")
		s.log.InfoStruct("Left", m)
		return
	}

	if len(rawData) > 4 {
		// просто выводим полученные данные как есть
		var otherData interface{}
		if err := json.Unmarshal(rawData, &otherData); err == nil {
			s.log.InfoStruct("	Received other data", otherData)
			c.JSON(http.StatusOK, "OK")
			return
		}
	}
	fmt.Println(len(rawData))
	fmt.Println(string(rawData))
	// Если ничего не удалось разобрать, выводим ошибку
	s.log.Error("Не удалось разобрать данные")
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
}

func (s *Server) ReadQueueTumcha(c *gin.Context) {
	s.PrintGoroutine()
	namesIds := GetDataSborkz.ReadQueueTumchaNameIds()
	if len(namesIds) == 0 {
		c.JSON(http.StatusNoContent, "nil")
		return
	}

	c.JSON(http.StatusOK, namesIds)
}
