package telegram

import (
	"fmt"
	"strings"
	"telegram/models"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func (t *Telegram) SaveMember(c *tgbotapi.Chat, user *tgbotapi.User) {
	t.mu.Lock()
	defer t.mu.Unlock()

	chat := models.Chat{
		ChatID:   c.ID,
		ChatName: c.Title,
	}

	for ch, _ := range t.ChatMembers {
		if ch.ChatID == chat.ChatID && ch.ChatName != chat.ChatName {
			ch.ChatName = chat.ChatName
		}
	}
	if t.ChatMembers[&chat] == nil {
		t.ChatMembers[&chat] = make(map[int64]tgbotapi.User)
	}

	// Обновляем информацию об участнике
	t.ChatMembers[&chat][user.ID] = *user
}

// Функция для получения всех отслеженных участников чата
func (t *Telegram) GetChatMembers(c *tgbotapi.Chat) []tgbotapi.User {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var members []tgbotapi.User
	for chat, m := range t.ChatMembers {
		if chat.ChatID == c.ID {
			for _, user := range m {
				members = append(members, user)
			}
		}
	}

	return members
}

// Функция для упоминания всех участников
func (t *Telegram) MentionAllMembers(c *tgbotapi.Chat, originalMessage *tgbotapi.Message) {
	if originalMessage.From.UserName != "mentalisit" {
		return
	}
	// Получаем отслеженных участников
	trackedMembers := t.GetChatMembers(c)

	var mentions []string
	mentionedUsers := make(map[int64]bool)

	// Затем остальные отслеженные участники
	for _, member := range trackedMembers {
		if !member.IsBot && !mentionedUsers[member.ID] {
			mentions = append(mentions, t.formatMention(member))
		}
	}

	// Формируем сообщение
	mentionText := "🔔 Упоминание всех участников:\n" + strings.Join(mentions, " ")
	fullMessage := fmt.Sprintf("%s\n\n%s %s", mentionText, originalMessage.From.String(), originalMessage.Text)

	msg := tgbotapi.NewMessage(c.ID, fullMessage)
	msg.ParseMode = "MarkdownV2"

	// Отправляем сообщение
	if _, err := t.t.Send(msg); err != nil {
		t.log.ErrorErr(err)
		return
	}

	// Удаляем оригинальное сообщение
	if originalMessage.MessageID != 0 {
		deleteMsg := tgbotapi.NewDeleteMessage(c.ID, originalMessage.MessageID)
		t.t.Send(deleteMsg)
	}
}

// Вспомогательная функция для форматирования упоминания
func (t *Telegram) formatMention(user tgbotapi.User) string {
	if user.UserName != "" {
		return "@" + EscapeMarkdownV2(user.UserName)
	}
	return fmt.Sprintf("[%s](tg://user?id=%d)", EscapeMarkdownV2(t.getUserName(&user)), user.ID)
}

func (t *Telegram) getUserName(user *tgbotapi.User) string {
	if user.FirstName != "" {
		if user.LastName != "" {
			return user.FirstName + " " + user.LastName
		}
		return user.FirstName
	}
	return "User"
}
