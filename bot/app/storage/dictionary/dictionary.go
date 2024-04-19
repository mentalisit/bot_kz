package dictionary

import (
	"encoding/json"
	"fmt"
	"github.com/mentalisit/logger"
	"os"
	"path/filepath"
	"strings"
)

type Dictionary struct {
	log        *logger.Logger
	dictionary map[string]map[string]string
}

func NewDictionary(log *logger.Logger) *Dictionary {
	dict := &Dictionary{
		log:        log,
		dictionary: make(map[string]map[string]string),
	}

	err := dict.readDirLocale()
	if err != nil {
		log.ErrorErr(err)
	}

	dict.setDictionaryJson(getDictionaryRuJson())
	dict.setDictionaryJson(getDictionaryEnJson())
	dict.setDictionaryJson(getDictionaryUaJson())

	return dict
}

// заготовка
func (dict *Dictionary) readDirLocale() error {
	// Получаем текущую директорию
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Формируем путь к папке "locale"
	localeDir := filepath.Join(currentDir, "locale")

	// Проверяем существование папки "locale" если её нет возвращаем nil
	if _, err = os.Stat(localeDir); os.IsNotExist(err) {
		return nil
	}

	// Получаем список файлов в папке "locale"
	files, err := os.ReadDir(localeDir)
	if err != nil {
		return err
	}

	// Итерируем по файлам и выводим файлы с расширением ".json"
	for _, file := range files {

		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {

			readFile, err := os.ReadFile(filepath.Join(localeDir, file.Name()))
			if err != nil {
				return err
			}

			var dictTemp map[string]map[string]string
			err = json.Unmarshal(readFile, &dictTemp)
			if err != nil {
				return err
			}

			for key, val := range dictTemp {
				dict.dictionary[key] = val
			}
		}
	}
	return nil
}

func (dict *Dictionary) setDictionaryJson(jsonByte []byte) {

	var dictTemp map[string]map[string]string

	err := json.Unmarshal(jsonByte, &dictTemp)
	if err != nil {
		dict.log.ErrorErr(err)
	}

	for key, val := range dictTemp {
		dict.dictionary[key] = val
	}
}
func (dict *Dictionary) GetText(lang string, key string) string {

	text := dict.dictionary[lang][key]

	if text == "" {
		if key == "" {
			text = "{'key' not specified}"
			dict.log.Error("{'key' not specified}")
		} else {
			dict.log.Error(fmt.Sprintf("GetText lang:%s  key:%s", lang, key))
			text = "{" + key + "}"
		}
	}

	return text
}
