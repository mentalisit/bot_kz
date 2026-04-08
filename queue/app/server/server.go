package server

import (
	"fmt"
	"net/http"
	"queue/config"
	"queue/kzbotdb"
	"queue/rsbotbd"
	"queue/server/getCountry"
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
	cache *getCountry.Cache
}

func NewServer(log *logger.Logger, cfg *config.ConfigBot) *Server {
	s := &Server{
		log:   log,
		queue: rsbotbd.NewQueue(log),
		kzbot: kzbotdb.NewDb(log, cfg),
		rs:    rs_bot.NewClient(log),
		cache: getCountry.NewCache(),
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
	router.Use(gin.LoggerWithFormatter(s.CustomLogFormatter))

	registerRoutes := func(r gin.IRouter) {
		r.GET("/queue", s.ReadAllQueue)
		r.POST("/api/queue", s.QueueApi)
		r.POST("/api/queue2", s.QueueApi2)
		r.GET("/api/webhooks", s.GetWebhooks)
		r.GET("/api/battles", s.GetBattlesAll)
		r.GET("/api/active0", s.ReadAllQueueActive0)
		r.POST("/api/left", s.Left)
		r.GET("/health", HealthCheckHandler)
	}

	registerRoutes(router)

	group := router.Group("/queue")
	registerRoutes(group)

	fmt.Println("Running port:" + port)
	//err := router.RunTLS(":"+port, "docker/cert/RSA-cert.pem", "docker/cert/RSA-privkey.pem")
	err := router.Run(":" + port)
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
	mdate := tm.Format("2006-01-02")
	mtime := tm.Format("15:04")
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

func (s *Server) CustomLogFormatter(param gin.LogFormatterParams) string {
	if param.Method == "OPTIONS" {
		return ""
	}

	latency := param.Latency.String()
	if param.Latency > time.Minute {
		latency = param.Latency.Truncate(time.Second).String()
	} else if param.Latency > time.Millisecond {
		latency = fmt.Sprintf("%.3fms", float64(param.Latency.Microseconds())/1000.0)
	} else if param.Latency > time.Microsecond {
		latency = fmt.Sprintf("%.3fµs", float64(param.Latency.Nanoseconds())/1000.0)
	}

	// Берём реальный IP из keys, если есть — иначе param.ClientIP
	clientIP := param.ClientIP
	if ip, ok := param.Keys["clientIP"]; ok {
		if ipStr, ok := ip.(string); ok && ipStr != "" {
			clientIP = ipStr
		}
	}

	arr, _ := s.cache.GetLocationInfo(clientIP)
	country := fmt.Sprintf("%s %s", clientIP, arr)

	return fmt.Sprintf("%s | %3d | %7v | %15s | %-5s | %#v\n%s",
		param.TimeStamp.Format("15:04"),
		param.StatusCode,
		latency,
		country,
		param.Method,
		param.Path,
		param.ErrorMessage,
	)
}
