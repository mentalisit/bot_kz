package bot

import (
	"fmt"
	"github.com/mentalisit/logger"
	"kz_bot/clients"
	"kz_bot/config"
	"kz_bot/models"
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
	storage    *storage.Storage
	client     *clients.Clients
	inbox      chan models.InMessage
	log        *logger.Logger
	debug      bool
	in         models.InMessage
	wg         sync.WaitGroup
	mu         sync.Mutex
	configCorp map[string]models.CorporationConfig
}

func NewBot(storage *storage.Storage, client *clients.Clients, log *logger.Logger, cfg *config.ConfigBot) *Bot {
	b := &Bot{
		storage:    storage,
		client:     client,
		log:        log,
		debug:      cfg.IsDebug,
		inbox:      make(chan models.InMessage, 10),
		configCorp: storage.CorpConfigRS,
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
		case in := <-b.client.Ds.ChanRsMessage:
			b.in = in
			b.LogicRs()
		case in := <-b.client.Tg.ChanRsMessage:
			b.in = in
			b.LogicRs()

		case in := <-b.inbox:
			b.in = in
			b.LogicRs()
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

// LogicRs логика игры
func (b *Bot) LogicRs() {
	if strings.HasPrefix(b.in.Mtext, ".") {
		b.accessChat()
		return
	}
	if len(b.in.Mtext) > 0 && b.in.Mtext != " `edit`" {
		if b.lRsPlus() {
		} else if b.lDarkRsPlus() {
		} else if b.lSubs() {
		} else if b.lDarkSubs() {
		} else if b.lQueue() {
		} else if b.lRsStart() {
		} else if b.lEvent() {
		} else if b.lTop() {
		} else if b.lEmoji() {
		} else if b.logicIfText() {
		} else if b.bridge() {
			//} else if b.lIfCommand() {
			//} else if b.SendALLChannel() {
		} else {
			b.cleanChat()
			//go b.Transtale()//нужно решить проблему с ошибками
		}

	} else if b.in.Option.MinusMin {
		b.CheckTimeQueue()
	} else if b.in.Option.Update {
		b.QueueLevel()
	}
}

func (b *Bot) cleanChat() {
	if b.in.Tip == ds && b.in.Config.DelMesComplite == 0 && !b.in.Option.Edit {
		go b.client.Ds.CleanChat(b.in.Config.DsChannel, b.in.Ds.Mesid, b.in.Mtext)
	}
	// if hs ua
	if b.in.Tip == tg && b.in.Config.TgChannel == "-1002116077159/44" {
		if !strings.HasPrefix(b.in.Mtext, ".") {
			go b.client.Tg.DelMessageSecond("-1002116077159/44", strconv.Itoa(b.in.Tg.Mesid), 600)
		}

	}
}

func (b *Bot) logicIfText() bool {
	iftext := true
	switch b.in.Mtext {
	case "+":
		if b.Plus() {
			return true
		}
	case "-":
		if b.Minus() {
			return true
		}
	case "Справка", "Help", "help":
		b.hhelp()
	case "update modules", "обновить модули":
		b.updateCompendiumModules()
		iftext = true
	case "OptimizationSborkz":
		go b.storage.DbFunc.OptimizationSborkz()
		b.iftipdelete()
	case "cleanrs":
		go b.client.Ds.CleanRsBotOtherMessage()
	default:
		iftext = false
	}
	return iftext
}

func (b *Bot) bridge() bool {
	if b.in.Config.Forward {
		//go b.Transtale()
		if b.in.Tip == ds {
			text := fmt.Sprintf("(DS)%s \n%s", b.in.Name, b.in.Mtext)
			b.client.Tg.SendChannelDelSecond(b.in.Config.TgChannel, text, 180)
			b.cleanChat()
		} else if b.in.Tip == tg {
			text := fmt.Sprintf("(TG)%s \n%s", b.in.Name, b.in.Mtext)
			b.client.Ds.SendChannelDelSecond(b.in.Config.DsChannel, text, 180)
			b.cleanChat()
		}
	}
	return b.in.Config.Forward
}
func (b *Bot) Autohelp() {
	tm := time.Now()
	mtime := tm.Format("15:04")
	EvenOrOdd, _ := strconv.Atoi((tm.Format("2006-01-02"))[8:])
	if mtime == "12:00" {
		a := b.storage.ConfigRs.AutoHelp()
		for _, s := range a {
			if s.DsChannel != "" {
				s.MesidDsHelp = b.client.Ds.HelpChannelUpdate(s)
			}
			if s.Forward && s.TgChannel != "" && EvenOrOdd%2 == 0 {
				text := fmt.Sprintf("%s \n%s", b.getLanguageText(s.Country, "info_bot_delete_msg"), b.getLanguageText(s.Country, "info_help_text"))
				if s.MesidTgHelp != "" {
					b.log.Info(s.MesidTgHelp)
					mID, err := strconv.Atoi(s.MesidTgHelp)
					if err != nil {
						return
					}
					b.log.Info(fmt.Sprintf("%s %d", s.MesidTgHelp, mID))
					go b.client.Tg.DelMessage(s.TgChannel, mID)
				}
				mid := b.client.Tg.SendHelp(s.TgChannel, strings.ReplaceAll(text, "3", "10"))
				b.log.Info(fmt.Sprintf("mid %d", mid))
				s.MesidTgHelp = strconv.Itoa(mid)

			}
			b.storage.ConfigRs.AutoHelpUpdateMesid(s)
		}
		time.Sleep(time.Minute)
		go b.client.Ds.CleanRsBotOtherMessage()
	} else if tm.Minute() == 0 {
		a := b.storage.ConfigRs.AutoHelp()
		for _, s := range a {
			if s.DsChannel != "" {
				MesidDsHelp := b.client.Ds.HelpChannelUpdate(s)
				if MesidDsHelp != s.MesidDsHelp {
					s.MesidDsHelp = MesidDsHelp
					b.storage.ConfigRs.AutoHelpUpdateMesid(s)
				}
			}
		}
	}
	time.Sleep(time.Minute)
}
