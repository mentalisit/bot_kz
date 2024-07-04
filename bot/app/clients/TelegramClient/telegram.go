package TelegramClient

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mentalisit/logger"
	"kz_bot/config"
	"kz_bot/models"
	"kz_bot/pkg/clientTelegram"
	"kz_bot/storage"
	"strconv"
)

type Telegram struct {
	ChanRsMessage chan models.InMessage
	t             *tgbotapi.BotAPI
	log           *logger.Logger
	storage       *storage.Storage
	debug         bool
	bridgeConfig  *map[string]models.BridgeConfig
	corpConfigRS  map[string]models.CorporationConfig
}

func NewTelegram(log *logger.Logger, st *storage.Storage, cfg *config.ConfigBot) *Telegram {
	client, err := clientTelegram.NewTelegram(log, cfg)
	if err != nil {
		return nil
	}

	tg := &Telegram{
		ChanRsMessage: make(chan models.InMessage, 10),
		t:             client,
		log:           log,
		storage:       st,
		debug:         cfg.IsDebug,
		bridgeConfig:  &st.BridgeConfigs,
		corpConfigRS:  st.CorpConfigRS,
	}

	go tg.update()

	return tg
}
func (t *Telegram) update() {
	ut := tgbotapi.NewUpdate(0)
	ut.Timeout = 60
	//получаем обновления от телеграм
	updates := t.t.GetUpdatesChan(ut)
	for update := range updates {
		//if update.InlineQuery != nil {
		//	t.handleInlineQuery(update.InlineQuery)
		//} else if update.ChosenInlineResult != nil {
		//go t.handleChosenInlineResult(update.ChosenInlineResult)
		//} else
		if update.CallbackQuery != nil {
			t.callback(update.CallbackQuery) //нажатия в чате
		} else if update.Message != nil {

			if update.Message.Chat.IsPrivate() { //если пишут боту в личку
				t.ifPrivatMesage(update.Message)
			} else if update.Message.IsCommand() {
				t.updatesComand(update.Message) //если сообщение является командой
			} else { //остальные сообщения
				t.logicMix(update.Message, false)
			}
		} else if update.EditedMessage != nil {
			t.logicMix(update.EditedMessage, true)
		} else if update.MyChatMember != nil {
			t.myChatMember(update.MyChatMember)

		} else if update.ChatMember != nil {
			t.chatMember(update.ChatMember)
		} else if update.ChatJoinRequest != nil {
			t.log.InfoStruct("ChatJoinRequest", update.ChatJoinRequest)
		} else {
			go func() {
				if update.Poll != nil {
					t.log.InfoStruct("pool ", update.Poll)
				} else if update.EditedChannelPost != nil {
					//t.log.InfoStruct("EditedChannelPost", update.EditedChannelPost)
				} else if update.ChannelPost != nil {
					//t.log.InfoStruct("update.ChannelPost.SenderChat", update.ChannelPost.SenderChat)
					//t.log.InfoStruct("update.ChannelPost.Chat", update.ChannelPost.Chat)
				} else {
					t.log.Info(fmt.Sprintf(" else update: %+v \n", update))
				}
			}()
		}
	}
}
func (t *Telegram) ifPrivatMesage(m *tgbotapi.Message) {
	if m.Text == "/start" {
		_, err := t.t.Send(tgbotapi.NewMessage(m.From.ID,
			"Возможность писать сообщения в личку активирована /"+
				" The ability to write private messages is activated"))
		if err != nil {
			t.log.ErrorErr(err)
			return
		}
	} else {
		//нужно решить что тут делать
		text := "эээ я же бот че ты мне пишешь тут, пиши в канале "
		ThreadID := m.MessageThreadID
		if !m.IsTopicMessage && ThreadID != 0 {
			ThreadID = 0
		}
		ChatId := strconv.FormatInt(m.Chat.ID, 10) + fmt.Sprintf("/%d", ThreadID)
		t.SendChannelDelSecond(ChatId, text, 600)
		t.DelMessageSecond(ChatId, strconv.Itoa(m.MessageID), 600)
		t.log.Info("DM " + m.From.String() + ": " + m.Text)
	}

}
func (t *Telegram) Shutdown() {
	t.t.StopReceivingUpdates()
}
