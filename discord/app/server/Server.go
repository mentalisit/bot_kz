package server

//
//import (
//	DiscordClient "discord/discord"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/mentalisit/logger"
//	"runtime"
//	"time"
//)
//
//type Server struct {
//	log *logger.Logger
//	ds  *DiscordClient.Discord
//}
//
//func NewServer(ds *DiscordClient.Discord, log *logger.Logger) *Server {
//	s := &Server{ds: ds, log: log}
//	go s.runServer()
//	return s
//}
//func (s *Server) runServer() {
//	gin.SetMode(gin.ReleaseMode)
//	r := gin.New()
//	r.POST("/func", s.funcInbox)
//	r.POST("/bridge", s.telegramSendBridge)
//
//	err := r.Run(":80")
//	if err != nil {
//		s.log.ErrorErr(err)
//		return
//	}
//}
//func (s *Server) PrintGoroutine() {
//	goroutine := runtime.NumGoroutine()
//	tm := time.Now()
//	mdate := (tm.Format("2006-01-02"))
//	mtime := (tm.Format("15:04"))
//	text := fmt.Sprintf(" %s %s Горутин  %d\n", mdate, mtime, goroutine)
//	if goroutine > 120 {
//		s.log.Info(text)
//		s.log.Panic(text)
//	} else if goroutine > 50 && goroutine%10 == 0 {
//		s.log.Info(text)
//	}
//
//	fmt.Println(text)
//}
