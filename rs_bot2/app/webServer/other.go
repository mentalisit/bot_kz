package webServer

import (
	"fmt"
	"net/http"
	"rs/models"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

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

// CORSMiddleware Выносим CORS в middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Добавляем CORS заголовки
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers",
			"Authorization, Content-Type, X-Sync-Mode, "+
				"X-Alt-Name, X-Corp-ID, X-Role-ID, Accept, Origin, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "false")

		// Обрабатываем preflight OPTIONS запросы
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

func (s *Server) checkConfigRs(inChannel string) (bool, models.CorporationConfigV2) {
	if inChannel != "" {
		for _, v2 := range s.db.ReadConfigV2() {
			for channel, _ := range v2.Channels {
				if channel == inChannel {
					return true, v2

				}
			}
		}

	}

	return false, models.CorporationConfigV2{}
}
