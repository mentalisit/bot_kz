package telegram

import (
	"fmt"
	"strconv"
	"strings"
	"telegram/models"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func (t *Telegram) logicMix(m *tgbotapi.Message, edit bool) {
	t.SaveMember(&m.Chat, m.From)
	//go t.imHere(m.Chat.ID, m.Chat)

	ThreadID := m.MessageThreadID
	if !m.IsTopicMessage && ThreadID != 0 {
		ThreadID = 0
	}
	ChatId := strconv.FormatInt(m.Chat.ID, 10) + fmt.Sprintf("/%d", ThreadID)

	////TODO на будущее если захочу реализацию
	//if m.From != nil && m.From.LanguageCode != "" {
	//	t.log.Info(m.From.LanguageCode)
	//}

	if m.Text == "@all" {
		t.MentionAllMembers(&m.Chat, m)
	}
	if m.Text == "roles" {
		t.SendWebAppButtonSmart(m.Chat.ID)
	}

	if strings.HasPrefix(m.Text, ".") {
		go t.ifPrefixPoint(m)
		return
	}

	//RsClient
	ok, config := t.checkChannelConfigTG(ChatId)
	if ok {
		go t.sendToRsFilter(m, config, ChatId)
		return
	}

	//bridge
	tg, bridgeConfig := t.bridgeCheckChannelConfigTg(ChatId)
	if tg {
		go t.sendToBridgeFilter(m, ChatId, bridgeConfig)
	}

	//compendium
	if strings.HasPrefix(m.Text, "%") {
		go t.sendToCompendiumFilter(m, ChatId)
	}
}

func (t *Telegram) sendToRsFilter(m *tgbotapi.Message, config models.CorporationConfig, ChatId string) {
	in := models.InMessage{
		Mtext:       m.Text,
		Tip:         "tg",
		Username:    m.From.String(),
		UserId:      strconv.FormatInt(m.From.ID, 10),
		NameNick:    "", //нет способа извлечь ник кроме member.CustomTitle
		NameMention: "@" + m.From.UserName,
		Tg: struct {
			Mesid int
		}{
			Mesid: m.MessageID,
		},
		Config: config,
		Option: models.Option{
			InClient: true,
		},
	}
	if in.Mtext == "" && (m.IsTopicMessage && m.MessageThreadID != 0) {
		t.DelMessageSecond(ChatId, strconv.Itoa(m.MessageID), 600)
	}

	t.api.SendRsBotAppRecover(in)
}
func (t *Telegram) sendToCompendiumFilter(m *tgbotapi.Message, ChatId string) {
	i := models.IncomingMessage{
		Text:        m.Text,
		DmChat:      strconv.FormatInt(m.From.ID, 10),
		Name:        m.From.String(),
		MentionName: "@" + m.From.UserName,
		NameId:      strconv.FormatInt(m.From.ID, 10),
		NickName:    "", //нет способа извлечь ник кроме member.CustomTitle
		Avatar:      t.loadAvatarIsExist(m.From.ID),
		ChannelId:   ChatId,
		GuildId:     strconv.FormatInt(m.Chat.ID, 10),
		GuildName:   m.Chat.Title,
		Type:        "tg",
	}
	chat, err := t.t.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: m.Chat.ChatConfig()})
	if err != nil {
		t.log.Error(err.Error())
	}
	if chat.Photo != nil {
		fileconfig := tgbotapi.FileConfig{FileID: chat.Photo.BigFileID}
		file, _ := t.t.GetFile(fileconfig)
		if file.FileID != "" {
			_, url := t.SaveAvatarLocalCache(strconv.FormatInt(m.Chat.ID, 10), "https://api.telegram.org/file/bot"+t.t.Token+"/"+file.FilePath)
			i.GuildAvatar = url
		}
	}
	member, _ := t.t.GetChatMember(tgbotapi.GetChatMemberConfig{ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: m.Chat.ID,
		},
		UserID: m.From.ID,
	}})
	if member.CustomTitle != "" {
		i.NickName = member.CustomTitle
	}

	if chat.Location != nil && chat.Location.Address != "" {
		t.log.Info(chat.Location.Address)
	}
	if m.From != nil && m.From.LanguageCode != "" {
		i.Language = m.From.LanguageCode
	} else {
		chatName := t.chatName(ChatId)
		if m.IsTopicMessage && m.ReplyToMessage != nil && m.ReplyToMessage.ForumTopicCreated != nil {
			chatName = fmt.Sprintf(" %s/%s", chatName, m.ReplyToMessage.ForumTopicCreated.Name)
		}
		i.Language = DetectLanguage(chatName)
	}

	t.api.SendCompendiumAppRecover(i)
}
func (t *Telegram) sendToBridgeFilter(m *tgbotapi.Message, ChatId string, config models.Bridge2Config) {
	if len(m.Text) < 3500 { //игнорируем сообщения большой длины
		t.handlePoll(m)
		if m.Text == "" {
			if m.NewChatMembers != nil {
				if len(m.NewChatMembers) != 1 {
					t.log.Info(strconv.Itoa(len(m.NewChatMembers)))
				}
				member0 := m.NewChatMembers[0].String()
				m.Text = member0 + " вступил(а) в группу"
			}
			if m.LeftChatMember != nil {
				m.Text = m.LeftChatMember.String() + " покинул(а) группу"
			}
		}
		mes := models.ToBridgeMessage{
			ChatId:        ChatId,
			Extra:         []models.FileInfo{},
			Config:        &config,
			Text:          m.Text,
			Tip:           "tg",
			MesId:         strconv.Itoa(m.MessageID),
			GuildId:       strconv.FormatInt(m.Chat.ID, 10),
			TimestampUnix: m.Time().Unix(),
			Sender:        ReplaceCyrillicToLatin(m.From.String()),
			Avatar:        t.getAvatarIsExist(m.From.ID),
		}

		if m.EditDate != 0 {
			mes.Tip = "tge"
		}

		err := t.handleDownloadBridge(&mes, m)
		if err != nil {
			t.log.ErrorErr(err)
		}

		// handle forwarded messages
		t.handleForwarded(&mes, m)

		// quote the previous message
		t.handleQuoting(&mes, m)

		if mes.Text != "" || len(mes.Extra) > 0 {
			t.api.SendBridgeAppRecover(mes)
		}
	}
}

func (t *Telegram) ifPrefixPoint(m *tgbotapi.Message) {
	ThreadID := m.MessageThreadID
	if !m.IsTopicMessage && m.MessageThreadID != 0 {
		ThreadID = 0
	}
	ChatId := strconv.FormatInt(m.Chat.ID, 10) + fmt.Sprintf("/%d", ThreadID)
	chatName := m.Chat.Title
	if m.IsTopicMessage && m.ReplyToMessage != nil && m.ReplyToMessage.ForumTopicCreated != nil {
		chatName = fmt.Sprintf("%s/%s", chatName, m.ReplyToMessage.ForumTopicCreated.Name)
	}

	good, config := t.checkChannelConfigTG(ChatId)
	in := models.InMessage{
		Mtext:       m.Text,
		Tip:         "tg",
		Username:    m.From.String(),
		UserId:      strconv.FormatInt(m.From.ID, 10),
		NameMention: "@" + m.From.UserName,
		Tg: struct {
			Mesid int
		}{
			Mesid: m.MessageID,
		},
		Option: models.Option{
			InClient: true,
		},
	}
	if good {
		in.Config = config
	} else {
		in.Config = models.CorporationConfig{
			CorpName:  chatName,
			TgChannel: ChatId,
			Guildid:   "",
		}
	}
	t.api.SendRsBotAppRecover(in)

	mes := models.ToBridgeMessage{
		Text:    m.Text,
		Sender:  m.From.String(),
		Tip:     "tg",
		ChatId:  ChatId,
		MesId:   strconv.Itoa(m.MessageID),
		GuildId: chatName,
		Config: &models.Bridge2Config{
			HostRelay: chatName,
		},
	}

	t.api.SendBridgeAppRecover(mes)

	//time.Sleep(5 * time.Second)
	//t.loadConfig()
}
func ReplaceCyrillicToLatin(input string) string {
	// Карта соответствия русских и украинских букв на латинские
	cyrillicToLatin := map[rune]string{
		'А': "A", 'Б': "B", 'В': "V", 'Г': "G", 'Д': "D", 'Е': "E", 'Ё': "Yo",
		'Ж': "Zh", 'З': "Z", 'И': "I", 'Й': "Y", 'К': "K", 'Л': "L", 'М': "M",
		'Н': "N", 'О': "O", 'П': "P", 'Р': "R", 'С': "S", 'Т': "T", 'У': "U",
		'Ф': "F", 'Х': "Kh", 'Ц': "Ts", 'Ч': "Ch", 'Ш': "Sh", 'Щ': "Shch", 'Ы': "Y",
		'Э': "E", 'Ю': "Yu", 'Я': "Ya", 'Ь': "", 'Ъ': "",

		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo",
		'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
		'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
		'ф': "f", 'х': "kh", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "shch", 'ы': "y",
		'э': "e", 'ю': "yu", 'я': "ya", 'ь': "", 'ъ': "",

		// Украинские буквы
		'Є': "Ye", 'І': "I", 'Ї': "Yi", 'Ґ': "G",
		'є': "ye", 'і': "i", 'ї': "yi", 'ґ': "g",
	}

	// Строим новую строку с заменой букв
	var result strings.Builder
	for _, char := range input {
		if latin, found := cyrillicToLatin[char]; found {
			result.WriteString(latin)
		} else {
			result.WriteRune(char) // Оставляем символ без изменений, если это не кириллица
		}
	}

	return result.String()
}
