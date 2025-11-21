package bot

import (
	"fmt"
	"rs/models"
	"rs/pkg/utils"
	"strconv"
	"time"
)

//lang ok
//wats lang not ok

func (b *Bot) MinusMin() {
	tt := b.storage.Timers.MinusMin()
	//go b.client.Ds.QueueSend(b.otherQueue.MyQueue())
	if time.Now().Minute()%5 == 0 {
		go func() {
			en, ru, ua := b.client.Ds.ReadNews()
			if en != "" && ru != "" && ua != "" {

				sendChannel := func(config models.CorporationConfig, text string) {
					text = fmt.Sprintf("%s \n%s", "Hades' Star Official", text)
					if config.TgChannel != "" {
						b.client.Tg.SendChannelDelSecond(config.TgChannel, text, 172800)
					}
					if config.DsChannel != "" {
						b.client.Ds.SendChannelDelSecond(config.DsChannel, text, 172800)
					}
				}

				for _, config := range b.storage.ConfigRs.ReadConfigRs() {
					if config.Country == "en" {
						sendChannel(config, en)
					}
					if config.Country == "ru" {
						sendChannel(config, ru)
					}
					if config.Country == "ua" {
						sendChannel(config, ua)
					}
				}
			}
		}()
	}

	if len(tt) > 0 {
		for _, t := range tt {
			if t.Corpname != "" {
				ok, config := b.CheckCorpNameConfig(t.Corpname)
				if ok {
					in := models.InMessage{
						Mtext:       "",
						Tip:         t.Tip,
						Username:    t.Name,
						UserId:      t.UserId,
						NameMention: t.Mention,
						RsTypeLevel: t.Lvlkz,
						TimeRs:      t.Timedown,
						Ds: struct {
							Mesid   string
							Guildid string
							Avatar  string
						}{
							Mesid:   t.Dsmesid,
							Guildid: config.Guildid,
						},
						Tg: struct {
							Mesid int
						}{
							Mesid: t.Tgmesid,
						},
						Config: config,
						Opt:    []string{models.OptionMinusMin},
					}

					if t.Timedown == 3 {
						text := in.GetNameMention() + b.getText(in, "info_time_almost_up")
						if in.Tip == ds {
							mID := b.client.Ds.SendEmbedTime(in.Config.DsChannel, text)
							go b.client.Ds.DeleteMessageSecond(in.Config.DsChannel, mID, 180)
						} else if in.Tip == tg {
							mID := b.client.Tg.SendEmbedTime(in.Config.TgChannel, text)
							go b.client.Tg.DelMessageSecond(in.Config.TgChannel, strconv.Itoa(mID), 180)
						}
						go b.ReadQueueLevel(in, 180)
					} else if t.Timedown == 0 || t.Timedown < -1 || t.Timedown < 0 {
						b.RsMinus(in)
					}

				}
			}
		}

		in := models.InMessage{
			Mtext: "",
			Opt:   []string{models.OptionMinusMinNext, models.OptionEdit},
		}
		b.Inbox <- in
	}

	go func() {
		timers := b.storage.TimeDeleteMessage.TimerMessage()
		for _, timer := range timers {
			if timer.Dsmesid != "" {
				b.client.Ds.DeleteMessage(timer.Dsmesid, timer.Dsmesid)
			} else if timer.Tgmesid != "" {
				atoi, _ := strconv.Atoi(timer.Tgmesid)
				b.client.Tg.DelMessage(timer.Tgchatid, atoi)
			}
			b.storage.TimeDeleteMessage.TimerDeleteMessage(timer)
		}
	}()

}
func (b *Bot) MinusMinMessageUpdate() {
	corpActive0 := b.storage.DbFunc.OneMinutsTimer()
	for _, corp := range corpActive0 {

		_, config := b.CheckCorpNameConfig(corp)

		dss, tgs := b.storage.DbFunc.MessageUpdateMin(corp)

		if config.DsChannel != "" {
			for _, d := range dss {
				b.Inbox <- b.storage.DbFunc.MessageUpdateDS(d, config)
			}
		}
		if config.TgChannel != "" {
			for _, t := range tgs {
				b.Inbox <- b.storage.DbFunc.MessageUpdateTG(t, config)
			}
		}
	}

}

func (b *Bot) ReadQueueLevel(in models.InMessage, second int) {
	ch := utils.WaitForMessage("ReadQueueLevel")
	defer close(ch)
	text, err := b.otherQueue.ReadingQueueByLevel(in.RsTypeLevel, in.Config.CorpName)
	if err != nil {
		b.log.ErrorErr(err)
		return
	}

	if text != "" {
		b.ifTipSendTextDelSecond(in, text, second)
	}
}
