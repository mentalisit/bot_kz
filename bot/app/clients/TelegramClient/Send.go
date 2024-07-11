package TelegramClient

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"kz_bot/models"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (t *Telegram) SendEmded(lvlkz string, chatid string, text string) int {
	a := strings.SplitN(chatid, "/", 2)
	chatId, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		t.log.ErrorErr(err)
	}
	ThreadID := 0
	if len(a) > 1 {
		ThreadID, err = strconv.Atoi(a[1])
		if err != nil {
			t.log.ErrorErr(err)
		}
	}
	var keyboardQueue = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(lvlkz+"+", lvlkz+"+"),
			tgbotapi.NewInlineKeyboardButtonData(lvlkz+"-", lvlkz+"-"),
			tgbotapi.NewInlineKeyboardButtonData(lvlkz+"++", lvlkz+"++"),
			tgbotapi.NewInlineKeyboardButtonData(lvlkz+"+30", lvlkz+"+++"),
		),
	)
	msg := tgbotapi.NewMessage(chatId, text)
	msg.MessageThreadID = ThreadID
	msg.ReplyMarkup = keyboardQueue
	message, _ := t.t.Send(msg)

	return message.MessageID

}
func (t *Telegram) SendEmbedTime(chatid string, text string) int {
	a := strings.SplitN(chatid, "/", 2)
	chatId, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		t.log.ErrorErr(err)
	}
	ThreadID := 0
	if len(a) > 1 {
		ThreadID, err = strconv.Atoi(a[1])
		if err != nil {
			t.log.ErrorErr(err)
		}
	}

	var keyboardQueue = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("+", "+"),
			tgbotapi.NewInlineKeyboardButtonData("-", "-"),
		),
	)
	msg := tgbotapi.NewMessage(chatId, text)
	msg.MessageThreadID = ThreadID
	msg.ReplyMarkup = keyboardQueue
	message, _ := t.t.Send(msg)

	return message.MessageID
}

// отправка сообщения в телегу
func (t *Telegram) SendChannel(chatid string, text string) int {
	a := strings.SplitN(chatid, "/", 2)
	chatId, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		t.log.ErrorErr(err)
	}
	ThreadID := 0
	if len(a) > 1 {
		ThreadID, err = strconv.Atoi(a[1])
		if err != nil {
			t.log.ErrorErr(err)
		}
	}
	m := tgbotapi.NewMessage(chatId, text)
	m.MessageThreadID = ThreadID
	tMessage, _ := t.t.Send(m)
	return tMessage.MessageID
}

//func (t *Telegram) SendText(chatid string, text string) error {
//	a := strings.SplitN(chatid, "/", 2)
//	chatId, err := strconv.ParseInt(a[0], 10, 64)
//	if err != nil {
//		t.log.ErrorErr(err)
//	}
//	ThreadID := 0
//	if len(a) > 1 {
//		ThreadID, err = strconv.Atoi(a[1])
//		if err != nil {
//			t.log.ErrorErr(err)
//		}
//	}
//	m := tgbotapi.NewMessage(chatId, text)
//	m.MessageThreadID = ThreadID
//	_, err = t.t.Send(m)
//	if err != nil {
//		return err
//	}
//	return nil
//}

func (t *Telegram) SendChannelDelSecond(chatid string, text string, second int) {
	a := strings.SplitN(chatid, "/", 2)
	chatId, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		t.log.ErrorErr(err)
	}
	ThreadID := 0
	if len(a) > 1 {
		ThreadID, err = strconv.Atoi(a[1])
		if err != nil {
			t.log.ErrorErr(err)
		}
	}
	m := tgbotapi.NewMessage(chatId, text)
	m.MessageThreadID = ThreadID
	tMessage, err1 := t.t.Send(m)
	if err1 != nil {
		t.log.Error(err1.Error())
	}
	if second <= 60 {
		go func() {
			time.Sleep(time.Duration(second) * time.Second)
			t.DelMessage(chatid, tMessage.MessageID)
		}()
	} else {
		t.storage.TimeDeleteMessage.TimerInsert(models.Timer{
			Tgmesid:  strconv.Itoa(tMessage.MessageID),
			Tgchatid: chatid,
			Timed:    second,
		})
	}
}

//func (t *Telegram) SendChannelAsync(chatid string, text string, resultChannel chan<- models.MessageIds, wg *sync.WaitGroup) {
//	defer wg.Done()
//	a := strings.SplitN(chatid, "/", 2)
//	chatId, err := strconv.ParseInt(a[0], 10, 64)
//	if err != nil {
//		t.log.ErrorErr(err)
//	}
//	ThreadID := 0
//	if len(a) > 1 {
//		ThreadID, err = strconv.Atoi(a[1])
//		if err != nil {
//			t.log.ErrorErr(err)
//		}
//	}
//	m := tgbotapi.NewMessage(chatId, text)
//	m.MessageThreadID = ThreadID
//	tMessage, _ := t.t.Send(m)
//	messageData := models.MessageIds{
//		MessageId: strconv.Itoa(tMessage.MessageID),
//		ChatId:    chatid,
//	}
//	resultChannel <- messageData
//}

//func (t *Telegram) SendFileFromURLAsync(chatid, text string, fileURL string, resultChannel chan<- models.MessageIds, wg *sync.WaitGroup) {
//	defer wg.Done()
//	fileURL = strings.TrimSpace(fileURL)
//	a := strings.SplitN(chatid, "/", 2)
//	chatId, err := strconv.ParseInt(a[0], 10, 64)
//	if err != nil {
//		t.log.ErrorErr(err)
//	}
//	ThreadID := 0
//	if len(a) > 1 {
//		ThreadID, err = strconv.Atoi(a[1])
//		if err != nil {
//			t.log.ErrorErr(err)
//		}
//	}
//
//	parsedURL, err := url.Parse(fileURL)
//	if err != nil {
//		t.log.ErrorErr(err)
//		return
//	}
//
//	// Используем path.Base для получения последней части URL, которая представляет собой имя файла
//	fileName := path.Base(parsedURL.Path)
//	parsedURL.RawQuery = ""
//	fileURL = parsedURL.String()
//
//	// Скачиваем файл по URL
//	resp, err := http.Get(fileURL)
//	if err != nil {
//		t.log.ErrorErr(err)
//		return
//	}
//	defer resp.Body.Close()
//
//	// Читаем содержимое файла
//	buffer := new(bytes.Buffer)
//	_, err = io.Copy(buffer, resp.Body)
//	if err != nil {
//		t.log.ErrorErr(err)
//		return
//	}
//	var media []interface{}
//
//	file := tgbotapi.FileBytes{
//		Name:  fileName,
//		Bytes: buffer.Bytes(),
//	}
//
//	switch filepath.Ext(fileName) {
//
//	case ".jpg", ".jpe", ".png":
//		pc := tgbotapi.NewInputMediaPhoto(file)
//		if text != "" {
//			pc.Caption = text
//		}
//		media = append(media, pc)
//	case ".mp4", ".m4v":
//		vc := tgbotapi.NewInputMediaVideo(file)
//		if text != "" {
//			vc.Caption = text
//		}
//		media = append(media, vc)
//	case ".mp3", ".oga":
//		ac := tgbotapi.NewInputMediaAudio(file)
//		if text != "" {
//			ac.Caption = text
//		}
//		media = append(media, ac)
//	default:
//		dc := tgbotapi.NewInputMediaDocument(file)
//		if text != "" {
//			dc.Caption = text
//		}
//		media = append(media, dc)
//	}
//
//	if len(media) == 0 {
//		return
//	}
//	mg := tgbotapi.MediaGroupConfig{
//		BaseChat: tgbotapi.BaseChat{
//			ChatID:          chatId,
//			MessageThreadID: ThreadID,
//			//ChannelUsername:  msg.Username,
//			//ReplyToMessageID: parentID,
//		},
//		Media: media,
//	}
//	m, err := t.t.SendMediaGroup(mg)
//	if err != nil {
//		t.log.ErrorErr(err)
//		return
//	}
//
//	messageData := models.MessageIds{
//		MessageId: strconv.Itoa(m[0].MessageID),
//		ChatId:    chatid,
//	}
//	resultChannel <- messageData
//}

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
func (t *Telegram) ChatTyping(chatId string) {
	chatid, threadID := t.chat(chatId)
	typingConfig := tgbotapi.NewChatAction(chatid, tgbotapi.ChatTyping)
	typingConfig.MessageThreadID = threadID
	_, _ = t.t.Send(typingConfig)
}

//	func (t *Telegram) SendBridgeAsync(chatid []string, text string, fileURL []string, resultChannel chan<- models.MessageIds, wg *sync.WaitGroup) {
//		for _, chat := range chatid {
//			chatId, threadID := t.chat(chat)
//
//			if len(fileURL) > 0 && fileURL[0] != "" {
//				go func() {
//					m, err := t.sendFileToTelegram(fileURL[0], text, chat)
//					if err != nil {
//						t.log.ErrorErr(err)
//					} else {
//						messageData := models.MessageIds{
//							MessageId: strconv.Itoa(m.MessageID),
//							ChatId:    chat,
//						}
//						resultChannel <- messageData
//					}
//					wg.Done()
//				}()
//
//			} else {
//				go func() {
//					defer wg.Done()
//					m := tgbotapi.NewMessage(chatId, text)
//					m.MessageThreadID = threadID
//					tMessage, _ := t.t.Send(m)
//					messageData := models.MessageIds{
//						MessageId: strconv.Itoa(tMessage.MessageID),
//						ChatId:    chat,
//					}
//					resultChannel <- messageData
//				}()
//			}
//		}
//	}
//

// Send func for compendium sendDM
func (t *Telegram) Send(chatid string, text string, parseMode string) (string, error) {
	a := strings.SplitN(chatid, "/", 2)
	chatId, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		t.log.ErrorErr(err)
	}
	ThreadID := 0
	if len(a) > 1 {
		ThreadID, err = strconv.Atoi(a[1])
		if err != nil {
			t.log.ErrorErr(err)
		}
	}
	m := tgbotapi.NewMessage(chatId, text)
	if parseMode != "" {
		m.ParseMode = parseMode
	}
	m.MessageThreadID = ThreadID
	message, errs := t.t.Send(m)
	if errs != nil {
		return "", errs
	}
	mid := strconv.Itoa(message.MessageID)
	return mid, nil
}

//func (t *Telegram) sendFileToTelegram(fileURL, text, chatID string) (*tgbotapi.Message, error) {
//	fileURL = strings.TrimSpace(fileURL)
//	parsedURL, err := url.Parse(fileURL)
//	if err != nil {
//		t.log.ErrorErr(err)
//	}
//	parsedURL.RawQuery = ""
//	// Определяем расширение файла
//	fileExtension := filepath.Ext(parsedURL.String())
//
//	file := tgbotapi.FileURL(fileURL)
//	chatId, threadID := t.chat(chatID)
//	var msg tgbotapi.Chattable
//	switch fileExtension {
//	case ".jpg", ".jpeg", ".png", ".gif":
//		{
//			ms := tgbotapi.NewPhoto(chatId, file)
//			ms.MessageThreadID = threadID
//			ms.Caption = text
//			msg = ms
//		}
//	case ".mp4", ".m4v":
//		{
//			ms := tgbotapi.NewVideo(chatId, file)
//			ms.MessageThreadID = threadID
//			ms.Caption = text
//			msg = ms
//		}
//	default:
//		{
//			ms := tgbotapi.NewDocument(chatId, file)
//			ms.MessageThreadID = threadID
//			ms.Caption = text
//			msg = ms
//		}
//	}
//
//	// Отправляем сообщение
//	m, err := t.t.Send(msg)
//	if err != nil {
//		t.log.Info(fileExtension)
//		t.log.Info(parsedURL.String())
//		t.log.Info(fileURL)
//		return nil, err
//	}
//
//	return &m, nil
//}

//func (t *Telegram) SendBridgeFunc(chatid []string, text string, Extra []models.FileInfo, resultChannel chan<- models.MessageIds, wg *sync.WaitGroup) {
//	for _, chat := range chatid {
//		chatId, threadID := t.chat(chat)
//
//		if len(Extra) > 0 {
//			go func() {
//				mid, err := t.sendFileExtra(Extra, text, chat)
//				if err != nil {
//					t.log.ErrorErr(err)
//				} else {
//					messageData := models.MessageIds{
//						MessageId: mid,
//						ChatId:    chat,
//					}
//					resultChannel <- messageData
//				}
//				wg.Done()
//			}()
//
//		} else {
//			go func() {
//				defer wg.Done()
//				m := tgbotapi.NewMessage(chatId, text)
//				m.MessageThreadID = threadID
//				tMessage, err := t.t.Send(m)
//				if err != nil {
//					t.log.ErrorErr(err)
//				}
//				messageData := models.MessageIds{
//					MessageId: strconv.Itoa(tMessage.MessageID),
//					ChatId:    chat,
//				}
//				resultChannel <- messageData
//			}()
//		}
//	}
//}

func (t *Telegram) sendFileExtra(extra []models.FileInfo, text, chatID string) (string, error) {
	if extra != nil {
		if len(extra) > 0 {
			chatId, threadID := t.chat(chatID)
			var media []interface{}
			for _, f := range extra {
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
					media = append(media, pc)
				case ".mp4", ".m4v":
					vc := tgbotapi.NewInputMediaVideo(fileRequestData)
					vc.Caption = text
					media = append(media, vc)
				case ".mp3", ".oga":
					ac := tgbotapi.NewInputMediaAudio(fileRequestData)
					ac.Caption = text
					media = append(media, ac)
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
					media = append(media, dc)
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

func (t *Telegram) SendBridgeFuncRest(in models.BridgeSendToMessenger) []models.MessageIds {
	var messageIds []models.MessageIds
	for _, chat := range in.ChannelId {
		chatId, threadID := t.chat(chat)

		if len(in.Extra) > 0 {
			mid, err := t.sendFileExtra(in.Extra, in.Text, chat)
			if err != nil {
				t.log.ErrorErr(err)
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
				t.log.ErrorErr(err)
			} else {
				messageData := models.MessageIds{
					MessageId: strconv.Itoa(tMessage.MessageID),
					ChatId:    chat,
				}
				messageIds = append(messageIds, messageData)
			}
		}
	}
	return messageIds
}

func (t *Telegram) SendHelp(chatid string, text string) int {
	a := strings.SplitN(chatid, "/", 2)
	chatId, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		t.log.ErrorErr(err)
	}
	ThreadID := 0
	if len(a) > 1 {
		ThreadID, err = strconv.Atoi(a[1])
		if err != nil {
			t.log.ErrorErr(err)
		}
	}

	_, config := t.checkChannelConfigTG(chatid)

	levels := t.storage.Count.ReadTop5Level(config.CorpName)
	var btt []tgbotapi.InlineKeyboardButton
	if len(levels) > 0 {
		for _, level := range levels {
			var bt tgbotapi.InlineKeyboardButton
			if level[:1] == "d" {
				bt = tgbotapi.NewInlineKeyboardButtonData(level[1:]+"*", level[1:]+"*")
			} else {
				bt = tgbotapi.NewInlineKeyboardButtonData(level+"+", level+"+")
			}
			btt = append(btt, bt)
		}
	}

	msg := tgbotapi.NewMessage(chatId, escapeMarkdownV2(text))
	msg.MessageThreadID = ThreadID
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	if len(btt) > 0 {
		var keyboardQueue = tgbotapi.NewInlineKeyboardMarkup(btt)
		msg.ReplyMarkup = keyboardQueue
	}

	message, err := t.t.Send(msg)
	if err != nil {
		t.log.Info(fmt.Sprintf("ERR chatid %s\n", chatid))
		t.log.ErrorErr(err)
		return 0
	}

	return message.MessageID

}

func escapeMarkdownV2(text string) string {
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
