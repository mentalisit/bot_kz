package tg

//func (t *Telegram) messagePrivatHandler(m *tgbotapi.Message) {
//
//	after, _ := strings.CutPrefix(m.Text, "%")
//	ChatId := strconv.FormatInt(m.Chat.ID, 10)
//
//	i := models.IncomingMessage{
//		Text:         after,
//		DmChat:       strconv.FormatInt(m.From.ID, 10),
//		Name:         m.From.String(),
//		MentionName:  "@" + m.From.String(),
//		NameId:       strconv.FormatInt(m.From.ID, 10),
//		Avatar:       t.GetAvatar(m.From.ID),
//		AvatarF:      "tg",
//		ChannelId:    ChatId,
//		GuildId:      strconv.FormatInt(m.Chat.ID, 10),
//		GuildName:    m.Chat.Title,
//		GuildAvatarF: "tg",
//
//		Type: "tg",
//	}
//	chat, err := t.t.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: m.Chat.ChatConfig()})
//	if err != nil {
//		t.log.Error(err.Error())
//	}
//	if chat.Photo != nil {
//		i.GuildAvatar = t.getFileLink(chat.Photo.BigFileID)
//	}
//
//	t.ChanMessage <- i
//
//}
