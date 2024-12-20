package logic

func (b *Bridge) Channels() (chatIdsTG, chatIdsDS []string) {
	for _, c := range b.in.Config.ChannelTg {
		if c.ChannelId != b.in.ChatId {
			if c.ChannelId != "" {
				chatIdsTG = append(chatIdsTG, c.ChannelId)
			}
		}
	}
	for _, d := range b.in.Config.ChannelDs {
		if d.ChannelId != b.in.ChatId {
			if d.ChannelId != "" {
				chatIdsDS = append(chatIdsDS, d.ChannelId)
			}
		}
	}
	return chatIdsTG, chatIdsDS
}
