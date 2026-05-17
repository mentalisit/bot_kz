package bot2

import (
	"fmt"
	"rs/models"
	"rs/pkg/utils"
	"strconv"
)

func (b *Bot) QueueAll(in *models.InMessageV2) {
	b.deleteInMessage(in)
	// Get all active queue levels for the corporation
	lvl, err := b.storage.GetActiveQueueLevelsByCorp(in.Config.Uid)
	if err != nil {
		b.log.ErrorErr(err)
		return
	}

	// Remove duplicates
	lvl = utils.RemoveDuplicates(lvl)

	if len(lvl) > 0 {
		for _, level := range lvl {
			if level != "" {
				// Create new Rs for each level
				levelRs := models.NewRs()
				levelRs.SetLevelRsOrDrs(in, level)

				// Set channel info from the first channel
				for _, channelInfo := range in.Config.Channels {
					levelRs.Info = channelInfo
					break
				}

				// Handle update auto help option
				if in.Options.Contains(models.OptionUpdateAutoHelp) {
					in.Options.Remove(models.OptionUpdateAutoHelp)
					if len(in.Options) > 0 {
						fmt.Printf("opt %+v\n", in.Options)
					}
				}

				// Add queue all option if not present
				if !in.Options.Contains(models.OptionQueueAll) {
					in.Options.Add(models.OptionQueueAll)
				}

				// Send to inbox for processing
				b.QueueLevel(in, levelRs)
			}
		}
	} else if in.Options.Contains(models.OptionInClient) {
		b.sendTextAfterDeleteSecond(in, b.getText(in, "no_active_queues"), 10)
		b.deleteInMessage(in)
	}
}

func (b *Bot) QueueLevel(in *models.InMessageV2, rs *models.Rs) {
	b.deleteInMessage(in)

	_, rs.CountQueue, _, rs.NumberLevel, _ = b.storage.GetQueueState(in.UserId, rs.RsTypeLevel, in.Config.Uid)

	// совподения количество условие
	if rs.CountQueue == 0 {
		if !in.Options.Contains(models.OptionQueue) && !in.Options.Contains(models.OptionQueueAll) {
			text := rs.GetTitle(b.getLanguageText) + b.getText(in, "empty")
			b.sendTextAfterDeleteSecond(in, text, 10)
		} else if in.Options.Contains(models.OptionQueue) || in.Options.Contains(models.OptionQueueAll) {
			b.sendTextAfterDeleteSecond(in, b.getText(in, "no_active_queues"), 10)
		}
		return
	}
	rs.U, _ = b.storage.GetActiveQueueByCorpAndLevel(in.Config.Uid, rs.RsTypeLevel)
	oldMessages, errId := b.storage.ReadQueueMessages(in.Config.Uid, rs.RsTypeLevel)
	if errId == nil {
		rs.QueueMessages = oldMessages
	}

	rs.MessageIdsChan = make(chan map[string]models.QueueMessages, 50)

	dsFunc := func(chId string, channelInfo models.Info) {
		defer b.wg.Done()
		chRs := *rs
		chRs.Ch = chId
		chRs.Info = &channelInfo
		chRs.NameRsLevel = chRs.GetRsNameLevel(b.Dictionary.GetText)
		chRs.LevelRsPingDs = b.getDsPing(&chRs)

		// Create a copy of discordMap for this channel using channel-specific language
		dsMap := b.getMapForDsChannel(chRs)
		dsMap["listUsers"] = b.getListUsers(&chRs.U, chRs.Info, false)

		// Handle edit option
		if in.Options.Contains(models.OptionEdit) && len(rs.U) > 0 {
			msgID, exists := rs.QueueMessages[chId]
			if exists {
				err := b.client.Ds.EditComplexButton(msgID.MessageID, chId, dsMap)
				if err != nil {
					b.log.Info(fmt.Sprintf("QueueLevelSendQueueMessage_ds %s %s %s\n%+v\n", msgID.MessageID, chId, in.Config.Uid, err))
					in.Options.Remove(models.OptionEdit)
				}
			}
		}

		// Send new message
		if !in.Options.Contains(models.OptionEdit) {
			newMessageId := b.client.Ds.SendComplex(chId, dsMap)
			if newMessageId != "" {
				rs.MessageIdsChan <- map[string]models.QueueMessages{
					chId: {TypeMessenger: channelInfo.TypeMessenger, MessageID: newMessageId},
				}
			}
			// Delete old message
			if len(rs.U) > 0 {
				if msgID, exists := rs.QueueMessages[chId]; exists {
					go b.client.Ds.DeleteMessage(chId, msgID.MessageID)
				}
			}
		}
	}
	tgFunc := func(chId string, channelInfo models.Info) {
		defer b.wg.Done()
		chRs := *rs
		chRs.Ch = chId
		chRs.Info = &channelInfo
		chRs.NameRsLevel = chRs.GetRsNameLevel(b.Dictionary.GetText)
		chRs.LevelRsPingDs = b.getDsPing(&chRs)

		_, level := chRs.TypeRedStar()
		textStart := fmt.Sprintf("%s (%d)\n\n", chRs.GetTitle(b.getLanguageText), chRs.NumberLevel)
		textEnd := fmt.Sprintf("\n%s++ - %s", level, b.getTextForInfo(chRs.Info, "forced_start"))
		users := b.getListUsers(&rs.U, chRs.Info, false)

		text := fmt.Sprintf("%s%s\n%s", textStart, users, textEnd)

		// Handle edit option
		if in.Options.Contains(models.OptionEdit) && len(rs.U) > 0 {
			msgID, exists := rs.QueueMessages[chId]
			if exists {
				messageID, _ := strconv.Atoi(msgID.MessageID)
				if messageID == 0 {
					in.Options.Remove(models.OptionEdit)
				} else {
					err := b.client.Tg.EditMessageTextKey(chId, messageID, text, chRs.GetLevelRs())
					if err != nil {
						b.log.Info(fmt.Sprintf("QueueLevelSendQueueMessage_tg %d %s %s\n%+v\n", messageID, chId, in.Config.Uid, err))
						in.Options.Remove(models.OptionEdit)
					}
				}
			}
		}

		// Send new message
		if !in.Options.Contains(models.OptionEdit) {
			newMessageId := b.client.Tg.SendEmbed(chRs.GetLevelRs(), chId, text)
			if newMessageId != 0 {
				rs.MessageIdsChan <- map[string]models.QueueMessages{
					chId: {TypeMessenger: channelInfo.TypeMessenger, MessageID: strconv.Itoa(newMessageId)},
				}
			}
			// Delete old message
			if len(rs.U) > 0 {
				if msgID, exists := rs.QueueMessages[chId]; exists {
					go b.client.Tg.DeleteMessage(chId, msgID.MessageID)
				}
			}
		}
	}

	if rs.CountQueue == 1 || rs.CountQueue == 2 || (rs.CountQueue == 3 && !rs.GetTypeRs()) {
		for channelId, info := range in.Config.Channels {
			if info.TypeMessenger == ds {
				b.wg.Add(1)
				go dsFunc(channelId, *info)

			} else if info.TypeMessenger == tg {
				b.wg.Add(1)
				go tgFunc(channelId, *info)
			} else {
				b.log.Info("todo for " + info.TypeMessenger)
			}
		}
	} else {
		// DRS queue with 3 players - complete
		b.log.Info(fmt.Sprintf("DRS queue complete with 3 players: %+v\n", rs.U))
	}

	b.wg.Wait()
	close(rs.MessageIdsChan)

	messageIds := make(map[string]models.QueueMessages)
	for mes := range rs.MessageIdsChan {
		for ch, m := range mes {
			messageIds[ch] = m
		}
	}
	if len(messageIds) != 0 {
		changed := b.UpdateQueueMessages(rs, messageIds)
		if changed {
			if err := b.storage.UpdateQueueMessages(in.Config.Uid, rs.RsTypeLevel, rs.QueueMessages); err != nil {
				b.log.ErrorErr(err)
			}
		}
	}
}
