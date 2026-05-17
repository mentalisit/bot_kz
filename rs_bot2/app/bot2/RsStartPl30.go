package bot2

import (
	"fmt"
	"rs/models"
)

//lang ok

func (b *Bot) RsStart(in *models.InMessageV2, rs *models.Rs) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.deleteInMessage(in)

	// Fetch all necessary queue state in one request
	active, count, _, numberkz, err := b.storage.GetQueueState(in.UserId, rs.RsTypeLevel, in.Config.Uid)
	if err != nil {
		b.log.ErrorErr(err)
		return
	}

	// 1. Проверка активности пользователя в очереди
	if !active {
		text := b.getText(in, "info_forced_start_available")
		b.sendTextAfterDeleteSecond(in, text, 10)
		return
	}

	// 2. Get queue data
	u, err := b.storage.GetActiveQueueByCorpAndLevel(in.Config.Uid, rs.RsTypeLevel)
	if err != nil {
		b.log.ErrorErr(err)
		return
	}
	if len(u) == 0 {
		return
	}
	rs.U = u

	messages, errId := b.storage.ReadQueueMessages(in.Config.Uid, rs.RsTypeLevel)
	if errId == nil {
		rs.QueueMessages = messages
	}

	rs.Info = in.Config.Channels[in.Messenger.ChannelId]
	rs.Ch = in.Messenger.ChannelId

	rs.CountQueue = count
	rs.NumberLevel = numberkz
	rs.TextQueueCompliteBonus = b.GetTextQueueComplite(&in.Config, rs.GetTypeRs())
	rs.NameRsLevel = rs.GetRsNameLevel(b.getLanguageText)

	if count > 0 {

		// Send start message to all channels
		for channelId, channelInfo := range in.Config.Channels {
			b.wg.Add(1)
			go func(chId string, info *models.Info) {
				defer b.wg.Done()
				chRs := *rs
				chRs.Ch = chId
				chRs.Info = info
				queue := rs.GetTitle(b.getLanguageText)
				textStart := fmt.Sprintf("🚀 %s (%d) %s\n\n",
					queue, numberkz, b.getTextForInfo(info, "was_launched_incomplete"))

				textEnd := fmt.Sprintf("\n%s \n%s", b.getTextForInfo(info, "go"), rs.TextQueueCompliteBonus)

				text := fmt.Sprintf("%s%s\n%s", textStart, b.getListUsers(&u, chRs.Info, true), textEnd)

				if chRs.Info.TypeMessenger == ds {
					b.client.Ds.SendWebhook(text, "RsBot", chRs.Ch, in.Messenger.UserAvatarUrl)
				} else if chRs.Info.TypeMessenger == tg {
					b.client.Tg.SendChannel(chRs.Ch, text)
				}

			}(channelId, channelInfo)
		}
		go func() {
			if len(rs.QueueMessages) > 0 {
				b.deleteQueueActiveMessage(rs.QueueMessages)
				b.storage.DeleteQueueMessages(in.Config.Uid, rs.RsTypeLevel)
			}
		}()

		// Wait for all messages to be sent
		b.wg.Wait()
		go b.SendOtherCorporationsPercent(in.Config)

		// Update queue as completed - move all users to completed queue
		_ = b.storage.MoveToCompleteQueue(rs)

		// Increment queue count for statistics
		err = b.storage.IncrementQueueCount(in.Config.Uid, rs.RsTypeLevel)
		if err != nil {
			b.log.ErrorErr(err)
		}

		fmt.Printf("RsStart users %+v\n", u)
		// Remove users from other queues
		go b.removeUsersFromOtherQueues(u)
	}
}

func (b *Bot) Pl30(in *models.InMessageV2, rs *models.Rs) {
	b.deleteInMessage(in)

	// Check if user is in any queue
	countName, err := b.storage.CountActiveQueueByUser(in.UserId, in.Config.Uid)
	if err != nil {
		b.log.ErrorErr(err)
		return
	}

	text := ""
	if countName == 0 {
		text = in.GetNameMention() + b.getText(in, "you_out_of_queue")
	} else if countName > 0 {
		// Get user's queue to check time
		u, err := b.storage.GetActiveQueueByUserAndCorp(in.UserId, in.Config.Uid)
		if err != nil || len(u) == 0 {
			text = in.GetNameMention() + b.getText(in, "you_out_of_queue")
		} else {
			// Check if user has enough time (30+ minutes)
			if u[0].RemainingTime >= 30 {
				text = fmt.Sprintf("%s %s",
					in.GetNameMention(), b.getText(in, "current_time_is_enough"))
			} else {
				text = in.GetNameMention() + b.getText(in, "timer_updated")

				// Update remaining time to 30 minutes
				err = b.storage.UpdateActiveQueueRemainingTime(u[0].ID, u[0].RemainingTime+30)
				if err != nil {
					b.log.ErrorErr(err)
				}

				// Show updated queue
				in.Options.Add(models.OptionPl30)
				in.Options.Add(models.OptionEdit)
				b.QueueLevel(in, rs)
			}
		}
	}

	b.sendTextAfterDeleteSecond(in, text, 20)
}
