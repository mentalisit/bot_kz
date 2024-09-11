package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mentalisit/logger"
	"queue/config"
	"queue/kzbotdb"
	"queue/rsbotbd"
)

type Server struct {
	log   *logger.Logger
	queue *rsbotbd.Queue
	kzbot *kzbotdb.Db
}

func NewServer(log *logger.Logger, cfg *config.ConfigBot) *Server {
	s := &Server{
		log:   log,
		queue: rsbotbd.NewQueue(log),
		kzbot: kzbotdb.NewDb(log, cfg),
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

	fmt.Println("Running port:" + port)
	err := router.RunTLS(":"+port, "cert/RSA-cert.pem", "cert/RSA-privkey.pem")
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
