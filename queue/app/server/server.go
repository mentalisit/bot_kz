package server

import (
	"fmt"
	"net/http"
	"queue/config"
	"queue/kzbotdb"
	"queue/rsbotbd"
	"queue/server/rs_bot"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mentalisit/logger"
)

type Server struct {
	log   *logger.Logger
	queue *rsbotbd.Queue
	kzbot *kzbotdb.Db
	rs    *rs_bot.Client
}

func NewServer(log *logger.Logger, cfg *config.ConfigBot) *Server {
	s := &Server{
		log:   log,
		queue: rsbotbd.NewQueue(log),
		kzbot: kzbotdb.NewDb(log, cfg),
		rs:    rs_bot.NewClient(log),
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
	router.GET("/api/webhooks", s.GetWebhooks)
	router.GET("/api/active0", s.ReadAllQueueActive0)
	router.POST("/api/left", s.Left)
	router.GET("/health", HealthCheckHandler)

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
	router.GET("/api/readsouzbot", s.ReadQueueTumcha)

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

// HealthCheckHandler проверяет здоровье сервиса
func HealthCheckHandler(c *gin.Context) {
	// Если все проверки пройдены успешно
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Service is healthy",
	})
}
