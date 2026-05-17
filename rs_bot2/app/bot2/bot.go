package bot2

import (
	"fmt"
	"rs/bot2/helpers"
	"rs/bot2/otherQueue"
	"rs/clients"
	"rs/models"
	"rs/storage"
	"rs/storage/dictionary"
	"rs/storage/postgresV2"
	"strings"
	"sync"
	"time"

	"github.com/mentalisit/logger"
)

const (
	ds = "ds"
	tg = "tg"
	wa = "wa"
)

type Bot struct {
	storage         *postgresV2.Db
	Dictionary      *dictionary.Dictionary
	client          *clients.Clients
	Inbox           chan models.InMessageV2
	log             *logger.Logger
	wg              sync.WaitGroup
	mu              sync.Mutex
	helpers         *helpers.Helpers
	otherQueue      *otherQueue.OtherQ
	AddLinkCodeFunc func(code, userID, username, provider string)
}

func NewBot(storage *storage.Storage, client *clients.Clients, log *logger.Logger) *Bot {
	b := &Bot{
		storage:    storage.V2,
		Dictionary: storage.Dictionary,
		client:     client,
		log:        log,
		Inbox:      make(chan models.InMessageV2, 30),
		helpers:    helpers.NewHelpers(log, storage.V2),
		otherQueue: otherQueue.NewOtherQ(log),
	}

	go b.loadInbox()
	go b.timerBot()

	return b
}

func (b *Bot) loadInbox() {
	//b.log.Info("Бот загружен и готов к работе ")

	for in := range b.Inbox {
		multiAccount, err := b.storage.FindMultiAccountByUserId(in.UserId)
		if err != nil || multiAccount == nil {
			multiAccount, _ = b.storage.CreateMultiAccountWithPlatform(in.UserId, in.Username, in.Messenger.TypeMessenger, in.Username)
		}
		if multiAccount != nil {
			in.MAcc = multiAccount
		}
		b.LogicRs(&in)
	}
}

func (b *Bot) timerBot() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		b.MinusMin()
		b.readNews()
		b.studyCompendium()

		tn := time.Now().UTC()

		//if now.Hour() == 23 && now.Minute() == 59 {
		//	b.storage.Battles.DeleteOldWebhooks()
		//}
		//
		if tn.Hour() == 12 && tn.Minute() == 0 {
			go b.AutoHelp()
		}
		switch tn.Minute() {

		case 0:
			//go b.AutoHelp() //автозапуск справки//todo need module telegram

		case 9, 19, 29, 39, 49, 59:
			//go b.ReadAndSendPic()

		}
	}
}

// LogicRs логика игры
func (b *Bot) LogicRs(in *models.InMessageV2) {
	if strings.HasPrefix(in.Text, ".") {
		b.accessChat(in)
		return
	}
	fmt.Printf("In %+v\n", in)

	if in.Text == "" {
		return
	}

	if b.lRsPlus(in) {
	} else if b.lSubs(in) {
	} else if b.lQueue(in) {
	} else if b.lRsStart(in) {
	} else if b.lTop(in) {
	} else if b.lEmoji(in) {
	} else if b.lHelp(in) {
	} else if b.bridge(in) {
	} else if strings.HasPrefix(in.Text, ".всем") {
		b.SendALLChannel(in)
	} else {
		b.cleanChat(in)
	}
}

func (b *Bot) cleanChat(in *models.InMessageV2) {
	if !strings.HasPrefix(in.Text, ".") && !in.Options.Contains(models.OptionEdit) {
		if in.Config.Channels[in.Messenger.ChannelId].Corp != nil && in.Config.Channels[in.Messenger.ChannelId].Corp.DeleteMessages {
			if in.Tip == ds {
				go b.client.Ds.DeleteMessageSecond(in.Messenger.ChannelId, in.Messenger.MessageId, 600)
			}
			if in.Tip == tg {
				go b.client.Tg.DelMessageSecond(in.Messenger.ChannelId, in.Messenger.MessageId, 600)
			}
		}

	}
}

func (b *Bot) bridge(in *models.InMessageV2) bool {
	if len(in.Config.Channels) != 1 {
		if in.Tip == ds {
			in.Text = b.client.Ds.ReplaceTextMessage(in.Text, in.Messenger.GuildId)
		}
		text := fmt.Sprintf("%s (%s) \n%s", in.Username, strings.ToUpper(in.Messenger.TypeMessenger), in.Text)
		for _, s := range in.Config.Channels {
			if in.Messenger.ChannelId != s.ChannelId {
				b.sendInfo(s, text)

				go b.cleanChat(in)
			}
		}
		return true
	}
	return false
}

func (b *Bot) SendALLChannel(in *models.InMessageV2) (bb bool) {
	text, found := strings.CutPrefix(in.Text, ".всем")
	if found && b.checkAdmin(in) {
		b.deleteInMessage(in)

		b.sendTextAfterDeleteSecond(in, "Начата рассылка.", 20)

		go func(textCopy string) {
			for _, config := range b.storage.ReadConfigV2() {
				for ch, i := range config.Channels {
					b.sendTypeMessenger(i.TypeMessenger, ch, textCopy, 86400)
				}
			}
		}(text)

		bb = true
	}

	return bb
}

func (b *Bot) sendTypeMessenger(TypeMessenger, ChannelId, text string, seconds int) {
	if TypeMessenger == ds {
		go b.client.Ds.SendChannelDelSecond(ChannelId, text, seconds)
	} else if TypeMessenger == tg {
		go b.client.Tg.SendChannelDelSecond(ChannelId, text, seconds)
	} else if TypeMessenger == wa {
		go b.client.Wa.SendChannelDelSecond(ChannelId, text, seconds)
	} else {
		b.log.Info("need implement for " + TypeMessenger)
	}
}

func (b *Bot) sendInfo(s *models.Info, text string) {
	if s.Corp != nil && s.Corp.DeleteMessages {
		timeDelete := s.Corp.DeleteMessagesDelay * 60
		b.sendTypeMessenger(s.TypeMessenger, s.ChannelId, text, timeDelete)
	} else {
		if s.TypeMessenger == ds {
			go b.client.Ds.Send(s.ChannelId, text)
		} else if s.TypeMessenger == tg {
			go b.client.Tg.SendChannel(s.ChannelId, text)
		} else if s.TypeMessenger == wa {
			go b.client.Wa.SendChannel(s.ChannelId, text)
		} else {
			b.log.Info("need implement for " + s.TypeMessenger)
		}
	}

}

func (b *Bot) studyCompendium() {
	// 1. Получаем всех, у кого есть хоть одна просроченная запись
	usersWithExpired, err := b.storage.GetAllExpiredStudy()
	if err != nil {
		b.log.ErrorErr(err)
		return
	}

	nowMs := time.Now().UnixMilli()

	for _, userRecord := range usersWithExpired {
		var stillActive []models.Studies
		var expired []models.Studies

		// 2. Сортируем модули пользователя на просроченные и активные
		for _, s := range userRecord.Studies {
			if s.EndTime < nowMs {
				expired = append(expired, s)
			} else {
				stillActive = append(stillActive, s)
			}
		}

		// 3. Обрабатываем просроченные
		for _, exp := range expired {
			multiAccount, _ := b.storage.FindMultiAccountByUUId(userRecord.Uuid.String())
			if multiAccount != nil {
				text := fmt.Sprintf("Модуль %s обновлен до %d уровня", exp.Name, exp.Level)
				if multiAccount.DiscordID != "" {
					go b.client.Ds.SendDmText(text, multiAccount.DiscordID)
				}
				if multiAccount.TelegramID != "" {
					go b.client.Tg.SendChannel(multiAccount.TelegramID, text)
				}

			}

			// Передаем изменение в основную базу модулей
			err = b.storage.SyncModuleStatus(userRecord, exp)
			if err != nil {
				b.log.ErrorErr(err)
			}
		}

		// 4. Очищаем базу: если ничего не осталось — удаляем строку, иначе обновляем JSON
		if len(expired) > 0 {
			if len(stillActive) == 0 {
				// Если активных модулей не осталось, удаляем всю запись пользователя
				err = b.storage.DeleteStudyRecord(userRecord)
			} else {
				// Если есть еще живые модули, перезаписываем JSON
				userRecord.Studies = stillActive
				err = b.storage.UpdateStudies(userRecord)
			}

			if err != nil {
				b.log.ErrorErr(err)
			}
		}
	}
}
