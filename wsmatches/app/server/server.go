package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mentalisit/logger"
	"net/http"
	"time"
	"ws/server/getCountry"
)

type Srv struct {
	log      *logger.Logger
	cache    *getCountry.Cache
	certFile string
	keyFile  string
}

func NewSrv(log *logger.Logger) *Srv {
	server := &Srv{
		log:      log,
		cache:    getCountry.NewCache(),
		certFile: "docker/cert/RSA-cert.pem",
		keyFile:  "docker/cert/RSA-privkey.pem",
	}
	server.runServer()
	return server
}

func (s *Srv) runServer() {
	port := "4443"
	//blockedIPs := []string{
	//	"34.42.221.171",
	//"34.72.76.85",}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	//	router.Use(BlockIPMiddleware(blockedIPs))

	router.Use(gin.LoggerWithFormatter(s.CustomLogFormatter))

	router.GET("/matches", s.getWsMatches)
	router.GET("/docs", s.docs)
	router.GET("/", s.docs)
	router.GET("/corps", s.getWsCorps)
	router.GET("/poll/:id", s.poll)

	fmt.Println("Running port:" + port)

	err := router.RunTLS(":"+port, s.certFile, s.keyFile)
	if err != nil {
		s.log.ErrorErr(err)
	}
}

func (s *Srv) CustomLogFormatter(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}

	country, err := s.cache.GetLocationInfo(param.ClientIP)
	if err != nil {
		fmt.Println(err)
	}
	addr := fmt.Sprintf("%s:%s", param.ClientIP, country)

	return fmt.Sprintf("[GIN]%v|%s %3d %s|%13v|%15s|%s %-3s %s %#v\n%s",
		param.TimeStamp.Format("2006/01/02-15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		addr,
		methodColor, param.Method, resetColor,
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
