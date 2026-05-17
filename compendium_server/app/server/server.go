package server

import (
	"compendium_s/config"
	"compendium_s/models"
	"compendium_s/server/getCountry"
	"compendium_s/storage"
	"compendium_s/storage/postgres"
	postgresv2 "compendium_s/storage/postgres/postgresV2"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mentalisit/logger"
)

type Server struct {
	log        *logger.Logger
	db         *postgres.Db
	dbV2       *postgresv2.Db
	roles      *Roles
	cache      *getCountry.Cache
	certFile   string
	keyFile    string
	cacheReq   map[string]cacheEntry
	cacheMutex sync.Mutex
}

func NewServer(log *logger.Logger, st *storage.Storage, cfg *config.ConfigBot) *Server {

	s := &Server{
		log:      log,
		db:       st.DB,
		dbV2:     st.DBv2,
		roles:    NewRoles(log),
		cache:    getCountry.NewCache(),
		certFile: "docker/cert/RSA-cert.pem",
		keyFile:  "docker/cert/RSA-privkey.pem",
		cacheReq: make(map[string]cacheEntry),
	}

	go s.RunServer()
	return s
}

func (s *Server) RunServer() {
	port := "28443"
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(CORSMiddleware())

	router.Use(gin.LoggerWithFormatter(s.CustomLogFormatter))

	registerRoutes := func(r gin.IRouter) {
		r.GET("/compendium/applink/identities", s.CheckIdentityHandler)
		r.POST("/compendium/applink/connect", s.CheckConnectHandler)
		r.POST("/compendium/cmd/syncTech/:mode", s.CheckSyncTechHandler)
		r.GET("/compendium/cmd/corpdata", s.CheckCorpDataHandler)
		r.GET("/compendium/applink/refresh", s.CheckRefreshHandler)
		r.GET("/compendium/user/corporations", s.CheckUserCorporationsHandler)
		r.GET("/compendium/api/tech", s.api)
		r.GET("/links", s.links)
		r.GET("/health", HealthCheckHandler)
		r.Static("/compendium/avatars", "docker/web/img/avatars")
	}

	registerRoutes(router)

	group := router.Group("/compendium")
	registerRoutes(group)

	//err := router.RunTLS(":"+port, s.certFile, s.keyFile)
	err := router.Run(":" + port)
	if err != nil {
		s.log.ErrorErr(err)
		//os.Exit(1)
	}
}

type db interface {
	CorpMembersRead(guildid string) ([]models.CorpMember, error)
	GuildRolesRead(guildid string) ([]models.CorpRole, error)
	GuildRolesExistSubscribe(guildid, RoleName, userid string) bool
	ListUserGetUserIdAndGuildId(token string) (userid string, guildid string, err error)
	ListUserGetByMatch(ttoken string) string
	ListUserUpdateToken(tokenOld, tokenNew string) error
	UsersGetByUserId(userid string) (*models.User, error)
	//GuildGet(guildid string) (*models.Guild, error)
	TechGet(username, userid string) ([]byte, error)
	TechUpdate(username, userid string, tech []byte) error
	CodeGet(code string) (*models.Code, error)
	CodeAllGet() []models.Code
	CodeDelete(code string)
	UserCorporationsGet(identity *models.Identity) ([]models.Guild, error)
	CorpMemberRead(userid string) ([]models.CorpMember, error)
	DeleteOldClient(userid string)
	SearchOldData(i models.Identity) (m models.Moving)
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

// HealthCheckHandler проверяет здоровье сервиса
func HealthCheckHandler(c *gin.Context) {
	// Если все проверки пройдены успешно
	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "Service is healthy",
		"timestamp": time.Now().Unix(),
		"cors":      "enabled",
	})
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
