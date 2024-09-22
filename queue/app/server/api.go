package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) QueueApi(c *gin.Context) {
	var rawData json.RawMessage

	// Читаем тело запроса как необработанные JSON-данные
	if err := c.ShouldBindJSON(&rawData); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("QueueApi", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if c.RemoteIP() != "192.99.246.144" && c.RemoteIP() != "172.18.0.1" {
		fmt.Println(c.RemoteIP())
		fmt.Println(string(rawData))
	}
	// Пробуем десериализовать как []DataCaprican
	var m []DataCaprican
	err := json.Unmarshal(rawData, &m)
	if err == nil {
		m1 := sortData(m)
		if len(m1) > 0 || string(rawData) == "[]" {
			fmt.Printf("	DataCaprican len%d data:%+v\n", len(m1), m1)
			merged1 = m1
			c.JSON(http.StatusOK, "OK")
			return
		}
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

	//var m []DataCaprican
	//if err := c.ShouldBindJSON(&m); err != nil {
	//	s.log.ErrorErr(err)
	//	s.log.InfoStruct("QueueApi", c.Request.Body)
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//merged1 = sortData(m)
	//
	//c.JSON(http.StatusOK, "OK")
}

func (s *Server) QueueApi2(c *gin.Context) {
	var m interface{}
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("QueueApi2", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.log.InfoStruct("QueueApi2", m)

	c.JSON(http.StatusOK, "OK")
}
