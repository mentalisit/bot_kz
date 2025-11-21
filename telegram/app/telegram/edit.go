package telegram

import (
	"telegram/models"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func (t *Telegram) EditMessageTextKey(chatid string, editMesId int, textEdit string, lvlkz string) error {
	chatId, _ := t.chat(chatid)
	var keyboardQueue = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(lvlkz+"+", lvlkz+"+"),
			tgbotapi.NewInlineKeyboardButtonData(lvlkz+"-", lvlkz+"-"),
			tgbotapi.NewInlineKeyboardButtonData(lvlkz+"++", lvlkz+"++"),
			tgbotapi.NewInlineKeyboardButtonData(lvlkz+"+30", lvlkz+"+++"),
		),
	)
	mes := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			BaseChatMessage: tgbotapi.BaseChatMessage{
				ChatConfig: tgbotapi.ChatConfig{
					ChatID: chatId,
				},
				MessageID: editMesId,
			},
			ReplyMarkup: &keyboardQueue,
		},
		Text: textEdit,
	}

	_, err := t.t.Send(mes)
	return err
}
func (t *Telegram) EditText(chatid string, editMesId int, textEdit, ParseMode string) error {
	chatId, _ := t.chat(chatid)
	msg := tgbotapi.NewEditMessageText(chatId, editMesId, textEdit)
	if ParseMode != "" {
		msg = tgbotapi.NewEditMessageText(chatId, editMesId, models.EscapeMarkdownV2ForLink(textEdit))
		msg.ParseMode = ParseMode
	}
	_, err := t.t.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
