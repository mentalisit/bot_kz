package bot2

import (
	"fmt"
	"rs/models"
)

// lang ok
func (b *Bot) emojiAdd(in *models.InMessageV2, slot, emo string) {
	b.deleteInMessage(in)
	t := b.storage.EmojiReadUUID(in.MAcc.UUID, in.Tip)
	if t == nil {
		b.storage.EmojiInsertEmptyUUID(in.MAcc.UUID, in.Tip)
	}
	text := b.storage.EmojiUpdateUUID(in.MAcc.UUID, in.Tip, slot, emo)
	b.sendTextAfterDeleteSecond(in, text, 20)
}
func (b *Bot) emojis(in *models.InMessageV2) {
	b.deleteInMessage(in)
	e := b.storage.EmojiReadUUID(in.MAcc.UUID, in.Tip)
	if e == nil {
		b.storage.EmojiInsertEmptyUUID(in.MAcc.UUID, in.Tip)
		e = &models.Emoji{}
	}

	channelInfo := in.Config.Channels[in.Messenger.ChannelId]
	if channelInfo == nil {
		return
	}

	text := b.getTextForInfo(channelInfo, "info_set_emoji") +
		"\n1 " + e.Em1 +
		"\n2 " + e.Em2 +
		"\n3 " + e.Em3 +
		"\n4 " + e.Em4

	if in.Tip == ds {
		text += fmt.Sprintf("\n %s %s %s %s", e.Em1, e.Em2, e.Em3, e.Em4)
	}
	b.sendTextAfterDeleteSecond(in, b.getTextForInfo(channelInfo, "your_emoji")+text, 60)
}
