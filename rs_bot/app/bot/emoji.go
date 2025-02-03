package bot

import (
	"fmt"
	"rs/models"
	"strings"
)

// lang ok
func (b *Bot) emodjiadd(in models.InMessage, slot, emo string) {
	b.iftipdelete(in)
	t := b.storage.Emoji.EmojiModuleReadUsers(in.Username, in.Tip)
	if len(t.Name) == 0 {
		b.storage.Emoji.EmInsertEmpty(in.Tip, in.Username)
	}
	text := b.storage.Emoji.EmojiUpdate(in.Username, in.Tip, slot, emo)
	b.ifTipSendTextDelSecond(in, text, 20)
}
func (b *Bot) emodjis(in models.InMessage) {
	b.iftipdelete(in)
	e := b.storage.Emoji.EmojiModuleReadUsers(in.Username, in.Tip)

	text := b.getText(in, "info_set_emoji") +
		"\n1 " + e.Em1 +
		"\n2 " + e.Em2 +
		"\n3 " + e.Em3 +
		"\n4 " + e.Em4
	if in.Tip == ds {
		text += fmt.Sprintf("\n %s %s %s %s", e.Module1, e.Module2, e.Module3, e.Weapon)
	}
	b.ifTipSendTextDelSecond(in, b.getText(in, "your_emoji")+text, 60)
}
func (b *Bot) instalNick(in models.InMessage, input string) (ok bool, nick string) {
	words := strings.Fields(input)
	if len(words) >= 2 && strings.ToLower(words[0]) == "nick" {
		nick = words[1]
		ok = true
		t := b.storage.Emoji.EmojiModuleReadUsers(in.Username, in.Tip)
		if len(t.Name) == 0 {
			b.storage.Emoji.EmInsertEmpty(in.Tip, in.Username)
		}
		go b.storage.Emoji.WeaponUpdate(in.Username, in.Tip, nick)
	} else if len(words) == 1 && strings.ToLower(words[0]) == "nick" {
		go b.storage.Emoji.WeaponUpdate(in.Username, in.Tip, "")
		return true, "удалено"
	}
	return ok, nick
}
