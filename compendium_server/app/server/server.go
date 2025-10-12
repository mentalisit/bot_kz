package server

import (
	"compendium_s/models"
	"compendium_s/server/getCountry"
	"compendium_s/storage"
	"compendium_s/storage/postgres/multi"
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
	db         db
	multi      *multi.Db
	roles      *Roles
	cache      *getCountry.Cache
	certFile   string
	keyFile    string
	cacheReq   map[string]cacheEntry
	cacheMutex sync.Mutex
}

type cacheEntry struct {
	data      any       // Сами данные, которые отправляем
	timestamp time.Time // Когда сохранили
}

func NewServer(log *logger.Logger, st *storage.Storage) *Server {
	s := &Server{
		log:      log,
		db:       st.DB,
		multi:    st.DB.Multi,
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
	port := "8443"
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.LoggerWithFormatter(s.CustomLogFormatter))

	router.OPTIONS("/compendium/applink/identities", s.Check)
	router.GET("/compendium/applink/identities", s.CheckIdentityHandler)

	router.OPTIONS("/compendium/applink/connect", s.Check)
	router.POST("/compendium/applink/connect", s.CheckConnectHandler)

	router.OPTIONS("/compendium/cmd/syncTech/:mode", s.Check)
	router.POST("/compendium/cmd/syncTech/:mode", s.CheckSyncTechHandler)

	router.OPTIONS("/compendium/cmd/corpdata", s.Check)
	router.GET("/compendium/cmd/corpdata", s.CheckCorpDataHandler)

	router.OPTIONS("/compendium/applink/refresh", s.Check)
	router.GET("/compendium/applink/refresh", s.CheckRefreshHandler)

	router.GET("/links", s.links)

	router.GET("/compendium/api/tech", s.api)

	router.Static("/compendium/avatars", "docker/compendium/avatars")
	router.Static("/docker/compendium/avatars", "docker/compendium/avatars")
	router.Static("/tv", "docker/compendium/tv")

	router.GET("/health", HealthCheckHandler)

	fmt.Println("Running port:" + port)
	err := router.RunTLS(":"+port, s.certFile, s.keyFile)
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
	GuildGet(guildid string) (*models.Guild, error)
	TechGet(username, userid, guildid string) ([]byte, error)
	TechUpdate(username, userid, guildid string, tech []byte) error
	CodeGet(code string) (*models.Code, error)
	CodeAllGet() []models.Code
	CodeDelete(code string)
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
func (s *Server) CustomLogFormatter(param gin.LogFormatterParams) string {
	if param.Method == "OPTIONS" {
		return fmt.Sprintf("")
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

	return fmt.Sprintf("[GIN]%v|%s %3d %s|%13v|%15s|%s%-3s%s|%#v\n%s",
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
		"message": "Service is healthy",
	})
}
