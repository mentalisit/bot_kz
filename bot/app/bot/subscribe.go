package bot

import (
	"context"
	"fmt"
	"time"
)

//lang ok

func (b *Bot) SubscribePing(tipPing int) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	men := b.storage.Subscribe.SubscribePing(ctx, b.in.NameMention, b.in.Lvlkz, b.in.Config.CorpName, tipPing, b.in.Config.TgChannel)
	if len(men) > 0 {
		men = fmt.Sprintf("%s%s\n%s", b.getText("SborNaKz"), b.in.Lvlkz, men)
		go b.client.Tg.SendChannelDelSecond(b.in.Config.TgChannel, men, 600)
	}
}

func (b *Bot) Subscribe(tipPing int) {
	if b.debug {
		fmt.Println("in Subscribe", b.in)
	}
	b.iftipdelete()
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	if b.in.Tip == ds {
		//go b.Ds.DeleteMessage(b.in.Config.DsChannel, b.in.Ds.Mesid)
		d, result := containsSymbolD(b.in.Lvlkz)
		argRoles := b.getText("kz") + b.in.Lvlkz
		if d {
			argRoles = b.getText("dkz") + result
		}
		if tipPing == 3 {
			argRoles = b.getText("kz") + b.in.Lvlkz + "+"
		}
		subscribeCode := b.client.Ds.Subscribe(b.in.Ds.Nameid, argRoles, b.in.Config.Guildid)
		var text string
		if subscribeCode == 0 {
			text = fmt.Sprintf("%s %s %s", b.in.NameMention, b.getText("teperViPodpisani"), argRoles)
		} else if subscribeCode == 1 {
			text = fmt.Sprintf("%s %s %s", b.in.NameMention, b.getText("ViUjePodpisan"), argRoles)
		} else if subscribeCode == 2 {
			text = b.getText("oshibkaNedostatochno") + argRoles
			b.log.Info(fmt.Sprintf("%+v %+v", b.in, b.in.Config))
		}
		b.client.Ds.SendChannelDelSecond(b.in.Config.DsChannel, text, 10)

	} else if b.in.Tip == tg {
		//проверка активной подписки
		counts := b.storage.Subscribe.CheckSubscribe(ctx, b.in.Name, b.in.Lvlkz, b.in.Config.TgChannel, tipPing)
		if counts == 1 {
			text := fmt.Sprintf("%s %s%s %d/4\n %s %s+",
				b.in.NameMention, b.getText("tiUjePodpisanNaKz"), b.in.Lvlkz, tipPing, b.getText("dlyaDobavleniyaVochered"), b.in.Lvlkz)
			go b.client.Tg.SendChannelDelSecond(b.in.Config.TgChannel, text, 10)
		} else {
			//добавление в оочередь пинга
			b.storage.Subscribe.Subscribe(ctx, b.in.Name, b.in.NameMention, b.in.Lvlkz, tipPing, b.in.Config.TgChannel)
			text := fmt.Sprintf("%s %s%s %d/4 \n %s %s+",
				b.in.NameMention, b.getText("viPodpisalisNaPing"), b.in.Lvlkz, tipPing, b.getText("dlyaDobavleniyaVochered"), b.in.Lvlkz)
			go b.client.Tg.SendChannelDelSecond(b.in.Config.TgChannel, text, 10)
		}
	}
}
func (b *Bot) Unsubscribe(tipPing int) {
	if b.debug {
		fmt.Println("in Unsubscribe", b.in)
	}
	b.iftipdelete()

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	if b.in.Tip == ds {
		d, result := containsSymbolD(b.in.Lvlkz)
		argRoles := b.getText("kz") + b.in.Lvlkz
		if d {
			argRoles = b.getText("dkz") + result
		}
		if tipPing == 3 {
			argRoles = b.getText("kz") + b.in.Lvlkz + "+"
		}
		unsubscribeCode := b.client.Ds.Unsubscribe(b.in.Ds.Nameid, argRoles, b.in.Config.Guildid)
		text := ""
		if unsubscribeCode == 0 {
			text = fmt.Sprintf("%s %s %s", b.in.NameMention, b.getText("viNePodpisani"), argRoles)
		} else if unsubscribeCode == 1 {
			text = fmt.Sprintf("%s %s %s", b.in.NameMention, argRoles, b.getText("netTakoiRoli"))
		} else if unsubscribeCode == 2 {
			text = fmt.Sprintf("%s %s %s", b.in.NameMention, b.getText("ViOtpisalis"), argRoles)
		} else if unsubscribeCode == 3 {
			text = b.getText("OshibkaNedostatochnadlyaS") + argRoles
			b.log.Info(fmt.Sprintf("%+v %+v", b.in, b.in.Config))
		}
		b.client.Ds.SendChannelDelSecond(b.in.Config.DsChannel, text, 10)
	} else if b.in.Tip == tg {
		//go b.Tg.DelMessage(b.in.Config.TgChannel, b.in.Tg.Mesid)
		//проверка активной подписки
		var text string
		counts := b.storage.Subscribe.CheckSubscribe(ctx, b.in.Name, b.in.Lvlkz, b.in.Config.TgChannel, tipPing)
		if counts == 0 {
			text = fmt.Sprintf("%s %s%s %d/4", b.in.NameMention, b.getText("tiNePodpisanNaPingKz"), b.in.Lvlkz, tipPing)
		} else if counts == 1 {
			//удаление с базы данных
			text = fmt.Sprintf("%s %s%s %d/4", b.in.NameMention, b.getText("otpisalsyaOtPingaKz"), b.in.Lvlkz, tipPing)
			b.storage.Subscribe.Unsubscribe(ctx, b.in.Name, b.in.Lvlkz, b.in.Config.TgChannel, tipPing)
		}
		b.client.Tg.SendChannelDelSecond(b.in.Config.TgChannel, text, 10)
	}
}
