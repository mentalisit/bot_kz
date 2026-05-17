package bot2

import (
	"fmt"
	"regexp"
	"rs/bot2/helpers"
	"rs/models"
	"strings"
)

func (b *Bot) EventStatistic(in *models.InMessageV2) (bb bool) {
	if strings.HasPrefix(in.Text, ".quality") {
		b.BattleStats(in)
		return false
	}
	var statsRegex = regexp.MustCompile(`^(?i)(\.статистика|\.statistic)\s+(.+)$`)
	// [1] - первая группа захвата (ключевое слово)
	// [2] - вторая группа захвата (имя)
	matches := statsRegex.FindStringSubmatch(in.Text)
	if len(matches) < 3 {
		return false
	}
	name := strings.TrimSpace(matches[2])

	b.deleteInMessage(in)

	statistics, err := b.storage.StatisticGetName(name)
	if err != nil || len(statistics) == 0 || statistics == nil {
		b.sendTextAfterDeleteSecond(in, fmt.Sprintf("Данные на имя %s не найдены.", name), 20)
		return false
	}

	if matches[1] == ".statistic" {
		picData := helpers.PicStatistic(name, statistics)
		if picData != nil {
			for ch, i := range in.Config.Channels {
				if i.TypeMessenger == ds {
					_ = b.client.Ds.SendChannelPic(ch, "", picData)
				} else if i.TypeMessenger == tg {
					_, _ = b.client.Tg.SendPic(ch, "", picData)
				} else {
					b.log.Info("need implement for " + i.TypeMessenger)
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
			go b.client.Ds.SendChannelDelSecond(in.Messenger.ChannelId, fmt.Sprintf("```%s```", text), 180)
		} else if in.Tip == tg {
			go b.client.Tg.SendChannelDelSecond(in.Messenger.ChannelId, text, 180)
		} else {
			b.log.Info(fmt.Sprintf("need make for %s \n", in.Tip))
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

func (b *Bot) BattleStats(in *models.InMessageV2) {
	//_, number := ParseQualityCommand(in.Text)
	//if number == 0 {
	//	number = 6
	//}
	//b.deleteInMessage(in)
	//
	//fmt.Println(in.Config.Uid)
	//info, err := b.storage.ReadCorpInfo(in.Config.Uid)
	//if err != nil || info == nil {
	//	return
	//}
	//
	//stats := b.storage.GetBattleStats(info.CorpName, number)
	//if stats == nil {
	//	return
	//}
	//picData := helpers.BattleStatsImage(info.CorpName, stats)
	//if picData != nil {
	//	for ch, i := range in.Config.Channels {
	//		if i.TypeMessenger == ds {
	//			_ = b.client.Ds.SendChannelPic(ch, "quality", picData)
	//		} else if i.TypeMessenger == tg {
	//			_, _ = b.client.Tg.SendPic(ch, "quality", picData)
	//		} else {
	//			b.log.Info("need implement for " + i.TypeMessenger)
	//		}
	//	}
	//}
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
