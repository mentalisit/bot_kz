package bot

import (
	"fmt"
	"kz_bot/models"
	"kz_bot/pkg/utils"
	"strconv"
	"time"
)

//lang ok

func (b *Bot) RsPlus(in models.InMessage) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.iftipdelete(in)
	CountName, err := b.storage.Count.СountName(in.UserId, in.Lvlkz, in.Config.CorpName)
	if err != nil {
		return
	}
	if CountName == 1 { //проверяем есть ли игрок в очереди
		b.ifTipSendMentionText(in, b.getText(in, "you_in_queue"))
	} else {
		countQueue, err1 := b.storage.Count.CountQueue(in.Lvlkz, in.Config.CorpName) //проверяем, есть ли кто-то в очереди
		if err1 != nil {
			return
		}
		numkzN, err2 := b.storage.Count.CountNumberNameActive1(in.Lvlkz, in.Config.CorpName, in.UserId) //проверяем количество боёв по уровню кз игрока
		if err2 != nil {
			return
		}
		numkzL, err3 := b.storage.DbFunc.NumberQueueLvl(in.Lvlkz, in.Config.CorpName) //проверяем какой номер боя определенной красной звезды
		if err3 != nil {
			return
		}

		dsmesid := ""
		tgmesid := 0

		n := b.getMap(in, numkzL)

		var u models.Users
		timekz, _ := strconv.Atoi(in.Timekz)
		UserIn := models.Sborkz{
			Name:     in.Username,
			UserId:   in.UserId,
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
					ch := utils.WaitForMessage("RsPlus73")
					u.User1 = UserIn
					n = b.helpers.GetQueueDiscord(n, u)
					dsmesid = b.client.Ds.SendComplexContent(in.Config.DsChannel,
						fmt.Sprintf(b.getText(in, "temp_queue_started"), in.Username, n["lvlkz"]))
					time.Sleep(1 * time.Second)
					b.client.Ds.EditComplexButton(dsmesid, in.Config.DsChannel, n)
					b.wg.Done()
					close(ch)
				}()
			}
			if in.Config.TgChannel != "" {
				b.wg.Add(1)
				go func() {
					ch := utils.WaitForMessage("RsPlus89")
					u.User1 = UserIn
					texttg = b.helpers.GetQueueTelegram(ntg, u)
					tgmesid = b.client.Tg.SendEmbed(in.Lvlkz, in.Config.TgChannel, texttg)
					b.SubscribePing(in, 1)
					b.wg.Done()
					close(ch)
				}()
			}
			go b.ReadQueueLevel(in)
		}
		u = b.storage.DbFunc.ReadAll(in.Lvlkz, in.Config.CorpName)

		if countQueue == 1 {
			dsmesid = u.User1.Dsmesid

			if in.Config.DsChannel != "" {
				b.wg.Add(1)
				go func() {
					ch := utils.WaitForMessage("RsPlus118")
					u.User2 = UserIn
					n = b.helpers.GetQueueDiscord(n, u)
					text := fmt.Sprintf("%s 2/4 %s %s", n["lvlkz"], in.Username, b.getText(in, "you_joined_queue"))
					//text := n["lvlkz"] + " 2/4 " + in.Username + b.getText(in, "you_joined_queue")
					go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 10)
					b.client.Ds.EditComplexButton(u.User1.Dsmesid, in.Config.DsChannel, n)
					b.wg.Done()
					close(ch)
				}()
			}
			if in.Config.TgChannel != "" {
				b.wg.Add(1)
				go func() {
					ch := utils.WaitForMessage("RsPlus133")
					u.User2 = UserIn
					texttg = b.helpers.GetQueueTelegram(ntg, u)

					tgmesid = b.client.Tg.SendEmbed(in.Lvlkz, in.Config.TgChannel, texttg)
					go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
					err = b.storage.Update.MesidTgUpdate(tgmesid, in.Lvlkz, in.Config.CorpName)
					if err != nil {
						err = b.storage.Update.MesidTgUpdate(tgmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							b.log.ErrorErr(err)
						}
					}
					b.wg.Done()
					close(ch)
				}()
			}
			go b.ReadQueueLevel(in)
		} else if countQueue == 2 {
			dsmesid = u.User1.Dsmesid

			if in.Config.DsChannel != "" {
				b.wg.Add(1)
				go func() {
					ch := utils.WaitForMessage("RsPlus164")
					u.User3 = UserIn
					n = b.helpers.GetQueueDiscord(n, u)
					lvlk3, err4 := b.client.Ds.RoleToIdPing(b.getText(in, "rs")+in.Lvlkz+"+", in.Config.Guildid)
					if err4 != nil {
						b.log.Info(fmt.Sprintf("RoleToIdPing %+v lvl %s", in.Config, in.Lvlkz[1:]))
					}
					text := fmt.Sprintf("%s  3/4 %s %s %s %s",
						n["lvlkz"], in.Username, b.getText(in, "you_joined_queue"), lvlk3, b.getText(in, "another_one_needed_to_complete_queue"))
					go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 10)
					b.client.Ds.EditComplexButton(u.User1.Dsmesid, in.Config.DsChannel, n)
					b.wg.Done()
					close(ch)
				}()
			}
			if in.Config.TgChannel != "" {
				b.wg.Add(1)
				go func() {
					ch := utils.WaitForMessage("RsPlus186")
					u.User3 = UserIn
					texttg = b.helpers.GetQueueTelegram(ntg, u)
					tgmesid = b.client.Tg.SendEmbed(in.Lvlkz, in.Config.TgChannel, texttg)
					go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
					err = b.storage.Update.MesidTgUpdate(tgmesid, in.Lvlkz, in.Config.CorpName)
					if err != nil {
						err = b.storage.Update.MesidTgUpdate(tgmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							b.log.ErrorErr(err)
						}
					}
					b.SubscribePing(in, 3)
					b.wg.Done()
					close(ch)
				}()
			}
		}
		if countQueue <= 2 {
			b.wg.Wait()
			b.storage.DbFunc.InsertQueue(dsmesid, "", in.Config.CorpName, in.Username, in.UserId, in.NameMention, in.Tip, in.Lvlkz, in.Timekz, tgmesid, numkzN)
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
					ch := utils.WaitForMessage("RsPlus235")
					n1, n2, n3, n4 := b.helpers.NameMention(u, ds)
					go b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
					go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel,
						fmt.Sprintf(" 4/4 %s %s", in.Username, b.getText(in, "you_joined_queue")), 10)
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
						dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Ds.Avatar)
					} else {
						dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, "")
					}
					err = b.storage.Update.MesidDsUpdate(dsmesid, in.Lvlkz, in.Config.CorpName)
					if err != nil {
						err = b.storage.Update.MesidDsUpdate(dsmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							b.log.ErrorErr(err)
						}
					}
					b.wg.Done()
					close(ch)
				}()
			}
			if in.Config.TgChannel != "" {
				b.wg.Add(1)
				go func() {
					ch := utils.WaitForMessage("RsPlus273")
					n1, n2, n3, n4 := b.helpers.NameMention(u, tg)
					go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
					go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel,
						in.Username+b.getText(in, "rs_queue_closed")+in.Lvlkz, 10)
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
					err = b.storage.Update.MesidTgUpdate(tgmesid, in.Lvlkz, in.Config.CorpName)
					if err != nil {
						err = b.storage.Update.MesidTgUpdate(tgmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							b.log.ErrorErr(err)
						}
					}
					b.wg.Done()
					close(ch)
				}()
			}

			b.wg.Wait()
			b.storage.DbFunc.InsertQueue(dsmesid, "", in.Config.CorpName, in.Username, in.UserId, in.NameMention, in.Tip, in.Lvlkz, in.Timekz, tgmesid, numkzN)
			err = b.storage.Update.UpdateCompliteRS(in.Lvlkz, dsmesid, tgmesid, "", numkzL, numberevent, in.Config.CorpName)
			if err != nil {
				err = b.storage.Update.UpdateCompliteRS(in.Lvlkz, dsmesid, tgmesid, "", numkzL, numberevent, in.Config.CorpName)
				if err != nil {
					b.log.ErrorErr(err)
				}
			}

			//отправляем сообщение о корпорациях с %
			go b.SendPercent(in.Config)

			//проверка есть ли игрок в других чатах
			user := []string{u.User1.UserId, u.User2.UserId, u.User3.UserId, in.UserId}
			go b.elseChat(user)
			go b.helpers.SaveUsersIdQueue(user, in.Config)

		}

	}
}
func (b *Bot) RsMinus(in models.InMessage) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.iftipdelete(in)

	CountNames, err := b.storage.Count.СountName(in.UserId, in.Lvlkz, in.Config.CorpName) //проверяем есть ли игрок в очереди
	if err != nil {
		b.log.ErrorErr(err)
		return
	}
	if CountNames == 0 {
		b.ifTipSendMentionText(in, b.getText(in, "you_out_of_queue"))
	} else if CountNames > 0 {
		//чтение айди очечреди
		u := b.storage.DbFunc.ReadAll(in.Lvlkz, in.Config.CorpName)
		//удаление с БД
		b.storage.DbFunc.DeleteQueue(in.UserId, in.Lvlkz, in.Config.CorpName)
		//проверяем очередь
		countQueue, err2 := b.storage.Count.CountQueue(in.Lvlkz, in.Config.CorpName)
		if err2 != nil {
			b.log.Error(err2.Error())
			return
		}

		darkStar, lvlkz := containsSymbolD(in.Lvlkz)
		var text string
		if darkStar {
			text = fmt.Sprintf("%s %s.", b.getText(in, "queue_drs")+lvlkz, b.getText(in, "was_deleted"))
		} else {
			text = fmt.Sprintf("%s %s.", b.getText(in, "rs_queue")+lvlkz, b.getText(in, "was_deleted"))
		}

		if in.Config.DsChannel != "" {
			go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, fmt.Sprintf("%s %s", in.Username, b.getText(in, "left_queue")), 10)
			if countQueue == 0 {
				go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 10)
				go b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
			}
		}
		if in.Config.TgChannel != "" {
			go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, fmt.Sprintf("%s %s", in.Username, b.getText(in, "left_queue")), 10)
			if countQueue == 0 {
				go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 10)
				go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
			}
		}
		if countQueue > 0 {

			b.QueueLevel(in)
		}
	}
}
