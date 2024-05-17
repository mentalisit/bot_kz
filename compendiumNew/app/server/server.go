package server

import (
	"compendium/config"
	"compendium/models"
	"compendium/storage"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mentalisit/logger"
)

type Server struct {
	log *logger.Logger
	db  db
	In  chan models.IncomingMessage
}

func NewServer(log *logger.Logger, cfg *config.ConfigBot, st *storage.Storage) *Server {
	s := &Server{
		log: log,
		db:  st.DB,
		In:  make(chan models.IncomingMessage, 10),
	}

	go s.RunServer(cfg.Port)
	go s.RunServerRestApi()
	return s
}

func (s *Server) RunServer(port string) {
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

	fmt.Println("Running port:" + port)
	err := router.RunTLS(":"+port, "cert/RSA-cert.pem", "cert/RSA-privkey.pem")
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
	UsersGetByUserId(userid string) (*models.User, error)
	GuildGet(guildid string) (*models.Guild, error)
	TechGet(username, userid, guildid string) ([]byte, error)
	TechUpdate(username, userid, guildid string, tech []byte) error
}
