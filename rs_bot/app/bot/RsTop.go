package bot

import (
	"fmt"
	"regexp"
	"rs/models"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GetCorpAlias(config models.CorporationConfig) string {
	corpName := ""
	switch config.CorpName {
	case "Корпорация  \"РУСЬ\".сбор-на-кз":
		corpName = "русь "
	case "IX Legion.сбор-на-кз-бот":
		corpName = "IX Легион"
	case "Неизбежный Рок/КЗ сбор":
		corpName = "Неизбежный Рок"
	case "Повстанцы Хаоса.кз-чат":
		corpName = "povstanci"
	case "ЛУННЫЙ ФЕНИКС/КЗ и новости":
		corpName = "Лунный Феникс"
	case "ДИВЕРСАНТЫ 2.0/Бот кз":
		corpName = "ДИВЕРСАНТЫ"

	default:
		return ""
	}
	return corpName
}

// lang ok
// нужно переделать полностью
func (b *Bot) Top(in models.InMessage) {
	corpAlias := GetCorpAlias(in.Config)
	if corpAlias != "" {
		b.TopGame(in)
		return
	}
	b.iftipdelete(in)
	number := 1
	message := ""
	message2 := ""
	var allpoints int
	var resultsTop []models.Top
	format := func(top models.Top) string {
		if top.Points == 0 {
			return fmt.Sprintf("%d. %s - %d \n", number, top.Name, top.Numkz)
		}
		allpoints += top.Points
		return fmt.Sprintf("%d. %s - %d (%s)\n", number, top.Name, top.Numkz, formatNumber(top.Points))
	}

	numEvent := b.storage.Event.NumActiveEvent(in.Config.CorpName)
	b.ifTipSendTextDelSecond(in, b.getText(in, "scan_db"), 5)

	_, level := in.TypeRedStar()
	if numEvent == 0 {
		if in.RsTypeLevel != "" {
			message = fmt.Sprintf("\xF0\x9F\x93\x96 %s %s%s:\n",
				b.getText(in, "top_participants"), b.getText(in, "rs"), level)
			resultsTop = b.storage.Top.TopLevelPerMonthNew(in.Config.CorpName, in.RsTypeLevel)
		} else {
			message = fmt.Sprintf("\xF0\x9F\x93\x96 %s:\n", b.getText(in, "top_participants"))
			resultsTop = b.storage.Top.TopAllPerMonthNew(in.Config.CorpName)
		}
	} else {
		if in.RsTypeLevel != "" {
			message = fmt.Sprintf("\xF0\x9F\x93\x96 %s %s %s%s\n     ",
				b.getText(in, "top_participants"), b.getText(in, "event"), b.getText(in, "rs"), level)
			resultsTop = b.storage.Top.TopEventLevelNew(in.Config.CorpName, in.RsTypeLevel, numEvent)
		} else {
			message = fmt.Sprintf("\xF0\x9F\x93\x96 %s %s:\n",
				b.getText(in, "top_participants"), b.getText(in, "event"))
			resultsTop = b.storage.Top.TopAllEventNew(in.Config.CorpName, numEvent)
		}
		resultsTop = mergeAndSumTops(resultsTop)
	}
	if len(resultsTop) == 0 {
		b.ifTipSendTextDelSecond(in, b.getText(in, "no_history"), 20)
		return
	} else if len(resultsTop) > 0 {
		b.ifTipSendTextDelSecond(in, b.getText(in, "form_list"), 5)
		for _, top := range resultsTop {
			message2 = message2 + format(top)
			number++
		}
	}
	if allpoints != 0 {
		message2 = fmt.Sprintf("%s\nTotal: %s", message2, formatNumber(allpoints))
	}

	if in.Tip == ds {
		mid := b.client.Ds.SendEmbedText(in.Config.DsChannel, message, message2)
		b.client.Ds.DeleteMessageSecond(in.Config.DsChannel, mid, 600)
	} else if in.Tip == tg {
		text := message + message2
		if in.Config.Guildid != "" {
			b.ifTipSendTextDelSecond(in, b.getText(in, "form_list"), 10)
			text = b.client.Ds.ReplaceTextMessage(text, in.Config.Guildid)
		}
		text = strings.ReplaceAll(text, "@", "")
		b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 600)
	}
	b.TopGame(in)
}

func (b *Bot) TopGame(in models.InMessage) {
	b.iftipdelete(in)

	corpName := GetCorpAlias(in.Config)
	if corpName == "" {
		b.log.Info("corpNameAlias not found ")
		return
	}

	number := 1
	title := ""
	message2 := ""
	var allpoints int
	var resultsTop []models.Top

	format := func(top models.Top) string {
		if top.Points == 0 {
			return fmt.Sprintf("%d. %s - %d \n", number, top.Name, top.Numkz)
		}
		allpoints += top.Points
		return fmt.Sprintf("%d. %s - %d (%s)\n", number, top.Name, top.Numkz, formatNumber(top.Points))
	}

	// Внутренняя функция для сбора статистики из списка корпораций
	collectStats := func(corps []string, eventID int) ([]models.PlayerStats, error) {
		aggregated := make(map[string]models.PlayerStats)
		for _, name := range corps {
			statsList, err := b.storage.Battles.BattlesGetAll(name, eventID)
			if err != nil {
				return nil, err
			}
			for _, s := range statsList {
				entry := aggregated[s.Player]
				entry.Player = s.Player
				entry.Points += s.Points
				entry.Runs += s.Runs
				aggregated[s.Player] = entry
			}
		}

		var final []models.PlayerStats
		for _, s := range aggregated {
			final = append(final, s)
		}
		return final, nil
	}
	// ---------------------------------------

	nextDateStart, nextDateStop, messageNews := b.storage.Battles.ReadEventScheduleAndMessage()
	date1 := time.Now().UTC().Format("02-01-2006")
	date2 := time.Now().UTC().Add(24 * time.Hour).Format("02-01-2006")
	numEvent := 0

	if date1 == nextDateStart || date2 == nextDateStop {
		numEvent = getSeasonNumber(messageNews)
		title = fmt.Sprintf("Сезон №%d %s - %s\n", numEvent, nextDateStart, nextDateStop)
	}

	b.ifTipSendTextDelSecond(in, b.getText(in, "scan_db"), 5)

	_, level := in.TypeRedStar()
	if numEvent == 0 {
		battlesTopGetAll, err := b.storage.Battles.BattlesTopGetAll(corpName)
		if err != nil {
			b.log.ErrorErr(err)
			return
		}
		if in.RsTypeLevel != "" {
			title = fmt.Sprintf("\xF0\x9F\x93\x96 %s %s%s:\n",
				b.getText(in, "top_participants"), b.getText(in, "rs"), level)

			levelInt, _ := strconv.Atoi(level)
			for _, top := range battlesTopGetAll {
				if top.Level == levelInt {
					resultsTop = append(resultsTop, models.Top{Name: top.Name, Numkz: top.Count, Points: 0})
				}
			}
		} else {
			title = fmt.Sprintf("\xF0\x9F\x93\x96 %s:\n", b.getText(in, "top_participants"))
			for _, t := range battlesTopGetAll {
				resultsTop = append(resultsTop, models.Top{Name: t.Name, Numkz: t.Count})
			}
		}
		sort.Slice(resultsTop, func(i, j int) bool { return resultsTop[i].Numkz > resultsTop[j].Numkz })
	} else {
		// Логика для сезона с объединением корпораций
		corpsToMerge := []string{corpName}

		// Если это "rusb", добавляем "best". Сюда можно дописать любое условие или доп. имена.
		if corpName == "русь " {
			corpsToMerge = append(corpsToMerge, "best")
			corpsToMerge = append(corpsToMerge, "IX Легион")
		}

		if corpName == "IX Легион" {
			corpsToMerge = append(corpsToMerge, "best")
			corpsToMerge = append(corpsToMerge, "русь ")
		}

		battlesGetAll, err := collectStats(corpsToMerge, numEvent)
		if err != nil {
			b.log.ErrorErr(err)
			return
		}

		for _, stats := range battlesGetAll {
			resultsTop = append(resultsTop, models.Top{
				Name:   stats.Player,
				Numkz:  stats.Runs,
				Points: stats.Points,
			})
		}
		resultsTop = mergeAndSumTops(resultsTop)
	}

	if len(resultsTop) == 0 {
		b.ifTipSendTextDelSecond(in, b.getText(in, "no_history"), 20)
		return
	}

	b.ifTipSendTextDelSecond(in, b.getText(in, "form_list"), 5)
	for _, top := range resultsTop {
		message2 += format(top)
		number++
	}

	if allpoints != 0 {
		message2 = fmt.Sprintf("%s\nTotal: %s", message2, formatNumber(allpoints))
	}

	if in.Tip == ds {
		mid := b.client.Ds.SendEmbedText(in.Config.DsChannel, title, message2)
		b.client.Ds.DeleteMessageSecond(in.Config.DsChannel, mid, 600)
	} else if in.Tip == tg {
		b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, title+message2, 600)
	}
}

func mergeAndSumTops(tops []models.Top) []models.Top {
	merged := make(map[string]models.Top)

	for _, top := range tops {
		// Удаляем знак $ из начала строки
		name := top.Name
		if strings.HasPrefix(name, "$") {
			name = strings.TrimPrefix(name, "$")
		}

		// Объединяем элементы с одинаковыми именами
		if existing, found := merged[name]; found {
			existing.Numkz += top.Numkz
			existing.Points += top.Points
			merged[name] = existing
		} else {
			top.Name = name
			merged[name] = top
		}
	}

	// Преобразуем карту обратно в срез
	var result []models.Top
	for _, top := range merged {
		result = append(result, top)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Points > result[j].Points
	})

	return result
}

func (b *Bot) UpdateTopEventForCorporation(Corp models.CorporationHistory, title string) {
	number := 1
	message2 := ""
	var allPoints int
	var resultsTop []models.Top
	format := func(top models.Top) string {
		if top.Points == 0 {
			return fmt.Sprintf("%d. %s - %d \n", number, top.Name, top.Numkz)
		}
		allPoints += top.Points
		return fmt.Sprintf("%d. %s - %d (%s)\n", number, top.Name, top.Numkz, formatNumber(top.Points))
	}

	numEvent := b.storage.Event.NumActiveEvent(Corp.CorpName)

	if numEvent == 0 {
		resultsTop = b.storage.Top.TopAllPerMonthNew(Corp.CorpName)
	} else {
		resultsTop = b.storage.Top.TopAllEventNew(Corp.CorpName, numEvent)
	}

	resultsTop = mergeAndSumTops(resultsTop)

	if len(resultsTop) == 0 {
		return
	} else if len(resultsTop) > 0 {
		for _, top := range resultsTop {
			message2 = message2 + format(top)
			number++
		}
	}
	if allPoints != 0 {
		message2 = fmt.Sprintf("%s\nВсего: %s\nОбновлено: <t:%d:R>",
			message2, formatNumber(allPoints), time.Now().UTC().Unix())
	}

	b.client.Ds.SendEmbedText(Corp.ChannelDs, title, message2)
}

func getSeasonNumber(text string) int {
	re := regexp.MustCompile(`Season (\d+)`)
	matches := re.FindStringSubmatch(text)

	if len(matches) < 2 {
		return 0
	}
	seasonNumber, _ := strconv.Atoi(matches[1])
	return seasonNumber
}

func formatNumber(n int) string {
	s := fmt.Sprintf("%d", n)
	var sb strings.Builder
	length := len(s)
	for i, c := range s {
		if i > 0 && (length-i)%3 == 0 {
			sb.WriteRune(' ')
		}
		sb.WriteRune(c)
	}
	return sb.String()
}
