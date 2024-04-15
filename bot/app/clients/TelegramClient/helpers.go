package TelegramClient

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

func (t *Telegram) nameOrNick(UserName, FirstName string) (name string) {
	if UserName != "" {
		name = UserName

	} else {
		name = FirstName
	}
	return name
}

//func (t *Telegram) GetAvatar(userid int64, name string) string {
//	AvatarTG := "https://thumb.cloud.mail.ru/weblink/thumb/xw1/VLES/v7tqy1nXQ/telegram.png"
//	userProfilePhotos, err := t.t.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{UserID: userid})
//	if err != nil || len(userProfilePhotos.Photos) == 0 {
//		AvatarTG = fmt.Sprintf("https://via.placeholder.com/128x128.png/%s/FFFFFF/?text=%s",
//			utils.GetRandomColor(), utils.ExtractUppercase(name))
//		return AvatarTG
//	}
//
//	fileconfig := tgbotapi.FileConfig{FileID: userProfilePhotos.Photos[0][0].FileID}
//	file, err := t.t.GetFile(fileconfig)
//	if err != nil {
//		t.log.ErrorErr(err)
//		return AvatarTG
//	}
//	return "https://api.telegram.org/file/bot" + t.t.Token + "/" + file.FilePath
//}

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
func (t *Telegram) chat(chatid string) (chatId int64, ThreadID int) {
	a := strings.SplitN(chatid, "/", 2)
	chatId, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		t.log.ErrorErr(err)
	}
	if len(a) > 1 {
		ThreadID, err = strconv.Atoi(a[1])
		if err != nil {
			t.log.ErrorErr(err)
		}
	}
	return chatId, ThreadID
}
func (t *Telegram) getLanguage(lang, key string) string {
	return t.storage.Dictionary.GetText(lang, key)
}
