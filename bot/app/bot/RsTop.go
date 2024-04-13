package bot

import (
	"context"
	"fmt"
	"time"
)

//lang ok
//нужно переделать полностью

func (b *Bot) TopLevel() {
	b.iftipdelete()
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	numEvent := b.storage.Event.NumActiveEvent(b.in.Config.CorpName)
	if numEvent == 0 {
		mesage := fmt.Sprintf("\xF0\x9F\x93\x96 %s %s%s:\n",
			b.getText("topUchastnikov"), b.getText("kz"), b.in.Lvlkz)

		b.ifTipSendTextDelSecond(b.getText("ScanDB"), 5)
		good := b.storage.Top.TopLevel(ctx, b.in.Config.CorpName, b.in.Lvlkz)
		if !good {
			b.ifTipSendTextDelSecond(b.getText("noHistory"), 20)
		} else if good {
			b.ifTipSendTextDelSecond(b.getText("formlist"), 5)
			mest := b.storage.Top.TopTemp(ctx)
			if b.in.Tip == ds {
				m := b.client.Ds.SendEmbedText(b.in.Config.DsChannel, mesage, mest)
				b.client.Ds.DeleteMesageSecond(b.in.Config.DsChannel, m.ID, 60)
			} else if b.in.Tip == tg {
				b.client.Tg.SendChannelDelSecond(b.in.Config.TgChannel, mesage+mest, 60)
			}
		}
	} else {
		mesage := fmt.Sprintf("\xF0\x9F\x93\x96 %s %s %s%s\n     ",
			b.getText("topUchastnikov"), b.getText("iventa"), b.getText("kz"), b.in.Lvlkz)
		b.ifTipSendTextDelSecond(b.getText("ScanDB"), 5)
		good := b.storage.Top.TopEventLevel(ctx, b.in.Config.CorpName, b.in.Lvlkz, numEvent)
		if !good {
			b.ifTipSendTextDelSecond(b.getText("noHistory"), 20)
		} else {
			b.ifTipSendTextDelSecond(b.getText("formlist"), 5)
			mest := b.storage.Top.TopTempEvent(ctx)
			if b.in.Tip == ds {
				m := b.client.Ds.SendEmbedText(b.in.Config.DsChannel, mesage, mest)
				b.client.Ds.DeleteMesageSecond(b.in.Config.DsChannel, m.ID, 60)
			} else if b.in.Tip == tg {
				b.client.Tg.SendChannelDelSecond(b.in.Config.TgChannel, mesage+mest, 60)
			}
		}
	}
}
func (b *Bot) TopAll() {
	b.iftipdelete()
	ctx := context.Background()
	numEvent := b.storage.Event.NumActiveEvent(b.in.Config.CorpName)
	if numEvent == 0 {
		mesage := fmt.Sprintf("\xF0\x9F\x93\x96 %s:\n", b.getText("topUchastnikov"))
		b.ifTipSendTextDelSecond(b.getText("ScanDB"), 5)
		good := b.storage.Top.TopAll(ctx, b.in.Config.CorpName)
		if good {
			b.ifTipSendTextDelSecond(b.getText("formlist"), 5)
			message2 := b.storage.Top.TopTemp(ctx)
			if b.in.Tip == ds {
				m := b.client.Ds.SendEmbedText(b.in.Config.DsChannel, mesage, message2)
				b.client.Ds.DeleteMesageSecond(b.in.Config.DsChannel, m.ID, 60)
			} else if b.in.Tip == tg {
				b.client.Tg.SendChannelDelSecond(b.in.Config.TgChannel, mesage+message2, 60)
			}
		} else if !good {
			b.ifTipSendTextDelSecond(b.getText("noHistory"), 10)
		}
	} else if numEvent > 0 {
		mesage := fmt.Sprintf("\xF0\x9F\x93\x96 %s %s:\n",
			b.getText("topUchastnikov"), b.getText("iventa"))
		b.ifTipSendTextDelSecond(b.getText("ScanDB"), 10)
		good := b.storage.Top.TopAllEvent(ctx, b.in.Config.CorpName, numEvent)
		if good {
			b.ifTipSendTextDelSecond(b.getText("formlist"), 5)
			message2 := b.storage.Top.TopTempEvent(ctx)
			if b.in.Tip == ds {
				m := b.client.Ds.SendEmbedText(b.in.Config.DsChannel, mesage, message2)
				go b.client.Ds.DeleteMesageSecond(b.in.Config.DsChannel, m.ID, 60)
			} else if b.in.Tip == tg {
				b.client.Tg.SendChannelDelSecond(b.in.Config.TgChannel, mesage+message2, 60)
			}
		} else if !good {
			b.ifTipSendTextDelSecond(b.getText("noHistory"), 10)
		}
	}
}
