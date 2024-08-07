package bot

import (
	"context"
	"fmt"
	"kz_bot/clients/restapi"
	conf "kz_bot/config"
	"kz_bot/models"
	"sort"
	"strconv"
	"time"
)

//lang ok
//wats lang not ok

func (b *Bot) MinusMin() {
	tt := b.storage.Timers.MinusMin(context.Background())

	if conf.Instance.BotMode == "dev" {
		return
	}

	go b.myQueue(tt)

	if len(tt) > 0 {
		for _, t := range tt {
			if t.Corpname != "" {
				ok, config := b.CheckCorpNameConfig(t.Corpname)
				if ok {
					timeDown := strconv.Itoa(t.Timedown)

					in := models.InMessage{
						Mtext:       "",
						Tip:         t.Tip,
						Username:    t.Name,
						UserId:      t.UserId,
						NameMention: t.Mention,
						Lvlkz:       t.Lvlkz,
						Timekz:      timeDown,
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
					b.inbox <- in

					if b.debug {
						fmt.Printf("\n  MinusMin []models.Sborkz %+v\n\n", t)
					}
				}
			}
		}
		b.UpdateMessage()
	}
}
func (b *Bot) UpdateMessage() {
	corpActive0 := b.storage.DbFunc.OneMinutsTimer(context.Background())
	for _, corp := range corpActive0 {

		_, config := b.CheckCorpNameConfig(corp)

		dss, tgs := b.storage.DbFunc.MessageUpdateMin(context.Background(), corp)

		if config.DsChannel != "" {
			for _, d := range dss {
				a := b.storage.DbFunc.MessageupdateDS(context.Background(), d, config)
				b.inbox <- a
			}
		}
		if config.TgChannel != "" {
			for _, t := range tgs {
				a := b.storage.DbFunc.MessageupdateTG(context.Background(), t, config)
				b.inbox <- a
			}
		}
	}
}

func (b *Bot) CheckTimeQueue(in models.InMessage) {
	atoi, err := strconv.Atoi(in.Timekz)
	if err != nil {
		b.log.ErrorErr(err)
	}
	if atoi == 3 {
		text := in.NameMention + b.getText(in, "info_time_almost_up")
		if in.Tip == ds {
			mID := b.client.Ds.SendEmbedTime(in.Config.DsChannel, text)
			go b.client.Ds.DeleteMesageSecond(in.Config.DsChannel, mID, 180)
		} else if in.Tip == tg {
			mID := b.client.Tg.SendEmbedTime(in.Config.TgChannel, text)
			go b.client.Tg.DelMessageSecond(in.Config.TgChannel, strconv.Itoa(mID), 180)
		}
	} else if atoi == 0 {
		b.RsMinus(in)
	} else if atoi < -1 {
		b.RsMinus(in)
	} else if atoi < 0 {
		b.RsMinus(in)
	}
}
func (b *Bot) myQueue(my []models.Sborkz) {
	time.Sleep(5 * time.Second)
	var text string
	if len(my) != 0 {
		sort.Slice(my, func(i, j int) bool {
			return my[i].Corpname < my[j].Corpname
		})
		var corpName string
		var level string
		var count int
		for _, sborkz := range my {
			if corpName != sborkz.Corpname {
				text += fmt.Sprintf("⚠️ сбор в %s \n", sborkz.Corpname)
				corpName = sborkz.Corpname
				level = ""
			}
			if sborkz.Corpname == corpName {
				if level != sborkz.Lvlkz {
					text += fmt.Sprintf("🔥 на кз%s \n", sborkz.Lvlkz)
					level = sborkz.Lvlkz
					count = 1
				}
				if level == sborkz.Lvlkz {
					text += fmt.Sprintf("%d. %s  %d\n", count, sborkz.Name, sborkz.Timedown)
					count += 1
				}
			}
		}
	}
	if text == "" {
		text = "нет активных очередей"
	}

	b.client.Ds.QueueSend(text, "RsBot")

	time.Sleep(5 * time.Second)

	queue, err := restapi.RsbotQueue()
	if err != nil {
		b.log.ErrorErr(err)
	}

	if queue == "" {
		queue = "нет активных очередей"
	}
	b.client.Ds.QueueSend(queue, "rssoyzbot")
}
