package telegram

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"telegram/models"
	"time"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func (t *Telegram) ChatTyping(chatId string) error {
	chatid, threadID := t.chat(chatId)
	typingConfig := tgbotapi.NewChatAction(chatid, tgbotapi.ChatTyping)
	typingConfig.MessageThreadID = threadID
	_, err := t.t.Send(typingConfig)
	return err
}
func (t *Telegram) SendChannelDelSecond(chatid string, text string, parsemode string, second int) (bool, error) {
	if strings.HasPrefix(text, "MENTION: UserId") {
		return t.SendChannelDelSecondRsMention(chatid, text, parsemode, second)
	}
	mId, err1 := t.SendChannel(chatid, text, parsemode)
	if err1 != nil {
		return false, err1
	}
	tu := int(time.Now().UTC().Unix())
	t.Storage.Db.TimerInsert(models.Timer{
		Tip:    "tg",
		ChatId: chatid,
		MesId:  mId,
		Timed:  tu + second,
	})

	if mId != "" {
		return true, nil
	}
	return false, nil
}
func (t *Telegram) SendChannel(chatid, text, parseMode string) (string, error) {
	chatId, threadID := t.chat(chatid)
	m := tgbotapi.NewMessage(chatId, text)
	m.MessageThreadID = threadID
	found, newText := models.FindTelegramMentions(text)
	if found || parseMode == "MarkdownV2" {
		m.ParseMode = "MarkdownV2"
		m.Text = newText
	}
	tMessage, err := t.t.Send(m)
	if err != nil {
		if err.Error() != tgbotapi.ErrAPIForbidden {
			fmt.Println(err)
			//t.log.ErrorErr(err)
			fmt.Printf("chatid '%s', text '%s'", chatid, m.Text)
		}
		return "", err
	}
	return strconv.Itoa(tMessage.MessageID), nil
}
func (t *Telegram) SendChannelReply(chatid, text, parseMode string, replyId int) (string, error) {
	chatId, threadID := t.chat(chatid)
	m := tgbotapi.NewMessage(chatId, text)
	m.MessageThreadID = threadID
	m.ReplyParameters.MessageID = replyId
	found, newText := models.FindTelegramMentions(text)
	if found || parseMode == "MarkdownV2" {
		m.ParseMode = "MarkdownV2"
		m.Text = newText
	}
	tMessage, err := t.t.Send(m)
	if err != nil {
		if err.Error() != tgbotapi.ErrAPIForbidden {
			fmt.Println(err)
			//t.log.ErrorErr(err)
			fmt.Printf("chatid '%s', text '%s'", chatid, m.Text)
		}
		return "", err
	}
	return strconv.Itoa(tMessage.MessageID), nil
}
func (t *Telegram) SendPic(chatID, text string, imageBytes []byte) (string, error) {
	chatid, threadID := t.chat(chatID)
	msg := tgbotapi.NewPhoto(chatid, tgbotapi.FileReader{
		Name:   "image.jpg",
		Reader: bytes.NewReader(imageBytes),
	})
	msg.MessageThreadID = threadID
	msg.Caption = text

	// Отправляем изображение
	m, err := t.t.Send(msg)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(m.MessageID), nil
}
func (t *Telegram) SendPicScoreboard(chatID, text, fileNameScoreboard string) (mid string, err error) {
	chatid, threadID := t.chat(chatID)
	open, err := os.Open("docker/scoreboard/" + fileNameScoreboard)
	if err != nil {
		return "", err
	}

	msg := tgbotapi.NewPhoto(chatid, tgbotapi.FileReader{
		Name:   fileNameScoreboard,
		Reader: open,
	})
	msg.MessageThreadID = threadID
	msg.Caption = text

	// Отправляем изображение
	mes, err := t.t.Send(msg)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(mes.MessageID), nil
}
func (t *Telegram) SendBridgeFuncRest(in models.BridgeSendToMessenger) []models.MessageIds {
	var messageIds []models.MessageIds
	for _, chat := range in.ChannelId {

		chatId, threadID := t.chat(chat)
		var replyId int
		if in.ReplyMap != nil && in.ReplyMap[chat] != "" {
			found := in.ReplyMap[chat]
			replyId, _ = strconv.Atoi(found)
		}

		if len(in.Extra) > 0 {
			mid, err := t.sendFileExtra(in.Extra, in.Text, chat, replyId)
			if err != nil {
				t.log.Info(fmt.Sprintf("err %s %s %s\n", err.Error(), in.Text, in.Sender))
			} else {
				messageData := models.MessageIds{
					MessageId: mid,
					ChatId:    chat,
				}
				messageIds = append(messageIds, messageData)
			}
		} else {
			if strings.Contains(in.Text, "@&") {
				in.Text = strings.ReplaceAll(in.Text, "&", "")
			}

			m := tgbotapi.NewMessage(chatId, in.Text)
			m.MessageThreadID = threadID
			m.ReplyParameters.MessageID = replyId
			tMessage, err := t.t.Send(m)
			if err != nil {
				t.log.InfoStruct(fmt.Sprintf("err %+v\n", err), in)
			} else {
				messageData := models.MessageIds{
					MessageId: strconv.Itoa(tMessage.MessageID),
					ChatId:    chat,
				}
				messageIds = append(messageIds, messageData)
			}

			re := regexp.MustCompile(`@\S+`)
			mentions := re.FindAllString(m.Text, -1)
			if len(mentions) > 0 {
				roles, _ := t.Storage.Db.GetChatsRoles(context.Background(), chatId)
				if roles != nil {
					rolesMap, rolesId := makeRolesMap(roles)
					us := make(map[string][]models.User)
					users, _ := t.Storage.Db.GetChatUsers(context.Background(), chatId)

					for _, mention := range mentions {
						roleName := strings.ToLower(mention[1:]) // убираем только "@" и регистр

						if roleID, exists := rolesMap[roleName]; exists {
							for _, user := range users {
								if _, hasRole := user.Roles[roleID]; hasRole {
									roleDisplayName := rolesId[roleID]
									if us[roleDisplayName] == nil {
										us[roleDisplayName] = []models.User{}
									}
									us[roleDisplayName] = append(us[roleDisplayName], user)
								}
							}
						}
					}

					if len(us) > 0 {
						t.MentionMembersRoles(chat, tMessage.MessageID, us)
					}
				}
			}

			if strings.Contains(in.Text, "@&Русь") || strings.Contains(in.Text, "@everyone") {
				pinConfig := tgbotapi.NewPinChatMessage(chatId, tMessage.MessageID, false)
				_, _ = t.t.Send(pinConfig)
			}

		}
	}
	return messageIds
}
func (t *Telegram) sendFileExtra(extra []models.FileInfo, text, chatID string, replyId int) (string, error) {
	if extra != nil {
		if len(extra) > 0 {
			chatId, threadID := t.chat(chatID)
			var media []tgbotapi.InputMedia
			for _, f := range extra {
				if f.URL == "" && len(f.Data) == 0 && f.FileID == "" {
					continue
				}
				var fileRequestData tgbotapi.RequestFileData

				if f.FileID != "" {
					fileRequestData = tgbotapi.FileID(f.FileID)
				} else if f.URL != "" {
					fileRequestData = tgbotapi.FileURL(f.URL)
				} else if len(f.Data) > 0 {
					fileRequestData = tgbotapi.FileBytes{
						Name:  f.Name,
						Bytes: f.Data,
					}
				}

				switch filepath.Ext(f.Name) {
				case ".jpg", ".jpe", ".png", ".jpeg":
					pc := tgbotapi.NewInputMediaPhoto(fileRequestData)
					pc.Caption = text
					media = append(media, &pc)
				case ".mp4", ".m4v":
					vc := tgbotapi.NewInputMediaVideo(fileRequestData)
					vc.Caption = text
					media = append(media, &vc)
				case ".mp3", ".oga":
					ac := tgbotapi.NewInputMediaAudio(fileRequestData)
					ac.Caption = text
					media = append(media, &ac)
				case ".ogg":
					chatid, _ := t.chat(chatID)
					voc := tgbotapi.NewVoice(chatid, fileRequestData)
					voc.Caption = text
					voc.MessageThreadID = threadID
					//voc.ReplyToMessageID = parentID
					res, err := t.t.Send(voc)
					if err != nil {
						return "", err
					}
					return strconv.Itoa(res.MessageID), nil
				case ".webp":
					dc := tgbotapi.NewInputMediaDocument(fileRequestData)
					dc.Caption = text
					media = append(media, &dc)
				default:
					dc := tgbotapi.NewInputMediaDocument(fileRequestData)
					dc.Caption = text
					media = append(media, &dc)
				}
			}

			if len(media) == 0 {
				return "", nil
			}

			mg := tgbotapi.MediaGroupConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatConfig: tgbotapi.ChatConfig{
						ChatID: chatId,
					},
					MessageThreadID: threadID,
					ReplyParameters: tgbotapi.ReplyParameters{
						MessageID:                replyId,
						ChatID:                   nil,
						AllowSendingWithoutReply: false,
						Quote:                    "",
						QuoteParseMode:           "",
						QuoteEntities:            nil,
						QuotePosition:            0,
					},
				},
				Media: media,
			}
			messages, err := t.t.SendMediaGroup(mg)
			if err != nil {
				// If the error indicates a sticker being sent as photo, retry as documents
				if strings.Contains(err.Error(), "can't use file of type Sticker as Photo") {
					t.log.Info("Retrying failed photo send as documents due to sticker detection")
					// Convert all photos in the media group to documents
					for i, m := range media {
						if pc, ok := m.(*tgbotapi.InputMediaPhoto); ok {
							dc := tgbotapi.NewInputMediaDocument(pc.Media)
							dc.Caption = pc.Caption
							media[i] = &dc
						}
					}
					// Retry with documents
					messages, retryErr := t.t.SendMediaGroup(mg)
					if retryErr != nil {
						return "", err // Return original error
					}
					return strconv.Itoa(messages[0].MessageID), nil
				}
				return "", err
			}
			// return first message id
			return strconv.Itoa(messages[0].MessageID), nil
		}
	}
	return "", nil
}

func (t *Telegram) SendEmbed(lvlkz string, chatid string, text string) (int, error) {
	chatId, threadID := t.chat(chatid)
	var keyboardQueue = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(lvlkz+"+", lvlkz+"+"),
			tgbotapi.NewInlineKeyboardButtonData(lvlkz+"-", lvlkz+"-"),
			tgbotapi.NewInlineKeyboardButtonData(lvlkz+"++", lvlkz+"++"),
			tgbotapi.NewInlineKeyboardButtonData(lvlkz+"+30", lvlkz+"+++"),
		),
	)
	msg := tgbotapi.NewMessage(chatId, text)
	msg.MessageThreadID = threadID
	msg.ReplyMarkup = keyboardQueue
	message, err := t.t.Send(msg)

	return message.MessageID, err

}
func (t *Telegram) SendEmbedTime(chatid string, text string) (int, error) {
	chatId, threadID := t.chat(chatid)
	var keyboardQueue = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("+", "+"),
			tgbotapi.NewInlineKeyboardButtonData("-", "-"),
		),
	)
	msg := tgbotapi.NewMessage(chatId, text)
	msg.MessageThreadID = threadID
	msg.ReplyMarkup = keyboardQueue
	message, err := t.t.Send(msg)

	return message.MessageID, err
}

func (t *Telegram) SendHelp(chatid string, text string, midHelpTgString string, ifUser bool) (string, error) {
	midHelpTg, err := strconv.Atoi(midHelpTgString)
	if err != nil {
		if !ifUser {
			t.log.Info(fmt.Sprintf("%s %s %d\n", chatid, midHelpTgString, midHelpTg))
		}
		midHelpTg = 0
	}

	var levels []string
	chatId, ThreadID := t.chat(chatid)
	ok, config := t.checkChannelConfigTG(chatid)

	if ok {
		levels = t.Storage.Db.ReadTop5Level(config.CorpName)
	}

	if !ifUser && config.TgChannel != "" {
		last := t.Storage.Db.ReadTelegramLastMessage(config.CorpName)
		if last == 0 {
			active := t.Storage.Db.ReadTelegramLastMessageActive(config.CorpName)
			if active == 0 {
				t.log.Info(fmt.Sprintf("удалена корпорация '%s' \n", config.CorpName))
				_, _ = t.SendChannel(config.TgChannel, "Бот отключен, для активации бота напишите команду \n.добавить", "")

				t.Storage.Db.DeleteConfigRs(config)
				return "Бот отключен", nil
			}
		}
		//
		if last-5 < midHelpTg {
			return midHelpTgString, nil
		} else {
			fmt.Printf("!ifUser if last-5 < midHelpTg last: %d midHelpTg: %d\n", last-5, midHelpTg)
		}
	}
	getButton := func(level string) string {
		after, found := strings.CutPrefix(level, "rs")
		if found {
			return after + "+"
		}
		after, found = strings.CutPrefix(level, "drs")
		if found {
			return after + "+"
		}
		after, found = strings.CutPrefix(level, "solo")
		if found {
			return after + "+"
		}
		return level
	}

	if midHelpTg != 0 {
		go func() {
			_ = t.DelMessage(chatid, midHelpTg)
		}()

	}

	var btt []tgbotapi.InlineKeyboardButton
	if len(levels) > 2 {
		for _, level := range levels {
			key := getButton(level)

			bt := tgbotapi.NewInlineKeyboardButtonData(key, key)

			btt = append(btt, bt)
		}
	} else {
		for i := 7; i < 12; i++ {
			l := strconv.Itoa(i)
			bt := tgbotapi.NewInlineKeyboardButtonData(l+"+", l+"+")
			btt = append(btt, bt)
		}
	}

	msg := tgbotapi.NewMessage(chatId, models.EscapeMarkdownV2ForHelp(text))

	msg.MessageThreadID = ThreadID
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	if len(btt) > 0 {
		var keyboardQueue = tgbotapi.NewInlineKeyboardMarkup(btt)
		msg.ReplyMarkup = keyboardQueue
	}

	message, err := t.t.Send(msg)
	if err != nil {
		t.log.Info(fmt.Sprintf("ERR chatid: %s\n err:%+v", chatid, err))
		return "", err
	}
	mid := strconv.Itoa(message.MessageID)

	return mid, nil
}

func (t *Telegram) SendPoll(m models.Request) string {
	chatid := m.Data["chatid"]
	question := m.Data["question"]
	url := m.Data["url"]
	createTime := m.Data["createTime"]
	description := ""
	for i, option := range m.Options {
		description += fmt.Sprintf("\n%d. %s", i+1, option)
	}
	title := fmt.Sprintf("Опрос от %s: \n\n  %s\n", m.Data["author"], question)

	chatId, ThreadID := t.chat(chatid)
	text := fmt.Sprintf("%s\n%s\n\n[результат](%s)", title, description, url)

	msg := tgbotapi.NewMessage(chatId, models.EscapeMarkdownV2ForLink(text))

	msg.MessageThreadID = ThreadID
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	btt := t.AddButtonPoll(createTime, m.Options)
	if len(btt) > 0 {
		var keyboardQueue = tgbotapi.NewInlineKeyboardMarkup(btt)
		msg.ReplyMarkup = keyboardQueue
	}

	message, err := t.t.Send(msg)
	if err != nil {
		t.log.Info(fmt.Sprintf("ERR chatid %s\n", chatid))
		t.log.ErrorErr(err)
		return ""
	}
	pinConfig := tgbotapi.PinChatMessageConfig{
		BaseChatMessage: tgbotapi.BaseChatMessage{
			ChatConfig: tgbotapi.ChatConfig{
				ChatID: chatId,
			},
			MessageID: message.MessageID,
		},
	}
	_, _ = t.t.Send(pinConfig)

	mid := strconv.Itoa(message.MessageID)
	return mid
}
func (t *Telegram) AddButtonPoll(createTime string, option []string) []tgbotapi.InlineKeyboardButton {
	var btt []tgbotapi.InlineKeyboardButton
	if len(option) > 0 {
		if len(option) > 0 && option[0] != "" {
			bt := tgbotapi.NewInlineKeyboardButtonData(emOne, createTime+".1")
			btt = append(btt, bt)
		}
		if len(option) > 1 && option[1] != "" {
			bt := tgbotapi.NewInlineKeyboardButtonData(emTwo, createTime+".2")
			btt = append(btt, bt)
		}
		if len(option) > 2 && option[2] != "" {
			bt := tgbotapi.NewInlineKeyboardButtonData(emTree, createTime+".3")
			btt = append(btt, bt)
		}
		if len(option) > 3 && option[3] != "" {
			bt := tgbotapi.NewInlineKeyboardButtonData(emFour, createTime+".4")
			btt = append(btt, bt)
		}
		if len(option) > 4 && option[4] != "" {
			bt := tgbotapi.NewInlineKeyboardButtonData(emFive, createTime+".5")
			btt = append(btt, bt)
		}
	}
	return btt
}

const (
	emOne  = "1️⃣"
	emTwo  = "2️⃣"
	emTree = "3️⃣"
	emFour = "4️⃣"
	emFive = "5️⃣"
)

func makeRolesMap(st []models.Role) (map[string]int64, map[int64]string) {
	rolesMap := make(map[string]int64)
	rolesId := make(map[int64]string)

	for _, role := range st {
		// Добавляем оригинальное название
		rolesMap[strings.ToLower(role.Name)] = role.ID
		rolesId[role.ID] = role.Name

		// Создаем варианты с заменой
		name := strings.ToLower(role.Name)
		mention := role.ID

		if strings.Contains(name, "ткз") {
			rolesMap[strings.ReplaceAll(name, "ткз", "drs")] = mention
			rolesMap[strings.ReplaceAll(name, "ткз", "тчз")] = mention
		}
		if strings.Contains(name, "drs") {
			rolesMap[strings.ReplaceAll(name, "drs", "ткз")] = mention
			rolesMap[strings.ReplaceAll(name, "drs", "тчз")] = mention
		}
		if strings.Contains(name, "тчз") {
			rolesMap[strings.ReplaceAll(name, "тчз", "ткз")] = mention
			rolesMap[strings.ReplaceAll(name, "тчз", "drs")] = mention
		}
	}

	return rolesMap, rolesId
}
