package logic

import (
	ds "bridge/Discord"
	tg "bridge/Telegram"
	"bridge/config"
	"bridge/models"
	"bridge/storage"
	"bridge/storage/postgres"
	wa "bridge/whatsapp"
	"fmt"
	"runtime"
	"time"

	"github.com/mentalisit/logger"
)

type Bridge struct {
	log      *logger.Logger
	in       models.ToBridgeMessage
	messages []models.BridgeTempMemory
	configs  map[string]models.Bridge2Config
	discord  *ds.Client
	telegram *tg.Client
	storage  BridgeConfig
	db       *postgres.Db
	whatsapp *wa.Client
}

func NewBridge(log *logger.Logger, st *storage.Storage, cfg *config.ConfigBot) *Bridge {
	bridge := &Bridge{
		log:      log,
		configs:  make(map[string]models.Bridge2Config),
		discord:  ds.NewClient(log),
		telegram: tg.NewClient(log),
		whatsapp: wa.NewClient(log),
		storage:  st.DB,
		db:       st.DB,
	}
	bridge.LoadConfig()

	//	go bridge.ServerRun()

	return bridge
}

type BridgeConfig interface {
	DBReadBridgeConfig2() []models.Bridge2Config
	UpdateBridge2Chat(br models.Bridge2Config)
	InsertBridge2Chat(br models.Bridge2Config)
	DeleteBridge2Chat(br models.Bridge2Config)

	SaveBridgeMap(msgMap map[string]string) error
	GetMapByLinkedID(msg map[string]string) (map[string]string, error)
}

func (b *Bridge) PrintGoroutine() {
	goroutine := runtime.NumGoroutine()
	tm := time.Now()
	mdate := (tm.Format("2006-01-02"))
	mtime := (tm.Format("15:04"))
	text := fmt.Sprintf(" %s %s Горутин  %d\n", mdate, mtime, goroutine)
	if goroutine > 120 {
		b.log.Info(text)
		b.log.Panic(text)
	} else if goroutine > 50 && goroutine%10 == 0 {
		b.log.Info(text)
	}

	fmt.Println(text)
}

// Shutdown корректно завершает все модули bridge
func (b *Bridge) Shutdown() {
	b.log.Info("Bridge shutting down...")

	// Close gRPC connections
	if b.discord != nil {
		if err := b.discord.Close(); err != nil {
			b.log.ErrorErr(err)
		}
	}

	if b.telegram != nil {
		if err := b.telegram.Close(); err != nil {
			b.log.ErrorErr(err)
		}
	}

	if b.whatsapp != nil {
		if err := b.whatsapp.Close(); err != nil {
			b.log.ErrorErr(err)
		}
	}

	b.log.Info("Bridge shutdown complete")
}
