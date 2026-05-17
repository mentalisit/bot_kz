package bot2

import (
	"fmt"
	"rs/models"
)

//lang ok

func (b *Bot) Plus(in *models.InMessageV2) bool {
	countName, _ := b.storage.CountActiveQueueByUser(in.UserId, in.Config.Uid)
	channelInfo := in.Config.Channels[in.Messenger.ChannelId]
	if channelInfo == nil {
		return false
	}

	ins := false
	rs := models.NewRs()
	checkFunc := func(t models.QueueActive) string {
		message := ""
		if t.RemainingTime > 3 {
			message = fmt.Sprintf("%s %s%s %s %d%s",
				t.Data.Mention, b.getTextForInfo(channelInfo, "info_cannot_click_plus"), t.Data.LvlRS,
				b.getTextForInfo(channelInfo, "you_will_still"), t.RemainingTime, b.getTextForInfo(channelInfo, "min"))
		} else if t.RemainingTime <= 3 {
			ins = true
			message = t.Data.Mention + b.getTextForInfo(channelInfo, "timer_updated")
			rs.SetLevelRsOrDrs(in, t.Data.LvlRS)
			in.Options.Remove(models.OptionReaction)
			in.Text = ""
			in.Options.Add(models.OptionPlus)
			b.Inbox <- *in
			b.deleteInMessage(in)
		}
		return message
	}
	if countName > 0 {
		tt, _ := b.storage.GetActiveQueueByUserAndCorp(in.UserId, in.Config.Uid)
		if len(tt) != 0 {
			minimalTime := tt[0].RemainingTime
			t := tt[0]
			for _, t1 := range tt {
				if t1.RemainingTime < minimalTime {
					minimalTime = t1.RemainingTime
					t = t1
				}
			}

			b.sendTextAfterDeleteSecond(in, checkFunc(t), 10)
		}

	}

	return ins
}
func (b *Bot) Minus(in *models.InMessageV2) bool {
	bb := false
	channelInfo := in.Config.Channels[in.Messenger.ChannelId]
	if channelInfo == nil {
		return false
	}

	checkFunc := func(t models.QueueActive) {
		if t.Data.UserID == in.UserId && t.RemainingTime > 3 {
			message := fmt.Sprintf("%s %s%s %s %d%s",
				t.Data.Mention, b.getTextForInfo(channelInfo, "info_cannot_click_minus"), t.Data.LvlRS,
				b.getTextForInfo(channelInfo, "you_will_still"), t.RemainingTime, b.getTextForInfo(channelInfo, "min"))
			b.sendTextAfterDeleteSecond(in, message, 10)
		} else if t.Data.UserID == in.UserId && t.RemainingTime <= 3 {
			rs := models.NewRs()
			rs.SetLevelRsOrDrs(in, t.Data.LvlRS)
			bb = true
			in.Options.Remove(models.OptionReaction)
			in.Options.Add(models.OptionUpdate)
			b.RsMinus(in, rs)
		}
	}
	countNames, _ := b.storage.CountActiveQueueByUser(in.UserId, in.Config.Uid)
	if countNames > 0 {
		tt, _ := b.storage.GetActiveQueueByUserAndCorp(in.UserId, in.Config.Uid)
		if len(tt) != 0 {
			minimalTime := tt[0].RemainingTime
			t := tt[0]
			for _, t1 := range tt {
				if t1.RemainingTime < minimalTime {
					minimalTime = t1.RemainingTime
					t = t1
				}
			}
			checkFunc(t)
		}
	}
	return bb
}
