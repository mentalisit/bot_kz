package bot

import (
	"fmt"
	"kz_bot/models"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	emOK      = "✅"
	emCancel  = "❎"
	emRsStart = "🚀"
	emPl30    = "⌛"
	emPlus    = "➕"
	emMinus   = "➖"
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
func (b *Bot) getMap(in models.InMessage, numkzl int) map[string]string {
	var n map[string]string
	n = make(map[string]string)

	if in.Config.DsChannel != "" {
		darkStar, lvlkz := containsSymbolD(in.Lvlkz)
		var err error
		if darkStar {
			n["lvlkz"], err = b.client.Ds.RoleToIdPing(b.getText(in, "drs")+lvlkz, in.Config.Guildid)
		} else {
			n["lvlkz"], err = b.client.Ds.RoleToIdPing(b.getText(in, "rs")+in.Lvlkz, in.Config.Guildid)
		}
		if err != nil {
			b.log.Info(fmt.Sprintf("RoleToIdPing lvl %s CorpName %s err: %+v", in.Lvlkz, in.Config.CorpName, err))
		}
	}

	n["lang"] = in.Config.Country
	n["title"] = b.getText(in, "rs_queue")
	if strings.HasPrefix(in.Lvlkz, "d") {
		n["title"] = b.getText(in, "queue_drs")
	}

	n["description"] = fmt.Sprintf("👇 %s <:rs:918545444425072671> %s (%d) ",
		b.getLanguageText(in.Config.Country, "wishing_to"), n["lvlkz"], numkzl)
	n["EmbedFieldName"] = fmt.Sprintf(" %s %s\n%s %s\n%s %s",
		emOK, b.getLanguageText(in.Config.Country, "to_add_to_queue"),
		emCancel, b.getLanguageText(in.Config.Country, "to_exit_the_queue"),
		emRsStart, b.getLanguageText(in.Config.Country, "forced_start"))
	n["EmbedFieldValue"] = b.getLanguageText(in.Config.Country, "data_updated") + ": "
	n["buttonLevel"] = in.Lvlkz
	return n
}
