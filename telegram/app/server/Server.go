package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mentalisit/logger"
	"runtime"
	"telegram/telegram"
	"time"
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

	r.POST("/send", s.Send)
	r.POST("/send_del", s.SendDel)
	r.POST("/send_help", s.SendHelp)
	r.POST("/send_embed", s.SendEmbed)
	r.POST("/send_embed_time", s.SendEmbedTime)
	r.POST("/chat_typing", s.SendChatTyping)

	r.POST("/edit_message", s.EditMessage)
	r.POST("/edit_message_text_key", s.EditMessageTextKey)

	r.POST("/delete_message", s.DeleteMessage)
	r.POST("/delete_message_second", s.DeleteMessageSecond)

	r.POST("/check_admin", s.CheckAdmin)
	r.POST("/get_avatar_url", s.GetAvatarUrl)
	r.POST("/send_poll", s.telegramSendPoll)

	err := r.Run(":80")
	if err != nil {
		s.log.ErrorErr(err)
		return
	}
}
func (s *Server) PrintGoroutine() {
	goroutine := runtime.NumGoroutine()
	tm := time.Now()
	mdate := (tm.Format("2006-01-02"))
	mtime := (tm.Format("15:04"))
	text := fmt.Sprintf(" %s %s Горутин  %d\n", mdate, mtime, goroutine)
	if goroutine > 120 {
		s.log.Info(text)
		s.log.Panic(text)
	} else if goroutine > 50 && goroutine%10 == 0 {
		s.log.Info(text)
	}

	fmt.Println(text)
}
