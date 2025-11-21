package telegram

import (
	"fmt"
	"net/http"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"

	"strconv"
	"telegram/models"
	"telegram/storage"
	"telegram/telegram/restapi"
	"time"

	"github.com/mentalisit/logger"
)

type Telegram struct {
	t                      *tgbotapi.BotAPI
	log                    *logger.Logger
	bridgeConfig           []models.Bridge2Config
	bridgeConfigUpdateTime int64
	Storage                *storage.Storage
	api                    *restapi.Recover
	//usernameMap            map[string]int
	server *http.Server
}

func NewTelegram(log *logger.Logger, token string, st *storage.Storage) *Telegram {
	botApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.ErrorErr(err)
		return nil
	}

	t := &Telegram{
		log:     log,
		t:       botApi,
		Storage: st,
		api:     restapi.NewRecover(log),
		//usernameMap: make(map[string]int),
	}

	fmt.Printf("Authorized on account %s\n", t.t.Self.UserName)

	go t.StartWebApp("8080")

	//t.loadConfig()
	go t.update()
	go t.DeleteMessageTimer()
	return t
}

func (t *Telegram) Close() {
	t.t.StopReceivingUpdates()
	t.api.Close()
}
func (t *Telegram) DeleteMessageTimer() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mes := t.Storage.Db.TimerReadMessage("tg")
			if len(mes) > 0 {
				for _, m := range mes {
					if m.MesId != "" {
						mid, _ := strconv.Atoi(m.MesId)
						_ = t.DelMessage(m.ChatId, mid)
						t.Storage.Db.TimerDeleteMessage(m)
					}
				}
			}
		}
	}
}
