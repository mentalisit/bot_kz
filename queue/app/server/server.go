package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mentalisit/logger"
	"queue/config"
	"queue/kzbotdb"
	"queue/rsbotbd"
	"runtime"
	"time"
)

type Server struct {
	log   *logger.Logger
	queue *rsbotbd.Queue
	kzbot *kzbotdb.Db
	ai    *Gemini
}

func NewServer(log *logger.Logger, cfg *config.ConfigBot) *Server {
	s := &Server{
		log:   log,
		queue: rsbotbd.NewQueue(log),
		kzbot: kzbotdb.NewDb(log, cfg),
		ai:    NewGemini(log),
	}
	go func() {
		err := s.runServer()
		if err != nil {
			recover()
		}
	}()
	go func() {
		err := s.runServerApi()
		if err != nil {
			recover()
		}
	}()

	return s
}
func (s *Server) runServer() error {
	port := "9443"

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/queue", s.ReadAllQueue)
	router.POST("/api/queue", s.QueueApi)
	router.POST("/api/queue2", s.QueueApi2)
	router.POST("/ai", s.GeminiAI)

	fmt.Println("Running port:" + port)
	err := router.RunTLS(":"+port, "docker/cert/RSA-cert.pem", "docker/cert/RSA-privkey.pem")
	if err != nil {
		return err
	}
	return nil
}
func (s *Server) runServerApi() error {
	port := "80"

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/queue", s.ReadAllQueue)

	fmt.Println("Running port:" + port)
	err := router.Run(":" + port)
	if err != nil {
		return err
	}
	return nil
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
