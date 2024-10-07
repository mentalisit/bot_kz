package bot

import (
	"fmt"
	"kz_bot/models"
	"sort"
	"strings"
)

// lang ok
// нужно переделать полностью
func (b *Bot) Top(in models.InMessage) {
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
		return fmt.Sprintf("%d. %s - %d (%d)\n", number, top.Name, top.Numkz, top.Points)
	}

	numEvent := b.storage.Event.NumActiveEvent(in.Config.CorpName)
	b.ifTipSendTextDelSecond(in, b.getText(in, "scan_db"), 5)

	if numEvent == 0 {
		if in.Lvlkz != "" {
			message = fmt.Sprintf("\xF0\x9F\x93\x96 %s %s%s:\n",
				b.getText(in, "top_participants"), b.getText(in, "rs"), in.Lvlkz)
			resultsTop = b.storage.Top.TopLevelPerMonthNew(in.Config.CorpName, in.Lvlkz)
		} else {
			message = fmt.Sprintf("\xF0\x9F\x93\x96 %s:\n", b.getText(in, "top_participants"))
			resultsTop = b.storage.Top.TopAllPerMonthNew(in.Config.CorpName)
		}
	} else {
		if in.Lvlkz != "" {
			message = fmt.Sprintf("\xF0\x9F\x93\x96 %s %s %s%s\n     ",
				b.getText(in, "top_participants"), b.getText(in, "event"), b.getText(in, "rs"), in.Lvlkz)
			resultsTop = b.storage.Top.TopEventLevelNew(in.Config.CorpName, in.Lvlkz, numEvent)
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
		message2 = fmt.Sprintf("%s\nTotal: %d", message2, allpoints)

	}

	if in.Tip == ds {
		mid := b.client.Ds.SendEmbedText(in.Config.DsChannel, message, message2)
		b.client.Ds.DeleteMesageSecond(in.Config.DsChannel, mid, 600)
	} else if in.Tip == tg {
		b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, message+message2, 600)
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
