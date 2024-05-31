package bot

import (
	"fmt"
	"github.com/mentalisit/logger"
	"kz_bot/bot/helpers"
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
}

func NewBot(storage *storage.Storage, client *clients.Clients, log *logger.Logger, cfg *config.ConfigBot) *Bot {
	b := &Bot{
		storage:    storage,
		client:     client,
		log:        log,
		debug:      cfg.IsDebug,
		inbox:      make(chan models.InMessage, 10),
		configCorp: storage.CorpConfigRS,
		helpers:    helpers.NewHelpers(log, storage),
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
			b.LogicRs(in)
		case in := <-b.client.Tg.ChanRsMessage:
			b.LogicRs(in)
		case in := <-b.inbox:
			b.LogicRs(in)
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
func (b *Bot) LogicRs(in models.InMessage) {
	if strings.HasPrefix(in.Mtext, ".") {
		b.accessChat(in)
		return
	}
	if len(in.Mtext) > 0 && in.Mtext != " `edit`" {
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
			//} else if b.SendALLChannel() {
		} else {
			b.cleanChat(in)
			//go b.Transtale()//нужно решить проблему с ошибками
		}

	} else if in.Option.MinusMin {
		b.CheckTimeQueue(in)
	} else if in.Option.Update {
		b.QueueLevel(in)
	}
}

func (b *Bot) cleanChat(in models.InMessage) {
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
			text := fmt.Sprintf("(DS)%s \n%s", in.Name, in.Mtext)
			go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 180)
			go b.cleanChat(in)
		} else if in.Tip == tg {
			text := fmt.Sprintf("(TG)%s \n%s", in.Name, in.Mtext)
			go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 180)
			go b.cleanChat(in)
		}
	}
	return in.Config.Forward
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
					mID, err := strconv.Atoi(s.MesidTgHelp)
					if err != nil {
						return
					}
					go b.client.Tg.DelMessage(s.TgChannel, mID)
				}
				mid := b.client.Tg.SendHelp(s.TgChannel, strings.Replace(text, "3", "10", 1))
				s.MesidTgHelp = strconv.Itoa(mid)

			} else if s.TgChannel != "" && !s.Forward {
				split := strings.Split(s.TgChannel, "/")
				if split[1] != "0" {
					text := fmt.Sprintf("%s\n%s ", b.getLanguageText(s.Country, "information"), b.getLanguageText(s.Country, "info_help_text"))
					if s.MesidTgHelp != "" {
						mID, err := strconv.Atoi(s.MesidTgHelp)
						if err != nil {
							return
						}
						go b.client.Tg.DelMessage(s.TgChannel, mID)
					}
					mid := b.client.Tg.SendHelp(s.TgChannel, text)
					s.MesidTgHelp = strconv.Itoa(mid)
				}
			}
			b.storage.ConfigRs.AutoHelpUpdateMesid(s)
		}
		time.Sleep(time.Minute)
		go b.client.Ds.CleanRsBotOtherMessage()
	} else if tm.Minute() == 0 {
		go func() {
			a := b.storage.ConfigRs.AutoHelp()
			for _, s := range a {
				if s.DsChannel != "" {
					MesidDsHelp := b.client.Ds.HelpChannelUpdate(s)
					if MesidDsHelp != s.MesidDsHelp {
						s.MesidDsHelp = MesidDsHelp
						b.storage.ConfigRs.AutoHelpUpdateMesid(s)
					}
					in := models.InMessage{Config: s}
					b.QueueAll(in)
				}
			}
		}()
	}
	time.Sleep(time.Minute)
}
