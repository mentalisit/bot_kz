package helpers

import "unicode"

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
