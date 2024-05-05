package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mentalisit/logger"
)

type Srv struct {
	log *logger.Logger
}

func NewSrv(log *logger.Logger, port string) *Srv {
	server := &Srv{log: log}
	server.runServer(port)
	return server
}

func (s *Srv) runServer(port string) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/matches", s.getWsMatches)
	router.GET("/docs", s.docs)
	router.GET("/", s.docs)
	router.GET("/corps", s.getWsCorps)

	fmt.Println("Running port:" + port)
	err := router.RunTLS(":"+port, "cert/RSA-cert.pem", "cert/RSA-privkey.pem")
	if err != nil {
		s.log.ErrorErr(err)
	}
}
