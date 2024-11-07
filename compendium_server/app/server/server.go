package server

import (
	"compendium_s/models"
	"compendium_s/storage"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mentalisit/logger"
	"runtime"
	"time"
)

type Server struct {
	log   *logger.Logger
	db    db
	roles *Roles
	//In  chan models.IncomingMessage
}

func NewServer(log *logger.Logger, st *storage.Storage) *Server {
	s := &Server{
		log:   log,
		db:    st.DB,
		roles: NewRoles(log),
	}

	go s.RunServer()
	return s
}

func (s *Server) RunServer() {
	port := "8443"
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

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

	fmt.Println("Running port:" + port)
	err := router.RunTLS(":"+port, "docker/cert/RSA-cert.pem", "docker/cert/RSA-privkey.pem")
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
