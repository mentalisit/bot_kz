package logic

import (
	"strings"
	"unicode"
)

//func (b *Bridge) Channels() (chatIdsTG, chatIdsDS []string) {
//	for _, c := range b.in.Config.Channel["tg"] {
//		if c.ChannelId != b.in.ChatId {
//			if c.ChannelId != "" {
//				chatIdsTG = append(chatIdsTG, c.ChannelId)
//			}
//		}
//	}
//	for _, d := range b.in.Config.Channel["ds"] {
//		if d.ChannelId != b.in.ChatId {
//			if d.ChannelId != "" {
//				chatIdsDS = append(chatIdsDS, d.ChannelId)
//			}
//		}
//	}
//	return chatIdsTG, chatIdsDS
//}

func (b *Bridge) GetChannels() (chatIdsTG, chatIdsDS, chatIdsWa []string) {
	for ClientType, Channels := range b.in.Config.Channel {
		for _, channel := range Channels {
			if channel.ChannelId != b.in.ChatId {
				switch ClientType {
				case "ds":
					chatIdsDS = append(chatIdsDS, channel.ChannelId)
				case "tg":
					chatIdsTG = append(chatIdsTG, channel.ChannelId)
				case "wa":
					chatIdsWa = append(chatIdsWa, channel.ChannelId)
				default:
					b.log.InfoStruct("bad type channel", ClientType)
				}
			}
		}
	}
	return chatIdsTG, chatIdsDS, chatIdsWa
}

// IsPurelyNumeric проверяет, содержит ли строка только цифры (0-9).
func IsPurelyNumeric(s string) bool {
	if len(s) == 0 {
		return false
	}
	return strings.IndexFunc(s, func(r rune) bool {
		return !unicode.IsDigit(r)
	}) == -1
}
