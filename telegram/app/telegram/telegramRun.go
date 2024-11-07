package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mentalisit/logger"
	"strconv"
	"telegram/models"
	"telegram/storage"
	"telegram/telegram/restapi"
	"time"
)

type Telegram struct {
	t                      *tgbotapi.BotAPI
	log                    *logger.Logger
	bridgeConfig           []models.BridgeConfig
	bridgeConfigUpdateTime int64
	Storage                *storage.Storage
	api                    *restapi.Recover
	usernameMap            map[string]int
}

func NewTelegram(log *logger.Logger, token string, st *storage.Storage) *Telegram {
	botApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.ErrorErr(err)
		return nil
	}
	t := &Telegram{
		log: log,
		t:   botApi,
		//bridgeConfig: make(map[string]models.BridgeConfig),
		//corpConfigRS: make(map[string]models.CorporationConfig),
		Storage:     st,
		api:         restapi.NewRecover(log),
		usernameMap: make(map[string]int),
	}

	fmt.Println(t.t.Self.UserName)

	//t.loadConfig()
	go t.update()
	go t.DeleteMessageTimer()
	return t
}

func (t *Telegram) Close() {
	t.api.Close()
}
func (t *Telegram) DeleteMessageTimer() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mes := t.Storage.Db.TimerReadMessage()
			if len(mes) > 0 {
				for _, m := range mes {
					if m.Tgmesid != "" {
						mid, _ := strconv.Atoi(m.Tgmesid)
						_ = t.DelMessage(m.Tgchatid, mid)
						t.Storage.Db.TimerDeleteMessage(m)
					}
				}
			}
		}
	}
}
