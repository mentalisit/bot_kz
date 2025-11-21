package bot

import (
	"fmt"
	"rs/bot/helpers"
	"rs/models"
	"strconv"
)

func (b *Bot) darkAlt(in models.InMessage, i int) {
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
	fmt.Println("RsDarkPlus " + alt)
	b.RsDarkPlus(in, alt)
}

func (b *Bot) RsDarkPlus(in models.InMessage, alt string) {
	b.helpers.ReadNameModules(in, alt)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.iftipdelete(in)
	CountName, err := b.storage.Count.СountName(in.UserId, in.RsTypeLevel, in.Config.CorpName)
	if err != nil {
		b.log.ErrorErr(err)
		return
	}
	if CountName == 1 { //проверяем есть ли игрок в очереди
		b.ifTipSendMentionText(in, b.getText(in, "you_in_queue"))
	} else {
		countQueue, numberName, numberLevel, errorsAll := b.storage.Count.CountQueueNumberNameActive1QueueLvl(in.RsTypeLevel, in.Config.CorpName, in.UserId)
		if errorsAll != nil {
			b.log.ErrorErr(errorsAll)
		}

		DsMessageId := ""
		TgMessageId := 0

		n := b.getMap(in, numberLevel)

		u := models.Users{}
		UserIn := models.Sborkz{
			Tip:      in.Tip,
			Name:     in.Username,
			UserId:   in.UserId,
			Mention:  in.NameMention,
			Numkzn:   numberName,
			Timedown: in.TimeRs,
			Wamesid:  alt,
		}

		darkOrRed, level := in.TypeRedStar()

		ntg := make(map[string]string)
		if in.IfTelegram() {
			ntg["text1"] = fmt.Sprintf("%s%s (%d)\n", b.getText(in, "queue_drs"), level, numberLevel)
			if !darkOrRed {
				ntg["text1"] = fmt.Sprintf("%s%s (%d)\n", b.getText(in, "rs_queue"), level, numberLevel)
			}
			ntg["text2"] = fmt.Sprintf("\n%s++ - %s", level, b.getText(in, "forced_start"))
			ntg["min"] = b.getText(in, "min")
		}

		if countQueue == 0 {
			if in.IfDiscord() {
				b.wg.Add(1)
				go func() {
					DsMessageId = b.client.Ds.SendComplexContent(in.Config.DsChannel,
						fmt.Sprintf(b.getText(in, "temp_queue_started"), in.Username, n["levelRs"]))
					b.wg.Done()
				}()
			}
			if in.IfTelegram() {
				b.wg.Add(1)
				go func() {
					text := fmt.Sprintf(b.getText(in, "temp_queue_started"), in.Username, b.getText(in, "drs")+level)
					if !darkOrRed {
						text = fmt.Sprintf(b.getText(in, "temp_queue_started"), in.Username, b.getText(in, "rs")+level)
					}
					TgMessageId = b.client.Tg.SendChannel(in.Config.TgChannel, text)

					go b.SubscribePing(in)
					b.wg.Done()
				}()
			}
			b.wg.Wait()
			b.storage.DbFunc.InsertQueue(DsMessageId, alt, in.Config.CorpName, in.Username, in.UserId, in.GetNameMention(), in.Tip, in.RsTypeLevel, in.TimeRs, TgMessageId, numberName)
			b.QueueLevel(in)
			go b.ReadQueueLevel(in, 30)
		} else {
			u = b.storage.DbFunc.ReadAll(in.RsTypeLevel, in.Config.CorpName)
			if u.User1.Dsmesid != "" {
				DsMessageId = u.User1.Dsmesid
			}
			if u.User1.Tgmesid != 0 {
				TgMessageId = u.User1.Tgmesid
			}

			if countQueue == 1 {
				if in.IfDiscord() {
					b.wg.Add(1)
					go func() {
						text := fmt.Sprintf("%s 2/3 %s %s \n%s",
							n["levelRs"], in.Username, b.getText(in, "you_joined_queue"), u.User1.Mention)
						if !darkOrRed {
							text = fmt.Sprintf("%s 2/4 %s %s",
								n["levelRs"], in.Username, b.getText(in, "you_joined_queue"))
						}
						go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 10)
						b.wg.Done()
					}()
				}
				if in.IfTelegram() {
					b.wg.Add(1)
					go func() {
						text := fmt.Sprintf("%s%s 2/3 %s %s \n%s", b.getText(in, "drs"), level,
							in.Username, b.getText(in, "you_joined_queue"), u.User1.Mention)
						if !darkOrRed {
							text = fmt.Sprintf("%s%s 2/4 %s %s", b.getText(in, "rs"), level,
								in.Username, b.getText(in, "you_joined_queue"))
						}
						go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 10)

						b.wg.Done()
					}()
				}
				b.wg.Wait()
				b.storage.DbFunc.InsertQueue(DsMessageId, alt, in.Config.CorpName, in.Username,
					in.UserId, in.GetNameMention(), in.Tip, in.RsTypeLevel, in.TimeRs, TgMessageId, numberName)
				b.QueueLevel(in)
				go b.ReadQueueLevel(in, 30)
			}
			if darkOrRed {
				if countQueue == 2 {
					//textEvent, numberRsEvent := b.EventText(in)
					//if textEvent == "" {
					textEvent := b.GetTextPercent(in.Config, true)
					//}
					numberEvent := b.storage.Event.NumActiveEvent(in.Config.CorpName) //получаем номер ивета если он активен
					//if numberEvent > 0 {
					//	numberLevel = numberRsEvent
					//}
					u.User3 = &UserIn

					if in.IfDiscord() {
						b.wg.Add(1)
						go func() {
							n1, n2, n3, _ := b.helpers.NameMention(u, ds)
							go b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
							go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel,
								fmt.Sprintf("🚀 3/3 %s %s", in.Username, b.getText(in, "you_joined_queue")), 10)
							text := fmt.Sprintf("3/3 %s%s %s\n %s\n %s\n %s\n%s %s",
								b.getText(in, "queue_drs"), level, b.getText(in, "queue_completed"),
								n1, n2, n3, b.getText(in, "go"), textEvent)

							DsMessageId = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Ds.Avatar)

							err = b.storage.Update.MesidDsUpdate(DsMessageId, in.RsTypeLevel, in.Config.CorpName)
							if err != nil {
								b.log.Error("MesidDsUpdate " + err.Error())
								err = b.storage.Update.MesidDsUpdate(DsMessageId, in.RsTypeLevel, in.Config.CorpName)
								if err != nil {
									b.log.ErrorErr(err)
									b.log.Warn("this problem")
								}

							}
							b.wg.Done()
						}()

					}
					if in.IfTelegram() {
						b.wg.Add(1)
						go func() {
							n1, n2, n3, _ := b.helpers.NameMention(u, tg)
							go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
							go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel,
								in.Username+b.getText(in, "drs_queue_closed")+level, 10)
							text := fmt.Sprintf("🚀 %s%s %s\n"+
								"%s\n%s\n%s\n %s \n%s",
								b.getText(in, "queue_drs"), level, b.getText(in, "queue_completed"),
								n1, n2, n3, b.getText(in, "go"), textEvent)
							fmt.Printf("MesidTgUpdate TgMessageId %d, level %s, in.Config.CorpName %s\n", TgMessageId, level, in.Config.CorpName)
							TgMessageId = b.client.Tg.SendChannel(in.Config.TgChannel, text)
							err = b.storage.Update.MesidTgUpdate(TgMessageId, in.RsTypeLevel, in.Config.CorpName)
							if err != nil {
								b.log.Error("MesidTgUpdate " + err.Error())
								err = b.storage.Update.MesidTgUpdate(TgMessageId, in.RsTypeLevel, in.Config.CorpName)
								if err != nil {
									b.log.ErrorErr(err)
									b.log.Warn("this problem")
								}
							}
							b.wg.Done()
						}()
					}

					b.wg.Wait()
					go b.SendLsNotification(in, u)
					b.storage.DbFunc.InsertQueue(DsMessageId, alt, in.Config.CorpName, in.Username, in.UserId, in.GetNameMention(), in.Tip, in.RsTypeLevel, in.TimeRs, TgMessageId, numberName)
					err = b.storage.Update.UpdateCompliteRS(in.RsTypeLevel, DsMessageId, TgMessageId, alt, numberLevel, numberEvent, in.Config.CorpName)
					if err != nil {
						b.log.Error("UpdateCompliteRS " + err.Error())
						err = b.storage.Update.UpdateCompliteRS(in.RsTypeLevel, DsMessageId, TgMessageId, alt, numberLevel, numberEvent, in.Config.CorpName)
						if err != nil {
							b.log.ErrorErr(err)
						}
					}

					//if numberRsEvent == 0 {
					//отправляем сообщение о корпорациях с %
					go b.SendPercent(in.Config)
					//}
					//проверка есть ли игрок в других чатах
					user, UserIdTg := u.GetAllUserId()
					go b.otherQueue.NeedRemoveOtherQueue(UserIdTg)
					go b.elseChat(user)
				}
			}

			if !darkOrRed {
				if countQueue == 2 {
					DsMessageId = u.User1.Dsmesid
					u.User3 = &UserIn

					if in.IfDiscord() {
						b.wg.Add(1)
						go func() {
							text := fmt.Sprintf("%s  3/4 %s %s %s",
								n["levelRs"], in.Username, b.getText(in, "you_joined_queue"),
								b.getText(in, "another_one_needed_to_complete_queue"))
							go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 10)
							b.wg.Done()
						}()
					}
					if in.IfTelegram() {
						b.wg.Add(1)
						go func() {
							text := fmt.Sprintf("%s  3/4 %s %s %s", b.getText(in, "rs")+level,
								in.Username, b.getText(in, "you_joined_queue"),
								b.getText(in, "another_one_needed_to_complete_queue"))

							go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 10)
							b.wg.Done()
						}()
					}
					b.wg.Wait()
					b.storage.DbFunc.InsertQueue(DsMessageId, "", in.Config.CorpName, in.Username,
						in.UserId, in.GetNameMention(), in.Tip, in.RsTypeLevel, in.TimeRs, TgMessageId, numberName)
					b.QueueLevel(in)
				}
				if countQueue == 3 {
					u.User4 = &UserIn
					//textEvent, numkzEvent := b.EventText(in)
					//if textEvent == "" {
					textEvent := b.GetTextPercent(in.Config, false)
					//}
					//numberevent := b.storage.Event.NumActiveEvent(in.Config.CorpName) //получаем номер ивета если он активен
					//if numberevent > 0 {
					//	numberLevel = numkzEvent
					//}
					if in.IfDiscord() {
						b.wg.Add(1)
						go func() {
							n1, n2, n3, n4 := b.helpers.NameMention(u, ds)
							go b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
							go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel,
								fmt.Sprintf(" 4/4 %s %s", in.Username, b.getText(in, "you_joined_queue")), 10)
							text := fmt.Sprintf("4/4 %s%s %s\n"+
								" %s\n %s\n %s\n %s\n%s %s",
								b.getText(in, "rs_queue"), level, b.getText(in, "queue_completed"),
								n1, n2, n3, n4, b.getText(in, "go"), textEvent)
							DsMessageId = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Ds.Avatar)
							err = b.storage.Update.MesidDsUpdate(DsMessageId, in.RsTypeLevel, in.Config.CorpName)
							if err != nil {
								b.log.ErrorErr(err)
							}
							b.wg.Done()
						}()
					}
					if in.IfTelegram() {
						b.wg.Add(1)
						go func() {
							n1, n2, n3, n4 := b.helpers.NameMention(u, tg)
							go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
							go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel,
								in.Username+b.getText(in, "rs_queue_closed")+level, 10)
							text := fmt.Sprintf("%s%s %s\n"+
								"%s\n%s\n%s\n%s\n %s \n%s",
								b.getText(in, "rs_queue"), level, b.getText(in, "queue_completed"),
								n1, n2, n3, n4, b.getText(in, "go"), textEvent)
							TgMessageId = b.client.Tg.SendChannel(in.Config.TgChannel, text)
							err = b.storage.Update.MesidTgUpdate(TgMessageId, in.RsTypeLevel, in.Config.CorpName)
							if err != nil {
								b.log.ErrorErr(err)
							}
							b.wg.Done()
						}()
					}

					b.wg.Wait()
					b.storage.DbFunc.InsertQueue(DsMessageId, "", in.Config.CorpName, in.Username, in.UserId, in.GetNameMention(),
						in.Tip, in.RsTypeLevel, in.TimeRs, TgMessageId, numberLevel)
					err = b.storage.Update.UpdateCompliteRS(in.RsTypeLevel, DsMessageId, TgMessageId, "", numberLevel, 0, in.Config.CorpName)
					if err != nil {
						err = b.storage.Update.UpdateCompliteRS(in.RsTypeLevel, DsMessageId, TgMessageId, "", numberLevel, 0, in.Config.CorpName)
						if err != nil {
							b.log.ErrorErr(err)
						}
					}

					//if numberevent == 0 {
					//отправляем сообщение о корпорациях с %
					go b.SendPercent(in.Config)
					//}

					//проверка есть ли игрок в других чатах
					user, UserIdTg := u.GetAllUserId()
					go b.otherQueue.NeedRemoveOtherQueue(UserIdTg)
					go b.elseChat(user)
				}
			}
		}
		//if darkOrRed {
		//	go b.CheckSubscribe(in)
		//}

	}
}

//func (b *Bot) RsSoloPlus(in models.InMessage) {
//	b.mu.Lock()
//	defer b.mu.Unlock()
//	b.iftipdelete(in)
//	numkzN, err2 := b.storage.Count.CountNumberNameActive1(in.RsTypeLevel, in.Config.CorpName, in.UserId) //проверяем количество боёв по уровню кз игрока
//	if err2 != nil {
//		return
//	}
//	numkzL, err3 := b.storage.DbFunc.NumberQueueLvl(in.RsTypeLevel, in.Config.CorpName) //проверяем какой номер боя определенной красной звезды
//	if err3 != nil {
//		return
//	}
//	dsmesid := ""
//	tgmesid := 0
//	textEvent, numkzEvent := b.EventText(in)
//	numberevent := b.storage.Event.NumActiveEvent(in.Config.CorpName) //получаем номер ивета если он активен
//	if numberevent > 0 {
//		numkzL = numkzEvent
//	} else {
//		b.ifTipSendTextDelSecond(in, "event not active ", 30)
//		return
//	}
//	_, level := in.TypeRedStar()
//	text := fmt.Sprintf("Соло 😱 %s \n🤘  %s \n%s%s", level, in.NameMention, b.getText(in, "go"), textEvent)
//	if in.IfDiscord() {
//		dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Ds.Avatar)
//	}
//	if in.IfTelegram() {
//		if in.IfDiscord() {
//			text = b.client.Ds.ReplaceTextMessage(text, in.Config.Guildid)
//		}
//		tgmesid = b.client.Tg.SendChannel(in.Config.TgChannel, text)
//	}
//
//	b.storage.DbFunc.InsertQueue(dsmesid, "", in.Config.CorpName, in.Username, in.UserId, in.NameMention, in.Tip, in.RsTypeLevel, in.TimeRs, tgmesid, numkzN)
//	err := b.storage.Update.UpdateCompliteSolo(in.RsTypeLevel, dsmesid, tgmesid, numkzL, numberevent, in.Config.CorpName)
//	if err != nil {
//		err = b.storage.Update.UpdateCompliteSolo(in.RsTypeLevel, dsmesid, tgmesid, numkzL, numberevent, in.Config.CorpName)
//		if err != nil {
//			b.log.ErrorErr(err)
//		}
//	}
//
//	//проверка есть ли игрок в других чатах
//	go b.elseChat([]string{in.UserId})
//	if in.Tip == tg {
//		go b.otherQueue.NeedRemoveOtherQueue([]string{in.UserId})
//	}
//
//}
//func (b *Bot) RsSoloPlusComplete(in models.InMessage, pointsStr string) {
//	b.mu.Lock()
//	defer b.mu.Unlock()
//	b.iftipdelete(in)
//	numkzN, _ := b.storage.Count.CountNumberNameActive1(in.RsTypeLevel, in.Config.CorpName, in.UserId) //проверяем количество боёв по уровню кз игрока
//	numberevent := b.storage.Event.NumActiveEvent(in.Config.CorpName)                                  //получаем номер ивета если он активен
//	if numberevent == 0 {
//		b.ifTipSendTextDelSecond(in, "event not active ", 30)
//		return
//	}
//	_, numkzEvent := b.EventText(in)
//	points, err := strconv.Atoi(pointsStr)
//	if err != nil || points > 99999 {
//		b.ifTipSendTextDelSecond(in, "не распознал очки  ", 30)
//		return
//	}
//	_, level := in.TypeRedStar()
//
//	mes1 := fmt.Sprintf("🔴 %s №%d (%s) ", b.getText(in, "event_game"), numkzEvent, level)
//	mesOld := fmt.Sprintf("🎉 %s %s %s\nㅤ\nㅤ", b.getText(in, "contributed"), in.NameMention, formatNumber(points))
//
//	dsmesid := ""
//	tgmesid := 0
//	text := fmt.Sprintf("%s %s \n%s", mes1, in.NameMention, mesOld)
//	if in.IfDiscord() {
//		dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Ds.Avatar)
//	}
//	if in.IfTelegram() {
//		tgmesid = b.client.Tg.SendChannel(in.Config.TgChannel, text)
//	}
//
//	b.storage.DbFunc.InsertQueueSolo(dsmesid, "", in.Config.CorpName, in.Username, in.UserId,
//		in.NameMention, in.Tip, in.RsTypeLevel, tgmesid, numberevent, numkzEvent, numkzN, points)
//}

func (b *Bot) SendLsNotification(in models.InMessage, u models.Users) {
	dmText := fmt.Sprintf("%s\n %s\n %s\n", b.getText(in, "go"), u.User2.Name, u.User3.Name)
	if u.User1.Tip == ds {
		go b.client.Ds.SendDmText(dmText, u.User1.UserId)
	} else if u.User1.Tip == tg {
		sendChannel := b.client.Tg.SendChannel(u.User1.UserId, dmText)
		b.client.Tg.DelMessageSecond(u.User1.UserId, strconv.Itoa(sendChannel), 1800)
	}
	dmText = fmt.Sprintf("%s\n %s\n %s\n", b.getText(in, "go"), u.User1.Name, u.User3.Name)
	if u.User2.Tip == ds {
		go b.client.Ds.SendDmText(dmText, u.User2.UserId)
	} else if u.User2.Tip == tg {
		sendChannel := b.client.Tg.SendChannel(u.User2.UserId, dmText)
		b.client.Tg.DelMessageSecond(u.User2.UserId, strconv.Itoa(sendChannel), 1800)
	}

}
