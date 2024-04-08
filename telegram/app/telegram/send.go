package telegram

import (
	"bytes"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"telegram/models"
	"telegram/telegram/restapi"
	"time"
)

func (t *Telegram) Send(chatid string, text string) (id string, err error) {
	chatId, threadID := t.chat(chatid)
	m := tgbotapi.NewMessage(chatId, text)
	m.MessageThreadID = threadID
	mmm, err := t.t.Send(m)
	if err != nil {
		return "", err
	}
	id = strconv.Itoa(mmm.MessageID)
	return id, nil
}

func (t *Telegram) ChatTyping(chatId string) {
	chatid, threadID := t.chat(chatId)
	typingConfig := tgbotapi.NewChatAction(chatid, tgbotapi.ChatTyping)
	typingConfig.MessageThreadID = threadID
	_, _ = t.t.Send(typingConfig)
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

func (t *Telegram) SendChannelDelSecond(chatid string, text string, second int) {
	chatId, threadID := t.chat(chatid)
	m := tgbotapi.NewMessage(chatId, text)
	m.MessageThreadID = threadID
	tMessage, err1 := t.t.Send(m)
	if err1 != nil {
		t.log.Error(err1.Error())
	}
	if second <= 60 {
		go func() {
			time.Sleep(time.Duration(second) * time.Second)
			_, _ = t.t.Request(tgbotapi.NewDeleteMessage(chatId, tMessage.MessageID))
		}()
	} else {
		err := restapi.SendInsertTimer(models.Timer{
			Tgmesid:  strconv.Itoa(tMessage.MessageID),
			Tgchatid: chatid,
			Timed:    second,
		})
		if err != nil {
			t.log.ErrorErr(err)
		}
	}
}

func (t *Telegram) SendEmbedTime(chatid string, text string, bt []models.Button) (mId string, err error) {
	chatId, threadID := t.chat(chatid)

	//var keyboardQueue = tgbotapi.NewInlineKeyboardMarkup(
	//	tgbotapi.NewInlineKeyboardRow(
	//		tgbotapi.NewInlineKeyboardButtonData("+", "+"),
	//		tgbotapi.NewInlineKeyboardButtonData("-", "-"),
	//	),
	//)
	//var keyboardQueue = tgbotapi.NewInlineKeyboardMarkup(
	//	tgbotapi.NewInlineKeyboardRow(
	//		tgbotapi.NewInlineKeyboardButtonData(lvlkz+"+", lvlkz+"+"),
	//		tgbotapi.NewInlineKeyboardButtonData(lvlkz+"-", lvlkz+"-"),
	//		tgbotapi.NewInlineKeyboardButtonData(lvlkz+"++", lvlkz+"++"),
	//		tgbotapi.NewInlineKeyboardButtonData(lvlkz+"+30", lvlkz+"+++"),
	//	),
	//)

	if len(bt) > 0 {
		var rows []tgbotapi.InlineKeyboardButton
		for _, button := range bt {
			rows = append(rows, tgbotapi.NewInlineKeyboardButtonData(button.Text, button.Data))
		}
		var keyboardQueue = tgbotapi.NewInlineKeyboardMarkup(rows)
		msg := tgbotapi.NewMessage(chatId, text)
		msg.MessageThreadID = threadID
		msg.ReplyMarkup = keyboardQueue
		message, err := t.t.Send(msg)
		if err != nil {
			return "", err
		}
		mId = strconv.Itoa(message.MessageID)
		return mId, nil
	}
	return "", errors.New("EmptyKeyboardButton")
}
