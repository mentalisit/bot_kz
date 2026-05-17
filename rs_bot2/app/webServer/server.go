package webServer

import (
	"rs/bot2"
	"rs/clients"
	"rs/models"
	"rs/storage"
	"rs/storage/postgresV2"
	"rs/webServer/getCountry"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mentalisit/logger"
)

type cacheEntry struct {
	data      any       // Сами данные, которые отправляем
	timestamp time.Time // Когда сохранили
}

type Server struct {
	log        *logger.Logger
	db         *postgresV2.Db
	cl         *clients.Clients
	bot        *bot2.Bot
	roles      *RolesHelper
	cache      *getCountry.Cache
	hub        *Hub
	certFile   string
	keyFile    string
	cacheReq   map[string]cacheEntry
	cacheMutex sync.Mutex
	LinkCodes  sync.Map // Карта для временных кодов вызова /link
	GameNames  map[string]models.GameAccountData
}

type LinkCodeData struct {
	UserID   string
	Username string
	Provider string
	Expires  time.Time
}

func NewServer(log *logger.Logger, st *storage.Storage, cl *clients.Clients, bot *bot2.Bot) *Server {

	s := &Server{
		log:       log,
		db:        st.V2,
		cl:        cl,
		bot:       bot,
		roles:     NewRolesHelper(log, cl.Ds),
		cache:     getCountry.NewCache(),
		hub:       NewHub(st.V2),
		certFile:  "docker/cert/RSA-cert.pem",
		keyFile:   "docker/cert/RSA-privkey.pem",
		cacheReq:  make(map[string]cacheEntry),
		GameNames: make(map[string]models.GameAccountData),
	}

	s.hub.onBroadcast = s.sendPushNotification

	go s.hub.Run()
	go s.RunServer()
	return s
}

func (s *Server) RunServer() {
	port := "8443"
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(CORSMiddleware())
	router.Use(gin.LoggerWithFormatter(s.CustomLogFormatter))

	registerRoutes := func(r gin.IRouter) {
		r.GET("/api/config", s.getMobileSetting)
		r.GET("/api/webhook", s.getWebhook)
		r.GET("/api/gameid", s.getGameData)
		r.POST("/api/config", s.postMobileSetting)
		r.DELETE("/api/config", s.deleteMobileSetting)
		r.Static("/api/setup", "/docker/rs_bot")

		r.Static("/settings", "/docker/settings")

		r.POST("/api/ds/link", s.postSettingDsLink)
		r.POST("/api/tg/link", s.postSettingTgLink)

		r.DELETE("/api/link", s.deleteSettingLink)

		r.GET("/health", HealthCheckHandler)

		r.GET("/api/events/get/seasons", s.getEvents)
		r.GET("/api/events/get/season", s.getEventId)
		r.GET("/api/events/get/user", s.getEventUser)

		r.GET("/api/ma/user", s.getMAcc)
		r.POST("/api/ma/user", s.postMAcc)
		r.POST("/api/ma/send", s.postMAccSendMessage)
		r.GET("/api/ma/get", s.getSecret)

		r.GET("/api/ga/check", s.checkGA)
		r.POST("/api/ga/user", s.postMAcc)

		r.GET("/api/guild/getData", s.GetGuildData)
		r.POST("/api/ws/universal", s.WsUniversal)

		//r.GET("/api/poll/getAll")
		//r.POST("/api/poll/create")
		//r.POST("/api/poll/votes")

		//compendium
		r.GET("/compendium/tech", s.compendiumTech)
		r.POST("/compendium/tech", s.compendiumTech)
		r.GET("/compendium/corps", s.compendiumCorps)
		r.GET("/compendium/accessible_corporations", s.compendiumAccessibleCorporations)
		r.GET("/compendium/multicorp", s.compendiumMultiCorp)
		r.POST("/compendium/multicorp_register", s.compendiumMultiCorpRegister)
		r.POST("/compendium/toggle_participation", s.compendiumToggleParticipation)
		r.GET("/compendium/corpdata", s.compendiumCorpData)

		r.POST("/compendium/study", s.postStudy)
		r.GET("/compendium/study", s.getStudy)

		// charts / графики
		r.GET("/api/chart/corps", s.getChartCorps)           // Список корп
		r.GET("/api/chart/levels", s.getChartLevels)         // Список уровней корпы
		r.GET("/api/chart/corp", s.getChartCorp)             // Данные корпы
		r.GET("/api/chart/user", s.getChartUser)             // Данные юзера
		r.GET("/api/chart/popularity", s.getChartPopularity) // Популярность по часам

		// roles logic
		r.GET("/api/user/chats", s.GetUserChats)

		// Chat roles
		r.GET("/api/chat/:chatId/roles", s.GetChatRoles)
		r.POST("/api/chat/:chatId/roles", s.CreateRole)
		r.PUT("/api/chat/:chatId/roles/:roleId", s.UpdateRole)
		r.DELETE("/api/chat/:chatId/roles/:roleId", s.DeleteRole)
		r.POST("/api/chat/:chatId/roles/:roleId/join", s.JoinRole)
		r.POST("/api/chat/:chatId/roles/:roleId/leave", s.LeaveRole)
		r.GET("/api/chat/:chatId/roles/:roleId/members", s.GetRoleMembers)

		// Chat users
		r.GET("/api/chat/:chatId/users", s.GetChatUsers)
		r.POST("/api/chat/:chatId/users/:userId/roles/:roleId", s.SetUserRole)
		r.DELETE("/api/chat/:chatId/users/:userId/roles/:roleId", s.SetUserRole)

		// Chat permissions
		r.GET("/api/chat/:chatId/permissions", s.GetUserPermissions)

		// Corporation members
		r.GET("/api/chat/:chatId/corp-members", s.GetCorpMembers)
		r.DELETE("/api/chat/:chatId/corp-members/:userId", s.RemoveCorpMember)

		// Push Notifications
		r.POST("/api/push/subscribe", s.postPushSubscribe)

		// Chat Channels
		r.GET("/api/chat/channels", s.getChatChannels)
		r.POST("/api/chat/channels", s.postChatChannel)
		r.GET("/api/chat/settings", s.getChatUserSettings)
		r.POST("/api/chat/settings", s.postChatUserSettings)

		// WebSocket Chat
		r.GET("/ws", s.ServeWs)
	}

	registerRoutes(router)

	group := router.Group("/rs")
	registerRoutes(group)

	//err := router.RunTLS(":"+port, s.certFile, s.keyFile)
	err := router.Run(":" + port)
	if err != nil {
		s.log.ErrorErr(err)
		//os.Exit(1)
	}
}
