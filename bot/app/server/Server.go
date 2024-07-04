package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mentalisit/logger"
	"kz_bot/clients"
	"net/http/pprof"
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
	ds := r.Group("/discord")
	{
		ds.POST("/send/bridge", s.discordSendBridge)
		ds.POST("/del", s.discordDel)
		ds.POST("/SendDel", s.discordSendDelSecond)
		ds.POST("/SendText", s.discordSendText)
		ds.POST("/edit", s.discordEditMessage)
		ds.POST("/GetRoles", s.discordGetRoles)
		ds.POST("/CheckRole", s.discordCheckRole)
		ds.POST("/SendPic", s.discordSendPic)
		ds.GET("/GetAvatarUrl", s.discordGetAvatarUrl)
	}

	//telegram
	tg := r.Group("/telegram")
	{
		tg.POST("/send/bridge", s.telegramSendBridge)
		tg.POST("/SendText", s.telegramSendText)
		tg.POST("/edit", s.telegramEditMessage)
		tg.POST("/SendPic", s.telegramSendPic)
		tg.POST("/del", s.telegramDel)
		tg.POST("/SendDel", s.telegramSendDelSecond)
		tg.GET("/GetAvatarUrl", s.telegramGetAvatarUrl)

	}

	//rsbot
	r.POST("/inbox", s.inboxRsBot)

	pprofRoutes(r)

	err := r.Run(":80")
	if err != nil {
		s.log.ErrorErr(err)
		return
	}
}

// pprofRoutes регистрирует обработчики pprof на маршрутизаторе Gin
func pprofRoutes(router *gin.Engine) {
	pprofGroup := router.Group("/debug/pprof")
	{
		pprofGroup.GET("/", gin.WrapF(pprof.Index))
		pprofGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
		pprofGroup.POST("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/trace", gin.WrapF(pprof.Trace))
		pprofGroup.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		pprofGroup.GET("/block", gin.WrapH(pprof.Handler("block")))
		pprofGroup.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		pprofGroup.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		pprofGroup.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		pprofGroup.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}
}
