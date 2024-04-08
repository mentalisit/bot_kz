package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"telegram/models"
	"telegram/telegram/restapi"
)

func (t *Telegram) logicMix(m *tgbotapi.Message, edit bool) {
	//go t.imHere(m.Chat.ID, m.Chat)
	if strings.HasPrefix(m.Text, ".") {
		t.accesChatTg(m) //это была начальная функция при добавлени бота в группу
	}

	ThreadID := m.MessageThreadID
	if !m.IsTopicMessage && ThreadID != 0 {
		ThreadID = 0
	}
	ChatId := strconv.FormatInt(m.Chat.ID, 10) + fmt.Sprintf("/%d", ThreadID)

	////TODO на будущее если захочу реализацию
	//if m.From != nil && m.From.LanguageCode != "" {
	//	t.log.Info(m.From.LanguageCode)
	//}

	//compendium
	if strings.HasPrefix(m.Text, "%") {
		go t.sendToCompendiumFilter(m, ChatId)
	}

	// RsClient
	ok, config := t.checkChannelConfigTG(ChatId)
	if ok {
		go t.sendToRsFilter(m, config, ChatId)
	}

	tg, bridgeConfig := t.bridgeCheckChannelConfigTg(ChatId)
	if tg {
		go t.sendToBridgeFilter(m, ChatId, bridgeConfig)
	}
}

func (t *Telegram) sendToRsFilter(m *tgbotapi.Message, config models.CorporationConfig, ChatId string) {
	in := models.InMessage{
		Mtext:       m.Text,
		Tip:         "tg",
		Name:        m.From.String(),
		NameMention: "@" + t.nickName(m.From, config.TgChannel),
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
	err := restapi.SendRsBotApp(in)
	if err != nil {
		t.log.ErrorErr(err)
	}
}
func (t *Telegram) sendToCompendiumFilter(m *tgbotapi.Message, ChatId string) {
	i := models.IncomingMessage{
		Text:         m.Text,
		DmChat:       strconv.FormatInt(m.From.ID, 10),
		Name:         m.From.String(),
		MentionName:  "@" + m.From.String(),
		NameId:       strconv.FormatInt(m.From.ID, 10),
		Avatar:       t.getAvatarIsExist(m.From.ID),
		AvatarF:      "tg",
		ChannelId:    ChatId,
		GuildId:      strconv.FormatInt(m.Chat.ID, 10),
		GuildName:    m.Chat.Title,
		GuildAvatarF: "tg",
		Type:         "tg",
	}
	chat, err := t.t.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: m.Chat.ChatConfig()})
	if err != nil {
		t.log.Error(err.Error())
	}
	if chat.Photo != nil {
		fileconfig := tgbotapi.FileConfig{FileID: chat.Photo.BigFileID}
		file, _ := t.t.GetFile(fileconfig)
		if file.FileID != "" {
			i.GuildAvatar = "https://api.telegram.org/file/bot" + t.t.Token + "/" + file.FilePath
		}

	}

	if chat.Location != nil && chat.Location.Address != "" {
		t.log.Info(chat.Location.Address)
	}
	err = restapi.SendCompendiumApp(i)
	if err != nil {
		t.log.InfoStruct("SendCompendiumApp", i)
		t.log.ErrorErr(err)
		return
	}
}
func (t *Telegram) sendToBridgeFilter(m *tgbotapi.Message, ChatId string, config models.BridgeConfig) {
	chatName := m.Chat.Title
	if m.IsTopicMessage && m.ReplyToMessage != nil && m.ReplyToMessage.ForumTopicCreated != nil {
		chatName = fmt.Sprintf("%s/%s", chatName, m.ReplyToMessage.ForumTopicCreated.Name)
	}

	if config.HostRelay != "" {
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
					ChatId:        ChatId,
					Extra:         []models.FileInfo{},
					Config:        &config,
					Text:          m.Text,
					Tip:           "tg",
					MesId:         strconv.Itoa(m.MessageID),
					GuildId:       strconv.FormatInt(m.Chat.ID, 10),
					TimestampUnix: m.Time().Unix(),
					Sender:        m.From.String(),
					Avatar:        t.getAvatarIsExist(m.From.ID),
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
					err = restapi.SendBridgeApp(mes)
					if err != nil {
						t.log.ErrorErr(err)
					}
				}
			}
		}()
	}
}
