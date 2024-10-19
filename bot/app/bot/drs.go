package bot

import (
	"fmt"
	"kz_bot/bot/helpers"
	"kz_bot/models"
	"kz_bot/pkg/utils"
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
					b.client.Ds.EditComplexButton(dsmesid, in.Config.DsChannel, n)
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
			go b.ReadQueueLevel(in)
		}

		u = b.storage.DbFunc.ReadAll(in.Lvlkz, in.Config.CorpName)

		if countQueue == 1 {
			dsmesid = u.User1.Dsmesid

			if in.Config.DsChannel != "" {
				b.wg.Add(1)
				go func() {
					ch := utils.WaitForMessage("RsDarkPlus223")
					u.User2 = UserIn
					n = b.helpers.GetQueueDiscord(n, u)
					text := n["lvlkz"] + fmt.Sprintf(" 2/3 %s %s \n%s", in.Username, b.getText(in, "you_joined_queue"), u.User1.Mention)
					go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 10)
					b.client.Ds.EditComplexButton(u.User1.Dsmesid, in.Config.DsChannel, n)
					b.wg.Done()
					close(ch)
				}()
			}
			if in.Config.TgChannel != "" {
				b.wg.Add(1)
				go func() {
					ch := utils.WaitForMessage("RsDarkPlus239")
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
		}
		if countQueue < 2 {
			b.wg.Wait()
			b.storage.DbFunc.InsertQueue(dsmesid, alt, in.Config.CorpName, in.Username, in.UserId, in.NameMention, in.Tip, in.Lvlkz, in.Timekz, tgmesid, numkzN)
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
					ch := utils.WaitForMessage("RsDarkPlus287")
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
						dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Ds.Avatar)
						if u.User1.Tip == ds {
							//go b.sendDmDark(b.getText(in, "go"), u.User1.Mention)
							go b.client.Ds.SendDmText(b.getText(in, "go"), u.User1.UserId)
						}
						if u.User2.Tip == ds {
							//go b.sendDmDark(b.getText(in, "go"), u.User2.Mention)
							go b.client.Ds.SendDmText(b.getText(in, "go"), u.User2.UserId)
						}
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
					ch := utils.WaitForMessage("RsDarkPlus330")
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
			b.storage.DbFunc.InsertQueue(dsmesid, alt, in.Config.CorpName, in.Username, in.UserId, in.NameMention, in.Tip, in.Lvlkz, in.Timekz, tgmesid, numkzN)
			err = b.storage.Update.UpdateCompliteRS(in.Lvlkz, dsmesid, tgmesid, alt, numkzL, numberevent, in.Config.CorpName)
			if err != nil {
				err = b.storage.Update.UpdateCompliteRS(in.Lvlkz, dsmesid, tgmesid, "", numkzL, numberevent, in.Config.CorpName)
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

func (b *Bot) RsSoloPlus(in models.InMessage) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.debug {
		fmt.Printf("\n\nin RsSoloPlus %+v\n", in)
	}
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
		if in.Tip == ds {
			dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, in.Ds.Avatar)
		} else {
			dsmesid = b.client.Ds.SendWebhook(text, "RsBot", in.Config.DsChannel, "")
			// dsmesid = b.client.Ds.Send(in.Config.DsChannel, text)
		}
	}
	if in.Config.TgChannel != "" {
		if in.Tip == ds {
			text = b.client.DS.ReplaceTextMessage(text, in.Config.Guildid)
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
