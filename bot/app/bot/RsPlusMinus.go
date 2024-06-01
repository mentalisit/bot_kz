package bot

import (
	"context"
	"fmt"
	"kz_bot/models"
	"strconv"
	"time"
)

//lang ok

func (b *Bot) RsPlus(in models.InMessage) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.iftipdelete(in)
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	CountName, err := b.storage.Count.СountName(ctx, in.Name, in.Lvlkz, in.Config.CorpName)
	if err != nil {
		return
	}
	if CountName == 1 { //проверяем есть ли игрок в очереди
		b.ifTipSendMentionText(in, b.getText(in, "you_in_queue"))
	} else {
		countQueue, err1 := b.storage.Count.CountQueue(ctx, in.Lvlkz, in.Config.CorpName) //проверяем, есть ли кто-то в очереди
		if err1 != nil {
			return
		}
		numkzN, err2 := b.storage.Count.CountNumberNameActive1(ctx, in.Lvlkz, in.Config.CorpName, in.Name) //проверяем количество боёв по уровню кз игрока
		if err2 != nil {
			return
		}
		numkzL, err3 := b.storage.DbFunc.NumberQueueLvl(ctx, in.Lvlkz, in.Config.CorpName) //проверяем какой номер боя определенной красной звезды
		if err3 != nil {
			return
		}

		dsmesid := ""
		tgmesid := 0
		var n map[string]string
		n = make(map[string]string)
		n["lang"] = in.Config.Country
		if in.Config.DsChannel != "" {
			n["lvlkz"], err = b.client.Ds.RoleToIdPing(b.getText(in, "rs")+in.Lvlkz, in.Config.Guildid)
			if err != nil {
				b.log.Info(fmt.Sprintf("RoleToIdPing %+v lvl %s", in.Config, in.Lvlkz[1:]))
			}
		}
		var u models.Users
		timekz, _ := strconv.Atoi(in.Timekz)
		UserIn := models.Sborkz{
			Name:     in.Name,
			Mention:  in.NameMention,
			Numkzn:   numkzN,
			Timedown: timekz,
		}

		texttg := ""
		ntg := make(map[string]string)
		if in.Config.TgChannel != "" {
			ntg["text1"] = fmt.Sprintf("%s%s (%d)\n", b.getText(in, "rs_queue"), in.Lvlkz, numkzN)
			ntg["text2"] = fmt.Sprintf("\n%s++ - %s", in.Lvlkz, b.getText(in, "forced_start"))
			ntg["min"] = b.getText(in, "min")
		}

		if countQueue == 0 {
			if in.Config.DsChannel != "" {
				b.wg.Add(1)
				go func() {
					u.User1 = UserIn
					n = b.helpers.GetQueueDiscord(n, u)
					//n["name1"] = fmt.Sprintf("%s  🕒  %s  (%d)", b.emReadName(in, in.Name, in.NameMention, ds), in.Timekz, numkzN)
					emb := b.client.Ds.EmbedDS(n, numkzL, 1, false)
					dsmesid = b.client.Ds.SendComplexContent(in.Config.DsChannel,
						fmt.Sprintf(b.getText(in, "temp_queue_started"), in.Name, n["lvlkz"]))
					time.Sleep(1 * time.Second)
					b.client.Ds.EditComplexButton(dsmesid, in.Config.DsChannel, emb, b.client.Ds.AddButtonsQueue(in.Lvlkz))
					b.wg.Done()
				}()
			}
			if in.Config.TgChannel != "" {
				b.wg.Add(1)
				go func() {
					//text := fmt.Sprintf("%s%s (%d)\n"+
					//	"1️⃣ %s - %s%s (%d) \n\n"+
					//	"%s++ - %s",
					//	b.getText(in, "rs_queue"), in.Lvlkz, numkzL,
					//	b.emReadName(in.Name, in.NameMention, tg), in.Timekz, b.getText(in, "min"), numkzN,
					//	in.Lvlkz, b.getText(in, "forced_start"))
					//text := fmt.Sprintf(b.getText(in, "temp1_queue"),
					//	in.Lvlkz, numkzL,
					//	b.emReadName(in, in.Name, in.NameMention, tg), in.Timekz, numkzN,
					//	in.Lvlkz)
					u.User1 = UserIn
					texttg = b.helpers.GetQueueTelegram(ntg, u)
					tgmesid = b.client.Tg.SendEmded(in.Lvlkz, in.Config.TgChannel, texttg)
					b.SubscribePing(in, 1)
					b.wg.Done()
				}()
			}
		}
		u = b.storage.DbFunc.ReadAll(ctx, in.Lvlkz, in.Config.CorpName)

		if countQueue == 1 {
			dsmesid = u.User1.Dsmesid

			if in.Config.DsChannel != "" {
				b.wg.Add(1)
				go func() {
					u.User2 = UserIn
					n = b.helpers.GetQueueDiscord(n, u)
					//n["name1"] = fmt.Sprintf("%s  🕒  %d  (%d)", b.emReadName(in, u.User1.Name, u.User1.Mention, ds), u.User1.Timedown, u.User1.Numkzn)
					//n["name2"] = fmt.Sprintf("%s  🕒  %s  (%d)", b.emReadName(in, in.Name, in.NameMention, ds), in.Timekz, numkzN)
					emb := b.client.Ds.EmbedDS(n, numkzL, 2, false)
					text := fmt.Sprintf("%s 2/4 %s %s", n["lvlkz"], in.Name, b.getText(in, "you_joined_queue"))
					//text := n["lvlkz"] + " 2/4 " + in.Name + b.getText(in, "you_joined_queue")
					go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 10)
					b.client.Ds.EditComplexButton(u.User1.Dsmesid, in.Config.DsChannel, emb, b.client.Ds.AddButtonsQueue(in.Lvlkz))
					b.wg.Done()
				}()
			}
			if in.Config.TgChannel != "" {
				b.wg.Add(1)
				go func() {
					//text1 := fmt.Sprintf("%s%s (%d)\n", b.getText(in, "rs_queue"), in.Lvlkz, numkzL)
					//name1 := fmt.Sprintf("1️⃣ %s - %d%s (%d) \n",
					//	b.emReadName(in, u.User1.Name, u.User1.Mention, tg), u.User1.Timedown, b.getText(in, "min"), u.User1.Numkzn)
					//name2 := fmt.Sprintf("2️⃣ %s - %s%s (%d) \n",
					//	b.emReadName(in, in.Name, in.NameMention, tg), in.Timekz, b.getText(in, "min"), numkzN)
					//text2 := fmt.Sprintf("\n%s++ - %s", in.Lvlkz, b.getText(in, "forced_start"))
					//text := fmt.Sprintf("%s %s %s %s", text1, name1, name2, text2)
					u.User2 = UserIn
					texttg = b.helpers.GetQueueTelegram(ntg, u)

					tgmesid = b.client.Tg.SendEmded(in.Lvlkz, in.Config.TgChannel, texttg)
					go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
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
		} else if countQueue == 2 {
			dsmesid = u.User1.Dsmesid

			if in.Config.DsChannel != "" {
				b.wg.Add(1)
				go func() {
					u.User3 = UserIn
					n = b.helpers.GetQueueDiscord(n, u)
					//n["name1"] = fmt.Sprintf("%s  🕒  %d  (%d)", b.emReadName(in, u.User1.Name, u.User1.Mention, in.Tip), u.User1.Timedown, u.User1.Numkzn)
					//n["name2"] = fmt.Sprintf("%s  🕒  %d  (%d)", b.emReadName(in, u.User2.Name, u.User2.Mention, in.Tip), u.User2.Timedown, u.User2.Numkzn)
					//n["name3"] = fmt.Sprintf("%s  🕒  %s  (%d)", b.emReadName(in, in.Name, in.NameMention, in.Tip), in.Timekz, numkzN)
					lvlk3, err4 := b.client.Ds.RoleToIdPing(b.getText(in, "rs")+in.Lvlkz+"+", in.Config.Guildid)
					if err4 != nil {
						b.log.Info(fmt.Sprintf("RoleToIdPing %+v lvl %s", in.Config, in.Lvlkz[1:]))
					}
					emb := b.client.Ds.EmbedDS(n, numkzL, 3, false)
					text := fmt.Sprintf("%s  3/4 %s %s %s %s",
						n["lvlkz"], in.Name, b.getText(in, "you_joined_queue"), lvlk3, b.getText(in, "another_one_needed_to_complete_queue"))
					go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 10)
					b.client.Ds.EditComplexButton(u.User1.Dsmesid, in.Config.DsChannel, emb, b.client.Ds.AddButtonsQueue(in.Lvlkz))
					b.wg.Done()
				}()
			}
			if in.Config.TgChannel != "" {
				b.wg.Add(1)
				go func() {
					//text1 := fmt.Sprintf("%s%s (%d)\n", b.getText(in, "rs_queue"), in.Lvlkz, numkzL)
					//name1 := fmt.Sprintf("1️⃣ %s - %d%s (%d) \n",
					//	b.emReadName(in, u.User1.Name, u.User1.Mention, tg), u.User1.Timedown, b.getText(in, "min"), u.User1.Numkzn)
					//name2 := fmt.Sprintf("2️⃣ %s - %d%s (%d) \n",
					//	b.emReadName(in, u.User2.Name, u.User2.Mention, tg), u.User2.Timedown, b.getText(in, "min"), u.User2.Numkzn)
					//name3 := fmt.Sprintf("3️⃣ %s - %s%s (%d) \n",
					//	b.emReadName(in, in.Name, in.NameMention, tg), in.Timekz, b.getText(in, "min"), numkzN)
					//text2 := fmt.Sprintf("\n%s++ - %s", in.Lvlkz, b.getText(in, "forced_start"))
					//text := fmt.Sprintf("%s %s %s %s %s", text1, name1, name2, name3, text2)

					u.User3 = UserIn
					texttg = b.helpers.GetQueueTelegram(ntg, u)
					tgmesid = b.client.Tg.SendEmded(in.Lvlkz, in.Config.TgChannel, texttg)
					go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
					err = b.storage.Update.MesidTgUpdate(ctx, tgmesid, in.Lvlkz, in.Config.CorpName)
					if err != nil {
						err = b.storage.Update.MesidTgUpdate(context.Background(), tgmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							b.log.ErrorErr(err)
						}
					}
					b.SubscribePing(in, 3)
					b.wg.Done()
				}()
			}
		}
		if countQueue <= 2 {
			b.wg.Wait()
			b.storage.DbFunc.InsertQueue(ctx, dsmesid, "", in.Config.CorpName, in.Name, in.NameMention, in.Tip, in.Lvlkz, in.Timekz, tgmesid, numkzN)
		}

		if countQueue == 3 {
			u.User4 = UserIn
			dsmesid = u.User1.Dsmesid

			textEvent, numkzEvent := b.EventText(in)
			if textEvent == "" {
				textEvent = b.GetTextPercent(in.Config, false)
			}
			numberevent := b.storage.Event.NumActiveEvent(in.Config.CorpName) //получаем номер ивета если он активен
			if numberevent > 0 {
				numkzL = numkzEvent
			}

			if in.Config.DsChannel != "" {
				b.wg.Add(1)
				go func() {
					n1, n2, n3, n4 := b.helpers.NameMention(in, u, ds)
					go b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
					go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel,
						fmt.Sprintf(" 4/4 %s %s", in.Name, b.getText(in, "you_joined_queue")), 10)
					//" 4/4 "+in.Name+" "+b.getText(in, "you_joined_queue"), 10)
					text := fmt.Sprintf("4/4 %s%s %s\n"+
						" %s\n"+
						" %s\n"+
						" %s\n"+
						" %s\n"+
						"%s %s",
						b.getText(in, "rs_queue"), in.Lvlkz, b.getText(in, "queue_completed"),
						n1,
						n2,
						n3,
						n4,
						b.getText(in, "go"), textEvent)

					if in.Tip == ds {
						dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Config.Guildid, in.Ds.Avatar)
					} else {
						dsmesid = b.client.Ds.Send(in.Config.DsChannel, text)
					}
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
					n1, n2, n3, n4 := b.helpers.NameMention(in, u, tg)
					go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
					go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel,
						in.Name+b.getText(in, "rs_queue_closed")+in.Lvlkz, 10)
					text := fmt.Sprintf("%s%s %s\n"+
						"%s\n"+
						"%s\n"+
						"%s\n"+
						"%s\n"+
						" %s \n"+
						"%s",
						b.getText(in, "rs_queue"), in.Lvlkz, b.getText(in, "queue_completed"),
						n1, n2, n3, n4,
						b.getText(in, "go"), textEvent)
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

			b.wg.Wait()
			b.storage.DbFunc.InsertQueue(ctx, dsmesid, "", in.Config.CorpName, in.Name, in.NameMention, in.Tip, in.Lvlkz, in.Timekz, tgmesid, numkzN)
			err = b.storage.Update.UpdateCompliteRS(ctx, in.Lvlkz, dsmesid, tgmesid, "", numkzL, numberevent, in.Config.CorpName)
			if err != nil {
				err = b.storage.Update.UpdateCompliteRS(context.Background(), in.Lvlkz, dsmesid, tgmesid, "", numkzL, numberevent, in.Config.CorpName)
				if err != nil {
					b.log.ErrorErr(err)
				}
			}

			//отправляем сообщение о корпорациях с %
			go b.SendPercent(in.Config)

			//проверка есть ли игрок в других чатах
			user := []string{u.User1.Name, u.User2.Name, u.User3.Name, in.Name}
			go b.elseChat(user)

		}

	}
}
func (b *Bot) RsMinus(in models.InMessage) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.iftipdelete(in)
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	CountNames, err := b.storage.Count.СountName(ctx, in.Name, in.Lvlkz, in.Config.CorpName) //проверяем есть ли игрок в очереди
	if err != nil {
		b.log.ErrorErr(err)
		return
	}
	if CountNames == 0 {
		b.ifTipSendMentionText(in, b.getText(in, "you_out_of_queue"))
	} else if CountNames > 0 {
		//чтение айди очечреди
		u := b.storage.DbFunc.ReadAll(ctx, in.Lvlkz, in.Config.CorpName)
		//удаление с БД
		b.storage.DbFunc.DeleteQueue(ctx, in.Name, in.Lvlkz, in.Config.CorpName)
		//проверяем очередь
		countQueue, err2 := b.storage.Count.CountQueue(ctx, in.Lvlkz, in.Config.CorpName)
		if err2 != nil {
			b.log.Error(err2.Error())
			return
		}

		if in.Config.DsChannel != "" {
			go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, fmt.Sprintf("%s %s", in.Name, b.getText(in, "left_queue")), 10)
			if countQueue == 0 {
				go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel,
					fmt.Sprintf("%s%s %s.", b.getText(in, "rs_queue"), in.Lvlkz, b.getText(in, "was_deleted")), 10)
				go b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
			}
		}
		if in.Config.TgChannel != "" {
			go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, fmt.Sprintf("%s %s", in.Name, b.getText(in, "left_queue")), 10)
			if countQueue == 0 {
				go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel,
					fmt.Sprintf("%s%s %s.", b.getText(in, "rs_queue"), in.Lvlkz, b.getText(in, "was_deleted")), 10)
				go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
			}
		}
		if countQueue > 0 {

			b.QueueLevel(in)
		}
	}
}
