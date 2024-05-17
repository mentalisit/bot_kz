package bot

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const dark = "d"

func (b *Bot) lDarkRsPlus() bool {
	var kzb string
	kz := false
	re := regexp.MustCompile(`^([7-9]|[1][0-2])([\*]|[-])(\+)?(\d|\d{2}|\d{3})$`) //три переменные
	arr := re.FindAllStringSubmatch(b.in.Mtext, -1)
	if len(arr) > 0 {
		kz = true
		b.in.Lvlkz = dark + arr[0][1]
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
			b.in.NameMention = "$" + b.in.NameMention
		}
		b.in.Timekz = strconv.Itoa(timekzz)
	}

	re2 := regexp.MustCompile(`^([7-9]|[1][0-2])([\*]|[-])(\+)?$`) // две переменные
	arr2 := (re2.FindAllStringSubmatch(b.in.Mtext, -1))
	if len(arr2) > 0 {
		kz = true
		b.in.Lvlkz = dark + arr2[0][1]
		kzb = arr2[0][2]
		b.in.Timekz = "30"
		if arr2[0][3] == "+" {
			b.in.NameMention = "$" + b.in.NameMention
		}
	}

	re2d := regexp.MustCompile(`^(d)([7-9]|[1][0-2])([\+]|[-])$`) // две переменные
	arr2d := (re2d.FindAllStringSubmatch(b.in.Mtext, -1))
	if len(arr2d) > 0 {
		kz = true
		b.in.Lvlkz = dark + arr2d[0][2]
		kzb = arr2d[0][3]
		b.in.Timekz = "30"
	}

	//solo
	re2s := regexp.MustCompile(`^([s]|[S])([7-9]|[1][0-2])(\+)$`) // две переменные
	arr2s := (re2s.FindAllStringSubmatch(b.in.Mtext, -1))
	if len(arr2s) > 0 {
		kz = true
		b.in.Lvlkz = arr2s[0][2]
		kzbs := arr2s[0][3]
		b.in.Timekz = "1"
		if kzbs == "+" {
			b.RsSoloPlus()
			return kz
		}
	}

	switch kzb {
	case "*":
		b.RsDarkPlus()
	case "+":
		b.RsDarkPlus()
	case "-":
		b.RsMinus()
	case "*+":
		b.RsDarkPlus()

	default:
		kz = false
	}
	return kz
}
func (b *Bot) RsDarkPlus() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.debug {
		fmt.Printf("\n\nin RsDarkPlus %+v\n", b.in)
	}
	b.iftipdelete()
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	CountName, err := b.storage.Count.СountName(ctx, b.in.Name, b.in.Lvlkz, b.in.Config.CorpName)
	if err != nil {
		return
	}
	if CountName == 1 { //проверяем есть ли игрок в очереди
		b.ifTipSendMentionText(b.getText("you_in_queue"))
	} else {
		countQueue, err1 := b.storage.Count.CountQueue(ctx, b.in.Lvlkz, b.in.Config.CorpName) //проверяем, есть ли кто-то в очереди
		if err1 != nil {
			return
		}
		numkzN, err2 := b.storage.Count.CountNumberNameActive1(ctx, b.in.Lvlkz, b.in.Config.CorpName, b.in.Name) //проверяем количество боёв по уровню кз игрока
		if err2 != nil {
			return
		}
		numkzL, err3 := b.storage.DbFunc.NumberQueueLvl(ctx, b.in.Lvlkz, b.in.Config.CorpName) //проверяем какой номер боя определенной красной звезды
		if err3 != nil {
			return
		}

		dsmesid := ""
		tgmesid := 0
		var n map[string]string
		n = make(map[string]string)
		n["lang"] = b.in.Config.Country
		if b.in.Config.DsChannel != "" {
			n["lvlkz"], err = b.client.Ds.RoleToIdPing(b.getText("drs")+b.in.Lvlkz[1:], b.in.Config.Guildid)
			if err != nil {
				b.log.Info(fmt.Sprintf("RoleToIdPing %+v lvl %s", b.in.Config, b.in.Lvlkz[1:]))
			}
		}
		if countQueue == 0 {
			if b.in.Config.DsChannel != "" {
				b.wg.Add(1)
				go func() {
					n["name1"] = fmt.Sprintf("%s  🕒  %s  (%d)", b.emReadName(b.in.Name, b.in.NameMention, ds), b.in.Timekz, numkzN)
					emb := b.client.Ds.EmbedDS(n, numkzL, 1, true)
					dsmesid = b.client.Ds.SendComplexContent(b.in.Config.DsChannel,
						fmt.Sprintf(b.getText("temp_queue_started"), b.in.Name, n["lvlkz"]))
					time.Sleep(1 * time.Second)
					b.client.Ds.EditComplexButton(dsmesid, b.in.Config.DsChannel, emb, b.client.Ds.AddButtonsQueue(b.in.Lvlkz))
					//b.client.Ds.AddEnojiRsQueue(b.in.Config.DsChannel, dsmesid)
					b.wg.Done()
				}()
			}
			if b.in.Config.TgChannel != "" {
				b.wg.Add(1)
				go func() {
					text := fmt.Sprintf("%s%s (%d)\n"+
						"1. %s - %s%s (%d) \n\n"+
						"%s++ - %s",
						b.getText("queue_drs"), b.in.Lvlkz[1:], numkzL,
						b.emReadName(b.in.Name, b.in.NameMention, tg), b.in.Timekz, b.getText("min"), numkzN,
						b.in.Lvlkz[1:], b.getText("forced_start"))
					tgmesid = b.client.Tg.SendEmded(b.in.Lvlkz, b.in.Config.TgChannel, text)
					b.SubscribePing(1)
					b.wg.Done()
				}()
			}
		}

		u := b.storage.DbFunc.ReadAll(ctx, b.in.Lvlkz, b.in.Config.CorpName)

		if countQueue == 1 {
			dsmesid = u.User1.Dsmesid

			if b.in.Config.DsChannel != "" {
				b.wg.Add(1)
				go func() {
					n["name1"] = fmt.Sprintf("%s  🕒  %d  (%d)", b.emReadName(u.User1.Name, u.User1.Mention, ds), u.User1.Timedown, u.User1.Numkzn)
					n["name2"] = fmt.Sprintf("%s  🕒  %s  (%d)", b.emReadName(b.in.Name, b.in.NameMention, ds), b.in.Timekz, numkzN)
					emb := b.client.Ds.EmbedDS(n, numkzL, 2, true)
					text := n["lvlkz"] + fmt.Sprintf(" 2/3 %s %s \n%s", b.in.Name, b.getText("you_joined_queue"), u.User1.Mention)
					//text := n["lvlkz"] + " 2/3 " + b.in.Name + b.getText("you_joined_queue") + "\n" + u.User1.Mention
					go b.client.Ds.SendChannelDelSecond(b.in.Config.DsChannel, text, 10)
					b.client.Ds.EditComplexButton(u.User1.Dsmesid, b.in.Config.DsChannel, emb, b.client.Ds.AddButtonsQueue(b.in.Lvlkz))
					b.wg.Done()
				}()
			}
			if b.in.Config.TgChannel != "" {
				b.wg.Add(1)
				go func() {
					text1 := fmt.Sprintf("%s%s (%d)\n", b.getText("queue_drs"), b.in.Lvlkz, numkzL)
					name1 := fmt.Sprintf("1. %s - %d%s (%d) \n",
						b.emReadName(u.User1.Name, u.User1.Mention, tg), u.User1.Timedown, b.getText("min"), u.User1.Numkzn)
					name2 := fmt.Sprintf("2. %s - %s%s (%d) \n",
						b.emReadName(b.in.Name, b.in.NameMention, tg), b.in.Timekz, b.getText("min"), numkzN)
					text2 := fmt.Sprintf("\n%s++ - %s", b.in.Lvlkz, b.getText("forced_start"))
					text := fmt.Sprintf("%s %s %s %s", text1, name1, name2, text2)
					tgmesid = b.client.Tg.SendEmded(b.in.Lvlkz, b.in.Config.TgChannel, text)
					go b.client.Tg.DelMessage(b.in.Config.TgChannel, u.User1.Tgmesid)
					err = b.storage.Update.MesidTgUpdate(ctx, tgmesid, b.in.Lvlkz, b.in.Config.CorpName)
					if err != nil {
						err = b.storage.Update.MesidTgUpdate(context.Background(), tgmesid, b.in.Lvlkz, b.in.Config.CorpName)
						if err != nil {
							b.log.ErrorErr(err)
						}
					}
					b.wg.Done()
				}()
			}
		}
		if countQueue < 2 {
			b.wg.Wait()
			b.storage.DbFunc.InsertQueue(ctx, dsmesid, "", b.in.Config.CorpName, b.in.Name, b.in.NameMention, b.in.Tip, b.in.Lvlkz, b.in.Timekz, tgmesid, numkzN)
		}

		if countQueue == 2 {
			dsmesid = u.User1.Dsmesid

			textEvent, numkzEvent := b.EventText()
			if textEvent == "" {
				textEvent = b.GetTextPercent(b.in.Config, true)
			}
			numberevent := b.storage.Event.NumActiveEvent(b.in.Config.CorpName) //получаем номер ивета если он активен
			if numberevent > 0 {
				numkzL = numkzEvent
			}

			if b.in.Config.DsChannel != "" {
				b.wg.Add(1)
				go func() {
					n1, n2, _, n3 := b.nameMention(u, ds)
					go b.client.Ds.DeleteMessage(b.in.Config.DsChannel, u.User1.Dsmesid)
					go b.client.Ds.SendChannelDelSecond(b.in.Config.DsChannel,
						//"🚀 3/3 "+b.in.Name+" "+b.getText("you_joined_queue"), 10)
						fmt.Sprintf("🚀 3/3 %s %s", b.in.Name, b.getText("you_joined_queue")), 10)
					text := fmt.Sprintf("3/3 %s%s %s\n"+
						" %s\n"+
						" %s\n"+
						" %s\n"+
						"%s %s",
						b.getText("queue_drs"), b.in.Lvlkz[1:], b.getText("queue_completed"),
						n1,
						n2,
						n3,
						b.getText("go"), textEvent)

					if b.in.Tip == ds {
						dsmesid = b.client.Ds.SendWebhook(text, "RsBot", b.in.Config.DsChannel, b.in.Config.Guildid, b.in.Ds.Avatar)
						if u.User1.Tip == ds {
							go b.sendDmDark(text, u.User1.Mention)
						}
						if u.User2.Tip == ds {
							go b.sendDmDark(text, u.User2.Mention)
						}
					} else {
						dsmesid = b.client.Ds.Send(b.in.Config.DsChannel, text)
					}
					err = b.storage.Update.MesidDsUpdate(ctx, dsmesid, b.in.Lvlkz, b.in.Config.CorpName)
					if err != nil {
						err = b.storage.Update.MesidDsUpdate(context.Background(), dsmesid, b.in.Lvlkz, b.in.Config.CorpName)
						if err != nil {
							b.log.ErrorErr(err)
						}
					}
					b.wg.Done()
				}()
			}
			if b.in.Config.TgChannel != "" {
				b.wg.Add(1)
				go func() {
					n1, n2, _, n3 := b.nameMention(u, tg)
					go b.client.Tg.DelMessage(b.in.Config.TgChannel, u.User1.Tgmesid)
					go b.client.Tg.SendChannelDelSecond(b.in.Config.TgChannel,
						b.in.Name+b.getText("drs_queue_closed")+b.in.Lvlkz[1:], 10)
					text := fmt.Sprintf("🚀 %s%s %s\n"+
						"%s\n"+
						"%s\n"+
						"%s\n"+
						" %s \n"+
						"%s",
						b.getText("queue_drs"), b.in.Lvlkz[1:], b.getText("queue_completed"),
						n1, n2, n3,
						b.getText("go"), textEvent)
					tgmesid = b.client.Tg.SendChannel(b.in.Config.TgChannel, text)
					err = b.storage.Update.MesidTgUpdate(ctx, tgmesid, b.in.Lvlkz, b.in.Config.CorpName)
					if err != nil {
						err = b.storage.Update.MesidTgUpdate(context.Background(), tgmesid, b.in.Lvlkz, b.in.Config.CorpName)
						if err != nil {
							b.log.ErrorErr(err)
						}
					}
					b.wg.Done()
				}()
			}

			b.wg.Wait()
			b.storage.DbFunc.InsertQueue(ctx, dsmesid, "", b.in.Config.CorpName, b.in.Name, b.in.NameMention, b.in.Tip, b.in.Lvlkz, b.in.Timekz, tgmesid, numkzN)
			err = b.storage.Update.UpdateCompliteRS(ctx, b.in.Lvlkz, dsmesid, tgmesid, "", numkzL, numberevent, b.in.Config.CorpName)
			if err != nil {
				err = b.storage.Update.UpdateCompliteRS(context.Background(), b.in.Lvlkz, dsmesid, tgmesid, "", numkzL, numberevent, b.in.Config.CorpName)
				if err != nil {
					b.log.ErrorErr(err)
				}
			}

			//отправляем сообщение о корпорациях с %
			go b.SendPercent(b.in.Config)

			//проверка есть ли игрок в других чатах
			user := []string{u.User1.Name, u.User2.Name, b.in.Name}
			go b.elseChat(user)
		}

	}
}
func (b *Bot) lDarkSubs() (bb bool) {
	bb = false
	var subs string
	re3 := regexp.MustCompile(`^([\+]|[-])(d)([7-9]|[1][0-2])$`) // две переменные для добавления или удаления подписок
	arr3 := (re3.FindAllStringSubmatch(b.in.Mtext, -1))
	if len(arr3) > 0 {
		b.in.Lvlkz = "d" + arr3[0][3]
		subs = arr3[0][1]
		bb = true
	}

	re6 := regexp.MustCompile(`^([\+][\+]|[-][-])(d)([7-9]|[1][0-2])$`) // две переменные
	arr6 := (re6.FindAllStringSubmatch(b.in.Mtext, -1))                 // для добавления или удаления подписок 2/3
	if len(arr6) > 0 {
		bb = true
		b.in.Lvlkz = "d" + arr6[0][3]
		subs = arr6[0][1]
	}

	switch subs {
	case "+":
		b.Subscribe(1)
	case "++":
		b.Subscribe(3)
	case "-":
		b.Unsubscribe(1)
	case "--":
		b.Unsubscribe(3)
	}
	return bb
}
func (b *Bot) RsSoloPlus() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.debug {
		fmt.Printf("\n\nin RsSoloPlus %+v\n", b.in)
	}
	b.iftipdelete()
	ctx := context.Background()
	numkzN, err2 := b.storage.Count.CountNumberNameActive1(ctx, b.in.Lvlkz[1:], b.in.Config.CorpName, b.in.Name) //проверяем количество боёв по уровню кз игрока
	if err2 != nil {
		return
	}
	numkzL, err3 := b.storage.DbFunc.NumberQueueLvl(ctx, b.in.Lvlkz[1:], b.in.Config.CorpName) //проверяем какой номер боя определенной красной звезды
	if err3 != nil {
		return
	}
	dsmesid := ""
	tgmesid := 0
	textEvent, numkzEvent := b.EventText()
	numberevent := b.storage.Event.NumActiveEvent(b.in.Config.CorpName) //получаем номер ивета если он активен
	if numberevent > 0 {
		numkzL = numkzEvent
	}
	text := fmt.Sprintf("Соло 😱 %s \n🤘  %s \n%s%s", b.in.Lvlkz, b.in.NameMention, b.getText("go"), textEvent)
	if b.in.Config.DsChannel != "" {
		if b.in.Tip == ds {
			dsmesid = b.client.Ds.SendWebhook(text, "RsBot", b.in.Config.DsChannel, b.in.Config.Guildid, b.in.Ds.Avatar)
		} else {
			dsmesid = b.client.Ds.Send(b.in.Config.DsChannel, text)
		}
	}
	if b.in.Config.TgChannel != "" {
		tgmesid = b.client.Tg.SendChannel(b.in.Config.TgChannel, text)
	}

	b.storage.DbFunc.InsertQueue(ctx, dsmesid, "", b.in.Config.CorpName, b.in.Name, b.in.NameMention, b.in.Tip, b.in.Lvlkz, b.in.Timekz, tgmesid, numkzN)
	err := b.storage.Update.UpdateCompliteRS(ctx, b.in.Lvlkz, dsmesid, tgmesid, "", numkzL, numberevent, b.in.Config.CorpName)
	if err != nil {
		err = b.storage.Update.UpdateCompliteRS(context.Background(), b.in.Lvlkz, dsmesid, tgmesid, "", numkzL, numberevent, b.in.Config.CorpName)
		if err != nil {
			b.log.ErrorErr(err)
		}
	}

	//проверка есть ли игрок в других чатах
	go b.elseChat([]string{b.in.Name})

}
func (b *Bot) sendDmDark(text, userMention string) {
	mentionRegexDs := regexp.MustCompile(`<@(\d+)>`)
	match := mentionRegexDs.FindStringSubmatch(userMention)
	if len(match) > 1 {
		id := match[1]
		b.client.Ds.SendDmText(text, id)
	}

}
