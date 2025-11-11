package wa

import (
	"fmt"
	"strings"
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
		fmt.Printf("text: %s %+v\n", text, messageInfo.Sender.String())

		g := b.getGroupCommunity(messageInfo)
		i := models.IncomingMessage{
			Text:        text,
			DmChat:      messageInfo.Sender.String(),
			Name:        b.getSenderName(messageInfo),
			MentionName: "@" + b.getSenderName(messageInfo),
			NameId:      messageInfo.Sender.String(),
			ChannelId:   g.ChannelId,
			GuildId:     g.GuildId,
			GuildName:   g.GuildName,
			GuildAvatar: g.GuildAvatar,
			Type:        "wa",
			Language:    "ru",
		}
		if avatarURL, exists := b.userAvatars[messageInfo.Sender.String()]; exists {
			_, url := b.SaveAvatarLocalCache(messageInfo.Sender.String(), avatarURL)
			i.Avatar = url
		}

		b.api.SendCompendiumAppRecover(i)

	}
	if strings.HasPrefix(text, ".") {
		g := b.getGroupCommunity(messageInfo)
		mes := models.ToBridgeMessage{
			Text:    text,
			Sender:  message.Info.Sender.String(),
			Tip:     "wa",
			ChatId:  g.ChannelId,
			MesId:   getMessageIdFormat(message.Info.Sender, message.Info.ID),
			GuildId: g.GuildId,
			Config: &models.Bridge2Config{
				HostRelay: g.GuildName,
			},
		}

		b.api.SendBridgeAppRecover(mes)
	}
}
