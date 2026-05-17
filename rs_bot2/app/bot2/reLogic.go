package bot2

import (
	"fmt"
	"regexp"
	"rs/models"
	"strconv"
)

func (b *Bot) lRsPlus(in *models.InMessageV2) (rsb bool) {
	var variableFound string
	reLevelTime := regexp.MustCompile(`^([3-9]|1[0-2])([+]|-|!|[*])([$]|[?])?(\d|\d{2}|\d{3})?$`)
	arr := reLevelTime.FindAllStringSubmatch(in.Text, -1)
	if len(arr) > 0 {
		rsb = true
		rs := models.NewRs()
		rs.TimeRs = arr[0][4]
		rs.SetLevelRsOrDrs(in, arr[0][1])
		variableFound = arr[0][2]

		if arr[0][3] == "$" {
			rs.Money = true
		} else if arr[0][3] == "?" {
			altNum, _ := strconv.Atoi(arr[0][4])
			if altNum == 0 || altNum > 5 {
				altNum = 1
			}
			rs.TimeRs = "30"

			if in.MAcc != nil {
				lenAlts := len(in.MAcc.Alts)
				if lenAlts > 0 {
					if lenAlts == 1 || altNum == 1 {
						rs.AltName = in.MAcc.Alts[0]
					} else if altNum > 1 {
						altNum = altNum - 1
						if lenAlts > altNum {
							rs.AltName = in.MAcc.Alts[altNum]
						}
					}
				}
			}

		}
		if rs.TimeRs == "" {
			rs.TimeRs = "30"
		}

		switch variableFound {
		case "+", "*":
			b.RsDarkPlus(in, rs)
		case "-":
			b.RsMinus(in, rs)
		case "!":
			{
				fmt.Println("using ! logicRs56")
				_, level := rs.TypeRedStar()
				rs.RsTypeLevel = "rs" + level
				b.RsDarkPlus(in, rs)
			}

		default:
			rsb = false
		}
	}

	return rsb
}

func (b *Bot) lQueue(in *models.InMessageV2) bool {
	re4 := regexp.MustCompile(`^([о]|[О]|[q]|[Q]|[Ч]|[ч])([3-9]|[1][0-2])$`) // две переменные для чтения  очереди
	arr4 := re4.FindAllStringSubmatch(in.Text, -1)
	rs := models.NewRs()
	if len(arr4) > 0 {
		rs.SetLevelRsOrDrs(in, arr4[0][2])
		b.QueueLevel(in, rs)
		return true
	}
	//rus
	if in.Text == "Очередь" || in.Text == "очередь" || in.Text == "Черга" || in.Text == "черга" || in.Text == "Queue" || in.Text == "queue" {
		b.QueueAll(in)
		return true
	}

	//todo придумать другую команду
	if in.Text == "Все очереди" {
		b.sendTextAfterDeleteSecond(in, "Поиск.... ", 10)
		b.deleteInMessage(in)
		go b.sendTextAfterDeleteSecond(in, b.otherQueue.MyQueue(), 30)
		return true
	}

	re4s := regexp.MustCompile(`^(Rs|rs)\s(Q|q)$`) // две переменные для чтения  очереди
	arr4s := re4s.FindAllStringSubmatch(in.Text, -1)
	if len(arr4s) > 0 {
		go b.QueueAll(in) //проверка совместимости
		return true
	}

	if in.Text == "+" && in.Options.Contains(models.OptionReaction) {
		b.Plus(in)
		return true
	} else if in.Text == "-" && in.Options.Contains(models.OptionReaction) {
		b.Minus(in)
		return true
	}

	return false
}

func (b *Bot) lSubs(in *models.InMessageV2) (bb bool) {
	bb = false
	var subs string
	rs := models.NewRs()
	re := regexp.MustCompile(`^([+]|[-])([3-9]|[1][0-2])$`) // две переменные для добавления или удаления подписок
	arr := re.FindAllStringSubmatch(in.Text, -1)
	if len(arr) > 0 {
		rs.SetLevelRsOrDrs(in, arr[0][2])
		subs = arr[0][1]
		bb = true
	}
	re1 := regexp.MustCompile(`^(Rs|rs)\s(S|s|u|U)\s([3-9]|[1][0-2])$`)
	arr1 := re1.FindAllStringSubmatch(in.Text, -1)
	if len(arr1) > 0 {
		rs.SetLevelRsOrDrs(in, arr1[0][3])
		subs = arr1[0][2]
		bb = true
		if subs == "S" || subs == "s" {
			subs = "+"
		} else if subs == "U" || subs == "u" {
			subs = "-"
		}
	}

	switch subs {
	case "+":
		go b.handleSubscription(in, rs, true)
	case "-":
		go b.handleSubscription(in, rs, false)
	}
	return bb
}

func (b *Bot) lRsStart(in *models.InMessageV2) (bb bool) {
	var rss string
	rs := models.NewRs()
	re5 := regexp.MustCompile(`^([3-9]|[1][0-2])([\+][\+])$`) //rs start
	arr5 := re5.FindAllStringSubmatch(in.Text, -1)
	if len(arr5) > 0 {
		bb = true
		rs.SetLevelRsOrDrs(in, arr5[0][1])
		rss = arr5[0][2]
	} else {
		re5 = regexp.MustCompile(`^(Rs|rs)\s(Start|start)\s([3-9]|[1][0-2])$`) //rs start
		arr5 = re5.FindAllStringSubmatch(in.Text, -1)
		if len(arr5) > 0 {
			bb = true
			rs.SetLevelRsOrDrs(in, arr5[0][3])
			rss = "++"
			//b.log.Println("Проверка совместимости принудительного старта ")
		}
	}
	reP := regexp.MustCompile(`^([3-9]|[1][0-2])([\+][\+][\+])$`) //p30pl
	arrP := reP.FindAllStringSubmatch(in.Text, -1)
	if len(arrP) > 0 {
		rs.SetLevelRsOrDrs(in, arrP[0][1])
		go b.Pl30(in, rs)
		return true
	}

	if rss == "++" {
		go b.RsStart(in, rs)
	}
	return bb
}

func (b *Bot) lTop(in *models.InMessageV2) (bb bool) {
	rs := models.NewRs()
	re8 := regexp.MustCompile(`^(Топ)\s([т]?[3-9]|[т]?[1][0-2])$`) // запрос топа по уровню
	arr8 := re8.FindAllStringSubmatch(in.Text, -1)
	if len(arr8) > 0 {
		rs.SetLevelRsOrDrs(in, arr8[0][2])
		go b.Top(in, rs)
		return true
	}

	//eng^(Топ)\s([d]?[3-9]|[d]?[1][0-2])$
	re8e := regexp.MustCompile(`^(Top)\s(d?[3-9]|d?[1][0-2])$`) // запрос топа по уровню
	arr8e := re8e.FindAllStringSubmatch(in.Text, -1)
	if len(arr8e) > 0 {
		rs.SetLevelRsOrDrs(in, arr8[0][2])
		go b.Top(in, rs)
		return true
	}

	switch in.Text {
	case "Топ", "Top":
		bb = true
		go b.Top(in, rs)
	case "Топ игры":
		bb = true
		go b.TopGame(in, rs)
	}

	return bb
}

func (b *Bot) lEmoji(in *models.InMessageV2) (bb bool) {
	var slot, emo string
	reEmoji := regexp.MustCompile("^(Эмоджи)\\s([1-4])\\s(<:\\w+:\\d+>)$") //добавления внутренних эмоджи
	arrEmoji := reEmoji.FindAllStringSubmatch(in.Text, -1)
	if len(arrEmoji) > 0 {
		slot = arrEmoji[0][2]
		emo = arrEmoji[0][3]
	}
	reEmoji = regexp.MustCompile("^(Эмоджи)\\s([1-4])\\s(\\P{Greek})$") //добавления эмоджи
	arrEmoji = reEmoji.FindAllStringSubmatch(in.Text, -1)
	if len(arrEmoji) > 0 {
		slot = arrEmoji[0][2]
		emo = arrEmoji[0][3]
	}
	reEmoji = regexp.MustCompile("^(Эмоджи)\\s([1-4])$") //удаление эмоджи с ячейки
	arrEmoji = reEmoji.FindAllStringSubmatch(in.Text, -1)
	if len(arrEmoji) > 0 {
		slot = arrEmoji[0][2]
		emo = ""
	}
	reEmoji = regexp.MustCompile("^(Rs|rs)\\s(icon)\\s([1-4])\\s(del)$") //удаление эмоджи с ячейки совместимость
	arrEmoji = reEmoji.FindAllStringSubmatch(in.Text, -1)
	if len(arrEmoji) > 0 {
		slot = arrEmoji[0][3]
		emo = ""
	}

	reEmoji = regexp.MustCompile("^(Rs|rs)\\s(icon)\\s([1-4])\\s(\\&\\#[0-9]+\\;)$") //Эмоджи совместимость
	arrEmoji = reEmoji.FindAllStringSubmatch(in.Text, -1)
	if len(arrEmoji) > 0 {
		slot = arrEmoji[0][3]
		emo = arrEmoji[0][4]
	}
	if slot != "" {
		go b.emojiAdd(in, slot, emo)
		bb = true
	}
	if in.Text == "Эмоджи" || in.Text == "Emoji" {
		bb = true
		go b.emojis(in)
	}

	return bb
}

func (b *Bot) lHelp(in *models.InMessageV2) bool {
	// Регулярное выражение для слов похожих на "Справка", "Help", "help" (русский, английский, украинский)
	reHelp := regexp.MustCompile(`(?i)^(справка|довідка|help|команди|commands|помощь|допомога|assist|info|інфо)$`)
	found := reHelp.MatchString(in.Text)
	if found {
		b.deleteInMessage(in)
		conf := b.SendHelpInMessenger(in)
		b.storage.UpdateConfigV2HelpMessage(conf)
	}
	return found
}
