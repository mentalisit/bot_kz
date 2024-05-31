package bot

import (
	"context"
	"fmt"
	"kz_bot/models"
	"time"
)

//lang ok
//нужно переделать полностью

func (b *Bot) TopLevel(in models.InMessage) {
	b.iftipdelete(in)
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	numEvent := b.storage.Event.NumActiveEvent(in.Config.CorpName)
	if numEvent == 0 {
		mesage := fmt.Sprintf("\xF0\x9F\x93\x96 %s %s%s:\n",
			b.getText(in, "top_participants"), b.getText(in, "rs"), in.Lvlkz)

		b.ifTipSendTextDelSecond(in, b.getText(in, "scan_db"), 5)
		good := b.storage.Top.TopLevelPerMonth(ctx, in.Config.CorpName, in.Lvlkz)
		if !good {
			b.ifTipSendTextDelSecond(in, b.getText(in, "no_history"), 20)
		} else if good {
			b.ifTipSendTextDelSecond(in, b.getText(in, "form_list"), 5)
			mest := b.storage.Top.TopTemp(ctx)
			if in.Tip == ds {
				m := b.client.Ds.SendEmbedText(in.Config.DsChannel, mesage, mest)
				b.client.Ds.DeleteMesageSecond(in.Config.DsChannel, m.ID, 60)
			} else if in.Tip == tg {
				b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, mesage+mest, 60)
			}
		}
	} else {
		mesage := fmt.Sprintf("\xF0\x9F\x93\x96 %s %s %s%s\n     ",
			b.getText(in, "top_participants"), b.getText(in, "event"), b.getText(in, "rs"), in.Lvlkz)
		b.ifTipSendTextDelSecond(in, b.getText(in, "scan_db"), 5)
		good := b.storage.Top.TopEventLevel(ctx, in.Config.CorpName, in.Lvlkz, numEvent)
		if !good {
			b.ifTipSendTextDelSecond(in, b.getText(in, "no_history"), 20)
		} else {
			b.ifTipSendTextDelSecond(in, b.getText(in, "form_list"), 5)
			mest := b.storage.Top.TopTempEvent(ctx)
			if in.Tip == ds {
				m := b.client.Ds.SendEmbedText(in.Config.DsChannel, mesage, mest)
				b.client.Ds.DeleteMesageSecond(in.Config.DsChannel, m.ID, 60)
			} else if in.Tip == tg {
				b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, mesage+mest, 60)
			}
		}
	}
}
func (b *Bot) TopAll(in models.InMessage) {
	b.iftipdelete(in)
	ctx := context.Background()
	numEvent := b.storage.Event.NumActiveEvent(in.Config.CorpName)
	if numEvent == 0 {
		mesage := fmt.Sprintf("\xF0\x9F\x93\x96 %s:\n", b.getText(in, "top_participants"))
		b.ifTipSendTextDelSecond(in, b.getText(in, "scan_db"), 5)
		good := b.storage.Top.TopAllPerMonth(ctx, in.Config.CorpName)
		if good {
			b.ifTipSendTextDelSecond(in, b.getText(in, "form_list"), 5)
			message2 := b.storage.Top.TopTemp(ctx)
			if in.Tip == ds {
				m := b.client.Ds.SendEmbedText(in.Config.DsChannel, mesage, message2)
				b.client.Ds.DeleteMesageSecond(in.Config.DsChannel, m.ID, 60)
			} else if in.Tip == tg {
				b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, mesage+message2, 60)
			}
		} else if !good {
			b.ifTipSendTextDelSecond(in, b.getText(in, "no_history"), 10)
		}
	} else if numEvent > 0 {
		mesage := fmt.Sprintf("\xF0\x9F\x93\x96 %s %s:\n",
			b.getText(in, "top_participants"), b.getText(in, "event"))
		b.ifTipSendTextDelSecond(in, b.getText(in, "scan_db"), 10)
		good := b.storage.Top.TopAllEvent(ctx, in.Config.CorpName, numEvent)
		if good {
			b.ifTipSendTextDelSecond(in, b.getText(in, "form_list"), 5)
			message2 := b.storage.Top.TopTempEvent(ctx)
			if in.Tip == ds {
				m := b.client.Ds.SendEmbedText(in.Config.DsChannel, mesage, message2)
				go b.client.Ds.DeleteMesageSecond(in.Config.DsChannel, m.ID, 60)
			} else if in.Tip == tg {
				b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, mesage+message2, 60)
			}
		} else if !good {
			b.ifTipSendTextDelSecond(in, b.getText(in, "no_history"), 10)
		}
	}
}
