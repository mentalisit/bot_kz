package TelegramClient

//
//func (t *Telegram) accessChatTg(m *tgbotapi.Message) {
//	after, res := strings.CutPrefix(m.Text, ".")
//	ThreadID := m.MessageThreadID
//	if !m.IsTopicMessage && m.MessageThreadID != 0 {
//		ThreadID = 0
//	}
//	ChatId := strconv.FormatInt(m.Chat.ID, 10) + fmt.Sprintf("/%d", ThreadID)
//	if res {
//		mId := strconv.Itoa(m.MessageID)
//		switch after {
//		case "add":
//			go t.DelMessageSecond(ChatId, mId, 10)
//			t.accessAddChannelTg(ChatId, "en", m)
//		case "добавить":
//			go t.DelMessageSecond(ChatId, mId, 10)
//			t.accessAddChannelTg(ChatId, "ru", m)
//		case "додати":
//			go t.DelMessageSecond(ChatId, mId, 10)
//			t.accessAddChannelTg(ChatId, "ua", m)
//		case "del":
//			go t.DelMessageSecond(ChatId, mId, 10)
//			t.accessDelChannelTg(ChatId, m)
//		case "удалить":
//			go t.DelMessageSecond(ChatId, mId, 10)
//			t.accessDelChannelTg(ChatId, m)
//		case "видалити":
//			go t.DelMessageSecond(ChatId, mId, 10)
//			t.accessDelChannelTg(ChatId, m)
//		case "паника":
//			t.log.Panic("перезагрузка по требованию")
//		default:
//			if t.setLang(m, ChatId) {
//				return
//			}
//		}
//	}
//}
//func (t *Telegram) accessAddChannelTg(chatid, lang string, m *tgbotapi.Message) { // внесение в дб и добавление в масив
//	ok, _ := t.checkChannelConfigTG(chatid)
//	if ok {
//		go t.SendChannelDelSecond(chatid, t.getLanguage(lang, "info_activation_not_required"), 20)
//	} else {
//		chatName := t.chatName(chatid)
//		if m.IsTopicMessage && m.ReplyToMessage != nil && m.ReplyToMessage.ForumTopicCreated != nil {
//			chatName = fmt.Sprintf(" %s/%s", chatName, m.ReplyToMessage.ForumTopicCreated.Name)
//		}
//		t.addTgCorpConfig(chatName, chatid, lang)
//		t.log.Info("новая активация корпорации " + chatName)
//		go t.SendChannelDelSecond(chatid, t.getLanguage(lang, "tranks_for_activation"), 60)
//	}
//}
//func (t *Telegram) accessDelChannelTg(chatid string, m *tgbotapi.Message) { //удаление с бд и масива для блокировки
//	ok, config := t.checkChannelConfigTG(chatid)
//	if !ok {
//		go t.SendChannelDelSecond(chatid, t.getLanguage("ru", "channel_not_connected"), 60)
//	} else {
//		t.storage.ConfigRs.DeleteConfigRs(config)
//		t.storage.ReloadDbArray()
//		t.corpConfigRS = t.storage.CorpConfigRS
//		chatName := t.chatName(chatid)
//		if m.IsTopicMessage && m.ReplyToMessage != nil && m.ReplyToMessage.ForumTopicCreated != nil {
//			chatName = fmt.Sprintf(" %s/%s", chatName, m.ReplyToMessage.ForumTopicCreated.Name)
//		}
//		t.log.Info("отключение корпорации " + chatName)
//		go t.SendChannelDelSecond(chatid, t.getLanguage(config.Country, "you_disabled_bot_functions"), 60)
//	}
//}
//func (t *Telegram) setLang(m *tgbotapi.Message, chatid string) bool {
//	re := regexp.MustCompile(`^\.set lang (ru|en|ua)$`)
//	matches := re.FindStringSubmatch(m.Text)
//	if len(matches) > 0 {
//		langUpdate := matches[1]
//		ok, config := t.checkChannelConfigTG(chatid)
//		if ok {
//			go t.DelMessageSecond(chatid, strconv.Itoa(m.MessageID), 10)
//			config.Country = langUpdate
//			t.corpConfigRS[config.CorpName] = config
//			t.storage.ConfigRs.AutoHelpUpdateMesid(config)
//			go t.SendChannelDelSecond(chatid, t.getLanguage(config.Country, "language_switched_to"), 20)
//			t.log.Info(fmt.Sprintf("замена языка в %s на %s", config.CorpName, config.Country))
//		}
//
//		return true
//	}
//	return false
//}
