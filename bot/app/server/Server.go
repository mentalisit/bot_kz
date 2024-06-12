package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mentalisit/logger"
	"kz_bot/clients"
)

type Server struct {
	cl  *clients.Clients
	log *logger.Logger
}

func NewServer(client *clients.Clients, log *logger.Logger) *Server {
	s := &Server{cl: client, log: log}
	go s.runServer()
	return s
}
func (s *Server) runServer() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Регистрация обработчика
	//discord
	r.POST("/discord/send/bridge", s.discordSendBridge)
	r.POST("/discord/del", s.discordDel)
	r.POST("/discord/SendDel", s.discordSendDelSecond)
	r.POST("/discord/SendText", s.discordSendText)
	r.POST("/discord/edit", s.discordEditMessage)
	r.POST("/discord/GetRoles", s.discordGetRoles)
	r.POST("/discord/CheckRole", s.discordCheckRole)
	r.POST("/discord/SendPic", s.discordSendPic)
	r.GET("/discord/GetAvatarUrl", s.discordGetAvatarUrl)

	//telegram
	r.POST("/telegram/send/bridge", s.telegramSendBridge)
	r.POST("/telegram/SendText", s.telegramSendText)
	r.POST("/telegram/edit", s.telegramEditMessage)
	r.POST("/telegram/SendPic", s.telegramSendPic)
	r.POST("/telegram/del", s.telegramDel)
	r.POST("/telegram/SendDel", s.telegramSendDelSecond)
	r.GET("/telegram/GetAvatarUrl", s.telegramGetAvatarUrl)

	//rsbot
	r.POST("/inbox", s.inboxRsBot)

	err := r.Run(":80")
	if err != nil {
		s.log.ErrorErr(err)
		return
	}
}
