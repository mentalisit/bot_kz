package telegram

import (
	"fmt"
	"sync"

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
	usernameMap            map[string]int
	chatMembers            map[int64]map[int64]tgbotapi.User // chatID -> userID -> User
	mu                     sync.RWMutex
}

func NewTelegram(log *logger.Logger, token string, st *storage.Storage) *Telegram {
	botApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.ErrorErr(err)
		return nil
	}
	t := &Telegram{
		log:         log,
		t:           botApi,
		Storage:     st,
		api:         restapi.NewRecover(log),
		usernameMap: make(map[string]int),
		chatMembers: make(map[int64]map[int64]tgbotapi.User),
	}

	fmt.Println(t.t.Self.UserName)

	//t.loadConfig()
	go t.update()
	go t.DeleteMessageTimer()
	go func() {
		t.chatMembers = t.Storage.Db.ReadAllMembers()
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				t.Storage.Db.UpsertChatData(t.chatMembers)
			}
		}
	}()
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
