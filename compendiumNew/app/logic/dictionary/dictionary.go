package dictionary

import (
	"encoding/json"
	"fmt"
	"github.com/mentalisit/logger"
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

	dict.setDictionaryJson(getDictionaryRuJson())
	dict.setDictionaryJson(getDictionaryUkJson())
	dict.setDictionaryJson(getDictionaryEnJson())

	fmt.Printf("Load language")
	for s, _ := range dict.dictionary {
		fmt.Printf(" %s", s)
	}
	fmt.Printf("\n")

	return dict
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
