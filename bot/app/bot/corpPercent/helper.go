package corpPercent

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

func percent(lvl int) int {
	p := 22
	for i := 2; i < lvl; i++ {
		p += 2
	}
	return p
}

func formatTime(ut int64) string {
	// Определите целевую дату
	targetDate := time.Unix(ut, 0)

	// Определите оставшееся время
	remainingTime := targetDate.Sub(time.Now().UTC())

	text := ""

	// Получите дни, часы и минуты из оставшегося времени
	days := remainingTime / (24 * time.Hour)
	if days > 0 {
		text += fmt.Sprintf("%dд ", days)
	}
	remainingTime = remainingTime % (24 * time.Hour)
	hours := remainingTime / time.Hour
	if hours > 0 {
		text += fmt.Sprintf("%dч ", hours)
	}
	remainingTime = remainingTime % time.Hour
	minutes := remainingTime / time.Minute
	if days == 0 && minutes > 0 {
		text += fmt.Sprintf("%dм", minutes)
	}
	return text
}

// Функция для сортировки среза строк по убыванию числовых значений первых двух символов
func sortByFirstTwoDigits(input []string) []string {
	// Создание кастомного типа для среза строк
	type sortableStrings []string

	// Реализация интерфейса sort.Interface для кастомного типа
	// Len возвращает длину среза
	// Less сравнивает строки по числовым значениям первых двух символов
	// Swap меняет местами элементы с указанными индексами
	var ss sortableStrings = input
	sort.Slice(ss, func(i, j int) bool {
		numI, _ := strconv.Atoi(ss[i][:2])
		numJ, _ := strconv.Atoi(ss[j][:2])
		return numI > numJ // сортировка по убыванию
	})

	// Преобразование кастомного типа обратно в срез строк
	return ss
}
