package bot

import (
	"fmt"
	"rs/models"
	"rs/pkg/utils"
	"strings"
)

//lang ok

func (b *Bot) RsStart(in models.InMessage) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.iftipdelete(in)
	countName, err := b.storage.Count.СountName(in.UserId, in.Lvlkz, in.Config.CorpName)
	if err != nil {
		return
	}
	if countName == 0 {
		text := b.getText(in, "info_forced_start_available")
		b.ifTipSendTextDelSecond(in, text, 10)
	} else if countName == 1 {
		numberkz, err1 := b.storage.DbFunc.NumberQueueLvl(in.Lvlkz, in.Config.CorpName)
		if err1 != nil {
			return
		}
		count, err2 := b.storage.Count.CountQueue(in.Lvlkz, in.Config.CorpName)
		if err2 != nil {
			return
		}

		dsmesid := ""
		tgmesid := 0

		d, level := containsSymbolD(in.Lvlkz)

		if count > 0 {
			u := b.storage.DbFunc.ReadAll(in.Lvlkz, in.Config.CorpName)
			textEvent, numkzEvent := b.EventText(in)
			if textEvent == "" {
				textEvent = b.GetTextPercent(in.Config, d)
			}
			numberevent := b.storage.Event.NumActiveEvent(in.Config.CorpName)
			if numberevent > 0 {
				numberkz = numkzEvent
			}
			textStart := fmt.Sprintf("🚀 %s%s (%d) %s\n\n", b.getText(in, "queue_drs"), level, numberkz, b.getText(in, "was_launched_incomplete"))
			if !d {
				textStart = fmt.Sprintf("🚀 %s%s (%d) %s\n\n", b.getText(in, "rs_queue"), in.Lvlkz, numberkz, b.getText(in, "was_launched_incomplete"))

			}
			textEnd := fmt.Sprintf("\n%s %s", b.getText(in, "go"), textEvent)

			if count == 1 {
				if in.Config.DsChannel != "" {
					b.wg.Add(1)
					go func() {
						ch := utils.WaitForMessage("RsStart54")
						name1, _, _, _ := b.helpers.NameMention(u, ds)
						text := fmt.Sprintf("%s1. %s%s", textStart, name1, textEnd)

						dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Ds.Avatar)

						go b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
						if err = b.storage.Update.MesidDsUpdate(dsmesid, in.Lvlkz, in.Config.CorpName); err != nil {
							b.log.ErrorErr(err)
						}
						b.wg.Done()
						close(ch)
					}()

				}
				if in.Config.TgChannel != "" {
					b.wg.Add(1)
					go func() {
						ch := utils.WaitForMessage("RsStart83")
						name1, _, _, _ := b.helpers.NameMention(u, tg)
						go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
						text := fmt.Sprintf("%s1. %s %s", textStart, name1, textEnd)
						tgmesid = b.client.Tg.SendChannel(in.Config.TgChannel, text)
						if err = b.storage.Update.MesidTgUpdate(tgmesid, in.Lvlkz, in.Config.CorpName); err != nil {
							b.log.ErrorErr(err)
						}
						b.wg.Done()
						close(ch)
					}()

				}
			} else if count == 2 {
				if in.Config.DsChannel != "" { //discord
					b.wg.Add(1)
					go func() {
						ch := utils.WaitForMessage("RsStart106")
						name1, name2, _, _ := b.helpers.NameMention(u, ds)
						text := fmt.Sprintf("%s1. %s\n2. %s %s", textStart, name1, name2, textEnd)
						dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Ds.Avatar)
						//go b.SendLsNotification(in, u)

						go b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
						if err = b.storage.Update.MesidDsUpdate(dsmesid, in.Lvlkz, in.Config.CorpName); err != nil {
							b.log.ErrorErr(err)
						}
						b.wg.Done()
						close(ch)
					}()

				}
				if in.Config.TgChannel != "" { //telegram
					b.wg.Add(1)
					go func() {
						ch := utils.WaitForMessage("RsStart137")
						name1, name2, _, _ := b.helpers.NameMention(u, tg)
						go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
						text := fmt.Sprintf("%s1. %s\n2. %s %s", textStart, name1, name2, textEnd)
						tgmesid = b.client.Tg.SendChannel(in.Config.TgChannel, text)
						if err = b.storage.Update.MesidTgUpdate(tgmesid, in.Lvlkz, in.Config.CorpName); err != nil {
							b.log.ErrorErr(err)
						}
						b.wg.Done()
						close(ch)
					}()
				}
			} else if count == 3 {
				if in.Config.DsChannel != "" { //discord
					b.wg.Add(1)
					go func() {
						ch := utils.WaitForMessage("RsStart161")
						name1, name2, name3, _ := b.helpers.NameMention(u, ds)
						text := fmt.Sprintf("%s1. %s\n2. %s\n3. %s %s", textStart, name1, name2, name3, textEnd)
						dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Ds.Avatar)

						//go b.SendLsNotification(in, u)
						go b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
						if err = b.storage.Update.MesidDsUpdate(dsmesid, in.Lvlkz, in.Config.CorpName); err != nil {
							b.log.ErrorErr(err)
						}
						b.wg.Done()
						close(ch)
					}()
				}
				if in.Config.TgChannel != "" { //telegram
					b.wg.Add(1)
					go func() {
						ch := utils.WaitForMessage("RsStart186")
						name1, name2, name3, _ := b.helpers.NameMention(u, tg)
						go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
						text := fmt.Sprintf("%s1. %s\n2. %s\n3. %s %s", textStart, name1, name2, name3, textEnd)
						tgmesid = b.client.Tg.SendChannel(in.Config.TgChannel, text)
						if err = b.storage.Update.MesidTgUpdate(tgmesid, in.Lvlkz, in.Config.CorpName); err != nil {
							b.log.ErrorErr(err)
						}
						b.wg.Done()
						close(ch)
					}()

				}
			}
			b.wg.Wait()
			err = b.storage.Update.UpdateCompliteRS(in.Lvlkz, dsmesid, tgmesid, "", numberkz, numberevent, in.Config.CorpName)
			if err != nil {
				b.log.ErrorErr(err)
				err = b.storage.Update.UpdateCompliteRS(in.Lvlkz, dsmesid, tgmesid, "", numberkz, numberevent, in.Config.CorpName)
				if err != nil {
					b.log.ErrorErr(err)
				}
			}

			if numberevent == 0 {
				//отправляем сообщение о корпорациях с %
				go b.SendPercent(in.Config)
			}

			user := u.GetAllUserId()
			go b.elseChat(user)
			go b.helpers.SaveUsersIdQueue(user, in.Config)
		}
	}
}

func (b *Bot) RsStartOld(in models.InMessage) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.iftipdelete(in)
	countName, err := b.storage.Count.СountName(in.UserId, in.Lvlkz, in.Config.CorpName)
	if err != nil {
		return
	}
	if countName == 0 {
		text := b.getText(in, "info_forced_start_available")
		b.ifTipSendTextDelSecond(in, text, 10)
	} else if countName == 1 {
		numberkz, err1 := b.storage.DbFunc.NumberQueueLvl(in.Lvlkz, in.Config.CorpName)
		if err1 != nil {
			return
		}
		count, err2 := b.storage.Count.CountQueue(in.Lvlkz, in.Config.CorpName)
		if err2 != nil {
			return
		}

		dsmesid := ""
		tgmesid := 0
		if count > 0 {
			u := b.storage.DbFunc.ReadAll(in.Lvlkz, in.Config.CorpName)
			textEvent, numkzEvent := b.EventText(in)
			if textEvent == "" {
				DarkFlag := strings.HasPrefix(in.Lvlkz, "d")
				textEvent = b.GetTextPercent(in.Config, DarkFlag)
			}
			numberevent := b.storage.Event.NumActiveEvent(in.Config.CorpName)
			if numberevent > 0 {
				numberkz = numkzEvent
			}
			if count == 1 {
				if in.Config.DsChannel != "" {
					b.wg.Add(1)
					go func() {
						ch := utils.WaitForMessage("RsStart54")
						name1, _, _, _ := b.helpers.NameMention(u, ds)
						text := fmt.Sprintf("🚀 %s%s (%d) %s \n\n1. %s\n%s %s",
							b.getText(in, "rs_queue"), in.Lvlkz, numberkz,
							b.getText(in, "was_launched_incomplete"), name1, b.getText(in, "go"), textEvent)

						dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Ds.Avatar)

						go b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
						err = b.storage.Update.MesidDsUpdate(dsmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							b.log.ErrorErr(err)
						}
						b.wg.Done()
						close(ch)
					}()

				}
				if in.Config.TgChannel != "" {
					b.wg.Add(1)
					go func() {
						ch := utils.WaitForMessage("RsStart83")
						name1, _, _, _ := b.helpers.NameMention(u, tg)
						go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
						text := fmt.Sprintf("🚀 %s%s (%d) %s \n\n1. %s\n%s %s",
							b.getText(in, "rs_queue"), in.Lvlkz, numberkz,
							b.getText(in, "was_launched_incomplete"), name1, b.getText(in, "go"), textEvent)
						tgmesid = b.client.Tg.SendChannel(in.Config.TgChannel, text)
						err = b.storage.Update.MesidTgUpdate(tgmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							b.log.ErrorErr(err)
						}
						b.wg.Done()
						close(ch)
					}()

				}
			} else if count == 2 {
				if in.Config.DsChannel != "" { //discord
					b.wg.Add(1)
					go func() {
						ch := utils.WaitForMessage("RsStart106")
						name1, name2, _, _ := b.helpers.NameMention(u, ds)
						text1 := fmt.Sprintf("🚀 %s%s (%d) %s \n\n",
							b.getText(in, "rs_queue"), in.Lvlkz, numberkz, b.getText(in, "was_launched_incomplete"))
						text2 := fmt.Sprintf("%s\n%s\n%s %s", name1, name2, b.getText(in, "go"), textEvent)
						text := text1 + text2
						dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Ds.Avatar)
						if u.User1.Tip == ds {
							go b.client.Ds.SendDmText(text, u.User1.UserId)
						}

						go b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
						err = b.storage.Update.MesidDsUpdate(dsmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							b.log.ErrorErr(err)
						}
						b.wg.Done()
						close(ch)
					}()

				}
				if in.Config.TgChannel != "" { //telegram
					b.wg.Add(1)
					go func() {
						ch := utils.WaitForMessage("RsStart137")
						name1, name2, _, _ := b.helpers.NameMention(u, tg)
						go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
						text1 := fmt.Sprintf("🚀 %s%s (%d) %s \n\n",
							b.getText(in, "rs_queue"), in.Lvlkz, numberkz, b.getText(in, "was_launched_incomplete"))
						text2 := fmt.Sprintf("%s\n%s\n%s %s", name1, name2, b.getText(in, "go"), textEvent)
						text := text1 + text2
						tgmesid = b.client.Tg.SendChannel(in.Config.TgChannel, text)
						err = b.storage.Update.MesidTgUpdate(tgmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							b.log.ErrorErr(err)
						}
						b.wg.Done()
						close(ch)
					}()

				}
			} else if count == 3 {
				if in.Config.DsChannel != "" { //discord
					b.wg.Add(1)
					go func() {
						ch := utils.WaitForMessage("RsStart161")
						name1, name2, name3, _ := b.helpers.NameMention(u, ds)
						text := fmt.Sprintf("🚀 %s%s (%d) %s \n\n%s\n%s\n%s\n%s %s",
							b.getText(in, "rs_queue"), in.Lvlkz, numberkz, b.getText(in, "was_launched_incomplete"),
							name1, name2, name3, b.getText(in, "go"), textEvent)
						dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Ds.Avatar)

						if u.User1.Tip == ds {
							go b.client.Ds.SendDmText(text, u.User1.UserId)
						}
						if u.User2.Tip == ds {
							go b.client.Ds.SendDmText(text, u.User1.UserId)
						}
						go b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
						err = b.storage.Update.MesidDsUpdate(dsmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							b.log.ErrorErr(err)
						}
						b.wg.Done()
						close(ch)
					}()
				}
				if in.Config.TgChannel != "" { //telegram
					b.wg.Add(1)
					go func() {
						ch := utils.WaitForMessage("RsStart186")
						name1, name2, name3, _ := b.helpers.NameMention(u, tg)
						go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
						text := fmt.Sprintf("🚀 %s%s (%d) %s \n\n%s\n%s\n%s\n%s %s",
							b.getText(in, "rs_queue"), in.Lvlkz, numberkz, b.getText(in, "was_launched_incomplete"),
							name1, name2, name3, b.getText(in, "go"), textEvent)
						tgmesid = b.client.Tg.SendChannel(in.Config.TgChannel, text)
						err = b.storage.Update.MesidTgUpdate(tgmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							b.log.ErrorErr(err)
						}
						b.wg.Done()
						close(ch)
					}()

				}
			}
			b.wg.Wait()
			err = b.storage.Update.UpdateCompliteRS(in.Lvlkz, dsmesid, tgmesid, "", numberkz, numberevent, in.Config.CorpName)
			if err != nil {
				err = b.storage.Update.UpdateCompliteRS(in.Lvlkz, dsmesid, tgmesid, "", numberkz, numberevent, in.Config.CorpName)
				if err != nil {
					b.log.ErrorErr(err)
				}
			}

			if numberevent == 0 {
				//отправляем сообщение о корпорациях с %
				go b.SendPercent(in.Config)
			}

			user := u.GetAllUserId()
			go b.elseChat(user)
			go b.helpers.SaveUsersIdQueue(user, in.Config)
		}
	}
}
func (b *Bot) Pl30(in models.InMessage) {
	countName := b.storage.Count.CountNameQueue(in.UserId)
	text := ""
	if countName == 0 {
		text = in.NameMention + b.getText(in, "you_out_of_queue")
	} else if countName > 0 {
		timedown := b.storage.DbFunc.P30Pl(in.Lvlkz, in.Config.CorpName, in.UserId)
		if timedown >= 150 {
			text = fmt.Sprintf("%s %s %d %s",
				in.NameMention, b.getText(in, "info_max_queue_time"), timedown, b.getText(in, "min"))
		} else {
			text = in.NameMention + b.getText(in, "timer_updated")
			b.storage.DbFunc.UpdateTimedown(in.Lvlkz, in.Config.CorpName, in.UserId)
			//in.O ption.Pl30 = true
			//in.O ption.Edit = true
			in.Opt.Add(models.OptionPl30)
			in.Opt.Add(models.OptionEdit)
			b.QueueLevel(in)
		}
	}
	b.ifTipSendTextDelSecond(in, text, 20)
}
