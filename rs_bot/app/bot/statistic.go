package bot

import (
	"fmt"
	"regexp"
	"rs/bot/helpers"
	"rs/models"
	"strconv"
	"strings"
)

func (b *Bot) EventStatistic(in models.InMessage) (bb bool) {
	if strings.HasPrefix(in.Mtext, ".quality") {
		b.BattleStats(in)
		return false
	}
	var statsRegex = regexp.MustCompile(`^(?i)(\.статистика|\.statistic)\s+(.+)$`)
	// [1] - первая группа захвата (ключевое слово)
	// [2] - вторая группа захвата (имя)
	matches := statsRegex.FindStringSubmatch(in.Mtext)
	if len(matches) < 3 {
		return false
	}
	name := strings.TrimSpace(matches[2])

	if in.Tip == ds {
		go b.client.Ds.DeleteMessageSecond(in.Config.DsChannel, in.Ds.Mesid, 60)
	} else if in.Tip == tg {
		go b.client.Tg.DelMessageSecond(in.Config.TgChannel, strconv.Itoa(in.Tg.Mesid), 60)
	}

	statistics, err := b.storage.Battles.StatisticGetName(name)
	if err != nil || len(statistics) == 0 || statistics == nil {
		b.ifTipSendTextDelSecond(in, fmt.Sprintf("Данные на имя %s не найдены.", name), 20)
		return false
	}

	if matches[1] == ".statistic" {
		picData := helpers.PicStatistic(name, statistics)
		if picData != nil {
			if in.IfDiscord() {
				b.client.Ds.SendChannelPic(in.Config.DsChannel, "", picData)
			}
			if in.IfTelegram() {
				id, err1 := b.client.Tg.SendPic(in.Config.TgChannel, "", picData)
				if err1 == nil {
					b.client.Tg.DelMessageSecond(in.Config.TgChannel, id, 300)
				}
			}
			return true
		}
	} else {
		const wEvent = 2
		const wLevel = 5
		const wRuns = 8
		const wPoints = 8

		// 1. Форматирование заголовка (для него используем обычные пробелы)
		header := fmt.Sprintf("%-*s %-*s %-*s%-s", wEvent, "Ивент", wLevel, "Уровень", wRuns, "Игры", "Очки")

		// Горизонтальная линия
		separator := strings.Repeat("-", 26)

		// Основная строка вывода
		text := fmt.Sprintf("Статистика игрока %s\n%s\n%s\n", name, separator, header)

		// 2. Цикл для вывода данных
		for _, s := range statistics {
			// Используем вспомогательную функцию для заполнения точек
			eventStr := padWithDots(s.EventId, wEvent)
			levelStr := padWithDots(s.Level, wLevel)
			runsStr := padWithDots(s.Runs, wRuns)
			pointsStr := padWithDots(s.Points, wPoints)

			// Объединяем строки, разделяя их пробелом (для наглядности)
			text += fmt.Sprintf("%s %s %s %s\n", eventStr, levelStr, runsStr, pointsStr)
		}

		fmt.Println(text)
		if in.Tip == ds {
			go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, fmt.Sprintf("```%s```", text), 180)
		} else if in.Tip == tg {
			go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 180)
		}

		//pic
		return true
	}
	return false
}

func padWithDots(value int, width int) string {
	// 1. Форматируем число с правым выравниванием
	// %*d: Правое выравнивание с заданной шириной
	paddedStr := fmt.Sprintf("%*d", width, value)

	// 2. Заменяем начальные пробелы точками
	// strings.ReplaceAll заменяет все пробелы, но нам нужно только те, что слева
	// Лучше найти первый непробельный символ и заменить все до него.

	// Находим индекс первого символа, который не является пробелом
	firstDigitIndex := strings.IndexFunc(paddedStr, func(r rune) bool {
		return r != ' '
	})

	if firstDigitIndex == -1 {
		// Если строка состоит только из пробелов (маловероятно), возвращаем ее.
		return paddedStr
	}

	// Создаем строку из точек
	dots := strings.Repeat(".", firstDigitIndex)

	// Возвращаем строку: точки + само число
	return dots + paddedStr[firstDigitIndex:]
}

func (b *Bot) BattleStats(in models.InMessage) {
	_, number := ParseQualityCommand(in.Mtext)
	if number == 0 {
		number = 6
	}
	if in.Tip == ds {
		go b.client.Ds.DeleteMessageSecond(in.Config.DsChannel, in.Ds.Mesid, 60)
	} else if in.Tip == tg {
		go b.client.Tg.DelMessageSecond(in.Config.TgChannel, strconv.Itoa(in.Tg.Mesid), 60)
	}
	corp := ""
	fmt.Println(in.Config.CorpName)
	if strings.HasPrefix(in.Config.CorpName, "IX Legion") {
		corp = "IX_Легион"
	} else if strings.HasPrefix(in.Config.CorpName, "Корпорация  \"РУСЬ\"") {
		corp = "rusb"
	}
	if corp == "" {
		return
	}

	stats := b.storage.Battles.GetBattleStats(corp, number)
	if stats == nil {
		return
	}
	picData := helpers.BattleStatsImage(corp, stats)
	if picData != nil {
		if in.IfDiscord() {
			b.client.Ds.SendChannelPic(in.Config.DsChannel, "quality", picData)
		}
		if in.IfTelegram() {
			id, err1 := b.client.Tg.SendPic(in.Config.TgChannel, "quality", picData)
			if err1 == nil {
				b.client.Tg.DelMessageSecond(in.Config.TgChannel, id, 300)
			}
		}
	}
}

func ParseQualityCommand(text string) (bool, int) {
	re := regexp.MustCompile(`^\.quality(?:\s+(\d{1,2}))?$`)
	matches := re.FindStringSubmatch(text)

	if matches == nil {
		return false, 0
	}

	if matches[1] == "" {
		return true, 0 // команда без числа
	}

	var number int
	fmt.Sscanf(matches[1], "%d", &number)
	return true, number
}
