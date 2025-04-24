package telegram

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"os"

	//tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"path/filepath"
	"strconv"
	"strings"
	"telegram/models"
	"time"
)

func (t *Telegram) ChatTyping(chatId string) error {
	chatid, threadID := t.chat(chatId)
	typingConfig := tgbotapi.NewChatAction(chatid, tgbotapi.ChatTyping)
	typingConfig.MessageThreadID = threadID
	_, err := t.t.Send(typingConfig)
	return err
}
func (t *Telegram) SendChannelDelSecond(chatid string, text string, second int) (bool, error) {
	chatId, threadID := t.chat(chatid)
	m := tgbotapi.NewMessage(chatId, text)
	m.MessageThreadID = threadID
	tMessage, err1 := t.t.Send(m)
	if err1 != nil {
		t.log.ErrorErr(err1)
		t.log.Info(fmt.Sprintf("chatid '%s', text %s, second %d", chatid, text, second))
		return false, err1
	}
	tu := int(time.Now().UTC().Unix())
	t.Storage.Db.TimerInsert(models.Timer{
		Tgmesid:  strconv.Itoa(tMessage.MessageID),
		Tgchatid: chatid,
		Timed:    tu + second,
	})

	if tMessage.MessageID != 0 {
		return true, nil
	}
	return false, nil
}
func (t *Telegram) SendChannel(chatid, text, parseMode string) (string, error) {
	chatId, threadID := t.chat(chatid)
	m := tgbotapi.NewMessage(chatId, text)
	m.MessageThreadID = threadID
	m.ParseMode = parseMode
	tMessage, err := t.t.Send(m)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return strconv.Itoa(tMessage.MessageID), nil
}
func (t *Telegram) SendPic(chatID, text string, imageBytes []byte) error {
	chatid, threadID := t.chat(chatID)
	msg := tgbotapi.NewPhoto(chatid, tgbotapi.FileReader{
		Name:   "image.jpg",
		Reader: bytes.NewReader(imageBytes),
	})
	msg.MessageThreadID = threadID
	msg.Caption = text

	// Отправляем изображение
	_, err := t.t.Send(msg)
	if err != nil {
		return err
	}

	return nil
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

		////бля текст для каждого канала нужен свой
		//func (b *Bridge) replaceText(text string,conf models.BridgeConfig) string {
		//	re := regexp.MustCompile("@&(\\w+)")
		//	mentionRole := re.FindAllStringSubmatch(text, -1)
		//	if len(mentionRole) > 0 {
		//	textReplace := mentionRole[0][0]
		//	roleName := mentionRole[0][1]
		//	mention := (roleName)
		//	text = strings.Replace(text, textReplace, mention, 1)
		//}
		//	return text
		//}

		chatId, threadID := t.chat(chat)

		if len(in.Extra) > 0 {
			mid, err := t.sendFileExtra(in.Extra, in.Text, chat)
			if err != nil {
				t.log.InfoStruct(fmt.Sprintf("err %+v\n", err), in)
			} else {
				messageData := models.MessageIds{
					MessageId: mid,
					ChatId:    chat,
				}
				messageIds = append(messageIds, messageData)
			}
		} else {

			m := tgbotapi.NewMessage(chatId, in.Text)
			m.MessageThreadID = threadID
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
			if strings.Contains(in.Text, "@&Русь") || strings.Contains(in.Text, "@everyone") {
				pinConfig := tgbotapi.NewPinChatMessage(chatId, tMessage.MessageID, false)
				_, _ = t.t.Send(pinConfig)
			}

		}
	}
	return messageIds
}
func (t *Telegram) sendFileExtra(extra []models.FileInfo, text, chatID string) (string, error) {
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
				case ".jpg", ".jpe", ".png":
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
				},
				Media: media,
			}
			messages, err := t.t.SendMediaGroup(mg)
			if err != nil {
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

	if !ifUser {
		last := t.Storage.Db.ReadTelegramLastMessage(config.CorpName)
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
		go t.DelMessage(chatid, midHelpTg)

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

	msg := tgbotapi.NewMessage(chatId, escapeMarkdownV2ForHelp(text))

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
func escapeMarkdownV2ForHelp(text string) string {
	var builder strings.Builder
	specialChars := "\\_[]()~>#+-=|{}.!"
	i := 0

	for i < len(text) {
		if i+1 < len(text) && text[i] == '*' && text[i+1] == '*' {
			i += 2
			continue
		}
		if i+1 < len(text) && text[i] == '7' && text[i+1] == '*' && text[i-1] == '*' {
			builder.WriteString("7\\*")
			i += 2
			continue
		}
		if text[i] == '*' {
			i++
			continue
		}

		// Проверяем, является ли текущий символ специальным
		if strings.ContainsRune(specialChars, rune(text[i])) {
			builder.WriteByte('\\') // Добавляем экранирующий символ
		}

		// Добавляем текущий символ в строку результата
		builder.WriteByte(text[i])
		i++
	}

	return builder.String()
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

	msg := tgbotapi.NewMessage(chatId, escapeMarkdownV2(text))

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
