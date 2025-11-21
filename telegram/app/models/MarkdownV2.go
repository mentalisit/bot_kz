package models

import (
	"fmt"
	"regexp"
	"strings"
)

// EscapeMarkdownV2 экранирует специальные символы для MarkdownV2
func EscapeMarkdownV2(text string) string {
	// Список специальных символов в MarkdownV2
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}

	escaped := text
	for _, char := range specialChars {
		escaped = strings.ReplaceAll(escaped, char, "\\"+char)
	}
	return escaped
}

// FindTelegramMentions находит упоминания в формате [%s](tg://user?id=%d)
func FindTelegramMentions(text string) (found bool, newText string) {
	// Регулярное выражение для поиска шаблона [текст](tg://user?id=число)
	re := regexp.MustCompile(`\[([^\]]+)\]\(tg://user\?id=(\d+)\)`)
	matches := re.FindAllStringSubmatchIndex(text, -1)
	if len(matches) == 0 {
		return false, text
	}

	var result strings.Builder
	lastPos := 0

	for _, match := range matches {
		if len(match) >= 6 {
			// Экранируем текст до упоминания
			if match[0] > lastPos {
				plainText := text[lastPos:match[0]]
				result.WriteString(EscapeMarkdownV2(plainText))
			}

			username := text[match[2]:match[3]]
			userIDStr := text[match[4]:match[5]] //на потом если захочу что то

			// Экранируем имя пользователя и создаем упоминание
			escapedUsername := EscapeMarkdownV2(username)
			mention := fmt.Sprintf("[%s](tg://user?id=%s)", escapedUsername, userIDStr)
			result.WriteString(mention)

			lastPos = match[1]
		}
	}

	// Экранируем оставшийся текст после последнего упоминания
	if lastPos < len(text) {
		plainText := text[lastPos:]
		result.WriteString(EscapeMarkdownV2(plainText))
	}

	return true, result.String()
}

func EscapeMarkdownV2ForLink(text string) string {
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

func EscapeMarkdownV2ForHelp(text string) string {
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
