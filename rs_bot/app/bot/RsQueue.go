package bot

import (
	"fmt"
	"rs/models"
	"rs/pkg/utils"
)

//lang ok

func (b *Bot) QueueLevel(in models.InMessage) {
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
	// совподения количество условие
	if count == 0 {
		if !in.Opt.Contains(models.OptionQueue) {
			text := b.getText(in, "rs_queue") + in.Lvlkz + b.getText(in, "empty")
			b.ifTipSendTextDelSecond(in, text, 10)
		} else if in.Opt.Contains(models.OptionQueue) {
			b.ifTipSendTextDelSecond(in, b.getText(in, "no_active_queues"), 10)
		}
		return
	}

	u := b.storage.DbFunc.ReadAll(in.Lvlkz, in.Config.CorpName)

	n := b.getMap(in, numberLvl)

	n = b.helpers.GetQueueDiscord(n, u)

	texttg := ""
	darkStar, level := containsSymbolD(in.Lvlkz)

	if in.Config.TgChannel != "" {
		ntg := make(map[string]string)
		if darkStar {
			ntg["text1"] = fmt.Sprintf("%s%s (%d)\n", b.getText(in, "queue_drs"), level, numberLvl)
		} else {
			ntg["text1"] = fmt.Sprintf("%s%s (%d)\n", b.getText(in, "rs_queue"), in.Lvlkz, numberLvl)
		}

		ntg["text2"] = fmt.Sprintf("\n%s++ - %s", in.Lvlkz, b.getText(in, "forced_start"))
		ntg["min"] = b.getText(in, "min")
		texttg = b.helpers.GetQueueTelegram(ntg, u)
	}

	fd := func(in models.InMessage) {
		if in.Opt.Contains(models.OptionEdit) {
			if err = b.client.Ds.EditComplexButton(u.User1.Dsmesid, in.Config.DsChannel, n); err != nil {
				b.log.Info(fmt.Sprintf("QueueLevel_ds %s %s %s\n%+v\n", u.User1.Dsmesid, in.Config.DsChannel, in.Config.CorpName, err))
				//in.O ption.Edit = false
				in.Opt.Remove(models.OptionEdit)
			}
		}
		if !in.Opt.Contains(models.OptionEdit) {
			dsmesid := b.client.Ds.SendComplex(in.Config.DsChannel, n)
			if dsmesid != "" {
				if err = b.storage.Update.MesidDsUpdate(dsmesid, in.Lvlkz, in.Config.CorpName); err != nil {
					b.log.ErrorErr(err)
				}
			}
			go b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
		}
	}
	ft := func(in models.InMessage) {
		if in.Opt.Contains(models.OptionEdit) {
			err = b.client.Tg.EditMessageTextKey(in.Config.TgChannel, u.User1.Tgmesid, texttg, in.Lvlkz)
			if err != nil {
				b.log.Info(fmt.Sprintf("QueueLevel_tg %d %s %s\n%+v\n", u.User1.Tgmesid, in.Config.TgChannel, in.Config.CorpName, err))
				//in.O ption.Edit = false
				in.Opt.Remove(models.OptionEdit)
			}
		}
		if !in.Opt.Contains(models.OptionEdit) {
			mesidTg := b.client.Tg.SendEmbed(in.Lvlkz, in.Config.TgChannel, texttg)
			if mesidTg != 0 {
				if err = b.storage.Update.MesidTgUpdate(mesidTg, in.Lvlkz, in.Config.CorpName); err != nil {
					b.log.ErrorErr(err)
				}
				go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
			}
		}
	}

	if count == 1 {

		if in.Config.DsChannel != "" {
			b.wg.Add(1)
			go func() {
				ch := utils.WaitForMessage("QueueLevel123")
				fd(in)
				b.wg.Done()
				close(ch)
			}()
		}
		if in.Config.TgChannel != "" {
			b.wg.Add(1)
			go func() {
				ch := utils.WaitForMessage("QueueLevel132")
				ft(in)
				b.wg.Done()
				close(ch)
			}()
		}
	} else if count == 2 {
		if in.Config.DsChannel != "" {
			b.wg.Add(1)
			go func() {
				ch := utils.WaitForMessage("QueueLevel142")
				fd(in)
				b.wg.Done()
				close(ch)
			}()
		}
		if in.Config.TgChannel != "" {
			b.wg.Add(1)
			go func() {
				ch := utils.WaitForMessage("QueueLevel151")
				ft(in)
				b.wg.Done()
				close(ch)
			}()
		}

	} else if count == 3 {
		if !darkStar {
			b.log.InfoStruct("QUEUE LEVEL 141 ", u)
			if in.Config.DsChannel != "" {
				b.wg.Add(1)
				go func() {
					ch := utils.WaitForMessage("QueueLevel163")
					fd(in)
					b.wg.Done()
					close(ch)
				}()
			}
			if in.Config.TgChannel != "" {
				b.wg.Add(1)
				go func() {
					ch := utils.WaitForMessage("QueueLevel172")
					ft(in)
					b.wg.Done()
					close(ch)
				}()
			}
		} else {
			b.log.InfoStruct("QUEUE LEVEL 161 ", u)
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
				if in.Opt.Contains(models.OptionUpdateAutoHelp) {
					in.Opt.Remove(models.OptionUpdateAutoHelp)
					if len(in.Opt) > 0 {
						fmt.Printf("opt %+v\n", in.Opt)
					}
				}
				in.Opt.Add(models.OptionQueueAll)
				in.Lvlkz = corp
				in.Mtext = ""
				b.Inbox <- in
				//in.Option.Queue = true
				//b.QueueLevel(in)
				//time.Sleep(1 * time.Second)

			}
		}
	} else if in.Opt.Contains(models.OptionInClient) {
		b.ifTipSendTextDelSecond(in, b.getText(in, "no_active_queues"), 10)
		b.iftipdelete(in)
	}

}
