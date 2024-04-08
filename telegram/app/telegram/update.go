package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *Telegram) update() {
	ut := tgbotapi.NewUpdate(0)
	ut.Timeout = 60
	updates := t.t.GetUpdatesChan(ut)

	for update := range updates {
		if update.CallbackQuery != nil {
			t.callback(update.CallbackQuery) //нажатия в чате
		} else if update.Message != nil {
			if update.Message.Chat.IsPrivate() { //если пишут боту в личку
				t.ifPrivatMesage(update.Message)
			} else if update.Message.IsCommand() {
				//t.updatesComand(update.Message) //если сообщение является командой
			} else { //остальные сообщения
				t.logicMix(update.Message, false)
			}
		} else if update.EditedMessage != nil {
			t.logicMix(update.EditedMessage, true)
		} else if update.MyChatMember != nil {
			t.myChatMember(update.MyChatMember)
		} else if update.ChatMember != nil {
			t.chatMember(update.ChatMember)
		} else if update.ChatJoinRequest != nil {
			t.log.InfoStruct("ChatJoinRequest", update.ChatJoinRequest)
		} else {
			go func() {
				if update.Poll != nil {
					t.log.InfoStruct("pool ", update.Poll)
				} else if update.EditedChannelPost != nil {
					//t.log.InfoStruct("EditedChannelPost", update.EditedChannelPost)
				} else if update.ChannelPost != nil {
					t.log.InfoStruct("update.ChannelPost.SenderChat", update.ChannelPost.SenderChat)
					t.log.InfoStruct("update.ChannelPost.Chat", update.ChannelPost.Chat)
				} else {
					t.log.Info(fmt.Sprintf(" else update: %+v \n", update))
				}
			}()
		}
	}
}
