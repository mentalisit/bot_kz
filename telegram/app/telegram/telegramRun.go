package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mentalisit/logger"
	"telegram/models"
	"telegram/storage"
	"telegram/telegram/restapi"
)

type Telegram struct {
	t            *tgbotapi.BotAPI
	log          *logger.Logger
	bridgeConfig map[string]models.BridgeConfig
	corpConfigRS map[string]models.CorporationConfig
	Storage      *storage.Storage
	api          *restapi.Recover
	usernameMap  map[string]int
}

func NewTelegram(log *logger.Logger, token string, st *storage.Storage) *Telegram {
	botApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.ErrorErr(err)
		return nil
	}
	t := &Telegram{
		log:          log,
		t:            botApi,
		bridgeConfig: make(map[string]models.BridgeConfig),
		corpConfigRS: make(map[string]models.CorporationConfig),
		Storage:      st,
		api:          restapi.NewRecover(log),
		usernameMap:  make(map[string]int),
	}

	fmt.Println(t.t.Self.UserName)

	t.loadConfig()
	go t.update()

	return t
}

func (t *Telegram) loadConfig() {
	bc := restapi.GetBridgeConfig()
	if len(bc) > 0 {
		t.bridgeConfig = bc
	}

	rs := t.Storage.Db.ReadConfigRs()
	if len(rs) > 0 {
		fmt.Printf("rsLoad %d\n", len(rs))
		for _, configRs := range rs {
			t.corpConfigRS[configRs.CorpName] = configRs
		}
	}

}
