package bot

import (
	"fmt"
	"rs/bot/helpers"
	"rs/models"
	"rs/pkg/utils"
	"strconv"
	"time"
)

const dark = "d"

func (b *Bot) darkAlt(in models.InMessage, i int) {
	if in.Tip == ds {
		alts := helpers.Get2AltsUserId(in.UserId)
		alt := ""
		lenAlts := len(alts)
		if lenAlts > 0 {
			if lenAlts == 1 || i == 1 {
				alt = alts[0]
			} else if i > 1 {
				i = i - 1
				if lenAlts > i {
					alt = alts[i]
				}
			}
		}
		b.RsDarkPlus(in, alt)
	}
}

func (b *Bot) RsDarkPlusOld(in models.InMessage, alt string) {
	b.helpers.ReadNameModules(in, alt)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.iftipdelete(in)
	CountName, err := b.storage.Count.СountName(in.UserId, in.Lvlkz, in.Config.CorpName)
	if err != nil {
		b.log.ErrorErr(err)
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

		u := models.Users{}
		timekz, _ := strconv.Atoi(in.Timekz)
		UserIn := models.Sborkz{
			Tip:      in.Tip,
			Name:     in.Username,
			UserId:   in.UserId,
			Mention:  in.NameMention,
			Numkzn:   numkzN,
			Timedown: timekz,
			Wamesid:  alt,
		}

		texttg := ""
		ntg := make(map[string]string)
		if in.Config.TgChannel != "" {
			ntg["text1"] = fmt.Sprintf("%s%s (%d)\n", b.getText(in, "queue_drs"), in.Lvlkz[1:], numkzL)
			ntg["text2"] = fmt.Sprintf("\n%s++ - %s", in.Lvlkz, b.getText(in, "forced_start"))
			ntg["min"] = b.getText(in, "min")
		}

		if countQueue == 0 {
			if in.Config.DsChannel != "" {
				b.wg.Add(1)
				go func() {
					ch := utils.WaitForMessage("RsDarkPlus179")
					u.User1 = UserIn
					n = b.helpers.GetQueueDiscord(n, u)
					dsmesid = b.client.Ds.SendComplexContent(in.Config.DsChannel,
						fmt.Sprintf(b.getText(in, "temp_queue_started"), in.Username, n["lvlkz"]))
					time.Sleep(1 * time.Second)
					err = b.client.Ds.EditComplexButton(dsmesid, in.Config.DsChannel, n)
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
					ch := utils.WaitForMessage("RsDarkPlus195")
					u.User1 = UserIn
					texttg = b.helpers.GetQueueTelegram(ntg, u)

					tgmesid = b.client.Tg.SendEmbed(in.Lvlkz, in.Config.TgChannel, texttg)
					b.SubscribePing(in, 1)
					b.wg.Done()
					close(ch)
				}()
			}
			b.wg.Wait()
			b.storage.DbFunc.InsertQueue(dsmesid, alt, in.Config.CorpName, in.Username, in.UserId, in.NameMention, in.Tip, in.Lvlkz, in.Timekz, tgmesid, numkzN)
			//b.ReadQueueLevel(in)
			return
		} else {
			u = b.storage.DbFunc.ReadAll(in.Lvlkz, in.Config.CorpName)
			if u.User1.Dsmesid != "" {
				dsmesid = u.User1.Dsmesid
			}
			if u.User1.Tgmesid != 0 {
				tgmesid = u.User1.Tgmesid
			}

			if countQueue == 1 {

				if in.Config.DsChannel != "" {
					b.wg.Add(1)
					go func() {
						ch := utils.WaitForMessage("RsDarkPlus223")
						u.User2 = &UserIn
						n = b.helpers.GetQueueDiscord(n, u)
						text := n["lvlkz"] + fmt.Sprintf(" 2/3 %s %s \n%s", in.Username, b.getText(in, "you_joined_queue"), u.User1.Mention)
						go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 10)
						err = b.client.Ds.EditComplexButton(dsmesid, in.Config.DsChannel, n)
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
						ch := utils.WaitForMessage("RsDarkPlus239")
						u.User2 = &UserIn
						texttg = b.helpers.GetQueueTelegram(ntg, u)

						tgmesid = b.client.Tg.SendEmbed(in.Lvlkz, in.Config.TgChannel, texttg)
						go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
						err = b.storage.Update.MesidTgUpdate(tgmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							b.log.ErrorErr(err)
						}
						b.wg.Done()
						close(ch)
					}()
				}
				//go b.ReadQueueLevel(in)
				b.wg.Wait()
				b.storage.DbFunc.InsertQueue(dsmesid, alt, in.Config.CorpName, in.Username, in.UserId, in.NameMention, in.Tip, in.Lvlkz, in.Timekz, tgmesid, numkzN)
				return
			}

			if countQueue == 2 {
				textEvent, numkzEvent := b.EventText(in)
				if textEvent == "" {
					textEvent = b.GetTextPercent(in.Config, true)
				}
				numberevent := b.storage.Event.NumActiveEvent(in.Config.CorpName) //получаем номер ивета если он активен
				if numberevent > 0 {
					numkzL = numkzEvent
				}
				u.User3 = &UserIn

				if in.Config.DsChannel != "" {
					b.wg.Add(1)
					go func() {
						ch := utils.WaitForMessage("RsDarkPlus287")
						n1, n2, n3, _ := b.helpers.NameMention(u, ds)
						go b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
						go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel,
							fmt.Sprintf("🚀 3/3 %s %s", in.Username, b.getText(in, "you_joined_queue")), 10)
						text := fmt.Sprintf("3/3 %s%s %s\n %s\n %s\n %s\n%s %s",
							b.getText(in, "queue_drs"), in.Lvlkz[1:], b.getText(in, "queue_completed"),
							n1, n2, n3, b.getText(in, "go"), textEvent)

						dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Ds.Avatar)

						err = b.storage.Update.MesidDsUpdate(dsmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							b.log.Error("MesidDsUpdate " + err.Error())
							err = b.storage.Update.MesidDsUpdate(dsmesid, in.Lvlkz, in.Config.CorpName)
							if err != nil {
								b.log.ErrorErr(err)
								b.log.Warn("this problem")
							}

						}
						b.wg.Done()
						close(ch)
					}()

				}
				if in.Config.TgChannel != "" {
					b.wg.Add(1)
					go func() {
						ch := utils.WaitForMessage("RsDarkPlus330")
						n1, n2, n3, _ := b.helpers.NameMention(u, tg)
						go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
						go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel,
							in.Username+b.getText(in, "drs_queue_closed")+in.Lvlkz[1:], 10)
						text := fmt.Sprintf("🚀 %s%s %s\n"+
							"%s\n%s\n%s\n %s \n%s",
							b.getText(in, "queue_drs"), in.Lvlkz[1:], b.getText(in, "queue_completed"),
							n1, n2, n3, b.getText(in, "go"), textEvent)
						fmt.Printf("MesidTgUpdate tgmesid %d, in.Lvlkz %s, in.Config.CorpName %s\n", tgmesid, in.Lvlkz, in.Config.CorpName)
						tgmesid = b.client.Tg.SendChannel(in.Config.TgChannel, text)
						err = b.storage.Update.MesidTgUpdate(tgmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							b.log.Error("MesidTgUpdate " + err.Error())
							err = b.storage.Update.MesidTgUpdate(tgmesid, in.Lvlkz, in.Config.CorpName)
							if err != nil {
								b.log.ErrorErr(err)
								b.log.Warn("this problem")
							}
						}
						b.wg.Done()
						close(ch)
					}()
				}

				b.wg.Wait()
				go b.SendLsNotification(in, u)
				b.storage.DbFunc.InsertQueue(dsmesid, alt, in.Config.CorpName, in.Username, in.UserId, in.NameMention, in.Tip, in.Lvlkz, in.Timekz, tgmesid, numkzN)
				err = b.storage.Update.UpdateCompliteRS(in.Lvlkz, dsmesid, tgmesid, alt, numkzL, numberevent, in.Config.CorpName)
				if err != nil {
					b.log.Error("UpdateCompliteRS " + err.Error())
					err = b.storage.Update.UpdateCompliteRS(in.Lvlkz, dsmesid, tgmesid, alt, numkzL, numberevent, in.Config.CorpName)
					if err != nil {
						b.log.ErrorErr(err)
					}
				}

				if numkzEvent == 0 {
					//отправляем сообщение о корпорациях с %
					go b.SendPercent(in.Config)
				}
				//проверка есть ли игрок в других чатах
				user := u.GetAllUserId()
				go b.elseChat(user)
				go b.helpers.SaveUsersIdQueue(user, in.Config)
			}
			go b.CheckSubscribe(in)
		}
	}
}
func (b *Bot) RsDarkPlus(in models.InMessage, alt string) {
	b.helpers.ReadNameModules(in, alt)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.iftipdelete(in)
	CountName, err := b.storage.Count.СountName(in.UserId, in.Lvlkz, in.Config.CorpName)
	if err != nil {
		b.log.ErrorErr(err)
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

		u := models.Users{}
		timekz, _ := strconv.Atoi(in.Timekz)
		UserIn := models.Sborkz{
			Tip:      in.Tip,
			Name:     in.Username,
			UserId:   in.UserId,
			Mention:  in.NameMention,
			Numkzn:   numkzN,
			Timedown: timekz,
			Wamesid:  alt,
		}

		//texttg := ""
		ntg := make(map[string]string)
		if in.Config.TgChannel != "" {
			ntg["text1"] = fmt.Sprintf("%s%s (%d)\n", b.getText(in, "queue_drs"), in.Lvlkz[1:], numkzL)
			ntg["text2"] = fmt.Sprintf("\n%s++ - %s", in.Lvlkz, b.getText(in, "forced_start"))
			ntg["min"] = b.getText(in, "min")
		}

		_, level := containsSymbolD(in.Lvlkz)

		if countQueue == 0 {
			if in.Config.DsChannel != "" {
				b.wg.Add(1)
				go func() {
					ch := utils.WaitForMessage("RsDarkPlus179")
					dsmesid = b.client.Ds.SendComplexContent(in.Config.DsChannel,
						fmt.Sprintf(b.getText(in, "temp_queue_started"), in.Username, n["lvlkz"]))
					b.wg.Done()
					close(ch)
				}()
			}
			if in.Config.TgChannel != "" {
				b.wg.Add(1)
				go func() {
					ch := utils.WaitForMessage("RsDarkPlus195")
					tgmesid = b.client.Tg.SendChannel(in.Config.TgChannel,
						fmt.Sprintf(b.getText(in, "temp_queue_started"),
							in.Username, b.getText(in, "drs")+level))

					b.SubscribePing(in, 1)
					b.wg.Done()
					close(ch)
				}()
			}
			b.wg.Wait()
			b.storage.DbFunc.InsertQueue(dsmesid, alt, in.Config.CorpName, in.Username, in.UserId, in.NameMention, in.Tip, in.Lvlkz, in.Timekz, tgmesid, numkzN)
			b.QueueLevel(in)
			b.ReadQueueLevel(in)
			return
		} else {
			u = b.storage.DbFunc.ReadAll(in.Lvlkz, in.Config.CorpName)
			if u.User1.Dsmesid != "" {
				dsmesid = u.User1.Dsmesid
			}
			if u.User1.Tgmesid != 0 {
				tgmesid = u.User1.Tgmesid
			}

			if countQueue == 1 {

				if in.Config.DsChannel != "" {
					b.wg.Add(1)
					go func() {
						ch := utils.WaitForMessage("RsDarkPlus223")
						//u.User2 = UserIn
						//n = b.helpers.GetQueueDiscord(n, u)
						text := n["lvlkz"] + fmt.Sprintf(" 2/3 %s %s \n%s",
							in.Username, b.getText(in, "you_joined_queue"), u.User1.Mention)
						go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 10)
						//err = b.client.Ds.EditComplexButton(dsmesid, in.Config.DsChannel, n)
						//if err != nil {
						//	b.log.ErrorErr(err)
						//}
						b.wg.Done()
						close(ch)
					}()
				}
				if in.Config.TgChannel != "" {
					b.wg.Add(1)
					go func() {
						ch := utils.WaitForMessage("RsDarkPlus239")
						//u.User2 = UserIn
						//texttg = b.helpers.GetQueueTelegram(ntg, u)

						text := fmt.Sprintf("%s%s 2/3 %s %s \n%s", b.getText(in, "drs"), level,
							in.Username, b.getText(in, "you_joined_queue"), u.User1.Mention)
						go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 10)

						//tgmesid = b.client.Tg.SendEmbed(in.Lvlkz, in.Config.TgChannel, texttg)
						//go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
						//err = b.storage.Update.MesidTgUpdate(tgmesid, in.Lvlkz, in.Config.CorpName)
						//if err != nil {
						//	b.log.ErrorErr(err)
						//}
						b.wg.Done()
						close(ch)
					}()
				}
				b.wg.Wait()
				b.storage.DbFunc.InsertQueue(dsmesid, alt, in.Config.CorpName, in.Username, in.UserId, in.NameMention, in.Tip, in.Lvlkz, in.Timekz, tgmesid, numkzN)
				b.QueueLevel(in)
				go b.ReadQueueLevel(in)
				return
			}

			if countQueue == 2 {
				textEvent, numkzEvent := b.EventText(in)
				if textEvent == "" {
					textEvent = b.GetTextPercent(in.Config, true)
				}
				numberevent := b.storage.Event.NumActiveEvent(in.Config.CorpName) //получаем номер ивета если он активен
				if numberevent > 0 {
					numkzL = numkzEvent
				}
				u.User3 = &UserIn

				if in.Config.DsChannel != "" {
					b.wg.Add(1)
					go func() {
						ch := utils.WaitForMessage("RsDarkPlus287")
						n1, n2, n3, _ := b.helpers.NameMention(u, ds)
						go b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
						go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel,
							fmt.Sprintf("🚀 3/3 %s %s", in.Username, b.getText(in, "you_joined_queue")), 10)
						text := fmt.Sprintf("3/3 %s%s %s\n %s\n %s\n %s\n%s %s",
							b.getText(in, "queue_drs"), level, b.getText(in, "queue_completed"),
							n1, n2, n3, b.getText(in, "go"), textEvent)

						dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Ds.Avatar)

						err = b.storage.Update.MesidDsUpdate(dsmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							b.log.Error("MesidDsUpdate " + err.Error())
							err = b.storage.Update.MesidDsUpdate(dsmesid, in.Lvlkz, in.Config.CorpName)
							if err != nil {
								b.log.ErrorErr(err)
								b.log.Warn("this problem")
							}

						}
						b.wg.Done()
						close(ch)
					}()

				}
				if in.Config.TgChannel != "" {
					b.wg.Add(1)
					go func() {
						ch := utils.WaitForMessage("RsDarkPlus330")
						n1, n2, n3, _ := b.helpers.NameMention(u, tg)
						go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
						go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel,
							in.Username+b.getText(in, "drs_queue_closed")+in.Lvlkz[1:], 10)
						text := fmt.Sprintf("🚀 %s%s %s\n"+
							"%s\n%s\n%s\n %s \n%s",
							b.getText(in, "queue_drs"), in.Lvlkz[1:], b.getText(in, "queue_completed"),
							n1, n2, n3, b.getText(in, "go"), textEvent)
						fmt.Printf("MesidTgUpdate tgmesid %d, in.Lvlkz %s, in.Config.CorpName %s\n", tgmesid, in.Lvlkz, in.Config.CorpName)
						tgmesid = b.client.Tg.SendChannel(in.Config.TgChannel, text)
						err = b.storage.Update.MesidTgUpdate(tgmesid, in.Lvlkz, in.Config.CorpName)
						if err != nil {
							b.log.Error("MesidTgUpdate " + err.Error())
							err = b.storage.Update.MesidTgUpdate(tgmesid, in.Lvlkz, in.Config.CorpName)
							if err != nil {
								b.log.ErrorErr(err)
								b.log.Warn("this problem")
							}
						}
						b.wg.Done()
						close(ch)
					}()
				}

				b.wg.Wait()
				go b.SendLsNotification(in, u)
				b.storage.DbFunc.InsertQueue(dsmesid, alt, in.Config.CorpName, in.Username, in.UserId, in.NameMention, in.Tip, in.Lvlkz, in.Timekz, tgmesid, numkzN)
				err = b.storage.Update.UpdateCompliteRS(in.Lvlkz, dsmesid, tgmesid, alt, numkzL, numberevent, in.Config.CorpName)
				if err != nil {
					b.log.Error("UpdateCompliteRS " + err.Error())
					err = b.storage.Update.UpdateCompliteRS(in.Lvlkz, dsmesid, tgmesid, alt, numkzL, numberevent, in.Config.CorpName)
					if err != nil {
						b.log.ErrorErr(err)
					}
				}

				if numkzEvent == 0 {
					//отправляем сообщение о корпорациях с %
					go b.SendPercent(in.Config)
				}
				//проверка есть ли игрок в других чатах
				user := u.GetAllUserId()
				go b.elseChat(user)
				go b.helpers.SaveUsersIdQueue(user, in.Config)
			}
			go b.CheckSubscribe(in)
		}
	}
}
func (b *Bot) RsSoloPlus(in models.InMessage) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.iftipdelete(in)
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
	textEvent, numkzEvent := b.EventText(in)
	numberevent := b.storage.Event.NumActiveEvent(in.Config.CorpName) //получаем номер ивета если он активен
	if numberevent > 0 {
		numkzL = numkzEvent
	} else {
		b.ifTipSendTextDelSecond(in, "event not active ", 30)
		return
	}
	text := fmt.Sprintf("Соло 😱 %s \n🤘  %s \n%s%s", in.Lvlkz, in.NameMention, b.getText(in, "go"), textEvent)
	if in.Config.DsChannel != "" {
		dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Ds.Avatar)
	}
	if in.Config.TgChannel != "" {
		if in.Config.DsChannel != "" {
			text = b.client.Ds.ReplaceTextMessage(text, in.Config.Guildid)
		}
		tgmesid = b.client.Tg.SendChannel(in.Config.TgChannel, text)
	}

	b.storage.DbFunc.InsertQueue(dsmesid, "", in.Config.CorpName, in.Username, in.UserId, in.NameMention, in.Tip, in.Lvlkz, in.Timekz, tgmesid, numkzN)
	err := b.storage.Update.UpdateCompliteSolo(in.Lvlkz, dsmesid, tgmesid, numkzL, numberevent, in.Config.CorpName)
	if err != nil {
		err = b.storage.Update.UpdateCompliteSolo(in.Lvlkz, dsmesid, tgmesid, numkzL, numberevent, in.Config.CorpName)
		if err != nil {
			b.log.ErrorErr(err)
		}
	}

	//проверка есть ли игрок в других чатах
	go b.elseChat([]string{in.UserId})

}

func (b *Bot) RsSoloPlusComplete(in models.InMessage, pointsStr string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.iftipdelete(in)
	numkzN, _ := b.storage.Count.CountNumberNameActive1(in.Lvlkz, in.Config.CorpName, in.UserId) //проверяем количество боёв по уровню кз игрока
	numberevent := b.storage.Event.NumActiveEvent(in.Config.CorpName)                            //получаем номер ивета если он активен
	if numberevent == 0 {
		b.ifTipSendTextDelSecond(in, "event not active ", 30)
		return
	}
	_, numkzEvent := b.EventText(in)
	points, err := strconv.Atoi(pointsStr)
	if err != nil {
		b.ifTipSendTextDelSecond(in, "не распознал очки  ", 30)
		return
	}

	mes1 := fmt.Sprintf("🔴 %s №%d (%s)\n", b.getText(in, "event_game"), numkzEvent, in.Lvlkz)
	mesOld := fmt.Sprintf("🎉 %s %s %d\nㅤ\nㅤ", b.getText(in, "contributed"), in.Username, points)

	dsmesid := ""
	tgmesid := 0
	text := fmt.Sprintf("%s %s \n%s", mes1, in.NameMention, mesOld)
	if in.Config.DsChannel != "" {
		dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Ds.Avatar)
	}
	if in.Config.TgChannel != "" {
		tgmesid = b.client.Tg.SendChannel(in.Config.TgChannel, text)
	}

	b.storage.DbFunc.InsertQueueSolo(dsmesid, "", in.Config.CorpName, in.Username, in.UserId,
		in.NameMention, in.Tip, in.Lvlkz, tgmesid, numberevent, numkzEvent, numkzN, points)
}
func (b *Bot) SendLsNotification(in models.InMessage, u models.Users) {
	dmText := fmt.Sprintf("%s\n%s %s\n", b.getText(in, "go"), u.User2.Name, u.User3.Name)
	if u.User1.Tip == ds {
		go b.client.Ds.SendDmText(dmText, u.User1.UserId)
	} else if u.User1.Tip == tg {
		b.client.Tg.SendChannel(u.User1.UserId, dmText)
	}
	dmText = fmt.Sprintf("%s\n%s %s\n", b.getText(in, "go"), u.User1.Name, u.User3.Name)
	if u.User2.Tip == ds {
		go b.client.Ds.SendDmText(dmText, u.User2.UserId)
	} else if u.User2.Tip == tg {
		b.client.Tg.SendChannel(u.User2.UserId, dmText)
	}

}
