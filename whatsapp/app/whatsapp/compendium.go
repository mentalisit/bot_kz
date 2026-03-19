package wa

import (
	"fmt"
	"whatsapp/models"

	"go.mau.fi/whatsmeow/types"
)

func (b *Whatsapp) compendium(text string, messageInfo types.MessageInfo) {
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
