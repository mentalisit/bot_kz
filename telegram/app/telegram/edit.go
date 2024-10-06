package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func (t *Telegram) EditMessageTextKey(chatid string, editMesId int, textEdit string, lvlkz string) {
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
	if err != nil {
		fmt.Println("EditMessageTextKey", mes, err)
	}
}
func (t *Telegram) EditTextParseMode(chatid string, editMesId int, textEdit, ParseMode string) error {
	chatId, _ := t.chat(chatid)
	msg := tgbotapi.NewEditMessageText(chatId, editMesId, textEdit)
	if ParseMode != "" {
		msg = tgbotapi.NewEditMessageText(chatId, editMesId, escapeMarkdownV2(textEdit))
		msg.ParseMode = ParseMode
	}
	_, err := t.t.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func escapeMarkdownV2(text string) string {
	// Специальные символы, которые нужно экранировать в MarkdownV2
	specialChars := "_*[]()~`>#+-=|{}.!"

	// Буфер для результата
	var builder strings.Builder

	// Переменные для отслеживания состояния
	var inLinkText bool
	var inLinkURL bool
	var linkTextBuffer strings.Builder
	var linkURLBuffer strings.Builder

	for i := 0; i < len(text); i++ {
		char := text[i]

		if char == '[' && !inLinkText && !inLinkURL {
			inLinkText = true
			builder.WriteByte(char)
			continue
		}

		if char == ']' && inLinkText && !inLinkURL {
			inLinkText = false
			builder.WriteString(linkTextBuffer.String())
			linkTextBuffer.Reset()
			builder.WriteByte(char)
			continue
		}

		if char == '(' && !inLinkText && !inLinkURL && i > 0 && text[i-1] == ']' {
			inLinkURL = true
			builder.WriteByte(char)
			continue
		}

		if char == ')' && !inLinkText && inLinkURL {
			inLinkURL = false
			builder.WriteString(linkURLBuffer.String())
			linkURLBuffer.Reset()
			builder.WriteByte(char)
			continue
		}

		if inLinkText {
			linkTextBuffer.WriteByte(char)
			continue
		}

		if inLinkURL {
			linkURLBuffer.WriteByte(char)
			continue
		}

		if strings.ContainsRune(specialChars, rune(char)) {
			builder.WriteByte('\\')
		}
		builder.WriteByte(char)
	}

	return builder.String()
}
