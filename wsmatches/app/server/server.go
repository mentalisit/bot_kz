package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mentalisit/logger"
)

type Srv struct {
	log *logger.Logger
}

func NewSrv(log *logger.Logger) *Srv {
	server := &Srv{log: log}
	server.runServer()
	return server
}

func (s *Srv) runServer() {
	port := "4443"
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/matches", s.getWsMatches)
	router.GET("/docs", s.docs)
	router.GET("/", s.docs)
	router.GET("/corps", s.getWsCorps)
	router.GET("/poll/:id", s.poll)

	fmt.Println("Running port:" + port)
	err := router.RunTLS(":"+port, "docker/cert/RSA-cert.pem", "docker/cert/RSA-privkey.pem")
	if err != nil {
		s.log.ErrorErr(err)
	}
}
