package serverV2

import (
	"compendium_s/server/getCountry"
	postgresv2 "compendium_s/storage/postgres/postgresV2"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mentalisit/logger"
)

type ServerV2 struct {
	log        *logger.Logger
	db         *postgresv2.Db
	cache      *getCountry.Cache
	cacheReq   map[string]cacheEntry
	cacheMutex sync.Mutex
	roles      *Roles
}

type cacheEntry struct {
	data      any       // Сами данные, которые отправляем
	timestamp time.Time // Когда сохранили
}

func NewServerV2(log *logger.Logger, db *postgresv2.Db) *ServerV2 {
	s := &ServerV2{
		log:      log,
		db:       db,
		cache:    getCountry.NewCache(),
		cacheReq: make(map[string]cacheEntry),
		roles:    NewRoles(log),
	}

	return s
}

func (s *ServerV2) RegisterV2Routes(router *gin.Engine) {

	router.Use(CORSMiddleware())

	//router.GET("/compendium/v2/identities", s.CheckIdentityHandler)

	router.POST("/compendium/v2/connect", s.CheckConnectHandler)

	router.POST("/compendium/v2/syncTech/:mode", s.CheckSyncTechHandler)

	router.GET("/compendium/v2/corpdata", s.CheckCorpDataHandler)

	router.GET("/compendium/v2/refresh", s.CheckRefreshHandler)

	router.GET("/compendium/v2/corporations", s.CheckUserCorporationsHandler)
	//router.GET("/compendium/v2/tech", s.api)

	//router.Static("/compendium/avatars", "docker/compendium/avatars")
	//router.Static("/docker/compendium/avatars", "docker/compendium/avatars")
	//router.Static("/tv", "docker/compendium/tv")
	//router.GET("/health", HealthCheckHandler)

}

func (s *ServerV2) CustomLogFormatter(param gin.LogFormatterParams) string {
	if param.Method == "OPTIONS" {
		return ""
	}
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

	return fmt.Sprintf("[GIN V2]%v|%s %3d %s|%13v|%15s|%s%-3s%s|%#v\n%s",
		param.TimeStamp.Format("2006/01/02-15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		addr,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

// HealthCheckHandler проверяет здоровье сервиса
func HealthCheckHandler(c *gin.Context) {
	// Если все проверки пройдены успешно
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Service V2 is healthy",
	})
}

// CORSMiddleware Выносим CORS в middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers",
			"Authorization, Content-Type, X-Sync-Mode, "+
				"X-Alt-Name, X-Corp-ID, X-Role-ID")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
