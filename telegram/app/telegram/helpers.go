package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"telegram/models"
	"telegram/telegram/restapi"
	"time"
)

const nickname = "Для того что бы БОТ мог Вас индентифицировать, создайте уникальный НикНей в настройках. Вы можете использовать a-z, 0-9 и символы подчеркивания. Минимальная длина - 5 символов."

var usernameMap map[string]int

func (t *Telegram) nickName(u *tgbotapi.User, channel string) string {
	if u.UserName == "" {
		if usernameMap[u.String()] == 0 {
			t.SendChannelDelSecond(channel, nickname, 60)
			usernameMap[u.String()] = 1
		} else if usernameMap[u.String()]%5 == 0 {
			t.SendChannelDelSecond(channel, nickname, 60)
			usernameMap[u.String()] += 1
		} else {
			usernameMap[u.String()] += 1
		}
	}
	return u.String()
}

func (t *Telegram) chat(chatid string) (chatId int64, threadID int) {
	a := strings.SplitN(chatid, "/", 2)
	chatId, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		t.log.ErrorErr(err)
	}
	if len(a) > 1 {
		threadID, err = strconv.Atoi(a[1])
		if err != nil {
			t.log.ErrorErr(err)
		}
	}
	return chatId, threadID
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
		errt := restapi.SendInsertTimer(models.Timer{
			Tgmesid:  strconv.Itoa(id),
			Tgchatid: chatid,
			Timed:    second,
		})
		if errt != nil {
			t.log.ErrorErr(err)
		}
	}
}
func (t *Telegram) DelMessage(chatid string, idSendMessage int) {
	a := strings.SplitN(chatid, "/", 2)
	chatId, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		t.log.ErrorErr(err)
	}
	_, _ = t.t.Request(tgbotapi.DeleteMessageConfig(tgbotapi.NewDeleteMessage(chatId, idSendMessage)))
}

func (t *Telegram) getAvatarIsExist(userid int64) string {

	userProfilePhotos, err := t.t.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{UserID: userid})
	if err != nil || len(userProfilePhotos.Photos) == 0 {
		return ""
	}

	fileconfig := tgbotapi.FileConfig{FileID: userProfilePhotos.Photos[0][0].FileID}
	file, err := t.t.GetFile(fileconfig)
	if err != nil {
		t.log.InfoStruct("userProfilePhotos.Photos", userProfilePhotos.Photos)
		t.log.ErrorErr(err)
		return ""
	}
	return "https://api.telegram.org/file/bot" + t.t.Token + "/" + file.FilePath
}
