package telegram

import (
	"fmt"
	"telegram/models2"
	"time"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func (t *Telegram) SaveMember(c *tgbotapi.Chat, u *tgbotapi.User) {
	chat := models2.Chat{
		ChatID:   c.ID,
		ChatName: c.Title,
	}
	if u.IsBot {
		return
	}

	user := models2.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		UserName:  u.UserName,
		IsAdmin:   false,
	}
	t.UpdateUserFromMessage(chat, user)
}

// UpdateUserFromMessage обновляет данные пользователя из сообщения
func (t *Telegram) UpdateUserFromMessage(chat models2.Chat, user models2.User) {
	// Сначала убедимся, что чат существует в бд
	_, err := t.Storage.Db.GetChat(chat.ChatID)
	if err != nil {
		// Если чат не существует, создаем его
		go func() {
			if chat.ChatName != "" {
				_ = t.Storage.Db.CreateOrUpdateChat(chat.ChatID, chat.ChatName)
			}
		}()
	}

	// Обновляем/добавляем пользователя в чат
	dbUser := models2.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UserName:  user.UserName,
		IsAdmin:   false, // По умолчанию не админ, будет обновлено отдельно
	}
	_ = t.Storage.Db.AddUserToChat(chat.ChatID, dbUser)
}

func (t *Telegram) TestFunc() {
	//t.CheckChatAvalibel()
	//t.GetChatsMember(392380978, false)
	//fmt.Println("next")
	//t.GetChatsMember(5977372281, false)
}

func (t *Telegram) CheckChatAvalibel() {
	chats, err := t.Storage.Db.GetChats()
	if err != nil {
		t.log.ErrorErr(err)
		return
	}
	deleteChatFull := func(chat models2.Chat) {
		err = t.Storage.Db.DeleteChatFull(chat.ChatID)
		if err != nil {
			t.log.ErrorErr(err)
		}
	}
	for _, chat := range chats {
		member, err := t.t.GetChatMember(tgbotapi.NewGetChatMember(chat.ChatID, t.t.Self.ID))

		if err != nil {
			if err.Error() == "Bad Request: group chat was upgraded to a supergroup chat" {
				fmt.Println("upgraded to a supergroup")
				deleteChatFull(chat)
			} else if err.Error() == "Bad Request: chat not found" {
				fmt.Println("chat not found")
				deleteChatFull(chat)
			} else {
				t.log.Error(fmt.Sprintln("Ошибка API", "chat", chat.ChatName, "err", err))
			}
			continue

		}

		if member.Status == "left" || member.Status == "kicked" {
			deleteChatFull(chat)
		}
	}
}

func (t *Telegram) GetChatsMember(m int64, full bool) []models2.ChatAccess {
	startTime := time.Now()
	defer func() {
		fmt.Printf("Команда GetChatsMember %v\n", time.Since(startTime))
	}()

	var chats []models2.Chat
	var err error

	if full {
		chats, err = t.Storage.Db.GetChats() // gel all chat in bot
	} else {
		chats, err = t.Storage.Db.GetUserChats(m)
	}
	if err != nil {
		t.log.ErrorErr(err)
	}

	var chatsArray []models2.ChatAccess
	for _, chat := range chats {
		member, err := t.t.GetChatMember(tgbotapi.NewGetChatMember(chat.ChatID, m))

		if err != nil {
			t.log.Error(fmt.Sprintln("Ошибка API", "chat", chat.ChatName, "err", err))
			continue
		}

		if member.Status == "creator" || member.Status == "administrator" || member.Status == "member" {
			c := models2.ChatAccess{
				ChatID:   chat.ChatID,
				ChatName: chat.ChatName,
				UserID:   m,
				Status:   member.Status,
			}

			chatsArray = append(chatsArray, c)
		}
	}
	return chatsArray
}
