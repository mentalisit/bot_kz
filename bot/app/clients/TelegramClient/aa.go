package TelegramClient

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"kz_bot/models"
	"strconv"
	"strings"
	"time"
)

func (t *Telegram) DelMessage(chatid string, idSendMessage int) {
	a := strings.SplitN(chatid, "/", 2)
	chatId, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		t.log.ErrorErr(err)
	}
	_, _ = t.t.Request(tgbotapi.DeleteMessageConfig(tgbotapi.NewDeleteMessage(chatId, idSendMessage)))
}

func (t *Telegram) DelMessageSecond(chatid string, idSendMessage string, second int) {
	id, err := strconv.Atoi(idSendMessage)
	if err != nil {
		return
	}
	if second <= 60 {
		go func() {
			time.Sleep(time.Duration(second) * time.Second)
			t.DelMessage(chatid, id)
		}()
	} else {
		t.storage.TimeDeleteMessage.TimerInsert(models.Timer{
			Tgmesid:  strconv.Itoa(id),
			Tgchatid: chatid,
			Timed:    second,
		})
	}
}
func (t *Telegram) EditMessageTextKey(chatid string, editMesId int, textEdit string, lvlkz string) {
	a := strings.SplitN(chatid, "/", 2)
	chatId, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		t.log.ErrorErr(err)
	}

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

	_, err = t.t.Send(mes)
	if err != nil {
		//t.log.InfoStruct("error EditMessageTextKey ", err)
	}
}
func (t *Telegram) EditText(chatid string, editMesId int, textEdit string) {
	a := strings.SplitN(chatid, "/", 2)
	chatId, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		t.log.ErrorErr(err)
	}
	_, err = t.t.Send(tgbotapi.NewEditMessageText(chatId, editMesId, textEdit))
	if err != nil {
		//t.log.Println("Ошибка редактирования EditText ", err)
	}
}
func (t *Telegram) EditTextParseMode(chatid string, editMesId int, textEdit, ParseMode string) error {
	a := strings.SplitN(chatid, "/", 2)
	chatId, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		t.log.ErrorErr(err)
	}
	msg := tgbotapi.NewEditMessageText(chatId, editMesId, textEdit)
	if ParseMode != "" {
		msg.ParseMode = ParseMode
	}
	_, err = t.t.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
func (t *Telegram) CheckAdminTg(chatid string, name string) bool {
	admin := false
	a := strings.SplitN(chatid, "/", 2)
	chatId, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		t.log.ErrorErr(err)
	}
	admins, err := t.t.GetChatAdministrators(
		tgbotapi.ChatAdministratorsConfig{
			ChatConfig: tgbotapi.ChatConfig{
				ChatID: chatId,
			},
		},
	)
	if err != nil {
		t.log.ErrorErr(err)
	}
	for _, ad := range admins {
		if name == ad.User.UserName && (ad.IsAdministrator() || ad.IsCreator()) {
			admin = true
			break
		}
	}
	return admin
}

func (t *Telegram) updatesComand(c *tgbotapi.Message) {
	ChatId := strconv.FormatInt(c.Chat.ID, 10) + fmt.Sprintf("/%d", c.MessageThreadID)
	if c.Command() == "chatid" {
		t.SendChannelDelSecond(ChatId, ChatId, 20)
	}
	ok, config := t.checkChannelConfigTG(ChatId)
	if ok {
		MessageID := strconv.Itoa(c.MessageID)
		switch c.Command() {
		case "help":
			t.help(config.TgChannel, MessageID)
		case "helpqueue":
			t.helpQueue(config.TgChannel, MessageID)
		case "helpnotification":
			t.helpNotification(config.TgChannel, MessageID)
		case "helpevent":
			t.helpEvent(config.TgChannel, MessageID)
		case "helptop":
			t.helpTop(config.TgChannel, MessageID)
		case "helpicon":
			t.helpIcon(config.TgChannel, MessageID)
		}
	} else {
		switch c.Command() {
		case "help":
			t.SendChannelDelSecond(ChatId, "Активируйте бота командой \n.добавить", 60)
		default:
			t.SendChannelDelSecond(ChatId, "Вам не доступна данная команда \n /help", 60)
		}
	}
}
func (t *Telegram) chatName(chatid string) (chatName string) {
	id, _ := t.chat(chatid)
	r, err := t.t.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: id}})
	if err != nil {
		t.log.ErrorErr(err)
	}
	return r.Title
}
func (t *Telegram) chatInviteLink(chatid int64) string {
	r, err := t.t.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: chatid}})
	if err != nil {
		t.log.ErrorErr(err)
	}
	return r.InviteLink
}

// Кеш для чатов проверки моего присутствия
var imHereChat []int64

func (t *Telegram) imHere(chatID int64, chat tgbotapi.Chat) {
	if chat.Type == "group" || chat.Type == "supergroup" {
		userID := int64(392380978)

		here := false
		for _, i := range imHereChat {
			if i == chatID {
				here = true
				break
			}
		}
		if !here {
			// Получаем информацию о членстве пользователя в группе
			m, err := t.t.GetChatMember(tgbotapi.GetChatMemberConfig{ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
				ChatConfig: tgbotapi.ChatConfig{ChatID: chatID},
				UserID:     userID,
			}})
			if err != nil {
				fmt.Println(err)
				return
			}
			if m.Status == "left" {
				t.log.Info("m.Status == left " + t.chatInviteLink(chatID))
			}
			imHereChat = append(imHereChat, chatID)
		}
	}
}

//func (t *Telegram) getChatPhoto(chatid string) string {
//	chatId, _ := t.chat(chatid)
//	chat, err := t.t.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: struct {
//		ChatID             int64
//		SuperGroupUsername string
//	}{ChatID: chatId}})
//	if err != nil {
//		return ""
//	}
//	corpAvatar := ""
//	if chat.Photo != nil {
//		fileconfig := tgbotapi.FileConfig{FileID: chat.Photo.SmallFileID}
//		file, err := t.t.GetFile(fileconfig)
//		if err != nil {
//			t.log.ErrorErr(err)
//		}
//
//		corpAvatar = "https://api.telegram.org/file/bot" + t.t.Token + "/" + file.FilePath
//	}
//	return corpAvatar
//}
