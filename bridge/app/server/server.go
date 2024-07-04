package server

import (
	ds "bridge/Discord"
	tg "bridge/Telegram"
	"bridge/models"
	"bridge/storage"
	"github.com/gin-gonic/gin"
	"github.com/mentalisit/logger"
	"net/http"
	"os"
)

type Bridge struct {
	log      *logger.Logger
	in       models.ToBridgeMessage
	messages []models.BridgeTempMemory
	configs  map[string]models.BridgeConfig
	discord  *ds.Discord
	telegram *tg.Telegram
	storage  BridgeConfig
}

func NewBridge(log *logger.Logger, st *storage.Storage) *Bridge {
	bridge := &Bridge{
		log:      log,
		configs:  make(map[string]models.BridgeConfig),
		discord:  ds.NewDiscord(log),
		telegram: tg.NewTelegram(log),
		storage:  st.DB,
	}
	bridge.LoadConfig()
	go bridge.ServerRun()

	return bridge
}

func (b *Bridge) ServerRun() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// Обработчик для принятия сообщений от DiscordService
	router.POST("/bridge/inbox", b.indoxBridge)
	//
	router.GET("/bridge/config", b.configBridge)

	err := router.Run(":80")
	if err != nil {
		b.log.ErrorErr(err)
		os.Exit(1)
	}
}
func (b *Bridge) indoxBridge(c *gin.Context) {
	var mes models.ToBridgeMessage

	if err := c.ShouldBindJSON(&mes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message received successfully"})

	go b.logic(mes)
}

func (b *Bridge) configBridge(c *gin.Context) {
	config := b.storage.DBReadBridgeConfig()
	c.JSON(http.StatusOK, config)
}

type BridgeConfig interface {
	DBReadBridgeConfig() []models.BridgeConfig
	UpdateBridgeChat(br models.BridgeConfig)
	InsertBridgeChat(br models.BridgeConfig)
	DeleteBridgeChat(br models.BridgeConfig)
}
