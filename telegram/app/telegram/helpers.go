package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"telegram/models"
	"time"
	"unicode"
)

const nickname = "Для того что бы БОТ мог Вас индентифицировать, создайте уникальный НикНей в настройках. Вы можете использовать a-z, 0-9 и символы подчеркивания. Минимальная длина - 5 символов."

func (t *Telegram) nickName(u *tgbotapi.User, channel string) string {
	if u.UserName == "" {
		if t.usernameMap[u.String()] == 0 {
			t.SendChannelDelSecond(channel, nickname, 60)
			t.usernameMap[u.String()] = 1
		} else if t.usernameMap[u.String()]%5 == 0 {
			t.SendChannelDelSecond(channel, nickname, 60)
			t.usernameMap[u.String()] += 1
		} else {
			t.usernameMap[u.String()] += 1
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
func (t *Telegram) DelMessageSecond(chatid string, idSendMessage string, second int) error {
	id, err := strconv.Atoi(idSendMessage)
	if err != nil {
		return err
	}
	if second <= 60 {
		go func() {
			time.Sleep(time.Duration(second) * time.Second)
			t.DelMessage(chatid, id)
		}()
	} else {
		t.Storage.Db.TimerInsert(models.Timer{
			Tgmesid:  strconv.Itoa(id),
			Tgchatid: chatid,
			Timed:    second,
		})
	}
	return nil
}
func (t *Telegram) DelMessage(chatid string, idSendMessage int) error {
	chatId, _ := t.chat(chatid)
	_, err := t.t.Request(tgbotapi.DeleteMessageConfig(tgbotapi.NewDeleteMessage(chatId, idSendMessage)))
	return err
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

func (t *Telegram) chatName(chatid string) (chatName string) {
	id, _ := t.chat(chatid)
	r, err := t.t.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: id}})
	if err != nil {
		t.log.ErrorErr(err)
	}
	return r.Title
}
func (t *Telegram) loadAvatarIsExist(userid int64) string {

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
	_, url := t.SaveAvatarLocalCache(strconv.FormatInt(userid, 10), "https://api.telegram.org/file/bot"+t.t.Token+"/"+file.FilePath)
	return url
}
func (t *Telegram) SaveAvatarLocalCache(userID, url string) (bool, string) {
	if url == "" {
		return false, ""
	}
	folder := "compendium/avatars"

	// Создаем HTTP-запрос
	resp, err := http.Get(url)
	if err != nil {
		t.log.ErrorErr(fmt.Errorf("ошибка при выполнении запроса: %v", err))
		return false, ""
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		t.log.ErrorErr(fmt.Errorf("неправильный статус-код: %d", resp.StatusCode))
		return false, ""
	}

	// Получаем расширение файла из URL
	fileExt := filepath.Ext(url)
	if fileExt == "" {
		t.log.ErrorErr(fmt.Errorf("не удалось определить расширение файла"))
		return false, ""
	}

	// Создаем директорию, если она не существует
	err = os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		t.log.ErrorErr(fmt.Errorf("ошибка при создании директории: %v", err))
		return false, ""
	}

	// Полный путь к файлу
	filename := userID + fileExt
	filePath := filepath.Join(folder, filename)
	// Используем path.Join для формирования URL
	newUrl := "https://compendiumnew.mentalisit.myds.me/" + path.Join(folder, filename)

	// Проверяем, существует ли файл
	if fileInfo, err := os.Stat(filePath); err == nil {
		// Файл существует, проверяем его размер
		if fileInfo.Size() == resp.ContentLength {
			// Размеры совпадают, пропускаем скачивание
			return true, newUrl
		}
	}

	// Создаем файл
	file, err := os.Create(filePath)
	if err != nil {
		t.log.ErrorErr(fmt.Errorf("ошибка при создании файла: %v", err))
		return false, ""
	}
	defer file.Close()

	// Копируем содержимое ответа в файл
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		t.log.ErrorErr(fmt.Errorf("ошибка при записи файла: %v", err))
		return false, ""
	}

	return true, newUrl
}

func DetectLanguage(s string) string {
	// Украинские символы, которых нет в русском языке
	ukrainianChars := "ґєіїҐЄІЇ"
	// Диапазон Unicode для русских символов
	var russianChars []*unicode.RangeTable
	russianChars = append(russianChars, unicode.Cyrillic)
	// Немецкие символы, которых нет в английском языке
	//germanChars := "ÄÖÜß"

	// Счетчики для каждого языка
	ukrainianCount := 0
	russianCount := 0
	//germanCount := 0

	contains := func(char rune, chars string) bool {
		for _, c := range chars {
			if c == char {
				return true
			}
		}
		return false
	}

	for _, char := range s {
		switch {
		case contains(char, ukrainianChars):
			ukrainianCount++
		case unicode.IsOneOf(russianChars, char):
			russianCount++
			//case contains(char, germanChars):
			//germanCount++
		}
	}

	if ukrainianCount > 0 {
		return "uk"
	} else if russianCount > 0 {
		return "ru"
		//} else if germanCount > 0 {
		//	return "de"
	} else {
		return "en"
	}
}

func (t *Telegram) CheckAdminTg(chatid string, name string) (admin bool) {
	chatId, _ := t.chat(chatid)
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

func (t *Telegram) GetAvatarUrl(userID string) string {
	userId, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		t.log.ErrorErr(err)
		return ""
	}
	photos, err := t.t.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{UserID: userId, Limit: 1})
	if err != nil || len(photos.Photos) == 0 {
		return ""
	}

	// Берем первый файл из списка фотографий профиля
	fileID := photos.Photos[0][0].FileID

	// Получаем файл по file_id
	file, err := t.t.GetFile(tgbotapi.FileConfig{
		FileID: fileID,
	})
	if err != nil {
		return ""
	}

	fileURL := file.Link(t.t.Token)

	_, newUrl := t.SaveAvatarLocalCache(strconv.FormatInt(userId, 10), fileURL)

	return newUrl
}
