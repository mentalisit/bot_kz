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

func (b *Bot) updateCompendiumModules(in models.InMessage) {
	b.iftipdelete(in)
	b.ifTipSendMentionText(in, b.helpers.UpdateCompendiumModules(in))
}
func (b *Bot) checkAdmin(in models.InMessage) bool {
	admin := false
	if in.Tip == ds {
		admin = b.client.Ds.CheckAdmin(in.Ds.Nameid, in.Config.DsChannel)
	} else if in.Tip == tg {
		admin = b.client.Tg.CheckAdminTg(in.Config.TgChannel, in.Username)
	}
	return admin
}

func (b *Bot) elseChat(user []string) { //проверяем всех игроков этой очереди на присутствие в других очередях или корпорациях
	user = utils.RemoveDuplicates(user)
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
				Username:    t.Name,
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
					m := b.client.Ds.SendWebhook(text2, in.Username, in.Config.DsChannel, in.Config.Guildid, in.Ds.Avatar)
					b.client.Ds.DeleteMesageSecond(in.Config.DsChannel, m, 90)
				}()
			} else if in.Tip == tg {
				go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text2, 90)
			}
		}
	}

}
