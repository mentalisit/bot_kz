package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mentalisit/logger"
	"kz_bot/clients"
	"net/http"
	"net/http/pprof"
	"time"
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

	r.Use(TimeoutMiddleware(1 * time.Minute))

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
		ds.POST("/GetMembersRoles", s.discordGetMembersRoles)
		ds.POST("/SendPic", s.discordSendPic)
		ds.GET("/GetAvatarUrl", s.discordGetAvatarUrl)
		ds.POST("/send_poll", s.discordSendPoll)
	}

	//telegram
	//tg := r.Group("/telegram")
	//{
	//tg.POST("/send/bridge", s.telegramSendBridge)
	//tg.POST("/SendText", s.telegramSendText)
	//tg.POST("/edit", s.telegramEditMessage)
	//tg.POST("/SendPic", s.telegramSendPic)
	//tg.POST("/del", s.telegramDel)
	//tg.POST("/SendDel", s.telegramSendDelSecond)
	//tg.GET("/GetAvatarUrl", s.telegramGetAvatarUrl)

	//}

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
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Создаём контекст с тайм-аутом
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Обновляем контекст запроса
		c.Request = c.Request.WithContext(ctx)

		// Канал завершения, который проверяет, завершён ли запрос
		finished := make(chan struct{})
		go func() {
			// Продолжение обработки запроса
			c.Next()
			close(finished)
		}()

		// Проверяем, не истекло ли время
		select {
		case <-ctx.Done():
			c.Writer.WriteHeader(http.StatusGatewayTimeout)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out"})
			c.Abort()
		case <-finished:
			// Если запрос завершился до тайм-аута, продолжаем
		}
	}
}
