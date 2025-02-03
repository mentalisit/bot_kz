package bot

import (
	"fmt"
	"github.com/mentalisit/logger"
	"rs/bot/helpers"
	"rs/bot/otherQueue"
	"rs/clients"
	"rs/models"
	"rs/pkg/utils"
	"rs/storage"
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
	storage    *storage.Storage
	client     *clients.Clients
	Inbox      chan models.InMessage
	log        *logger.Logger
	wg         sync.WaitGroup
	mu         sync.Mutex
	helpers    *helpers.Helpers
	otherQueue *otherQueue.OtherQ
}

func NewBot(storage *storage.Storage, client *clients.Clients, log *logger.Logger) *Bot {
	b := &Bot{
		storage:    storage,
		client:     client,
		log:        log,
		Inbox:      make(chan models.InMessage, 30),
		helpers:    helpers.NewHelpers(log, storage),
		otherQueue: otherQueue.NewOtherQ(log),
	}

	go b.loadInbox()
	go b.timerBot()

	return b
}

func (b *Bot) loadInbox() {
	b.log.Info("Бот загружен и готов к работе ")

	for {
		//ПОЛУЧЕНИЕ СООБЩЕНИЙ
		select {
		case in := <-b.Inbox:
			b.PrepareLogicRs(in)

		}
	}
}
func (b *Bot) timerBot() { //цикл для удаления сообщений
	for {
		now := time.Now()
		if now.Second() == 0 {
			b.MinusMin()

			if now.Minute() == 0 {
				go b.AutoHelp() //автозапуск справки
			}
		}
		time.Sleep(1 * time.Second)
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
	case <-time.After(15 * time.Second):
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
		utils.PrintGoroutine(b.log)
		if strings.HasPrefix(in.Mtext, "!") {
			answer := b.helpers.GeminiSay(in.Mtext, in.Username)
			for _, text := range answer {
				fmt.Printf("answer gemini %s\n", text)
				if text == "" {
					return
				}
				if in.Config.DsChannel != "" {
					go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 600)
					go b.client.Ds.DeleteMessageSecond(in.Config.DsChannel, in.Ds.Mesid, 600)
				} else if in.Config.TgChannel != "" {
					go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 600)
					go b.client.Tg.DelMessageSecond(in.Config.TgChannel, strconv.Itoa(in.Tg.Mesid), 600)
				}
			}
		} else {
			text := "эээ я же бот че ты мне пишешь тут, пиши в канале "
			if in.Config.DsChannel != "" {
				go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 600)
				go b.client.Ds.DeleteMessageSecond(in.Config.DsChannel, in.Ds.Mesid, 600)
			} else if in.Config.TgChannel != "" {
				go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 600)
				go b.client.Tg.DelMessageSecond(in.Config.TgChannel, strconv.Itoa(in.Tg.Mesid), 600)
			}

			if conf.DsChannel != "" {
				go b.client.Ds.SendWebhook(in.Mtext, in.Username, conf.DsChannel, in.Ds.Avatar)
			}
			if conf.TgChannel != "" {
				go b.client.Tg.SendChannel(conf.TgChannel, fmt.Sprintf("%s: %s", in.Username, in.Mtext))
			}
		}

		return
	}

	if in.Config.CorpName == "Корпорация  \"РУСЬ\".сбор-на-кз" {
		b.insertUserAccount(in)
	}

	ch := utils.WaitForMessage("LogicRs ")
	if len(in.Mtext) > 0 && in.Mtext != " `edit`" {
		utils.PrintGoroutine(b.log)
		fmt.Printf("%+v LogicRs %s %s %s\n", time.Now().Format(time.DateTime), in.Config.CorpName, in.Username, in.Mtext)
		if b.lRsPlus(in) {
		} else if b.lSubs(in) {
		} else if b.lQueue(in) {
		} else if b.lRsStart(in) {
		} else if b.lEvent(in) {
		} else if b.lTop(in) {
		} else if b.lEmoji(in) {
		} else if b.logicIfText(in) {
		} else if b.bridge(in) {
		} else {
			b.cleanChat(in)
		}
	} else if in.Opt.Contains(models.OptionMinusMinNext) {
		b.MinusMinMessageUpdate()
	} else if in.Opt.Contains(models.OptionMessageUpdateDS) {
		b.QueueLevel(in)
	} else if in.Opt.Contains(models.OptionMessageUpdateTG) {
		b.QueueLevel(in)
	} else if in.Opt.Contains(models.OptionUpdateAutoHelp) {
		b.QueueAll(in)
	} else if in.Opt.Contains(models.OptionQueueAll) {
		b.QueueLevel(in)
	} else if in.Opt.Contains(models.OptionPlus) {
		b.QueueLevel(in)
	}
	close(ch)
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
	case "OptimizationSborkz":
		go b.OptimizationSborkz()
		b.iftipdelete(in)
	case "cleanrs":
		go b.client.Ds.CleanRsBotOtherMessage()
	default:
		iftext = false
	}
	return iftext
}

func (b *Bot) cleanChat(in models.InMessage) {
	ch := utils.WaitForMessage("cleanChat")
	defer close(ch)
	if !strings.HasPrefix(in.Mtext, ".") && !in.Opt.Contains(models.OptionEdit) {
		if in.Tip == ds && in.Config.DelMesComplite == 0 {
			go b.client.Ds.DeleteMessageSecond(in.Config.DsChannel, in.Ds.Mesid, 600)
		}
		if in.Tip == tg && IsThisTopicTG(in.Config.TgChannel) {
			go b.client.Tg.DelMessageSecond(in.Config.TgChannel, strconv.Itoa(in.Tg.Mesid), 600)
		}
	}
}

func (b *Bot) bridge(in models.InMessage) bool {
	if in.Config.Forward {
		if in.Tip == ds {
			mText := b.client.Ds.ReplaceTextMessage(in.Mtext, in.Config.Guildid)
			text := fmt.Sprintf("(DS)%s \n%s", in.Username, mText)
			go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 600)
			go b.cleanChat(in)
		} else if in.Tip == tg {
			text := fmt.Sprintf("(TG)%s \n%s", in.Username, in.Mtext)
			go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 600)
			go b.cleanChat(in)
		}
	}
	return in.Config.Forward
}
