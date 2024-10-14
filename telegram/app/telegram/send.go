package telegram

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"path/filepath"
	"strconv"
	"strings"
	"telegram/models"
	"time"
)

func (t *Telegram) ChatTyping(chatId string) {
	chatid, threadID := t.chat(chatId)
	typingConfig := tgbotapi.NewChatAction(chatid, tgbotapi.ChatTyping)
	typingConfig.MessageThreadID = threadID
	_, _ = t.t.Send(typingConfig)
}
func (t *Telegram) SendChannelDelSecond(chatid string, text string, second int) bool {
	chatId, threadID := t.chat(chatid)
	m := tgbotapi.NewMessage(chatId, text)
	m.MessageThreadID = threadID
	tMessage, err1 := t.t.Send(m)
	if err1 != nil {
		t.log.ErrorErr(err1)
		t.log.Info(fmt.Sprintf("chatid '%s', text %s, second %d", chatid, text, second))
		return false
	}
	if second <= 60 {
		go func() {
			time.Sleep(time.Duration(second) * time.Second)
			_, _ = t.t.Request(tgbotapi.NewDeleteMessage(chatId, tMessage.MessageID))
		}()
	} else {
		t.Storage.Db.TimerInsert(models.Timer{
			Tgmesid:  strconv.Itoa(tMessage.MessageID),
			Tgchatid: chatid,
			Timed:    second,
		})
	}

	if tMessage.MessageID != 0 {
		return true
	}
	return false
}
func (t *Telegram) SendChannel(chatid string, text string) (string, error) {
	chatId, threadID := t.chat(chatid)
	m := tgbotapi.NewMessage(chatId, text)
	m.MessageThreadID = threadID
	tMessage, err := t.t.Send(m)
	if err != nil {
		t.log.ErrorErr(err)
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
func (t *Telegram) SendEmbed(lvlkz string, chatid string, text string) int {
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
	message, _ := t.t.Send(msg)

	return message.MessageID

}
func (t *Telegram) SendEmbedTime(chatid string, text string) int {
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
	message, _ := t.t.Send(msg)

	return message.MessageID
}

func (t *Telegram) SendHelp(chatid string, text string, midHelpTgString string) string {
	midHelpTg, err := strconv.Atoi(midHelpTgString)
	if err != nil {
		midHelpTg = 0
	}

	var levels []string
	chatId, ThreadID := t.chat(chatid)
	_, config := t.checkChannelConfigTG(chatid)

	levels = t.Storage.Db.ReadTop5Level(config.CorpName)
	last := t.Storage.Db.ReadTelegramLastMessage(config.CorpName)

	if last < midHelpTg {
		return midHelpTgString
	}

	t.DelMessage(chatid, midHelpTg)

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
	} else {
		for i := 7; i < 12; i++ {
			var bt tgbotapi.InlineKeyboardButton
			l := strconv.Itoa(i)
			bt = tgbotapi.NewInlineKeyboardButtonData(l+"*", l+"*")
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
		t.log.Info(fmt.Sprintf("ERR chatid %s\n", chatid))
		t.log.ErrorErr(err)
		return ""
	}
	mid := strconv.Itoa(message.MessageID)

	return mid
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
