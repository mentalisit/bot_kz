package bot2

import (
	"rs/models"
	"strconv"
)

func (b *Bot) MinusMin() {
	tt, err := b.storage.MinusMin()
	if err != nil || len(tt) == 0 {
		return
	}

	for _, t := range tt {
		if t.Data.CorporationUuid != "" {
			// Get config by corporation name
			config := b.storage.ReadConfigV2Uid(t.Data.CorporationUuid)
			if config == nil {
				continue
			}
			// Create InMessageV2 from QueueActive
			in := models.InMessageV2{
				Text:        "",
				Tip:         t.Data.Tip,
				Username:    t.Data.Name,
				UserId:      t.Data.UserID,
				NameMention: t.Data.Mention,
				Messenger: models.Info{
					TypeMessenger: t.Data.Tip,
				},
				Config: *config,
			}
			in.Options.Add(models.OptionMinusMin)
			in.Options.Add(models.OptionEdit)

			// Find the correct channel ID for this messenger type
			var channelId string
			for _, channel := range config.Channels {
				if channel.TypeMessenger == t.Data.Tip {
					channelId = channel.ChannelId
					break
				}
			}
			in.Messenger.ChannelId = channelId

			// Create Rs from QueueData
			rs := models.NewRs()
			rs.RsTypeLevel = t.Data.LvlRS
			rs.TimeRs = t.Data.Time

			if t.RemainingTime == 3 {
				for _, info := range config.Channels {
					if info.Game != nil && info.Game.Alias == t.Data.Alias || info.GuildName == t.Data.Alias {
						text := in.GetNameMention() + b.getLanguageText(info.Language, "info_time_almost_up")
						if t.Data.Tip == "ds" {
							mID := b.client.Ds.SendEmbedTime(info.ChannelId, text)
							go b.client.Ds.DeleteMessageSecond(info.ChannelId, mID, 180)
						} else if t.Data.Tip == "tg" {
							mID := b.client.Tg.SendEmbedTime(info.ChannelId, text)
							go b.client.Tg.DelMessageSecond(info.ChannelId, strconv.Itoa(mID), 180)
						}
					}
				}
			} else if t.RemainingTime <= 0 {
				b.RsMinus(&in, rs)
			}

		}
	}

	// Update all active queues after time decrement
	b.updateAllActiveQueues()
}

// updateAllActiveQueues обновляет все активные очереди
func (b *Bot) updateAllActiveQueues() {
	// Get all active queues
	activeQueues, err := b.storage.GetActiveQueue()
	if err != nil || len(activeQueues) == 0 {
		return
	}

	// Get unique corporation names
	corpMap := make(map[string]bool)
	for _, queue := range activeQueues {
		if queue.Data.CorporationUuid != "" {
			corpMap[queue.Data.CorporationUuid] = true
		}
	}

	// Update each corporation's queues
	for corpName := range corpMap {
		// Get config for corporation
		config := b.storage.ReadConfigV2Uid(corpName)
		if config == nil {
			continue
		}

		// Create InMessageV2 for queue update
		in := &models.InMessageV2{
			Text:    "",
			Options: []string{models.OptionMinusMinNext, models.OptionEdit},
			Config:  *config,
		}

		// Get all active queue levels for this corporation
		levels, err := b.storage.GetActiveQueueLevelsByCorp(corpName)
		if err != nil {
			continue
		}

		// Update each queue level
		for _, level := range levels {
			if level != "" {
				rs := models.NewRs()
				rs.RsTypeLevel = level

				// Set channel info from the first channel
				for _, channelInfo := range config.Channels {
					rs.Info = channelInfo
					break
				}

				// Update queue display
				go b.QueueLevel(in, rs)
			}
		}
	}
}

//func (b *Bot) MinusMinMessageUpdate() {
//	corpActive0 := b.storage.DbFunc.OneMinutsTimer()
//	for _, corp := range corpActive0 {
//
//		_, config := b.CheckCorpNameConfig(corp)
//
//		dss, tgs := b.storage.DbFunc.MessageUpdateMin(corp)
//
//		if config.DsChannel != "" {
//			for _, d := range dss {
//				b.Inbox <- b.storage.DbFunc.MessageUpdateDS(d, config)
//			}
//		}
//		if config.TgChannel != "" {
//			for _, t := range tgs {
//				b.Inbox <- b.storage.DbFunc.MessageUpdateTG(t, config)
//			}
//		}
//	}
//
//}

//func (b *Bot) ReadQueueLevel(in *models.InMessageV2, second int) {
//	ch := utils.WaitForMessage("ReadQueueLevel")
//	defer close(ch)
//	text, err := b.otherQueue.ReadingQueueByLevel(in.RsTypeLevel, in.Config.CorpName)
//	if err != nil {
//		b.log.ErrorErr(err)
//		return
//	}
//
//	if text != "" {
//		b.sendTextAfterDeleteSecond(in, text, second)
//	}
//}
