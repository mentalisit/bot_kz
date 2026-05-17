package wa

import (
	"strings"
	"time"
	"whatsapp/models"

	"go.mau.fi/whatsmeow/types/events"
)

func (b *Whatsapp) filter(message *events.Message) {
	msg := message.Message

	// Фильтрация: nil, от меня, старые сообщения
	if msg == nil || message.Info.IsFromMe || message.Info.Timestamp.Before(b.startedAt) {
		return
	}

	var text string
	if msg.GetExtendedTextMessage() == nil {
		text = msg.GetConversation()
	} else {
		text = msg.GetExtendedTextMessage().GetText()
	}
	messageInfo := message.Info

	if strings.HasPrefix(text, "%") {
		b.compendium(text, messageInfo)
	}

	if strings.HasPrefix(text, ".") {
		b.bridgePoint(text, messageInfo)

		good2, config2 := b.checkChannelConfig2(message.Info.Chat.String())
		g := b.getGroupCommunity(messageInfo)
		in2 := models.InMessageV2{
			Text:        text,
			Tip:         "wa",
			NameNick:    "",
			Username:    b.getSenderName(messageInfo),
			UserId:      messageInfo.Sender.String(),
			NameMention: "@" + b.getSenderName(messageInfo),
			Messenger: models.Info{
				TypeMessenger:  "wa",
				MessageId:      getMessageIdFormat(messageInfo.Sender, messageInfo.ID),
				ChannelId:      messageInfo.Chat.String(),
				GuildId:        g.GuildId,
				GuildName:      g.GuildName,
				GuildAvatarUrl: g.GuildAvatar,
				Language:       "ru",
				CreatedAt:      time.Now(),
			},
			Config:  models.CorporationConfigV2{},
			Options: models.Options{models.OptionInClient},
		}
		if avatarURL, exists := b.userAvatars[messageInfo.Sender.String()]; exists {
			_, url := b.SaveAvatarLocalCache(messageInfo.Sender.String(), avatarURL)
			in2.Messenger.UserAvatarUrl = url
		}

		if good2 {
			in2.Config = config2
		}
		if text == ".setup" || strings.HasPrefix(text, ".invite ") {

		}

		b.api.SendRsBotV2AppRecover(in2)
	}
}
