package dictionary

import (
	"encoding/json"
	"fmt"
	gt "github.com/bas24/googletranslatefree"
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
	} else {
		for s, _ := range dict.dictionary {
			fmt.Println(s)
		}
	}

	dict.setDictionaryJson(getDictionaryRuJson())
	dict.setDictionaryJson(getDictionaryEnJson())
	dict.setDictionaryJson(getDictionaryUaJson())

	//dict.TranslateViaGoogle("de")
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

func (dict *Dictionary) TranslateViaGoogle(lang string) error {
	langName := nameMap[lang]

	dictTemp := make(map[string]map[string]string)
	dictTemp[lang] = make(map[string]string)
	fmt.Print("Create translate " + lang)
	for key, s := range dict.dictionary["en"] {
		if key == "language_switched_to" && langName != "" {
			s = strings.Replace(s, "English", langName, 1)
		}
		text, err := gt.Translate(s, "auto", lang)
		if err != nil {
			dict.log.ErrorErr(err)
			return err
		}
		dictTemp[lang][key] = text
		fmt.Print(".")
	}
	fmt.Print("\n")
	// Сериализуем структуру в JSON
	jsonData, err := json.MarshalIndent(dictTemp, "", "  ")
	if err != nil {
		fmt.Println("Ошибка сериализации JSON:", err)
		return err
	}

	// Получаем текущую директорию
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Записываем JSON в файл
	fileName := lang + ".json"
	os.Mkdir("locale", 0750)
	path := filepath.Join(currentDir, "locale", fileName)
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Ошибка записи в файл:", err)
		return err
	}

	fmt.Println("JSON успешно записан в файл", fileName)
	dict.readDirLocale()
	return nil
}
func (dict *Dictionary) CheckTranslateLanguage(lang string) bool {
	for s, _ := range dict.dictionary {
		if s == lang {
			return true
		}
	}
	return false
}
