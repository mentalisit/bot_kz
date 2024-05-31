package bot

import (
	"context"
	"fmt"
	gt "github.com/bas24/googletranslatefree"
	"kz_bot/models"
	"kz_bot/pkg/utils"
	"strings"
)

// lang
func (b *Bot) iftipdelete(in models.InMessage) {
	if in.Tip == ds && !in.Option.Reaction && !in.Option.Update && !in.Option.Edit {
		go b.client.Ds.DeleteMessage(in.Config.DsChannel, in.Ds.Mesid)
		go b.client.Ds.ChannelTyping(in.Config.DsChannel)
	} else if in.Tip == tg && !in.Option.Reaction && !in.Option.Update {
		go b.client.Tg.ChatTyping(in.Config.TgChannel)
		go b.client.Tg.DelMessage(in.Config.TgChannel, in.Tg.Mesid)
	}
}
func (b *Bot) ifTipSendMentionText(in models.InMessage, text string) {
	text = fmt.Sprintf("%s %s", in.NameMention, text)
	if in.Tip == ds {
		go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 10)
	} else if in.Tip == tg {
		go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 10)
	}
}
func (b *Bot) ifTipSendTextDelSecond(in models.InMessage, text string, time int) {
	if in.Tip == ds {
		go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, time)
	} else if in.Tip == tg {
		go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, time)
	}
}

//func (b *Bot) emReadName(in models.InMessage, name, nameMention, tip string) string { // склеиваем имя и эмоджи
//	t := b.storage.Emoji.EmojiModuleReadUsers(context.Background(), name, tip)
//	newName := name
//	if tip == ds {
//		newName = nameMention
//	} else {
//		newName = name
//	}
//
//	if len(t.Name) > 0 {
//		if tip == ds && tip == t.Tip {
//			newName = fmt.Sprintf("%s %s %s %s %s %s%s%s%s", nameMention, t.Module1, t.Module2, t.Module3, t.Weapon, t.Em1, t.Em2, t.Em3, t.Em4)
//		} else if tip == tg && tip == t.Tip {
//			newName = fmt.Sprintf("%s %s%s%s%s", name, t.Em1, t.Em2, t.Em3, t.Em4)
//			if t.Weapon != "" {
//				newName = fmt.Sprintf("%s [%s] %s%s%s%s", name, t.Weapon, t.Em1, t.Em2, t.Em3, t.Em4)
//			}
//		}
//	} else if in.Tip == ds && in.Config.Guildid == "716771579278917702" && in.Name == name {
//		genesis, enrich, rsextender := helpers.GetTechDataUserId(in.Ds.Nameid)
//		b.storage.Emoji.EmInsertEmpty(context.Background(), "ds", name)
//		one := fmt.Sprintf("<:rse:1199068829511335946> %d ", rsextender)
//		two := fmt.Sprintf("<:genesis:1199068748280242237> %d ", genesis)
//		three := fmt.Sprintf("<:enrich:1199068793633251338> %d ", enrich)
//		newName = fmt.Sprintf("%s ", nameMention)
//		if rsextender != 0 {
//			b.storage.Emoji.ModuleUpdate(context.Background(), name, "ds", "1", one)
//			newName += one
//		}
//		if genesis != 0 {
//			b.storage.Emoji.ModuleUpdate(context.Background(), name, "ds", "2", two)
//			newName += two
//		}
//		if enrich != 0 {
//			b.storage.Emoji.ModuleUpdate(context.Background(), name, "ds", "3", three)
//			newName += three
//		}
//
//	}
//	return newName
//}

func (b *Bot) emReadMention(name, nameMention, tip string) string { // склеиваем имя и эмоджи
	t := b.storage.Emoji.EmojiModuleReadUsers(context.Background(), name, tip)
	newName := nameMention

	if len(t.Name) > 0 {
		if tip == ds && tip == t.Tip {
			newName = fmt.Sprintf("%s %s %s %s %s %s%s%s%s", nameMention, t.Module1, t.Module2, t.Module3, t.Weapon, t.Em1, t.Em2, t.Em3, t.Em4)
		} else if tip == tg && tip == t.Tip {
			newName = fmt.Sprintf("%s %s%s%s%s", nameMention, t.Em1, t.Em2, t.Em3, t.Em4)
			if t.Weapon != "" {
				newName = fmt.Sprintf("%s [%s] %s%s%s%s", nameMention, t.Weapon, t.Em1, t.Em2, t.Em3, t.Em4)
			}
		}
	}
	return newName
}

func (b *Bot) updateCompendiumModules(in models.InMessage) {
	b.iftipdelete(in)
	b.ifTipSendMentionText(in, b.helpers.UpdateCompendiumModules(in))
}
func (b *Bot) checkAdmin(in models.InMessage) bool {
	admin := false
	if in.Tip == ds {
		admin = b.client.Ds.CheckAdmin(in.Ds.Nameid, in.Config.DsChannel)
	} else if in.Tip == tg {
		admin = b.client.Tg.CheckAdminTg(in.Config.TgChannel, in.Name)
	}
	return admin
}

func (b *Bot) nameMention(in models.InMessage, u models.Users, tip string) (n1, n2, n3, n4 string) {
	if u.User1.Tip == tip {
		n1 = b.emReadMention(u.User1.Name, u.User1.Mention, tip)
	} else {
		n1 = u.User1.Name
	}
	if u.User2.Tip == tip {
		n2 = b.emReadMention(u.User2.Name, u.User2.Mention, tip)
	} else {
		n2 = u.User2.Name
	}
	if u.User3.Tip == tip {
		n3 = b.emReadMention(u.User3.Name, u.User3.Mention, tip)
	} else {
		n3 = u.User3.Name
	}
	if in.Tip == tip {
		n4 = b.emReadMention(in.Name, in.NameMention, tip)
	} else {
		n4 = in.Name
	}
	return
}

func (b *Bot) elseChat(user []string) { //проверяем всех игроков этой очереди на присутствие в других очередях или корпорациях
	user = utils.RemoveDuplicateElementString(user)
	for _, u := range user {
		if b.storage.Count.CountNameQueue(context.Background(), u) > 0 {
			b.elsetrue(u)
		}
	}
}
func (b *Bot) elsetrue(name string) { //удаляем игрока с очереди
	tt := b.storage.DbFunc.ElseTrue(context.Background(), name)
	for _, t := range tt {
		ok, config := b.CheckCorpNameConfig(t.Corpname)
		if ok {
			in := models.InMessage{
				Mtext:       t.Lvlkz + "-",
				Tip:         t.Tip,
				Name:        t.Name,
				NameMention: t.Mention,
				Lvlkz:       t.Lvlkz,
				Ds: struct {
					Mesid   string
					Nameid  string
					Guildid string
					Avatar  string
				}{
					Mesid:   t.Dsmesid,
					Nameid:  "",
					Guildid: ""},
				Tg: struct {
					Mesid int
				}{
					Mesid: t.Tgmesid,
				},
				Config: config,
				Option: models.Option{Elsetrue: true},
			}
			b.inbox <- in
		}
	}
}

func (b *Bot) hhelp(in models.InMessage) {
	b.iftipdelete(in)
	if in.Tip == ds {
		go b.client.Ds.Help(in.Config.DsChannel, in.Config.Country)
	} else if in.Tip == tg {
		go b.client.Tg.Help(in.Config.TgChannel, in.Config.Country)
	}
}

func (b *Bot) getText(in models.InMessage, key string) string {
	return b.storage.Dictionary.GetText(in.Config.Country, key)
}

func (b *Bot) getLanguageText(lang, key string) string {
	return b.storage.Dictionary.GetText(lang, key)
}

func containsSymbolD(s string) (dark bool, result string) {
	for _, char := range s {
		if char == 'd' {
			dark = true
		}
	}
	if dark {
		result = strings.Replace(s, "d", "", -1)
	}

	return dark, result
}

func (b *Bot) Transtale(in models.InMessage) {
	text2, err := gt.Translate(in.Mtext, "auto", in.Config.Country)
	if err == nil {
		if in.Mtext != text2 {
			if in.Tip == ds {
				go func() {
					m := b.client.Ds.SendWebhook(text2, in.Name, in.Config.DsChannel, in.Config.Guildid, in.Ds.Avatar)
					b.client.Ds.DeleteMesageSecond(in.Config.DsChannel, m, 90)
				}()
			} else if in.Tip == tg {
				go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text2, 90)
			}
		}
	}

}
