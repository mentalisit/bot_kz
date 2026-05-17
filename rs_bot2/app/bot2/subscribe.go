package bot2

import (
	"fmt"
	"rs/models"
)

// lang ok

func (b *Bot) SubscribePing(in *models.InMessageV2, rs *models.Rs) {
	tt := fmt.Sprintf("MENTION: UserId %s RsTypeLevel %s", in.UserId, rs.RsTypeLevel)
	b.client.Tg.SendChannelDelSecondRsMention(rs.Ch, tt, 1800)
}

func (b *Bot) handleSubscription(in *models.InMessageV2, rs *models.Rs, isSubscribe bool) {
	b.deleteInMessage(in)

	channelInfo := in.Config.Channels[in.Messenger.ChannelId]
	if channelInfo == nil {
		return
	}

	argRoles := b.getTextForInfo(channelInfo, "rs") + rs.GetLevelRs()
	if rs.GetTypeRs() {
		argRoles = b.getTextForInfo(channelInfo, "drs") + rs.GetLevelRs()
	}

	// Generate response messages based on subscription type
	codes := func(code int) string {
		if isSubscribe {
			switch code {
			case 0:
				return fmt.Sprintf("%s %s %s", in.GetNameMention(), b.getTextForInfo(channelInfo, "you_subscribed_to"), argRoles)
			case 1:
				return fmt.Sprintf("%s %s %s", in.GetNameMention(), b.getTextForInfo(channelInfo, "you_already_subscribed_to"), argRoles)
			case 2:
				return b.getTextForInfo(channelInfo, "error_rights_assign") + argRoles
			}
		} else {
			switch code {
			case 0:
				return fmt.Sprintf("%s %s %s", in.GetNameMention(), b.getTextForInfo(channelInfo, "you_not_subscribed_to_role"), argRoles)
			case 1:
				return fmt.Sprintf("%s %s %s", in.GetNameMention(), argRoles, b.getTextForInfo(channelInfo, "role_not_exist"))
			case 2:
				return fmt.Sprintf("%s %s %s", in.GetNameMention(), b.getTextForInfo(channelInfo, "you_unsubscribed"), argRoles)
			case 3:
				return b.getTextForInfo(channelInfo, "error_rights_remove") + argRoles
			}
		}
		return ""
	}

	// Execute subscription/unsubscription based on messenger type
	if in.Tip == ds {
		var code int
		if isSubscribe {
			code = b.client.Ds.Subscribe(in.UserId, argRoles, in.Messenger.GuildId)
		} else {
			code = b.client.Ds.Unsubscribe(in.UserId, argRoles, in.Messenger.GuildId)
		}
		b.client.Ds.SendChannelDelSecond(in.Messenger.ChannelId, codes(code), 10)
	} else if in.Tip == tg {
		var code int
		if isSubscribe {
			code = b.client.Tg.Subscribe(in.UserId, rs.RsTypeLevel, in.Messenger.GuildId)
		} else {
			code = b.client.Tg.Unsubscribe(in.UserId, rs.RsTypeLevel, in.Messenger.ChannelId)
		}
		go b.client.Tg.SendChannelDelSecond(in.Messenger.ChannelId, codes(code), 10)
	}
}
