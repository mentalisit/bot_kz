package bot

import (
	"context"
	"fmt"
	"kz_bot/models"
	"time"
)

//lang ok

func (b *Bot) SubscribePing(in models.InMessage, tipPing int) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	men := b.storage.Subscribe.SubscribePing(ctx, in.NameMention, in.Lvlkz, in.Config.CorpName, tipPing, in.Config.TgChannel)
	if len(men) > 0 {
		men = fmt.Sprintf("%s%s\n%s", b.getText(in, "call_rs"), in.Lvlkz, men)
		go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, men, 600)
	}
}

func (b *Bot) Subscribe(in models.InMessage, tipPing int) {
	b.iftipdelete(in)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
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
		subscribeCode := b.client.Ds.Subscribe(in.Ds.Nameid, argRoles, in.Config.Guildid)
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
		counts := b.storage.Subscribe.CheckSubscribe(ctx, in.Name, in.Lvlkz, in.Config.TgChannel, tipPing)
		if counts == 1 {
			text := fmt.Sprintf("%s %s%s %d/4\n %s %s+",
				in.NameMention, b.getText(in, "you_subscribed_to_rs"), in.Lvlkz, tipPing, b.getText(in, "to_add_to_queue_post"), in.Lvlkz)
			go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 10)
		} else {
			//добавление в оочередь пинга
			b.storage.Subscribe.Subscribe(ctx, in.Name, in.NameMention, in.Lvlkz, tipPing, in.Config.TgChannel)
			text := fmt.Sprintf("%s %s%s %d/4 \n %s %s+",
				in.NameMention, b.getText(in, "you_subscribed_to_rs_ping"), in.Lvlkz, tipPing, b.getText(in, "to_add_to_queue_post"), in.Lvlkz)
			go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 10)
		}
	}
}
func (b *Bot) Unsubscribe(in models.InMessage, tipPing int) {
	b.iftipdelete(in)

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	if in.Tip == ds {
		d, result := containsSymbolD(in.Lvlkz)
		argRoles := b.getText(in, "rs") + in.Lvlkz
		if d {
			argRoles = b.getText(in, "drs") + result
		}
		if tipPing == 3 {
			argRoles = b.getText(in, "rs") + in.Lvlkz + "+"
		}
		unsubscribeCode := b.client.Ds.Unsubscribe(in.Ds.Nameid, argRoles, in.Config.Guildid)
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
		counts := b.storage.Subscribe.CheckSubscribe(ctx, in.Name, in.Lvlkz, in.Config.TgChannel, tipPing)
		if counts == 0 {
			text = fmt.Sprintf("%s %s%s %d/4", in.NameMention, b.getText(in, "you_not_subscribed_to_rs_ping"), in.Lvlkz, tipPing)
		} else if counts == 1 {
			//удаление с базы данных
			text = fmt.Sprintf("%s %s%s %d/4", in.NameMention, b.getText(in, "you_unsubscribed_from_rs_ping"), in.Lvlkz, tipPing)
			b.storage.Subscribe.Unsubscribe(ctx, in.Name, in.Lvlkz, in.Config.TgChannel, tipPing)
		}
		b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 10)
	}
}
