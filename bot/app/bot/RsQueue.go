package bot

import (
	"context"
	"fmt"
	"kz_bot/config"
	"kz_bot/models"
	"kz_bot/pkg/utils"
	"time"
)

//lang ok

func (b *Bot) QueueLevel(in models.InMessage) {
	if in.Config.DsChannel != "1210280495238090782" && config.Instance.BotMode == "dev" {
		return
	}
	b.iftipdelete(in)
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	count, err := b.storage.Count.CountQueue(ctx, in.Lvlkz, in.Config.CorpName)
	if err != nil {
		b.log.ErrorErr(err)
		return
	}
	if count == 0 {
		in.Lvlkz = "d" + in.Lvlkz
		count, err = b.storage.Count.CountQueue(ctx, in.Lvlkz, in.Config.CorpName)
		if err != nil {
			b.log.ErrorErr(err)
			return
		}
		if count == 0 {
			in.Lvlkz = in.Lvlkz[1:]
		}
	}
	numberLvl, err2 := b.storage.DbFunc.NumberQueueLvl(ctx, in.Lvlkz, in.Config.CorpName)
	if err2 != nil {
		b.log.ErrorErr(err)
		return
	}
	// совподения количество  условие
	if count == 0 {
		if !in.Option.Queue {
			text := b.getText(in, "rs_queue") + in.Lvlkz + b.getText(in, "empty")
			b.ifTipSendTextDelSecond(in, text, 10)
		} else if in.Option.Queue {
			b.ifTipSendTextDelSecond(in, b.getText(in, "no_active_queues"), 10)
		}
		return
	}

	u := b.storage.DbFunc.ReadAll(ctx, in.Lvlkz, in.Config.CorpName)
	var n map[string]string
	n = make(map[string]string)
	n["lang"] = in.Config.Country
	n = b.helpers.GetQueueDiscord(n, u)

	texttg := ""
	if in.Config.TgChannel != "" {
		ntg := make(map[string]string)
		ntg["text1"] = fmt.Sprintf("%s%s (%d)\n", b.getText(in, "rs_queue"), in.Lvlkz, numberLvl)
		ntg["text2"] = fmt.Sprintf("\n%s++ - %s", in.Lvlkz, b.getText(in, "forced_start"))
		ntg["min"] = b.getText(in, "min")
		texttg = b.helpers.GetQueueTelegram(ntg, u)
	}
	darkStar, lvlkz := containsSymbolD(in.Lvlkz)
	if in.Config.DsChannel != "" {
		if darkStar {
			n["lvlkz"], err = b.client.Ds.RoleToIdPing(b.getText(in, "drs")+lvlkz, in.Config.Guildid)
		} else {
			n["lvlkz"], err = b.client.Ds.RoleToIdPing(b.getText(in, "rs")+in.Lvlkz, in.Config.Guildid)
		}
		if err != nil {
			b.log.Info(fmt.Sprintf("RoleToIdPing %+v lvl %s", in.Config, in.Lvlkz))
		}
	}

	fd := func() {
		emb := b.client.Ds.EmbedDS(n, numberLvl, count, darkStar)
		if in.Option.Edit {
			errr := b.client.Ds.EditComplexButton(u.User1.Dsmesid, in.Config.DsChannel, emb, b.client.Ds.AddButtonsQueue(in.Lvlkz))
			if errr != nil {
				b.log.Info(fmt.Sprintf("QueueLevel %s %s", u.User1.Dsmesid, in.Config.DsChannel))
				b.log.ErrorErr(errr)
				in.Option.Edit = false
			}
		}
		if !in.Option.Edit {
			b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
			dsmesid := b.client.Ds.SendComplex(in.Config.DsChannel, emb, b.client.Ds.AddButtonsQueue(in.Lvlkz))

			err = b.storage.Update.MesidDsUpdate(ctx, dsmesid, in.Lvlkz, in.Config.CorpName)
			if err != nil {
				err = b.storage.Update.MesidDsUpdate(context.Background(), dsmesid, in.Lvlkz, in.Config.CorpName)
				if err != nil {
					b.log.ErrorErr(err)
				}
			}
		}
	}
	ft := func() {
		if in.Option.Edit {
			b.client.Tg.EditMessageTextKey(in.Config.TgChannel, u.User1.Tgmesid, texttg, in.Lvlkz)
		} else if !in.Option.Edit {
			mesidTg := b.client.Tg.SendEmded(in.Lvlkz, in.Config.TgChannel, texttg)
			err = b.storage.Update.MesidTgUpdate(ctx, mesidTg, in.Lvlkz, in.Config.CorpName)
			if err != nil {
				err = b.storage.Update.MesidTgUpdate(context.Background(), mesidTg, in.Lvlkz, in.Config.CorpName)
				if err != nil {
					b.log.ErrorErr(err)
				}
			}
			b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
		}
	}
	if count == 1 {

		if in.Config.DsChannel != "" {
			b.wg.Add(1)
			go func() {
				fd()
				b.wg.Done()
			}()
		}
		if in.Config.TgChannel != "" {
			b.wg.Add(1)
			go func() {
				ft()
				b.wg.Done()
			}()
		}
	} else if count == 2 {
		if in.Config.DsChannel != "" {
			b.wg.Add(1)
			go func() {
				fd()
				b.wg.Done()
			}()
		}
		if in.Config.TgChannel != "" {
			b.wg.Add(1)
			go func() {
				ft()
				b.wg.Done()
			}()
		}

	} else if count == 3 {

		if in.Config.DsChannel != "" {
			b.wg.Add(1)
			go func() {
				fd()
				b.wg.Done()
			}()
		}
		if in.Config.TgChannel != "" {
			b.wg.Add(1)
			go func() {
				ft()
				b.wg.Done()
			}()
		}
	}
	b.wg.Wait()
}
func (b *Bot) QueueAll(in models.InMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	lvl := b.storage.DbFunc.Queue(ctx, in.Config.CorpName)
	lvlk := utils.RemoveDuplicates(lvl)
	if len(lvlk) > 0 {
		for _, corp := range lvlk {
			if corp != "" {
				in.Option.Queue = true
				in.Lvlkz = corp
				b.QueueLevel(in)

			}
		}
	} else if in.Option.InClient {
		b.ifTipSendTextDelSecond(in, b.getText(in, "no_active_queues"), 10)
		b.iftipdelete(in)
	}

}
