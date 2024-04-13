package dictionary

import (
	"github.com/mentalisit/logger"
)

type Dictionary struct {
	ua  map[string]string
	en  map[string]string
	ru  map[string]string
	log *logger.Logger
}

func NewDictionary(log *logger.Logger) *Dictionary {
	dict := &Dictionary{
		ua:  make(map[string]string),
		en:  make(map[string]string),
		ru:  make(map[string]string),
		log: log,
	}

	dict.setDictionary()

	return dict
}

func (dict *Dictionary) setDictionary() {

	dict.setDictionaryUa()
	dict.setDictionaryRu()
	//dict.setDictionaryEn()
	dict.setDictionaryEnJson()
}

func (dict *Dictionary) GetText(lang string, key string) string {

	var text string

	if lang == "ru" {
		text = dict.ru[key]
	} else if lang == "ua" {
		text = dict.ua[key]
	} else {
		text = dict.en[key]
	}

	if text == "" {
		if key == "" {
			text = "{'key' not specified}"
			dict.log.Error("{'key' not specified}")
		} else {
			//dict.log.Error(fmt.Sprintf("GetText lang:%s  key:%s", lang, key))
			//text = "{" + key + "}"
		}
	}

	return text
}
