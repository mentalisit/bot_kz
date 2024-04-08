package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"path/filepath"
	"strconv"
	"telegram/models"
)

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
