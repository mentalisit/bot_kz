package server

import (
	DiscordClient "discord/discord"
	"github.com/gin-gonic/gin"
	"github.com/mentalisit/logger"
)

type Server struct {
	log *logger.Logger
	ds  *DiscordClient.Discord
}

func NewServer(ds *DiscordClient.Discord, log *logger.Logger) *Server {
	s := &Server{ds: ds, log: log}
	go s.runServer()
	return s
}
func (s *Server) runServer() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.POST("/func", s.funcInbox)
	r.POST("/bridge", s.telegramSendBridge)

	err := r.Run(":80")
	if err != nil {
		s.log.ErrorErr(err)
		return
	}
}
