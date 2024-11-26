package bot

import (
	"fmt"
	"rs/models"
	"strings"
)

//lang ok

func (b *Bot) SubscribePing(in models.InMessage, tipPing int) {
	men := b.storage.Subscribe.SubscribePing(in.NameMention, in.Lvlkz, tipPing, in.Config.TgChannel)
	if len(men) > 0 {
		d, result := containsSymbolD(in.Lvlkz)
		lvl := b.getText(in, "rs") + result
		if d {
			lvl = b.getText(in, "drs") + result
		}
		text1 := fmt.Sprintf(b.getText(in, "call_rs"), lvl)
		text := fmt.Sprintf("%s\n%s", text1, men)
		go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 600)
	}
}

func (b *Bot) Subscribe(in models.InMessage, tipPing int) {
	b.iftipdelete(in)
	if in.Tip == ds {
		//go b.Ds.DeleteMessage(in.Config.DsChannel, in.Ds.Mesid)
		d, result := containsSymbolD(in.Lvlkz)
		argRoles := b.getText(in, "rs") + in.Lvlkz
		if d {
			argRoles = b.getText(in, "drs") + result
		}
		if tipPing == 3 {
			argRoles = b.getText(in, "rs") + in.Lvlkz + "+"
		}
		subscribeCode := b.client.Ds.Subscribe(in.UserId, argRoles, in.Config.Guildid)
		var text string
		if subscribeCode == 0 {
			text = fmt.Sprintf("%s %s %s", in.NameMention, b.getText(in, "you_subscribed_to"), argRoles)
		} else if subscribeCode == 1 {
			text = fmt.Sprintf("%s %s %s", in.NameMention, b.getText(in, "you_already_subscribed_to"), argRoles)
		} else if subscribeCode == 2 {
			text = b.getText(in, "error_rights_assign") + argRoles
			b.log.Info(fmt.Sprintf("%+v %+v", in, in.Config))
		}
		b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 10)

	} else if in.Tip == tg {
		//проверка активной подписки
		counts := b.storage.Subscribe.CheckSubscribe(in.Username, in.Lvlkz, in.Config.TgChannel, tipPing)
		if counts == 1 {
			text := fmt.Sprintf("%s %s%s %d/4\n %s %s+",
				in.NameMention, b.getText(in, "you_subscribed_to_rs"), in.Lvlkz, tipPing, b.getText(in, "to_add_to_queue_post"), in.Lvlkz)
			go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 10)
		} else {
			//добавление в оочередь пинга
			b.storage.Subscribe.Subscribe(in.Username, in.NameMention, in.Lvlkz, tipPing, in.Config.TgChannel)
			text := fmt.Sprintf("%s %s%s %d/4 \n %s %s+",
				in.NameMention, b.getText(in, "you_subscribed_to_rs_ping"), in.Lvlkz, tipPing, b.getText(in, "to_add_to_queue_post"), in.Lvlkz)
			go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 10)
		}
	}
}
func (b *Bot) Unsubscribe(in models.InMessage, tipPing int) {
	b.iftipdelete(in)

	if in.Tip == ds {
		d, result := containsSymbolD(in.Lvlkz)
		argRoles := b.getText(in, "rs") + in.Lvlkz
		if d {
			argRoles = b.getText(in, "drs") + result
		}
		if tipPing == 3 {
			argRoles = b.getText(in, "rs") + in.Lvlkz + "+"
		}
		unsubscribeCode := b.client.Ds.Unsubscribe(in.UserId, argRoles, in.Config.Guildid)
		text := ""
		if unsubscribeCode == 0 {
			text = fmt.Sprintf("%s %s %s", in.NameMention, b.getText(in, "you_not_subscribed_to_role"), argRoles)
		} else if unsubscribeCode == 1 {
			text = fmt.Sprintf("%s %s %s", in.NameMention, argRoles, b.getText(in, "role_not_exist"))
		} else if unsubscribeCode == 2 {
			text = fmt.Sprintf("%s %s %s", in.NameMention, b.getText(in, "you_unsubscribed"), argRoles)
		} else if unsubscribeCode == 3 {
			text = b.getText(in, "error_rights_remove") + argRoles
			b.log.Info(fmt.Sprintf("%+v %+v", in, in.Config))
		}
		b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 10)
	} else if in.Tip == tg {
		//go b.Tg.DelMessage(in.Config.TgChannel, in.Tg.Mesid)
		//проверка активной подписки
		var text string
		counts := b.storage.Subscribe.CheckSubscribe(in.Username, in.Lvlkz, in.Config.TgChannel, tipPing)
		if counts == 0 {
			text = fmt.Sprintf("%s %s%s %d/4", in.NameMention, b.getText(in, "you_not_subscribed_to_rs_ping"), in.Lvlkz, tipPing)
		} else if counts == 1 {
			//удаление с базы данных
			text = fmt.Sprintf("%s %s%s %d/4", in.NameMention, b.getText(in, "you_unsubscribed_from_rs_ping"), in.Lvlkz, tipPing)
			b.storage.Subscribe.Unsubscribe(in.Username, in.Lvlkz, in.Config.TgChannel, tipPing)
			//внесение информации об отказе от авто подписки tipPing 0
			b.storage.Subscribe.Subscribe(in.Username, in.NameMention, in.Lvlkz, 0, in.Config.TgChannel)
		}
		b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 10)
	}
}
func (b *Bot) CheckSubscribe(in models.InMessage) {
	if in.Tip == "ds" {
		return
	}
	if strings.HasPrefix(in.NameMention, "@@") {
		return
	}
	if strings.HasPrefix(in.Username, "$") {
		in.Username, _ = strings.CutPrefix(in.Username, "$")
	}
	drs, result := containsSymbolD(in.Lvlkz)
	if !drs {
		return
	}
	argRoles := b.getText(in, "drs") + result
	if in.Tip == tg {
		//проверка отписки после авто подписки
		counts2 := b.storage.Subscribe.CheckSubscribe(in.Username, in.Lvlkz, in.Config.TgChannel, 0)
		if counts2 != 0 {
			return
		}

		//проверка активной подписки
		counts := b.storage.Subscribe.CheckSubscribe(in.Username, in.Lvlkz, in.Config.TgChannel, 1)

		if counts > 0 {
			return
		} else {
			//добавление в оочередь пинга
			b.storage.Subscribe.Subscribe(in.Username, in.NameMention, in.Lvlkz, 1, in.Config.TgChannel)
			text := fmt.Sprintf(b.getText(in, "you_subscribed_automated"), in.NameMention, argRoles)
			go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 10)
		}
	}
}
