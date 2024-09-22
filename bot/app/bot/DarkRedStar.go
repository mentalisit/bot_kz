package bot

import (
	"context"
	"fmt"
	"kz_bot/bot/helpers"
	"kz_bot/models"
	"regexp"
	"strconv"
	"time"
)

const dark = "d"

func (b *Bot) lDarkRsPlus(in models.InMessage) bool {
	var kzb string
	kz := false
	re := regexp.MustCompile(`^([7-9]|[1][0-2])([\*]|[-])(\+)?(\d|\d{2}|\d{3})$`) //три переменные
	arr := re.FindAllStringSubmatch(in.Mtext, -1)
	if len(arr) > 0 {
		kz = true
		in.Lvlkz = dark + arr[0][1]
		kzb = arr[0][2]
		timekzz, err := strconv.Atoi(arr[0][4])
		if err != nil {
			b.log.ErrorErr(err)
			timekzz = 0
		}
		if timekzz > 180 {
			timekzz = 180
		}
		if arr[0][3] == "+" {
			in.NameMention = "$" + in.NameMention
		}
		in.Timekz = strconv.Itoa(timekzz)
	}

	re2 := regexp.MustCompile(`^([7-9]|[1][0-2])([\*]|[-])([\+]|[\?])?([1-5])?$`) // две переменные
	arr2 := (re2.FindAllStringSubmatch(in.Mtext, -1))
	if len(arr2) > 0 {
		kz = true
		in.Lvlkz = dark + arr2[0][1]
		kzb = arr2[0][2]
		in.Timekz = "30"
		if arr2[0][3] == "+" {
			in.NameMention = "$" + in.NameMention
		}
		if arr2[0][3] == "?" && arr2[0][4] != "" {
			atoi, _ := strconv.Atoi(arr2[0][4])
			b.darkAlt(in, atoi)
			return true
		} else if arr2[0][3] == "?" {
			b.darkAlt(in, 1)
			return true
		}
	}

	re2d := regexp.MustCompile(`^(d)([7-9]|[1][0-2])([\+]|[-])$`) // две переменные
	arr2d := (re2d.FindAllStringSubmatch(in.Mtext, -1))
	if len(arr2d) > 0 {
		kz = true
		in.Lvlkz = dark + arr2d[0][2]
		kzb = arr2d[0][3]
		in.Timekz = "30"
	}

	//solo
	re2s := regexp.MustCompile(`^([s]|[S])([7-9]|[1][0-2])(\+)$`) // две переменные
	arr2s := (re2s.FindAllStringSubmatch(in.Mtext, -1))
	if len(arr2s) > 0 {
		kz = true
		in.Lvlkz = "d" + arr2s[0][2]
		kzbs := arr2s[0][3]
		in.Timekz = "1"
		if kzbs == "+" {
			b.RsSoloPlus(in)
			return kz
		}
	}

	switch kzb {
	case "*":
		b.RsDarkPlus(in, "")
	case "+":
		b.RsDarkPlus(in, "")
	case "-":
		b.RsMinus(in)
	case "*+":
		b.RsDarkPlus(in, "")

	default:
		kz = false
	}
	return kz
}
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

func (b *Bot) RsDarkPlus(in models.InMessage, alt string) {
	b.helpers.ReadNameModules(in, alt)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.iftipdelete(in)
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	CountName, err := b.storage.Count.СountName(ctx, in.UserId, in.Lvlkz, in.Config.CorpName)
	if err != nil {
		b.log.ErrorErr(err)
		return
	}
	if CountName == 1 { //проверяем есть ли игрок в очереди
		b.ifTipSendMentionText(in, b.getText(in, "you_in_queue"))
	} else {
		countQueue, err1 := b.storage.Count.CountQueue(ctx, in.Lvlkz, in.Config.CorpName) //проверяем, есть ли кто-то в очереди
		if err1 != nil {
			return
		}
		numkzN, err2 := b.storage.Count.CountNumberNameActive1(ctx, in.Lvlkz, in.Config.CorpName, in.UserId) //проверяем количество боёв по уровню кз игрока
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
			n["lvlkz"], err = b.client.Ds.RoleToIdPing(b.getText(in, "drs")+in.Lvlkz[1:], in.Config.Guildid)
			if err != nil {
				b.log.Info(fmt.Sprintf("RoleToIdPing %+v lvl %s", in.Config, in.Lvlkz[1:]))
			}
		}
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
					u.User1 = UserIn
					n = b.helpers.GetQueueDiscord(n, u)
					//n["name1"] = fmt.Sprintf("%s  🕒  %s  (%d)", b.emReadName(in, in.Username, in.NameMention, ds), in.Timekz, numkzN)
					emb := b.client.Ds.EmbedDS(n, numkzL, 1, true)
					dsmesid = b.client.Ds.SendComplexContent(in.Config.DsChannel,
						fmt.Sprintf(b.getText(in, "temp_queue_started"), in.Username, n["lvlkz"]))
					time.Sleep(1 * time.Second)
					b.client.Ds.EditComplexButton(dsmesid, in.Config.DsChannel, emb, b.client.Ds.AddButtonsQueue(in.Lvlkz))
					b.wg.Done()
				}()
			}
			if in.Config.TgChannel != "" {
				b.wg.Add(1)
				go func() {
					//text := fmt.Sprintf("%s%s (%d)\n"+
					//	"1. %s - %s%s (%d) \n\n"+
					//	"%s++ - %s",
					//	b.getText(in, "queue_drs"), in.Lvlkz[1:], numkzL,
					//	b.emReadName(in, in.Username, in.NameMention, tg), in.Timekz, b.getText(in, "min"), numkzN,
					//	in.Lvlkz[1:], b.getText(in, "forced_start"))

					u.User1 = UserIn
					texttg = b.helpers.GetQueueTelegram(ntg, u)

					tgmesid = b.client.Tg.SendEmded(in.Lvlkz, in.Config.TgChannel, texttg)
					b.SubscribePing(in, 1)
					b.wg.Done()
				}()
			}
			go b.ReadQueueLevel(in)
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
					//n["name2"] = fmt.Sprintf("%s  🕒  %s  (%d)", b.emReadName(in, in.Username, in.NameMention, ds), in.Timekz, numkzN)
					emb := b.client.Ds.EmbedDS(n, numkzL, 2, true)
					text := n["lvlkz"] + fmt.Sprintf(" 2/3 %s %s \n%s", in.Username, b.getText(in, "you_joined_queue"), u.User1.Mention)
					go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 10)
					b.client.Ds.EditComplexButton(u.User1.Dsmesid, in.Config.DsChannel, emb, b.client.Ds.AddButtonsQueue(in.Lvlkz))
					b.wg.Done()
				}()
			}
			if in.Config.TgChannel != "" {
				b.wg.Add(1)
				go func() {
					//text1 := fmt.Sprintf("%s%s (%d)\n", b.getText(in, "queue_drs"), in.Lvlkz, numkzL)
					//name1 := fmt.Sprintf("1. %s - %d%s (%d) \n",
					//	b.emReadName(in, u.User1.Name, u.User1.Mention, tg), u.User1.Timedown, b.getText(in, "min"), u.User1.Numkzn)
					//name2 := fmt.Sprintf("2. %s - %s%s (%d) \n",
					//	b.emReadName(in, in.Username, in.NameMention, tg), in.Timekz, b.getText(in, "min"), numkzN)
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
			go b.ReadQueueLevel(in)
		}
		if countQueue < 2 {
			b.wg.Wait()
			b.storage.DbFunc.InsertQueue(ctx, dsmesid, alt, in.Config.CorpName, in.Username, in.UserId, in.NameMention, in.Tip, in.Lvlkz, in.Timekz, tgmesid, numkzN)
		}

		if countQueue == 2 {
			dsmesid = u.User1.Dsmesid

			textEvent, numkzEvent := b.EventText(in)
			if textEvent == "" {
				textEvent = b.GetTextPercent(in.Config, true)
			}
			numberevent := b.storage.Event.NumActiveEvent(in.Config.CorpName) //получаем номер ивета если он активен
			if numberevent > 0 {
				numkzL = numkzEvent
			}
			u.User3 = UserIn

			if in.Config.DsChannel != "" {
				b.wg.Add(1)
				go func() {
					n1, n2, n3, _ := b.helpers.NameMention(u, ds)
					go b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
					go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel,
						fmt.Sprintf("🚀 3/3 %s %s", in.Username, b.getText(in, "you_joined_queue")), 10)
					text := fmt.Sprintf("3/3 %s%s %s\n"+
						" %s\n"+
						" %s\n"+
						" %s\n"+
						"%s %s",
						b.getText(in, "queue_drs"), in.Lvlkz[1:], b.getText(in, "queue_completed"),
						n1,
						n2,
						n3,
						b.getText(in, "go"), textEvent)

					if in.Tip == ds {
						dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Config.Guildid, in.Ds.Avatar)
						if u.User1.Tip == ds {
							//go b.sendDmDark(b.getText(in, "go"), u.User1.Mention)
							go b.client.Ds.SendDmText(b.getText(in, "go"), u.User1.UserId)
						}
						if u.User2.Tip == ds {
							//go b.sendDmDark(b.getText(in, "go"), u.User2.Mention)
							go b.client.Ds.SendDmText(b.getText(in, "go"), u.User2.UserId)
						}
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
					n1, n2, n3, _ := b.helpers.NameMention(u, tg)
					go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
					go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel,
						in.Username+b.getText(in, "drs_queue_closed")+in.Lvlkz[1:], 10)
					text := fmt.Sprintf("🚀 %s%s %s\n"+
						"%s\n"+
						"%s\n"+
						"%s\n"+
						" %s \n"+
						"%s",
						b.getText(in, "queue_drs"), in.Lvlkz[1:], b.getText(in, "queue_completed"),
						n1, n2, n3,
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
			b.storage.DbFunc.InsertQueue(ctx, dsmesid, alt, in.Config.CorpName, in.Username, in.UserId, in.NameMention, in.Tip, in.Lvlkz, in.Timekz, tgmesid, numkzN)
			err = b.storage.Update.UpdateCompliteRS(ctx, in.Lvlkz, dsmesid, tgmesid, alt, numkzL, numberevent, in.Config.CorpName)
			if err != nil {
				err = b.storage.Update.UpdateCompliteRS(context.Background(), in.Lvlkz, dsmesid, tgmesid, "", numkzL, numberevent, in.Config.CorpName)
				if err != nil {
					b.log.ErrorErr(err)
				}
			}

			//отправляем сообщение о корпорациях с %
			go b.SendPercent(in.Config)

			//проверка есть ли игрок в других чатах
			user := []string{u.User1.UserId, u.User2.UserId, in.UserId}
			go b.elseChat(user)
			go b.helpers.SaveUsersIdQueue(user, in.Config)
		}

	}
}
func (b *Bot) lDarkSubs(in models.InMessage) (bb bool) {
	bb = false
	var subs string
	re3 := regexp.MustCompile(`^([\+]|[-])(d)([7-9]|[1][0-2])$`) // две переменные для добавления или удаления подписок
	arr3 := (re3.FindAllStringSubmatch(in.Mtext, -1))
	if len(arr3) > 0 {
		in.Lvlkz = "d" + arr3[0][3]
		subs = arr3[0][1]
		bb = true
	}

	re6 := regexp.MustCompile(`^([\+][\+]|[-][-])(d)([7-9]|[1][0-2])$`) // две переменные
	arr6 := (re6.FindAllStringSubmatch(in.Mtext, -1))                   // для добавления или удаления подписок 2/3
	if len(arr6) > 0 {
		bb = true
		in.Lvlkz = "d" + arr6[0][3]
		subs = arr6[0][1]
	}

	switch subs {
	case "+":
		b.Subscribe(in, 1)
	case "++":
		b.Subscribe(in, 3)
	case "-":
		b.Unsubscribe(in, 1)
	case "--":
		b.Unsubscribe(in, 3)
	}
	return bb
}
func (b *Bot) RsSoloPlus(in models.InMessage) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.debug {
		fmt.Printf("\n\nin RsSoloPlus %+v\n", in)
	}
	b.iftipdelete(in)
	ctx := context.Background()
	numkzN, err2 := b.storage.Count.CountNumberNameActive1(ctx, in.Lvlkz, in.Config.CorpName, in.UserId) //проверяем количество боёв по уровню кз игрока
	if err2 != nil {
		return
	}
	numkzL, err3 := b.storage.DbFunc.NumberQueueLvl(ctx, in.Lvlkz, in.Config.CorpName) //проверяем какой номер боя определенной красной звезды
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
		//todo send not event
		b.ifTipSendTextDelSecond(in, "event not active ", 30)
		return
	}
	text := fmt.Sprintf("Соло 😱 %s \n🤘  %s \n%s%s", in.Lvlkz, in.NameMention, b.getText(in, "go"), textEvent)
	if in.Config.DsChannel != "" {
		if in.Tip == ds {
			dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Config.Guildid, in.Ds.Avatar)
		} else {
			dsmesid = b.client.Ds.Send(in.Config.DsChannel, text)
		}
	}
	if in.Config.TgChannel != "" {
		tgmesid = b.client.Tg.SendChannel(in.Config.TgChannel, text)
	}

	b.storage.DbFunc.InsertQueue(ctx, dsmesid, "", in.Config.CorpName, in.Username, in.UserId, in.NameMention, in.Tip, in.Lvlkz, in.Timekz, tgmesid, numkzN)
	err := b.storage.Update.UpdateCompliteSolo(ctx, in.Lvlkz, dsmesid, tgmesid, numkzL, numberevent, in.Config.CorpName)
	if err != nil {
		err = b.storage.Update.UpdateCompliteSolo(context.Background(), in.Lvlkz, dsmesid, tgmesid, numkzL, numberevent, in.Config.CorpName)
		if err != nil {
			b.log.ErrorErr(err)
		}
	}

	//проверка есть ли игрок в других чатах
	go b.elseChat([]string{in.UserId})

}

//func (b *Bot) sendDmDark(text, userMention string) {
//	mentionRegexDs := regexp.MustCompile(`<@(\d+)>`)
//	match := mentionRegexDs.FindStringSubmatch(userMention)
//	if len(match) > 1 {
//		id := match[1]
//		b.client.Ds.SendDmText(text, id)
//	}
//
//}
