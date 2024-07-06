package bot

import (
	"context"
	"fmt"
	"kz_bot/models"
	"strings"
	"time"
)

//lang ok

func (b *Bot) RsStart(in models.InMessage) {
	b.mu.Lock()
	defer b.mu.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	b.iftipdelete(in)
	countName, err := b.storage.Count.СountName(ctx, in.UserId, in.Lvlkz, in.Config.CorpName)
	if err != nil {
		return
	}
	if countName == 0 {
		text := b.getText(in, "info_forced_start_available")
		b.ifTipSendTextDelSecond(in, text, 10)
	} else if countName == 1 {
		numberkz, err1 := b.storage.DbFunc.NumberQueueLvl(ctx, in.Lvlkz, in.Config.CorpName)
		if err1 != nil {
			return
		}
		count, err2 := b.storage.Count.CountQueue(ctx, in.Lvlkz, in.Config.CorpName)
		if err2 != nil {
			return
		}

		dsmesid := ""
		tgmesid := 0
		if count > 0 {
			u := b.storage.DbFunc.ReadAll(ctx, in.Lvlkz, in.Config.CorpName)
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
						name1, _, _, _ := b.helpers.NameMention(u, ds)
						text := fmt.Sprintf("🚀 %s%s (%d) %s \n\n1. %s\n%s %s",
							b.getText(in, "rs_queue"), in.Lvlkz, numberkz,
							b.getText(in, "was_launched_incomplete"), name1, b.getText(in, "go"), textEvent)

						if in.Tip == ds {
							dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Config.Guildid, in.Ds.Avatar)

						} else {
							dsmesid = b.client.Ds.Send(in.Config.DsChannel, text)
						}

						go b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
						err = b.storage.Update.MesidDsUpdate(ctx, dsmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							err = b.storage.Update.MesidDsUpdate(context.Background(), dsmesid, in.Lvlkz, in.Config.CorpName)
							if err != nil {
								b.log.ErrorErr(err)
							}
						}
						b.wg.Done()
					}()

				}
				if in.Config.TgChannel != "" {
					b.wg.Add(1)
					go func() {
						name1, _, _, _ := b.helpers.NameMention(u, tg)
						go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
						text := fmt.Sprintf("🚀 %s%s (%d) %s \n\n1. %s\n%s %s",
							b.getText(in, "rs_queue"), in.Lvlkz, numberkz,
							b.getText(in, "was_launched_incomplete"), name1, b.getText(in, "go"), textEvent)
						tgmesid = b.client.Tg.SendChannel(in.Config.TgChannel, text)
						err = b.storage.Update.MesidTgUpdate(ctx, tgmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							err = b.storage.Update.MesidTgUpdate(context.Background(), tgmesid, in.Lvlkz, in.Config.CorpName)
							if err != nil {
								b.log.ErrorErr(err)
							}
						}
						b.wg.Done()
					}()

				}
			} else if count == 2 {
				if in.Config.DsChannel != "" { //discord
					b.wg.Add(1)
					go func() {
						name1, name2, _, _ := b.helpers.NameMention(u, ds)
						text1 := fmt.Sprintf("🚀 %s%s (%d) %s \n\n",
							b.getText(in, "rs_queue"), in.Lvlkz, numberkz, b.getText(in, "was_launched_incomplete"))
						text2 := fmt.Sprintf("%s\n%s\n%s %s", name1, name2, b.getText(in, "go"), textEvent)
						text := text1 + text2
						if in.Tip == ds {
							dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Config.Guildid, in.Ds.Avatar)
							if u.User1.Tip == ds {
								//go b.sendDmDark(text, u.User1.Mention)
								go b.client.Ds.SendDmText(text, u.User1.UserId)
							}
						} else {
							dsmesid = b.client.Ds.Send(in.Config.DsChannel, text)
						}
						go b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
						err = b.storage.Update.MesidDsUpdate(ctx, dsmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							err = b.storage.Update.MesidDsUpdate(context.Background(), dsmesid, in.Lvlkz, in.Config.CorpName)
							if err != nil {
								b.log.ErrorErr(err)
							}
						}
						b.wg.Done()
					}()

				}
				if in.Config.TgChannel != "" { //telegram
					b.wg.Add(1)
					go func() {
						name1, name2, _, _ := b.helpers.NameMention(u, tg)
						go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
						text1 := fmt.Sprintf("🚀 %s%s (%d) %s \n\n",
							b.getText(in, "rs_queue"), in.Lvlkz, numberkz, b.getText(in, "was_launched_incomplete"))
						text2 := fmt.Sprintf("%s\n%s\n%s %s", name1, name2, b.getText(in, "go"), textEvent)
						text := text1 + text2
						tgmesid = b.client.Tg.SendChannel(in.Config.TgChannel, text)
						err = b.storage.Update.MesidTgUpdate(ctx, tgmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							err = b.storage.Update.MesidTgUpdate(context.Background(), tgmesid, in.Lvlkz, in.Config.CorpName)
							if err != nil {
								b.log.ErrorErr(err)
							}
						}
						b.wg.Done()
					}()

				}
			} else if count == 3 {
				if in.Config.DsChannel != "" { //discord
					b.wg.Add(1)
					go func() {
						name1, name2, name3, _ := b.helpers.NameMention(u, ds)
						text := fmt.Sprintf("🚀 %s%s (%d) %s \n\n%s\n%s\n%s\n%s %s",
							b.getText(in, "rs_queue"), in.Lvlkz, numberkz, b.getText(in, "was_launched_incomplete"),
							name1, name2, name3, b.getText(in, "go"), textEvent)
						if in.Tip == ds {
							dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Config.Guildid, in.Ds.Avatar)
						} else {
							dsmesid = b.client.Ds.Send(in.Config.DsChannel, text)
						}
						go b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
						err = b.storage.Update.MesidDsUpdate(ctx, dsmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							err = b.storage.Update.MesidDsUpdate(context.Background(), dsmesid, in.Lvlkz, in.Config.CorpName)
							if err != nil {
								b.log.ErrorErr(err)
							}
						}
						b.wg.Done()
					}()
				}
				if in.Config.TgChannel != "" { //telegram
					b.wg.Add(1)
					go func() {
						name1, name2, name3, _ := b.helpers.NameMention(u, tg)
						go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
						text := fmt.Sprintf("🚀 %s%s (%d) %s \n\n%s\n%s\n%s\n%s %s",
							b.getText(in, "rs_queue"), in.Lvlkz, numberkz, b.getText(in, "was_launched_incomplete"),
							name1, name2, name3, b.getText(in, "go"), textEvent)
						tgmesid = b.client.Tg.SendChannel(in.Config.TgChannel, text)
						err = b.storage.Update.MesidTgUpdate(ctx, tgmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							err = b.storage.Update.MesidTgUpdate(context.Background(), tgmesid, in.Lvlkz, in.Config.CorpName)
							if err != nil {
								b.log.ErrorErr(err)
							}
						}
						b.wg.Done()
					}()

				}
			}
			b.wg.Wait()
			err = b.storage.Update.UpdateCompliteRS(ctx, in.Lvlkz, dsmesid, tgmesid, "", numberkz, numberevent, in.Config.CorpName)
			if err != nil {
				err = b.storage.Update.UpdateCompliteRS(context.Background(), in.Lvlkz, dsmesid, tgmesid, "", numberkz, numberevent, in.Config.CorpName)
				if err != nil {
					b.log.ErrorErr(err)
				}
			}

			//отправляем сообщение о корпорациях с %
			go b.SendPercent(in.Config)

			user := []string{u.User1.UserId, u.User2.UserId, u.User3.UserId, in.UserId}
			b.elseChat(user)
		}
	}
}
func (b *Bot) Pl30(in models.InMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	countName := b.storage.Count.CountNameQueue(ctx, in.UserId)
	text := ""
	if countName == 0 {
		text = in.NameMention + b.getText(in, "you_out_of_queue")
	} else if countName > 0 {
		timedown := b.storage.DbFunc.P30Pl(ctx, in.Lvlkz, in.Config.CorpName, in.UserId)
		if timedown >= 150 {
			text = fmt.Sprintf("%s %s %d %s",
				in.NameMention, b.getText(in, "info_max_queue_time"), timedown, b.getText(in, "min"))
		} else {
			text = in.NameMention + b.getText(in, "timer_updated")
			b.storage.DbFunc.UpdateTimedown(ctx, in.Lvlkz, in.Config.CorpName, in.UserId)
			in.Option.Pl30 = true
			in.Option.Edit = true
			b.QueueLevel(in)
		}
	}
	b.ifTipSendTextDelSecond(in, text, 20)
}
