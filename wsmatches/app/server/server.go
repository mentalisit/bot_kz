package server

import (
	"fmt"
	"net/http"
	"time"
	"ws/dbpostgres"
	"ws/server/getCountry"

	"github.com/gin-gonic/gin"
	"github.com/mentalisit/logger"
)

type Srv struct {
	log      *logger.Logger
	cache    *getCountry.Cache
	Db       *dbpostgres.Db
	certFile string
	keyFile  string
}

func NewSrv(log *logger.Logger, Db *dbpostgres.Db) *Srv {
	server := &Srv{
		log:      log,
		cache:    getCountry.NewCache(),
		Db:       Db,
		certFile: "docker/cert/RSA-cert.pem",
		keyFile:  "docker/cert/RSA-privkey.pem",
	}
	server.runServer()
	return server
}

func (s *Srv) runServer() {
	port := "24443"

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(CORSMiddleware())

	router.Use(gin.LoggerWithFormatter(s.CustomLogFormatter))

	// Регистрируем маршруты через функцию
	registerRoutes := func(r gin.IRouter) {
		r.GET("/matches", s.getWsMatches)
		r.GET("/docs", s.docs)
		r.GET("/", s.docs)
		r.GET("/corps", s.getWsCorps)
		r.GET("/poll/:id", s.poll)
		r.GET("/api/poll/:id", s.pollAPI)
		r.GET("/health", HealthCheckHandler)
	}

	// Маршруты без префикса: /matches, /docs, ...
	registerRoutes(router)

	// Маршруты с префиксом /ws: /ws/matches, /ws/docs, ...
	wsGroup := router.Group("/ws")
	registerRoutes(wsGroup)

	fmt.Println("Running port:" + port)

	err := router.Run(":" + port)
	if err != nil {
		s.log.ErrorErr(err)
	}
}

func (s *Srv) CustomLogFormatter(param gin.LogFormatterParams) string {
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

func BlockIPMiddleware(blockedIPs []string) gin.HandlerFunc {
	blockedSet := make(map[string]struct{}, len(blockedIPs))
	for _, ip := range blockedIPs {
		blockedSet[ip] = struct{}{}
	}

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		if _, blocked := blockedSet[clientIP]; blocked {
			// Если IP-адрес заблокирован, возвращаем ошибку
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort() // Прекращаем дальнейшую обработку запроса
			return
		}
		c.Next() // Продолжаем обработку для других IP-адресов
	}
}

// HealthCheckHandler проверяет здоровье сервиса
func HealthCheckHandler(c *gin.Context) {
	// Если все проверки пройдены успешно
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Service is healthy",
	})
}

// CORSMiddleware Выносим CORS в middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Добавляем CORS заголовки
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")

		// Обрабатываем preflight OPTIONS запросы
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
