package bot2

import (
	"fmt"
	"rs/models"
	"strconv"
	"sync"
)

func (b *Bot) RsDarkPlus(in *models.InMessageV2, rs *models.Rs) {
	b.deleteInMessage(in)

	if in.MAcc != nil && in.MAcc.ActiveAccount != "" && in.MAcc.ActiveAccount != in.MAcc.Nickname {
		rs.AltName = in.MAcc.ActiveAccount
	}

	// 1. Single DB Read for initial state
	active, countQueue, numberName, numberLevel, err := b.storage.GetQueueState(in.UserId, rs.RsTypeLevel, in.Config.Uid)
	if err != nil {
		b.log.ErrorErr(err)
		return
	}

	if active {
		b.SendMentionTextDel(in, b.getText(in, "you_in_queue"))
		return
	}

	rs.CountQueue = countQueue
	rs.NumberName = numberName
	rs.NumberLevel = numberLevel

	// 2. Prepare user record and fetch current queue members if any
	userIn := b.createQueueUser(in, rs)
	if rs.CountQueue != 0 {
		rs.U, _ = b.storage.GetActiveQueueByCorpAndLevel(in.Config.Uid, rs.RsTypeLevel)
		messages, errId := b.storage.ReadQueueMessages(in.Config.Uid, rs.RsTypeLevel)
		if errId == nil {
			rs.QueueMessages = messages
		}
	}
	if rs.CountQueue == 2 {
		rs.TextQueueCompliteBonus = b.GetTextQueueComplite(&in.Config, true)
	} else if rs.CountQueue == 3 {
		rs.TextQueueCompliteBonus = b.GetTextQueueComplite(&in.Config, false)
	}

	// 3. Process each channel in parallel
	var localWg sync.WaitGroup

	for ch, info := range in.Config.Channels {
		// Scoped copies for goroutines
		chRs := *rs
		chRs.Ch = ch
		chRs.Info = info
		chRs.NameRsLevel = chRs.GetRsNameLevel(b.Dictionary.GetText)
		chRs.LevelRsPingDs = b.getDsPing(&chRs)

		localWg.Add(1)
		go func(r models.Rs, u *models.QueueActive) {
			defer localWg.Done()
			b.processQueueLogic(in, &r, u)
		}(chRs, userIn)
	}

	localWg.Wait()
	close(rs.MessageIdsChan)
	for msg := range rs.MessageIdsChan {
		for ch, m := range msg {
			if rs.QueueMessages == nil {
				rs.QueueMessages = make(map[string]models.QueueMessages)
			}
			rs.QueueMessages[ch] = m
		}
	}

	if rs.CountQueue == 0 || rs.CountQueue == 1 || (rs.CountQueue == 2 && !rs.GetTypeRs()) {
		if err = b.storage.InsertActiveQueue(*userIn); err != nil {
			b.log.ErrorErr(err)
			return
		}
		if len(rs.QueueMessages) != 0 {
			err = b.storage.SaveQueueMessages(in.Config.Uid, rs.RsTypeLevel, rs.QueueMessages)
			if err != nil {
				b.log.ErrorErr(err)
			}
		}
		// Refresh view
		b.QueueLevel(in, rs)
	} else {
		// Queue complete (DRS 3/3 or RS 4/4)
		go b.SendOtherCorporationsPercent(in.Config)
		b.storage.MoveToCompleteQueue(rs)
		b.SendLsNotification(in, rs.U)
		fmt.Printf("RsDarkPlus users %+v\n", rs.U)
		_ = b.storage.IncrementQueueCount(in.Config.Uid, rs.RsTypeLevel)
		go b.removeUsersFromOtherQueues(rs.U)
	}
}

// processQueueLogic determines which handler to call based on current queue status
func (b *Bot) processQueueLogic(in *models.InMessageV2, rs *models.Rs, userIn *models.QueueActive) {
	if rs.Info == nil {
		rs.Info = in.Config.Channels[in.Messenger.ChannelId]
	}

	youJoinedQueue := b.getTextForInfo(rs.Info, "you_joined_queue")

	var text string

	sendDel10second := func(text string) {
		if rs.Info.TypeMessenger == ds {
			go b.client.Ds.SendChannelDelSecond(rs.Ch, text, 10)
		} else if rs.Info.TypeMessenger == tg {
			go b.client.Tg.SendChannelDelSecond(rs.Ch, text, 10)
		}
	}

	if rs.CountQueue == 0 {
		text = fmt.Sprintf(b.getTextForInfo(rs.Info, "temp_queue_started"), in.Username, rs.LevelRsPingDs)
		var mid string
		if rs.Info != nil && rs.Info.TypeMessenger == ds {
			mid = b.client.Ds.SendComplexContent(rs.Ch, text)
		} else if rs.Info != nil && rs.Info.TypeMessenger == tg {
			b.client.Tg.SendChannelDelSecond(rs.Ch, text, 10)
			go b.SubscribePing(in, rs)
		}
		if mid != "" {
			if rs.MessageIdsChan != nil {
				rs.MessageIdsChan <- map[string]models.QueueMessages{
					rs.Ch: {TypeMessenger: rs.Info.TypeMessenger, MessageID: mid},
				}
			}
		}
		return
	}
	if rs.CountQueue == 1 {
		if !rs.GetTypeRs() {
			text = fmt.Sprintf("%s 2/4 %s %s", rs.LevelRsPingDs, in.Username, youJoinedQueue)
		} else {
			text = fmt.Sprintf("%s 2/3 %s %s", rs.LevelRsPingDs, in.Username, youJoinedQueue)
		}
		sendDel10second(text)
		return
	}

	sendComplete := func(text1, text string) {
		rs.U = append(rs.U, *userIn)
		if rs.Info != nil {
			b.deleteQueueActiveMessage(rs.QueueMessages)

			var mid string
			if rs.Info.TypeMessenger == ds {
				b.client.Ds.SendChannelDelSecond(rs.Ch, text1, 10)
				mid = b.client.Ds.SendWebhook(text, "RsBot", rs.Ch, in.Messenger.UserAvatarUrl)
			} else if rs.Info.TypeMessenger == tg {
				b.client.Tg.SendChannelDelSecond(rs.Ch, text1, 10)
				midInt := b.client.Tg.SendChannel(rs.Ch, text)
				mid = strconv.Itoa(midInt)
			}
			if mid != "" {
				if rs.MessageIdsChan != nil {
					rs.MessageIdsChan <- map[string]models.QueueMessages{
						rs.Ch: {TypeMessenger: rs.Info.TypeMessenger, MessageID: mid},
					}
				}
			}
		}
	}

	if rs.CountQueue == 2 {
		if rs.GetTypeRs() {
			text1 := fmt.Sprintf("🚀 3/3 %s %s", in.Username, youJoinedQueue)
			text = fmt.Sprintf("🚀 3/3 %s %s\n\n%s\n\n %s\n %s",
				rs.NameRsLevel, b.getTextForInfo(rs.Info, "queue_completed"),
				b.getListUsers(&rs.U, rs.Info, true), b.getTextForInfo(rs.Info, "go"),
				rs.TextQueueCompliteBonus)

			sendComplete(text1, text)
		} else {
			text = fmt.Sprintf("%s 3/4 %s %s", rs.LevelRsPingDs, in.Username, youJoinedQueue)
			sendDel10second(text)
		}
	}
	if rs.CountQueue == 3 {
		if rs.GetTypeRs() {
			b.log.InfoStruct("fack ", rs)
			return
		}

		n1, n2, n3, n4 := b.helpers.NameMention(&rs.U, rs.Info.TypeMessenger)
		text1 := fmt.Sprintf(" 4/4 %s %s", in.Username, youJoinedQueue)
		text = fmt.Sprintf("4/4 %s%s %s\n %s\n %s\n %s\n%s %s\n%s",
			b.getTextForInfo(rs.Info, "rs_queue"), rs.GetRsNameLevel(b.getLanguageText),
			b.getTextForInfo(rs.Info, "queue_completed"),
			n1, n2, n3, n4, b.getTextForInfo(rs.Info, "go"), rs.TextQueueCompliteBonus)

		sendComplete(text1, text)
	}
}

func (b *Bot) SendLsNotification(in *models.InMessageV2, u []models.QueueActive) {
	channelInfo := in.Config.Channels[in.Messenger.ChannelId]
	if channelInfo == nil {
		return
	}

	dmText := fmt.Sprintf("%s\n", b.getTextForInfo(channelInfo, "go"))
	for _, user := range u {
		dmText += user.Data.Name + "\n"
	}
	for _, user := range u {
		if user.Data.MAcc != nil &&
			user.Data.MAcc.Data != nil &&
			user.Data.MAcc.Data.NotifyPM &&
			user.Data.UserID != in.UserId {
			if user.Data.Tip == ds {
				go b.client.Ds.SendDmText(dmText, user.Data.UserID)
			} else if user.Data.Tip == tg {
				go b.client.Tg.SendChannelDelSecond(user.Data.UserID, dmText, 1800)
			}
		}
	}
}
