package bot

import (
	"fmt"
	"kz_bot/config"
	"kz_bot/models"
	"kz_bot/pkg/utils"
)

//lang ok

func (b *Bot) QueueLevel(in models.InMessage) {
	if in.Config.DsChannel != "1210280495238090782" && config.Instance.BotMode == "dev" {
		return
	}
	b.iftipdelete(in)

	count, err := b.storage.Count.CountQueue(in.Lvlkz, in.Config.CorpName)
	if err != nil {
		b.log.ErrorErr(err)
		return
	}
	if count == 0 {
		in.Lvlkz = "d" + in.Lvlkz
		count, err = b.storage.Count.CountQueue(in.Lvlkz, in.Config.CorpName)
		if err != nil {
			b.log.ErrorErr(err)
			return
		}
		if count == 0 {
			in.Lvlkz = in.Lvlkz[1:]
		}
	}
	numberLvl, err2 := b.storage.DbFunc.NumberQueueLvl(in.Lvlkz, in.Config.CorpName)
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

	u := b.storage.DbFunc.ReadAll(in.Lvlkz, in.Config.CorpName)

	n := b.getMap(in, numberLvl)

	n = b.helpers.GetQueueDiscord(n, u)

	texttg := ""
	if in.Config.TgChannel != "" {
		ntg := make(map[string]string)
		ntg["text1"] = fmt.Sprintf("%s%s (%d)\n", b.getText(in, "rs_queue"), in.Lvlkz, numberLvl)
		ntg["text2"] = fmt.Sprintf("\n%s++ - %s", in.Lvlkz, b.getText(in, "forced_start"))
		ntg["min"] = b.getText(in, "min")
		texttg = b.helpers.GetQueueTelegram(ntg, u)
	}

	fd := func() {
		//emb := b.client.Ds.EmbedDS(n, numberLvl, count, darkStar)
		if in.Option.Edit {
			errr := b.client.Ds.EditComplexButton(u.User1.Dsmesid, in.Config.DsChannel, n)
			if errr != nil {
				b.log.Info(fmt.Sprintf("QueueLevel %s %s", u.User1.Dsmesid, in.Config.DsChannel))
				b.log.ErrorErr(errr)
				in.Option.Edit = false
			}
		}
		if !in.Option.Edit {
			b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
			dsmesid := b.client.Ds.SendComplex(in.Config.DsChannel, n)

			err = b.storage.Update.MesidDsUpdate(dsmesid, in.Lvlkz, in.Config.CorpName)
			if err != nil {
				err = b.storage.Update.MesidDsUpdate(dsmesid, in.Lvlkz, in.Config.CorpName)
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
			mesidTg := b.client.Tg.SendEmbed(in.Lvlkz, in.Config.TgChannel, texttg)
			err = b.storage.Update.MesidTgUpdate(mesidTg, in.Lvlkz, in.Config.CorpName)
			if err != nil {
				err = b.storage.Update.MesidTgUpdate(mesidTg, in.Lvlkz, in.Config.CorpName)
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
				ch := utils.WaitForMessage("QueueLevel123")
				fd()
				b.wg.Done()
				close(ch)
			}()
		}
		if in.Config.TgChannel != "" {
			b.wg.Add(1)
			go func() {
				ch := utils.WaitForMessage("QueueLevel132")
				ft()
				b.wg.Done()
				close(ch)
			}()
		}
	} else if count == 2 {
		if in.Config.DsChannel != "" {
			b.wg.Add(1)
			go func() {
				ch := utils.WaitForMessage("QueueLevel142")
				fd()
				b.wg.Done()
				close(ch)
			}()
		}
		if in.Config.TgChannel != "" {
			b.wg.Add(1)
			go func() {
				ch := utils.WaitForMessage("QueueLevel151")
				ft()
				b.wg.Done()
				close(ch)
			}()
		}

	} else if count == 3 {

		if in.Config.DsChannel != "" {
			b.wg.Add(1)
			go func() {
				ch := utils.WaitForMessage("QueueLevel163")
				fd()
				b.wg.Done()
				close(ch)
			}()
		}
		if in.Config.TgChannel != "" {
			b.wg.Add(1)
			go func() {
				ch := utils.WaitForMessage("QueueLevel172")
				ft()
				b.wg.Done()
				close(ch)
			}()
		}
	}
	b.wg.Wait()
}
func (b *Bot) QueueAll(in models.InMessage) {
	lvl := b.storage.DbFunc.Queue(in.Config.CorpName)
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
