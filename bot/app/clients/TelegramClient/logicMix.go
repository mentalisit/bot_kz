package TelegramClient

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"kz_bot/clients/helper"
	"kz_bot/clients/restapi"
	"kz_bot/models"
	"strconv"
	"strings"
	"time"
)

func (t *Telegram) logicMix(m *tgbotapi.Message, edit bool) {
	go t.imHere(m.Chat.ID, m.Chat)
	//t.accessChatTg(m) //это была начальная функция при добавлени бота в группу
	ThreadID := m.MessageThreadID
	if !m.IsTopicMessage && ThreadID != 0 {
		ThreadID = 0
	}
	ChatId := strconv.FormatInt(m.Chat.ID, 10) + fmt.Sprintf("/%d", ThreadID)

	////TODO на будущее если захочу реализацию
	//if m.From != nil && m.From.LanguageCode != "" {
	//	t.log.Info(m.From.LanguageCode)
	//}

	if strings.HasPrefix(m.Text, "%") {
		t.sendToCompendiumFilter(m, ChatId)
		return
	}

	if strings.HasPrefix(m.Text, ".") {
		t.ifPrefixPoint(m)
		return
	}

	// RsClient
	ok, config := t.checkChannelConfigTG(ChatId)
	if ok {
		t.sendToRsFilter(m, config, ChatId)
	}

	tg, bridgeConfig := t.bridgeCheckChannelConfigTg(ChatId)
	if tg {
		chatName := m.Chat.Title
		if m.IsTopicMessage && m.ReplyToMessage != nil && m.ReplyToMessage.ForumTopicCreated != nil {
			chatName = fmt.Sprintf("%s/%s", chatName, m.ReplyToMessage.ForumTopicCreated.Name)
		}

		go func() {
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
					ChatId: ChatId,
					Extra:  []models.FileInfo{},
					Config: &bridgeConfig,
				}
				t.filterNewBridge(m, mes)
			}
		}()

	}
}
func (t *Telegram) sendToRsFilter(m *tgbotapi.Message, config models.CorporationConfig, ChatId string) {
	name := t.nameNick(m.From.UserName, m.From.FirstName, config.TgChannel)
	in := models.InMessage{
		Mtext:       m.Text,
		Tip:         "tg",
		Username:    name,
		UserId:      strconv.FormatInt(m.From.ID, 10),
		NameNick:    "", //нет способа извлечь ник кроме member.CustomTitle
		NameMention: "@" + name,
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
	if in.Mtext == "" && config.Forward {
		t.DelMessageSecond(ChatId, strconv.Itoa(m.MessageID), 180)
	}

	t.ChanRsMessage <- in
}
func (t *Telegram) sendToCompendiumFilter(m *tgbotapi.Message, ChatId string) {
	i := models.IncomingMessage{
		Text:        m.Text,
		DmChat:      strconv.FormatInt(m.From.ID, 10),
		Name:        m.From.String(),
		MentionName: "@" + m.From.String(),
		NameId:      strconv.FormatInt(m.From.ID, 10),
		NickName:    "", //нет способа извлечь ник кроме member.CustomTitle
		Avatar:      t.loadAvatarIsExist(m.From.ID),
		//AvatarF:      "tg",
		ChannelId: ChatId,
		GuildId:   strconv.FormatInt(m.Chat.ID, 10),
		GuildName: m.Chat.Title,
		//GuildAvatarF: "tg",
		Type: "tg",
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
		i.Language = helper.DetectLanguage(chatName)
	}
	err = restapi.SendCompendiumApp(i)
	if err != nil {
		t.log.InfoStruct("SendCompendiumApp", i)
		t.log.ErrorErr(err)
		return
	}
}
func (t *Telegram) filterNewBridge(m *tgbotapi.Message, mes models.ToBridgeMessage) {
	mes.Text = m.Text
	mes.Tip = "tg"
	mes.MesId = strconv.Itoa(m.MessageID)
	mes.GuildId = strconv.FormatInt(m.Chat.ID, 10)
	mes.TimestampUnix = m.Time().Unix()
	mes.Sender = m.From.String()
	mes.Avatar = t.getAvatarIsExist(m.From.ID)

	err := t.handleDownloadBridge(&mes, m)
	if err != nil {
		t.log.ErrorErr(err)
	}

	// handle forwarded messages
	t.handleForwarded(&mes, m)

	// quote the previous message
	t.handleQuoting(&mes, m)

	if mes.Text != "" || len(mes.Extra) > 0 {
		err = restapi.SendBridgeApp(mes)
		if err != nil {
			t.log.ErrorErr(err)
		}
	}
}
func (t *Telegram) ifPrefixPoint(m *tgbotapi.Message) {
	ThreadID := m.MessageThreadID
	if !m.IsTopicMessage && m.MessageThreadID != 0 {
		ThreadID = 0
	}
	ChatId := strconv.FormatInt(m.Chat.ID, 10) + fmt.Sprintf("/%d", ThreadID)
	chatName := t.chatName(ChatId)
	if m.IsTopicMessage && m.ReplyToMessage != nil && m.ReplyToMessage.ForumTopicCreated != nil {
		chatName = fmt.Sprintf(" %s/%s", chatName, m.ReplyToMessage.ForumTopicCreated.Name)
	}
	good, config := t.checkChannelConfigTG(ChatId)
	in := models.InMessage{
		Mtext:       m.Text,
		Tip:         "tg",
		Username:    m.From.String(),
		UserId:      strconv.FormatInt(m.From.ID, 10),
		NameMention: "@" + m.From.String(),
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
	t.ChanRsMessage <- in
	go func() {
		time.Sleep(5 * time.Second)
		t.corpConfigRS = t.storage.CorpConfigRS
	}()
	go func() {
		mes := models.ToBridgeMessage{
			Text:    m.Text,
			Sender:  m.From.String(),
			Tip:     "tg",
			ChatId:  ChatId,
			MesId:   strconv.Itoa(m.MessageID),
			GuildId: chatName,
			Config: &models.BridgeConfig{
				HostRelay: chatName,
			},
		}
		err := restapi.SendBridgeApp(mes)
		if err != nil {
			t.log.ErrorErr(err)
			return
		}
		time.Sleep(3 * time.Second)
		t.storage.ReloadDbArray()
	}()

}
