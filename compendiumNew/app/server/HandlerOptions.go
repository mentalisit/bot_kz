package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) Check(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	// Получение запрашиваемого метода из заголовка Access-Control-Request-Method
	//requestedMethod := c.GetHeader("Access-Control-Request-Method")

	// Получение запрашиваемых заголовков из заголовка Access-Control-Request-Headers
	requestedHeaders := c.GetHeader("Access-Control-Request-Headers")

	// Вывод полученных опций в консоль для отладки
	//fmt.Println("Requested method:", requestedMethod)
	//fmt.Println("Requested headers:", requestedHeaders)
	if requestedHeaders != "" {
		c.Header("Access-Control-Allow-Headers", requestedHeaders)
	}
	c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
	// Отправка пустого ответа с кодом состояния 200 (OK)
	c.Status(http.StatusOK)
}
