package server

import "C"
import (
	"encoding/json"
	"fmt"
	"net/http"
	"queue/rsbotbd/GetDataSborkz"

	"github.com/gin-gonic/gin"
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
	type NameID struct {
		NameID string `json:"nameid"`
	}

	var rawData json.RawMessage

	// Читаем тело запроса как необработанные JSON-данные
	if err := c.ShouldBindJSON(&rawData); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("Left", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Пробуем десериализовать
	var nameIDs []NameID
	var namesIds []string

	err := json.Unmarshal(rawData, &nameIDs)
	if err != nil {
		err = json.Unmarshal(rawData, &namesIds)
		if err != nil {
			var id string
			err = json.Unmarshal(rawData, &id)
			if err != nil {
				s.log.ErrorErr(err)
				s.log.InfoStruct("Left", string(rawData))
				fmt.Println("Error parsing JSON:", err)

			}
			namesIds = append(namesIds, id)
		}
	} else {
		for _, nameID := range nameIDs {
			namesIds = append(namesIds, nameID.NameID)
		}
	}

	if len(namesIds) > 0 {
		fmt.Println("Left ", namesIds)
		s.rs.SendOtherQueue(namesIds)
	}
}

func (s *Server) ReadQueueTumcha(c *gin.Context) {
	s.PrintGoroutine()
	namesIds := GetDataSborkz.ReadQueueTumchaNameIds() //namesIds []int64
	fmt.Printf("namesIds: %+v\n", namesIds)
	if len(namesIds) == 0 {
		c.JSON(http.StatusNoContent, "nil")
		return
	}

	c.JSON(http.StatusOK, namesIds)
}
