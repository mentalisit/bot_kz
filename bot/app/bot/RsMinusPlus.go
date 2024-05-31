package bot

import (
	"context"
	"fmt"
	"kz_bot/models"
	"time"
)

//lang ok

func (b *Bot) Plus(in models.InMessage) bool {
	b.iftipdelete(in)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	countName := b.storage.Count.CountNameQueueCorp(ctx, in.Name, in.Config.CorpName)
	message := ""
	ins := false
	if countName > 0 && in.Option.Reaction {
		b.iftipdelete(in)
		t := b.storage.Timers.UpdateMitutsQueue(ctx, in.Name, in.Config.CorpName)
		if t.Timedown > 3 {
			message = fmt.Sprintf("%s %s%s %s %d%s",
				t.Mention, b.getText(in, "info_cannot_click_plus"), t.Lvlkz, b.getText(in, "you_will_still"), t.Timedown, b.getText(in, "min"))
		} else if t.Timedown <= 3 {
			ins = true
			message = t.Mention + b.getText(in, "timer_updated")
			in.Lvlkz = t.Lvlkz
			in.Option.Reaction = false
			go b.QueueLevel(in)
		}
		b.ifTipSendTextDelSecond(in, message, 10)
		//if in.Tip == ds {
		//	b.client.Ds.DeleteMessage(in.Config.DsChannel, in.Ds.Mesid)
		//} else if in.Tip == tg {
		//	b.client.Tg.DelMessage(in.Config.TgChannel, in.Tg.Mesid)
		//}
	}
	return ins
}
func (b *Bot) Minus(in models.InMessage) bool {
	bb := false
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	countNames := b.storage.Count.CountNameQueueCorp(ctx, in.Name, in.Config.CorpName)
	if countNames > 0 && in.Option.Reaction {
		b.iftipdelete(in)
		t := b.storage.Timers.UpdateMitutsQueue(ctx, in.Name, in.Config.CorpName)
		if t.Name == in.Name && t.Timedown > 3 {
			message := fmt.Sprintf("%s %s%s %s %d%s",
				t.Mention, b.getText(in, "info_cannot_click_minus"), t.Lvlkz, b.getText(in, "you_will_still"), t.Timedown, b.getText(in, "min"))
			b.ifTipSendTextDelSecond(in, message, 10)
		} else if t.Name == in.Name && t.Timedown <= 3 {
			in.Lvlkz = t.Lvlkz
			bb = true
			in.Option.Reaction = false
			in.Option.Update = true
			go b.RsMinus(in)
		}
		if in.Tip == ds {
			b.client.Ds.DeleteMessage(in.Config.DsChannel, in.Ds.Mesid)
		} else if in.Tip == tg {
			b.client.Tg.DelMessage(in.Config.TgChannel, in.Tg.Mesid)
		}

	}
	return bb
}
