package words

import "kz_bot/storage/dictionary"

type Words struct {
	ua            map[string]string
	en            map[string]string
	ru            map[string]string
	newDictionary *dictionary.Dictionary
}

func NewWords(dictionary *dictionary.Dictionary) *Words {
	w := &Words{
		ua:            make(map[string]string),
		en:            make(map[string]string),
		ru:            make(map[string]string),
		newDictionary: dictionary,
	}
	w.setWords()

	return w
}

func (w *Words) setWords() {
	w.setWordsUa()
	w.setWordsRu()
	w.setWordsEn()

}
func (w *Words) GetWords(lang string, key string) string {
	text := w.newDictionary.GetText(lang, key)
	if text != "" {
		return text
	}

	if lang == "ru" {
		return w.ru[key]
	} else if lang == "ua" {
		return w.ua[key]
	} else {
		return w.en[key]
	}
}
