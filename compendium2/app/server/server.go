package server

import (
	"compendium/config"
	"compendium/models"
	"compendium/storage"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mentalisit/logger"
)

type Server struct {
	log *logger.Logger
	db  *storage.Storage
	In  chan models.IncomingMessage
}

func NewServer(log *logger.Logger, cfg *config.ConfigBot, st *storage.Storage) *Server {
	s := &Server{
		log: log,
		db:  st,
		In:  make(chan models.IncomingMessage, 10),
	}

	go s.RunServer(cfg.Port)
	go s.RunServerRestApi()
	return s
}

func (s *Server) RunServer(port string) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.OPTIONS("/compendium/applink/identities", s.Check)
	router.GET("/compendium/applink/identities", s.CheckIdentityHandler)

	router.OPTIONS("/compendium/applink/connect", s.Check)
	router.POST("/compendium/applink/connect", s.CheckConnectHandler)

	router.OPTIONS("/compendium/cmd/syncTech/:mode", s.Check)
	router.POST("/compendium/cmd/syncTech/:mode", s.CheckSyncTechHandler)

	router.OPTIONS("/compendium/cmd/corpdata", s.Check)
	router.GET("/compendium/cmd/corpdata", s.CheckCorpDataHandler)

	router.OPTIONS("/compendium/applink/refresh", s.Check)
	router.GET("/compendium/applink/refresh", s.CheckRefreshHandler)

	router.GET("/ws/matches", s.getWsMatches)
	//router.GET("/ws/matchesall", s.getWsMatchesAll)
	router.GET("/ws/corps", s.getWsCorps)

	fmt.Println("Running port:" + port)
	err := router.RunTLS(":"+port, "cert/RSA-cert.pem", "cert/RSA-privkey.pem")
	if err != nil {
		s.log.ErrorErr(err)
		//os.Exit(1)
	}
}