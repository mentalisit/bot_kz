package bot

import (
	"fmt"
	gt "github.com/bas24/googletranslatefree"
	"rs/models"
	"rs/pkg/utils"
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
	var err error
	if in.Tip == ds {
		admin = b.client.Ds.CheckAdmin(in.UserId, in.Config.DsChannel)
	} else if in.Tip == tg {
		admin, err = b.client.Tg.CheckAdminTg(in.Config.TgChannel, in.Username)
		if err != nil {
			b.log.ErrorErr(err)
		}
	} else if in.Username == "Mentalisit" || in.Username == "mentalisit" {
		admin = true
	}
	return admin
}

func (b *Bot) elseChat(user []string) { //проверяем всех игроков этой очереди на присутствие в других очередях или корпорациях
	ch := utils.WaitForMessage("SendChannelDelSecond")
	defer close(ch)
	user = utils.RemoveDuplicates(user)
	for _, u := range user {
		if b.storage.Count.CountNameQueue(u) > 0 {
			b.elsetrue(u)
		}
	}
}
func (b *Bot) elsetrue(userid string) { //удаляем игрока с очереди
	tt := b.storage.DbFunc.ElseTrue(userid)
	for _, t := range tt {
		ok, config := b.CheckCorpNameConfig(t.Corpname)
		if ok {
			in := models.InMessage{
				Mtext:       t.Lvlkz + "-",
				Tip:         t.Tip,
				Username:    t.Name,
				UserId:      t.UserId,
				NameMention: t.Mention,
				Lvlkz:       t.Lvlkz,
				Ds: struct {
					Mesid   string
					Guildid string
					Avatar  string
				}{
					Mesid:   t.Dsmesid,
					Guildid: ""},
				Tg: struct {
					Mesid int
				}{
					Mesid: t.Tgmesid,
				},
				Config: config,
				Option: models.Option{Elsetrue: true},
			}
			b.Inbox <- in
		}
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
	} else {
		result = s
	}

	return dark, result
}

func (b *Bot) Transtale(in models.InMessage) {
	text2, err := gt.Translate(in.Mtext, "auto", in.Config.Country)
	if err == nil {
		if in.Mtext != text2 {
			if in.Tip == ds {
				go func() {
					ch := utils.WaitForMessage("Translate")
					m := b.client.Ds.SendWebhook(text2, in.Username, in.Config.DsChannel, in.Ds.Avatar)
					b.client.Ds.DeleteMessageSecond(in.Config.DsChannel, m, 90)
					close(ch)
				}()
			} else if in.Tip == tg {
				go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text2, 90)
			}
		}
	}

}
