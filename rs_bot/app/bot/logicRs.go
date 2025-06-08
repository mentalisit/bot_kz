package bot

import (
	"fmt"
	"regexp"
	"rs/models"
	"strconv"
	"strings"
)

// lang ok
// ivent not lang
func (b *Bot) lRsPlus(in models.InMessage) (rs bool) {
	var variableFound string
	reLevelTime := regexp.MustCompile(`^([3-9]|1[0-2])([+]|-|!|[*])([$]|[?])?(\d|\d{2}|\d{3})?$`)
	arr := reLevelTime.FindAllStringSubmatch(in.Mtext, -1)
	if len(arr) > 0 {
		rs = true
		in.SetLevelRsOrDrs(arr[0][1])
		variableFound = arr[0][2]
		in.SetTimeRs(arr[0][4])

		if arr[0][3] == "$" {
			in.NameMention = "$" + in.NameMention
		} else if arr[0][3] == "?" {
			alt, _ := strconv.Atoi(arr[0][4])
			if alt == 0 || alt > 5 {
				alt = 1
			}
			in.SetTimeRs("30")
			b.darkAlt(in, alt)
			return rs
		}
	}

	switch variableFound {
	case "+":
		b.RsDarkPlus(in, "")
	case "-":
		b.RsMinus(in)
	case "*":
		{
			fmt.Println("using * logicRs50")
			b.RsDarkPlus(in, "")
		}

	case "!":
		{
			fmt.Println("using ! logicRs56")
			_, level := in.TypeRedStar()
			in.RsTypeLevel = "rs" + level
			b.RsDarkPlus(in, "")
		}

	default:
		rs = false
	}

	//solo
	//reSolo := regexp.MustCompile(`^([sSсС])([7-9]|1[0-2])(\+)?\s*(\d*)$`)
	//arrSolo := reSolo.FindAllStringSubmatch(in.Mtext, -1)
	//if len(arrSolo) > 0 {
	//	in.SetLevelRsOrDrs("solo" + arrSolo[0][2])
	//	in.SetTimeRs("1")
	//
	//	points := arrSolo[0][4]
	//	if points != "" {
	//		b.RsSoloPlusComplete(in, points)
	//		return true
	//	}
	//
	//	if arrSolo[0][3] == "+" {
	//		b.RsSoloPlus(in)
	//		return true
	//	}
	//}

	return rs
}

func (b *Bot) lSubs(in models.InMessage) (bb bool) {
	bb = false
	var subs string
	re := regexp.MustCompile(`^([+]|[-])([3-9]|[1][0-2])$`) // две переменные для добавления или удаления подписок
	arr := re.FindAllStringSubmatch(in.Mtext, -1)
	if len(arr) > 0 {
		in.SetLevelRsOrDrs(arr[0][2])
		subs = arr[0][1]
		bb = true

	}
	re1 := regexp.MustCompile(`^(Rs|rs)\s(S|s|u|U)\s([3-9]|[1][0-2])$`)
	arr1 := re1.FindAllStringSubmatch(in.Mtext, -1)
	if len(arr1) > 0 {
		in.SetLevelRsOrDrs(arr1[0][3])
		subs = arr1[0][2]
		bb = true
		if subs == "S" || subs == "s" {
			subs = "+"
		} else if subs == "U" || subs == "u" {
			subs = "-"
		}
	}

	reAdd := regexp.MustCompile(`^(подписать)\s([3-9]|[1][0-2])\s(@\w+(\s@\w+)*)$`)
	arrAdd := reAdd.FindAllStringSubmatch(in.Mtext, -1)
	if len(arrAdd) > 0 {
		adminTg, err := b.client.Tg.CheckAdminTg(in.Config.TgChannel, in.Username)
		if err != nil {
			b.log.ErrorErr(err)
		}
		if adminTg {
			bb = true
			splitNames := strings.Split(arrAdd[0][3], " ")
			for _, s := range splitNames {
				in.NameMention = s
				in.Username = s[1:]
				in.SetLevelRsOrDrs(arrAdd[0][2])
				go b.Subscribe(in)
			}
		}
	}

	switch subs {
	case "+":
		go b.Subscribe(in)
	case "-":
		go b.Unsubscribe(in)
	}
	return bb
}

func (b *Bot) lQueue(in models.InMessage) (bb bool) {
	re4 := regexp.MustCompile(`^([о]|[О]|[q]|[Q]|[Ч]|[ч])([3-9]|[1][0-2])$`) // две переменные для чтения  очереди
	arr4 := (re4.FindAllStringSubmatch(in.Mtext, -1))
	if len(arr4) > 0 {
		in.SetLevelRsOrDrs(arr4[0][2])
		bb = true
		b.QueueLevel(in)
	}
	//rus
	if in.Mtext == "Очередь" || in.Mtext == "очередь" || in.Mtext == "Черга" || in.Mtext == "черга" || in.Mtext == "Queue" || in.Mtext == "queue" {
		bb = true
		b.iftipdelete(in)
		b.QueueAll(in)
	}

	//todo придумать другую команду
	if in.Mtext == "Все очереди" {
		bb = true
		b.ifTipSendTextDelSecond(in, "Поиск.... ", 10)
		b.iftipdelete(in)
		b.ifTipSendTextDelSecond(in, b.otherQueue.MyQueue(), 30)
	}

	re4s := regexp.MustCompile(`^(Rs|rs)\s(Q|q)$`) // две переменные для чтения  очереди
	arr4s := (re4s.FindAllStringSubmatch(in.Mtext, -1))
	if len(arr4s) > 0 {
		bb = true
		go b.QueueAll(in) //проверка совместимости
	}
	return bb
}

func (b *Bot) lRsStart(in models.InMessage) (bb bool) {
	var rss string
	re5 := regexp.MustCompile(`^([3-9]|[1][0-2])([\+][\+])$`) //rs start
	arr5 := re5.FindAllStringSubmatch(in.Mtext, -1)
	if len(arr5) > 0 {
		bb = true
		in.SetLevelRsOrDrs(arr5[0][1])
		rss = arr5[0][2]
	} else {
		re5 = regexp.MustCompile(`^(Rs|rs)\s(Start|start)\s([3-9]|[1][0-2])$`) //rs start
		arr5 = re5.FindAllStringSubmatch(in.Mtext, -1)
		if len(arr5) > 0 {
			bb = true
			in.SetLevelRsOrDrs(arr5[0][3])
			rss = "++"
			//b.log.Println("Проверка совместимости принудительного старта ")
		}
	}
	reP := regexp.MustCompile(`^([3-9]|[1][0-2])([\+][\+][\+])$`) //p30pl
	arrP := reP.FindAllStringSubmatch(in.Mtext, -1)
	if len(arrP) > 0 {
		in.SetLevelRsOrDrs(arrP[0][1])
		bb = true
		go b.Pl30(in)
	}

	if rss == "++" {
		go b.RsStart(in)
	}
	return bb
}

// ivent not lang
func (b *Bot) lEvent(in models.InMessage) (bb bool) {
	//re7 := regexp.MustCompile(`^(["К"]|["к"]|["K"]|["k"])\s+([0-9]+)\s+([0-9]+)$`) // ивент
	//arr7 := (re7.FindAllStringSubmatch(in.Mtext, -1))
	//if len(arr7) > 0 {
	//	bb = true
	//	points, err := strconv.Atoi(arr7[0][3])
	//	if err != nil {
	//		b.log.ErrorErr(err)
	//	}
	//	numkz, err := strconv.Atoi(arr7[0][2])
	//	if err != nil {
	//		b.log.ErrorErr(err)
	//	}
	//	go b.EventPoints(in, numkz, points)
	//}
	//re7s := regexp.MustCompile(`^(rs|Rs)\s(p|P)\s([0-9]+)\s([0-9]+)$`)
	//arr7s := (re7s.FindAllStringSubmatch(in.Mtext, -1))
	//if len(arr7s) > 0 {
	//	bb = true
	//	points, err := strconv.Atoi(arr7[0][4])
	//	if err != nil {
	//		b.log.ErrorErr(err)
	//	}
	//	numkz, err := strconv.Atoi(arr7[0][3])
	//	if err != nil {
	//		b.log.ErrorErr(err)
	//	}
	//	go b.EventPoints(in, numkz, points)
	//}
	switch in.Mtext {
	case "Ивент старт":
		go b.EventStart(in)
		bb = true
	case "event add corp":
		go b.EventPreStart(in)
		bb = true

	case "Ивент стоп":
		go b.EventStop(in)
		bb = true
	}
	return bb
}

func (b *Bot) lTop(in models.InMessage) (bb bool) {

	re8 := regexp.MustCompile(`^(Топ)\s([т]?[3-9]|[т]?[1][0-2])$`) // запрос топа по уровню
	arr8 := (re8.FindAllStringSubmatch(in.Mtext, -1))
	if len(arr8) > 0 {
		in.SetLevelRsOrDrs(arr8[0][2])
		go b.Top(in)
		bb = true
		return bb
	}

	//eng^(Топ)\s([d]?[3-9]|[d]?[1][0-2])$
	re8e := regexp.MustCompile(`^(Top)\s(d?[3-9]|d?[1][0-2])$`) // запрос топа по уровню
	arr8e := (re8e.FindAllStringSubmatch(in.Mtext, -1))
	if len(arr8e) > 0 {
		in.SetLevelRsOrDrs(arr8[0][2])
		go b.Top(in)
		bb = true
		return bb
	}

	switch in.Mtext {
	case "Топ", "Top":
		bb = true
		go b.Top(in)
	case "Топ игры":
		bb = true
		go b.TopGame(in)
	}

	return bb
}

func (b *Bot) lEmoji(in models.InMessage) (bb bool) {
	var slot, emo string
	reEmodji := regexp.MustCompile("^(Эмоджи)\\s([1-4])\\s(<:\\w+:\\d+>)$") //добавления внутрених эмоджи
	arrEmodji := (reEmodji.FindAllStringSubmatch(in.Mtext, -1))
	if len(arrEmodji) > 0 {
		slot = arrEmodji[0][2]
		emo = arrEmodji[0][3]
	}
	reEmodji = regexp.MustCompile("^(Эмоджи)\\s([1-4])\\s(\\P{Greek})$") //добавления эмоджи
	arrEmodji = (reEmodji.FindAllStringSubmatch(in.Mtext, -1))
	if len(arrEmodji) > 0 {
		slot = arrEmodji[0][2]
		emo = arrEmodji[0][3]
	}
	reEmodji = regexp.MustCompile("^(Эмоджи)\\s([1-4])$") //удаление эмоджи с ячейки
	arrEmodji = (reEmodji.FindAllStringSubmatch(in.Mtext, -1))
	if len(arrEmodji) > 0 {
		slot = arrEmodji[0][2]
		emo = ""
	}
	reEmodji = regexp.MustCompile("^(Rs|rs)\\s(icon)\\s([1-4])\\s(del)$") //удаление эмоджи с ячейки совместимость
	arrEmodji = (reEmodji.FindAllStringSubmatch(in.Mtext, -1))
	if len(arrEmodji) > 0 {
		slot = arrEmodji[0][3]
		emo = ""
	}

	reEmodji = regexp.MustCompile("^(Rs|rs)\\s(icon)\\s([1-4])\\s(\\&\\#[0-9]+\\;)$") //Эмоджи совместимость
	arrEmodji = (reEmodji.FindAllStringSubmatch(in.Mtext, -1))
	if len(arrEmodji) > 0 {
		slot = arrEmodji[0][3]
		emo = arrEmodji[0][4]
	}
	if slot != "" {
		go b.emodjiadd(in, slot, emo)
		bb = true
	}
	if in.Mtext == "Эмоджи" || in.Mtext == "Emoji" {
		bb = true
		go b.emodjis(in)
	}
	//if in.Tip == "tg" {
	//	ok, n := b.instalNick(in, in.Mtext)
	//	if ok {
	//		go b.client.Tg.DelMessageSecond(in.Config.TgChannel, strconv.Itoa(in.Tg.Mesid), 20)
	//		bb = true
	//		if n == "удалено" {
	//			go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, fmt.Sprintf("Удалено дополнительное имя "), 20)
	//		} else {
	//			go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, fmt.Sprintf("Установлено дополнительное имя %s", n), 20)
	//		}
	//	}
	//}

	return bb
}

func (b *Bot) SendALLChannel(in models.InMessage) (bb bool) {
	text, found := strings.CutPrefix(in.Mtext, ".всем")
	if found && b.checkAdmin(in) {
		if in.Tip == ds {
			go b.client.Ds.DeleteMessage(in.Config.DsChannel, in.Ds.Mesid)
		} else if in.Tip == tg {
			go b.client.Tg.DelMessage(in.Config.TgChannel, in.Tg.Mesid)
		}

		b.ifTipSendTextDelSecond(in, "Начата рассылка.", 20)

		go func() {
			for _, config := range b.storage.ConfigRs.ReadConfigRs() {
				if config.DsChannel != "" {
					b.client.Ds.SendChannelDelSecond(config.DsChannel, text, 86400)
				}
				if config.TgChannel != "" {
					b.client.Tg.SendChannelDelSecond(config.TgChannel, text, 86400)
				}
			}
		}()

		bb = true
	}

	return bb
}
