package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mentalisit/logger"
	"telegram/models"
	"telegram/telegram/restapi"
)

type Telegram struct {
	t            *tgbotapi.BotAPI
	log          *logger.Logger
	bridgeConfig map[string]models.BridgeConfig
	corpConfigRS map[string]models.CorporationConfig
}

func NewTelegram(log *logger.Logger, token string) *Telegram {
	botApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.ErrorErr(err)
		return nil
	}
	t := &Telegram{
		log: log,
		t:   botApi,
	}

	go t.update()

	return t
}

func (t *Telegram) LoadConfig() {
	bc, _ := restapi.GetBridgeConfig()
	for _, configBridge := range bc {
		t.bridgeConfig[configBridge.NameRelay] = configBridge
	}
	rs, _ := restapi.GetRsBotConfig()
}
