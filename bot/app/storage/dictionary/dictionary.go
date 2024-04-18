package dictionary
//
import (
	"fmt"
	"github.com/mentalisit/logger"
	"encoding/json"
)

type Dictionary struct {
	//ua  map[string]string
	//en  map[string]string
	//ru  map[string]string
	dictionary  map[string]map[string]string
	log *logger.Logger
}

func NewDictionary(log *logger.Logger) *Dictionary {
	dict := &Dictionary{
		//ua:  make(map[string]string),
		//en:  make(map[string]string),
		//ru:  make(map[string]string),
		dictionary:  make(map[string]map[string]string),
		log: log,
	}

	dict.setDictionary()

	return dict
}

func (dict *Dictionary) setDictionary() {
	dict.setDictionaryJson(getDictionaryEnJson())
	dict.setDictionaryJson(getDictionaryRuJson())
	dict.setDictionaryJson(getDictionaryUaJson())
}

func (dict *Dictionary) GetText(lang string, key string) string {

	//var text string
	var text = dict.dictionary[lang][key]

	//if lang == "ru" {
	//	text = dict.ru[key]
	//} else if lang == "ua" {
	//	text = dict.ua[key]
	//} else {
	//	text = dict.en[key]
	//}

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

func (dict *Dictionary) setDictionaryJson(jsonText string) {

	var dictTemp map[string]map[string]string

	err := json.Unmarshal([]byte(jsonText), &dictTemp)
	if err != nil {
		dict.log.ErrorErr(err)
	}

	for key, val := range dictTemp {
       dict.dictionary[key] = val
  	}
}