package dictionary

import (
	"fmt"
	"github.com/mentalisit/logger"
)

type Dictionary struct {
	//ua  map[string]string
	//en  map[string]string
	//ru  map[string]string
	log        *logger.Logger
	dictionary map[string]map[string]string
}

func NewDictionary(log *logger.Logger) *Dictionary {
	dict := &Dictionary{
		//ua:  make(map[string]string),
		//en:  make(map[string]string),
		//ru:  make(map[string]string),
		log:        log,
		dictionary: make(map[string]map[string]string),
	}

	dict.setDictionary()

	return dict
}

func (dict *Dictionary) setDictionary() {
	dict.setDictionaryUaJson()
	dict.setDictionaryRuJson()
	dict.setDictionaryEnJson()
}

func (dict *Dictionary) GetText(lang string, key string) string {

	var text string

	text = dict.dictionary[lang][key]

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
