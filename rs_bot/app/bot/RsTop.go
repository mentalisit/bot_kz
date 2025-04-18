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
}

func (b *Bot) TopForGame(in models.InMessage) {
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

	b.ifTipSendTextDelSecond(in, b.getText(in, "scan_db"), 5)
	message = fmt.Sprintf("\xF0\x9F\x93\x96 %s %s:\n",
		b.getText(in, "top_participants"), b.getText(in, "event"))
	star, _ := b.storage.Top.RedStarFightGetStar()
	if star == nil {
		return
	}
	rr := make(map[string][]models.Top)
	for _, fight := range star {
		if fight.CountParticipants() != 0 {
			for _, participants := range fight.ParticipantsSlice {
				name := combineNames(participants.PlayerName)
				rr[name] = append(rr[name], models.Top{
					Name:   name,
					Numkz:  1,
					Points: fight.Points / fight.CountParticipants(),
				})
			}
		}
	}
	rrr := make(map[string]models.Top)

	for s, tops := range rr {
		var tr models.Top
		tr.Name = s
		for _, top := range tops {
			tr.Points += top.Points
			tr.Numkz += top.Numkz
		}
		rrr[s] = tr
	}

	for _, top := range rrr {
		resultsTop = append(resultsTop, top)
	}

	sort.Slice(resultsTop, func(i, j int) bool {
		return resultsTop[i].Points > resultsTop[j].Points
	})

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
}

func combineNames(r string) string {
	switch r {
	case "Mchuleft", "Valenvaryon":
		return "Mchuleft"
	case "Overturned", "Overturned-1.1":
		return "Overturned"
	case "RedArrow", "Light Matter", "Dark Matter", "Drake":
		return "RedArrow"
	case "Коньячный ЗАВОД", "falcon_2":
		return "falcon_2"
	case "Silent_Noise", "WarySamurai1055":
		return "Silent_Noise"
	case "arsenium23", "Tabu 666", "Psyker":
		return "Tabu"
	case "delov@r", "delovar", "Plague":
		return "delovar"
	case "iVanCoMik", "eVanCoMik", "VanCoMik":
		return "VanCoMik"
	case "Альтаир", "АЛЬТАИР", "Storm":
		return "Альтаир"
	case "Гэндальф серый", "Ёжик71":
		return "Ёжик71"
	case "Джонни_De", "JonnyDe":
		return "JonnyDe"
	case "N@N", "ChubbChubbs":
		return "ChubbChubbs"
	case "Nixonblade", "TimA":
		return "Nixonblade"

	default:
		return r
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

func (b *Bot) UpdateTopEvent() {
	corp1 := models.CorporationHistory{
		CorpName:  "Корпорация  \"РУСЬ\".сбор-на-кз",
		ChannelDs: "1198012575615561979",
	}
	corp2 := models.CorporationHistory{
		CorpName:  "IX Legion.сбор-на-кз-бот",
		ChannelDs: "1253851338408857641",
	}

	nextDateStart, nextDateStop, messagee := b.storage.Event.ReadEventScheduleAndMessage()
	date1 := time.Now().UTC().Format("02-01-2006")
	date2 := time.Now().UTC().Add(24 * time.Hour).Format("02-01-2006")
	if date1 == nextDateStart || date2 == nextDateStop {
		title := fmt.Sprintf("Сезон №%d   %s  -  %s", getSeasonNumber(messagee), nextDateStart, nextDateStop)
		b.UpdateTopEventForCorporation(corp1, title)
		b.UpdateTopEventForCorporation(corp2, title)
	}
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
