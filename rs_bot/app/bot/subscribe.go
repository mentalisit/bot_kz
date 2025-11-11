package bot

import (
	"fmt"
	"rs/models"
	"strings"
)

// lang ok

func (b *Bot) SubscribePing(in models.InMessage) {
	subscribes := b.storage.Subscribe.SubscribePing(models.Subscribe{
		Name:    in.Username,
		Lvlkz:   in.RsTypeLevel,
		Tip:     1,
		ChatId:  in.Config.TgChannel,
		UserId:  in.UserId,
		Mention: in.NameMention,
	})

	var MentionText string

	darkOrRed, level := in.TypeRedStar()
	lvl := b.getText(in, "rs") + level
	if darkOrRed {
		lvl = b.getText(in, "drs") + level
	}
	text1 := fmt.Sprintf(b.getText(in, "call_rs"), lvl)

	for _, s := range subscribes {
		if s.Mention == "@" {
			MentionText += fmt.Sprintf("[%s](tg://user?id=%s), ", s.Name, s.UserId)
		} else {
			MentionText = MentionText + s.Mention + ", "
		}
	}
	MentionText = strings.TrimSuffix(MentionText, ", ")

	if len(subscribes) > 0 {

		text := fmt.Sprintf("%s\n%s", text1, MentionText)
		go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 1800)
	}
}

func (b *Bot) Subscribe(in models.InMessage) {
	b.iftipdelete(in)
	darkOrRed, level := in.TypeRedStar()
	var text string
	if in.Tip == ds {
		argRoles := b.getText(in, "rs") + level
		if darkOrRed {
			argRoles = b.getText(in, "drs") + level
		}
		subscribeCode := b.client.Ds.Subscribe(in.UserId, argRoles, in.Config.Guildid)
		if subscribeCode == 0 {
			text = fmt.Sprintf("%s %s %s", in.GetNameMention(), b.getText(in, "you_subscribed_to"), argRoles)
		} else if subscribeCode == 1 {
			text = fmt.Sprintf("%s %s %s", in.GetNameMention(), b.getText(in, "you_already_subscribed_to"), argRoles)
		} else if subscribeCode == 2 {
			text = b.getText(in, "error_rights_assign") + argRoles
			b.log.Info(fmt.Sprintf("%+v %+v", in, in.Config))
		}
		b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 10)

	} else if in.Tip == tg {
		suds := models.Subscribe{
			Name:    in.Username,
			Lvlkz:   in.RsTypeLevel,
			Tip:     1,
			ChatId:  in.Config.TgChannel,
			UserId:  in.UserId,
			Mention: in.NameMention,
		}
		//проверка активной подписки
		counts := b.storage.Subscribe.CheckSubscribe(suds)
		if counts == 1 {
			text = fmt.Sprintf("%s %s%s %d/4\n %s %s+",
				in.GetNameMention(), b.getText(in, "you_subscribed_to_rs"), level, 1,
				b.getText(in, "to_add_to_queue_post"), level)
		} else {
			//добавление в оочередь пинга
			b.storage.Subscribe.Subscribe(suds)
			text = fmt.Sprintf("%s %s%s %d/4 \n %s %s+",
				in.GetNameMention(), b.getText(in, "you_subscribed_to_rs_ping"),
				level, 1, b.getText(in, "to_add_to_queue_post"), level)
		}
		go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 10)
	}
}
func (b *Bot) Unsubscribe(in models.InMessage) {
	b.iftipdelete(in)
	darkOrRed, level := in.TypeRedStar()

	if in.Tip == ds {
		argRoles := b.getText(in, "rs") + level
		if darkOrRed {
			argRoles = b.getText(in, "drs") + level
		}

		unsubscribeCode := b.client.Ds.Unsubscribe(in.UserId, argRoles, in.Config.Guildid)
		text := ""
		if unsubscribeCode == 0 {
			text = fmt.Sprintf("%s %s %s", in.GetNameMention(), b.getText(in, "you_not_subscribed_to_role"), argRoles)
		} else if unsubscribeCode == 1 {
			text = fmt.Sprintf("%s %s %s", in.GetNameMention(), argRoles, b.getText(in, "role_not_exist"))
		} else if unsubscribeCode == 2 {
			text = fmt.Sprintf("%s %s %s", in.GetNameMention(), b.getText(in, "you_unsubscribed"), argRoles)
		} else if unsubscribeCode == 3 {
			text = b.getText(in, "error_rights_remove") + argRoles
			b.log.Info(fmt.Sprintf("%+v %+v", in, in.Config))
		}
		b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 10)
	} else if in.Tip == tg {
		suds := models.Subscribe{
			Name:    in.Username,
			Lvlkz:   in.RsTypeLevel,
			Tip:     1,
			ChatId:  in.Config.TgChannel,
			UserId:  in.UserId,
			Mention: in.NameMention,
		}
		var text string
		counts := b.storage.Subscribe.CheckSubscribe(suds)
		if counts == 0 {
			text = fmt.Sprintf("%s %s%s %d/4", in.GetNameMention(),
				b.getText(in, "you_not_subscribed_to_rs_ping"), in.RsTypeLevel, 1)
		} else if counts == 1 {
			//удаление с базы данных
			text = fmt.Sprintf("%s %s%s %d/4", in.GetNameMention(),
				b.getText(in, "you_unsubscribed_from_rs_ping"), in.RsTypeLevel, 1)
			b.storage.Subscribe.Unsubscribe(suds)
			//внесение информации об отказе от авто подписки tipPing 0
			suds.Tip = 0
			b.storage.Subscribe.Subscribe(suds)
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
	drs, result := in.TypeRedStar()
	if !drs {
		return
	}
	argRoles := b.getText(in, "drs") + result
	if in.Tip == tg {
		suds := models.Subscribe{
			Name:    in.Username,
			Lvlkz:   in.RsTypeLevel,
			Tip:     0,
			ChatId:  in.Config.TgChannel,
			UserId:  in.UserId,
			Mention: in.NameMention,
		}
		//проверка отписки после авто подписки
		counts2 := b.storage.Subscribe.CheckSubscribe(suds)
		if counts2 != 0 {
			return
		}

		//проверка активной подписки
		suds.Tip = 1
		counts := b.storage.Subscribe.CheckSubscribe(suds)

		if counts > 0 {
			return
		} else {
			//добавление в оочередь пинга
			b.storage.Subscribe.Subscribe(suds)
			text := fmt.Sprintf(b.getText(in, "you_subscribed_automated"), in.GetNameMention(), argRoles)
			go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 10)
		}
	}
}
