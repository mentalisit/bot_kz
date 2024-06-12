package TelegramClient

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
)

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
			fmt.Println("Файл уже существует и имеет тот же размер, пропуск скачивания.")
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
