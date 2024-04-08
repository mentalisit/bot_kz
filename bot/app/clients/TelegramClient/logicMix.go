package TelegramClient

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"kz_bot/clients/restapi"
	"kz_bot/models"
	"strconv"
	"strings"
	"time"
)

func (t *Telegram) logicMix(m *tgbotapi.Message, edit bool) {
	go t.imHere(m.Chat.ID, m.Chat)
	t.accesChatTg(m) //это была начальная функция при добавлени бота в группу
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
	}

	// RsClient
	ok, config := t.checkChannelConfigTG(ChatId)
	if ok {
		t.sendToRsFilter(m, config, ChatId)
	}

	tg, bridgeConfig := t.bridgeCheckChannelConfigTg(ChatId)
	if tg || strings.HasPrefix(m.Text, ".") {

		username := t.nameOrNick(m.From.UserName, m.From.FirstName)
		chatName := m.Chat.Title
		if m.IsTopicMessage && m.ReplyToMessage != nil && m.ReplyToMessage.ForumTopicCreated != nil {
			chatName = fmt.Sprintf("%s/%s", chatName, m.ReplyToMessage.ForumTopicCreated.Name)
		}

		if tg {
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
		//new
		if strings.HasPrefix(m.Text, ".") {
			go func() {
				mes := models.ToBridgeMessage{
					Text:    m.Text,
					Sender:  username,
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
				t.bridgeConfig = t.storage.BridgeConfigs
			}()
		}
	}
}
func (t *Telegram) sendToRsFilter(m *tgbotapi.Message, config models.CorporationConfig, ChatId string) {
	name := t.nameNick(m.From.UserName, m.From.FirstName, config.TgChannel)
	in := models.InMessage{
		Mtext:       m.Text,
		Tip:         "tg",
		Name:        name,
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
