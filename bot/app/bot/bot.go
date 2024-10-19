package bot

import (
	"fmt"
	"github.com/mentalisit/logger"
	"kz_bot/bot/helpers"
	"kz_bot/bot/otherQueue"
	"kz_bot/clients"
	"kz_bot/config"
	"kz_bot/models"
	"kz_bot/pkg/utils"
	"kz_bot/storage"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	ds = "ds"
	tg = "tg"
)

type Bot struct {
	storage *storage.Storage
	client  *clients.Clients
	inbox   chan models.InMessage
	log     *logger.Logger
	debug   bool
	//in         models.InMessage
	wg         sync.WaitGroup
	mu         sync.Mutex
	configCorp map[string]models.CorporationConfig
	helpers    *helpers.Helpers
	otherQueue *otherQueue.OtherQ
}

func NewBot(storage *storage.Storage, client *clients.Clients, log *logger.Logger, cfg *config.ConfigBot) *Bot {
	b := &Bot{
		storage:    storage,
		client:     client,
		log:        log,
		debug:      cfg.IsDebug,
		inbox:      make(chan models.InMessage, 30),
		configCorp: storage.CorpConfigRS,
		helpers:    helpers.NewHelpers(log, storage),
		otherQueue: otherQueue.NewOtherQ(log),
	}
	go func() {
		for {
			if len(storage.CorpConfigRS) > 0 {
				b.configCorp = storage.CorpConfigRS
				fmt.Println("Bot Loaded")
				break
			} else {
				time.Sleep(1 * time.Second)
			}
		}
	}()

	go b.loadInbox()
	go b.RemoveMessage()

	return b
}

func (b *Bot) loadInbox() {
	b.log.Info("Бот загружен и готов к работе ")

	for {
		//ПОЛУЧЕНИЕ СООБЩЕНИЙ
		select {
		case in := <-b.client.DS.ChanRsMessage:
			b.PrepareLogicRs(in)
			if len(b.client.DS.ChanRsMessage) > 5 {
				b.log.Info(fmt.Sprintf("len(b.client.Ds.ChanRsMessage) = %d", len(b.client.DS.ChanRsMessage)))
			}
		case in := <-b.client.Tg.ChanRsMessage:
			b.PrepareLogicRs(in)
			if len(b.client.Tg.ChanRsMessage) > 5 {
				b.log.Info(fmt.Sprintf("len(b.client.Tg.ChanRsMessage) = %d", len(b.client.Tg.ChanRsMessage)))
			}
		case in := <-b.inbox:
			b.PrepareLogicRs(in)
			if len(b.inbox) > 15 {
				b.log.Info(fmt.Sprintf("len(b.inbox) = %d\n %+v\n", len(b.inbox), in))
			}
		}
	}
}
func (b *Bot) RemoveMessage() { //цикл для удаления сообщений
	for {
		now := time.Now()
		if now.Second() == 0 {
			b.MinusMin() //ежеминутное обновление активной очереди
			b.client.DeleteMessageTimer()

			if now.Minute() == 0 {
				b.Autohelp() //автозапуск справки
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func (b *Bot) PrepareLogicRs(in models.InMessage) {
	// Канал для отслеживания завершения запроса
	done := make(chan struct{})

	go func() {
		ch := utils.WaitForMessage("PrepareLogicRs")
		b.LogicRs(in)
		close(done)
		close(ch)
	}()

	select {
	case <-done:
		// Запрос завершился до истечения таймаута
	case <-time.After(10 * time.Second):
		// Логируем, если запрос завис
		b.log.InfoStruct("PrepareLogicRs", in)
	}

}

// LogicRs логика игры
func (b *Bot) LogicRs(in models.InMessage) {
	if strings.HasPrefix(in.Mtext, ".") {
		b.accessChat(in)
		return
	}

	dm, conf := b.helpers.IfMessageDM(in)
	if dm {
		text := "эээ я же бот че ты мне пишешь тут, пиши в канале "
		if in.Config.DsChannel != "" {
			b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 600)
			b.client.Ds.DeleteMesageSecond(in.Config.DsChannel, in.Ds.Mesid, 600)
		} else if in.Config.TgChannel != "" {
			b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 600)
			b.client.Tg.DelMessageSecond(in.Config.TgChannel, strconv.Itoa(in.Tg.Mesid), 600)
		}

		if conf.DsChannel != "" {
			b.client.Ds.SendWebhook(in.Mtext, in.Username, conf.DsChannel, in.Ds.Avatar)
		}
		if conf.TgChannel != "" {
			b.client.Tg.SendChannel(conf.TgChannel, fmt.Sprintf("%s: %s", in.Username, in.Mtext))
		}

		return
	}

	ch := utils.WaitForMessage("LogicRs ")
	if len(in.Mtext) > 0 && in.Mtext != " `edit`" {
		utils.PrintGoroutine(b.log)
		fmt.Printf("%+v LogicRs %s %s %s\n", time.Now().Format(time.DateTime), in.Config.CorpName, in.Username, in.Mtext)
		if b.lRsPlus(in) {
		} else if b.lDarkRsPlus(in) {
		} else if b.lSubs(in) {
		} else if b.lDarkSubs(in) {
		} else if b.lQueue(in) {
		} else if b.lRsStart(in) {
		} else if b.lEvent(in) {
		} else if b.lTop(in) {
		} else if b.lEmoji(in) {
		} else if b.logicIfText(in) {
		} else if b.bridge(in) {
			//} else if b.lIfCommand() {
		} else {
			b.cleanChat(in)
			//go b.Transtale()//нужно решить проблему с ошибками
		}

	} else if in.Option.MinusMin {
		b.CheckTimeQueue(in)
	} else if in.Option.Update {
		b.QueueLevel(in)
	}
	close(ch)
}

func (b *Bot) cleanChat(in models.InMessage) {
	ch := utils.WaitForMessage("cleanChat")
	defer close(ch)
	if in.Tip == ds && in.Config.DelMesComplite == 0 && !in.Option.Edit {
		go b.client.Ds.CleanChat(in.Config.DsChannel, in.Ds.Mesid, in.Mtext)
	}
	// if hs ua
	if in.Tip == tg && in.Config.TgChannel == "-1002116077159/44" {
		if !strings.HasPrefix(in.Mtext, ".") {
			go b.client.Tg.DelMessageSecond("-1002116077159/44", strconv.Itoa(in.Tg.Mesid), 600)
		}

	}
}

func (b *Bot) logicIfText(in models.InMessage) bool {
	iftext := true
	switch in.Mtext {
	case "+":
		if b.Plus(in) {
			return true
		}
	case "-":
		if b.Minus(in) {
			return true
		}
	case "Справка", "Help", "help":
		b.hhelp(in)
	case "update modules", "обновить модули":
		go b.updateCompendiumModules(in)
		iftext = true
	case "OptimizationSborkz":
		go b.storage.DbFunc.OptimizationSborkz()
		b.iftipdelete(in)
	case "cleanrs":
		go b.client.Ds.CleanRsBotOtherMessage()
	default:
		iftext = false
	}
	return iftext
}

func (b *Bot) bridge(in models.InMessage) bool {
	if in.Config.Forward {
		if in.Tip == ds {
			text := fmt.Sprintf("(DS)%s \n%s", in.Username, in.Mtext)
			go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 180)
			go b.cleanChat(in)
		} else if in.Tip == tg {
			text := fmt.Sprintf("(TG)%s \n%s", in.Username, in.Mtext)
			go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 180)
			go b.cleanChat(in)
		}
	}
	return in.Config.Forward
}
