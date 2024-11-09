package bot

import (
	"rs/models"
	"rs/pkg/utils"
	"strconv"
)

//lang ok
//wats lang not ok

func (b *Bot) MinusMin() {
	tt := b.storage.Timers.MinusMin()
	go b.client.Ds.QueueSend(b.otherQueue.MyQueue())

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
						Lvlkz:       t.Lvlkz,
						Timekz:      strconv.Itoa(t.Timedown),
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
						Option: models.Option{
							MinusMin: true,
							Edit:     true},
					}

					if t.Timedown == 3 {
						text := in.NameMention + b.getText(in, "info_time_almost_up")
						if in.Tip == ds {
							mID := b.client.Ds.SendEmbedTime(in.Config.DsChannel, text)
							go b.client.Ds.DeleteMessageSecond(in.Config.DsChannel, mID, 180)
						} else if in.Tip == tg {
							mID := b.client.Tg.SendEmbedTime(in.Config.TgChannel, text)
							go b.client.Tg.DelMessageSecond(in.Config.TgChannel, strconv.Itoa(mID), 180)
						}
					} else if t.Timedown == 0 || t.Timedown < -1 || t.Timedown < 0 {
						b.RsMinus(in)
					}

				}
			}
		}
		corpActive0 := b.storage.DbFunc.OneMinutsTimer()
		for _, corp := range corpActive0 {

			_, config := b.CheckCorpNameConfig(corp)

			dss, tgs := b.storage.DbFunc.MessageUpdateMin(corp)

			if config.DsChannel != "" && config.TgChannel != "" {
				for _, d := range dss {
					in := b.storage.DbFunc.MessageupdateDS(d, config)
					b.QueueLevel(in)
				}
				for _, t := range tgs {
					in := b.storage.DbFunc.MessageupdateTG(t, config)
					b.QueueLevel(in)
				}
			} else if config.DsChannel != "" && config.TgChannel == "" {
				for _, d := range dss {
					in := b.storage.DbFunc.MessageupdateDS(d, config)
					go b.QueueLevel(in)
				}
			} else if config.TgChannel != "" && config.DsChannel == "" {
				for _, t := range tgs {
					in := b.storage.DbFunc.MessageupdateTG(t, config)
					go b.QueueLevel(in)
				}
			}
		}
	}
}

func (b *Bot) ReadQueueLevel(in models.InMessage) {
	ch := utils.WaitForMessage("ReadQueueLevel")
	defer close(ch)
	text, err := b.otherQueue.ReadingQueueByLevel(in.Lvlkz, in.Config.CorpName)
	if err != nil {
		b.log.ErrorErr(err)
		return
	}

	if text != "" {
		b.ifTipSendTextDelSecond(in, text, 30)
	}
}