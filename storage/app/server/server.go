package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mentalisit/logger"
	"storage/mongodb"
)

type Server struct {
	db  *mongodb.DB
	log *logger.Logger
}

func NewServer(db *mongodb.DB, log *logger.Logger) *Server {
	s := &Server{db: db, log: log}
	go func() {
		err := s.runServer()
		if err != nil {
			recover()
		}
	}()

	return s
}
func (s *Server) runServer() error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Регистрация обработчиков bridge
	r.GET("/storage/bridge/read", s.DBReadBridgeConfig)
	r.POST("/storage/bridge/update", s.UpdateBridgeChat)
	r.POST("/storage/bridge/insert", s.InsertBridgeChat)

	//r.GET("/storage/timer/delete", s.DeleteMessageTimer)
	//r.POST("/storage/timer/insert", s.InsertTimer)

	r.GET("/storage/rsbot/readqueue", s.DBReadRsBotMySQL)

	r.GET("/storage/rsbot/read", s.DBReadRsConfig)
	r.POST("/storage/rsbot/update", s.UpdateRsConfig)
	r.POST("/storage/rsbot/insert", s.InsertRsConfig)
	r.DELETE("/storage/rsbot/delete", s.DeleteRsConfig)

	err := r.Run(":80")
	if err != nil {
		s.log.ErrorErr(err)
		return err
	}
	return nil
}
