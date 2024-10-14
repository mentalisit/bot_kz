package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mentalisit/logger"
	"telegram/telegram"
)

type Server struct {
	log *logger.Logger
	tg  *telegram.Telegram
}

func NewServer(tg *telegram.Telegram, log *logger.Logger) *Server {
	s := &Server{tg: tg, log: log}
	go s.runServer()
	return s
}
func (s *Server) runServer() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.POST("/func", s.funcInbox)
	r.POST("/sendPic", s.telegramSendPic)
	r.POST("/bridge", s.telegramSendBridge)
	r.GET("/GetAvatarUrl", s.telegramGetAvatarUrl)

	err := r.Run(":80")
	if err != nil {
		s.log.ErrorErr(err)
		return
	}
}
