package bot

import (
	"fmt"
	"rs/models"
	"strconv"
	"strings"
	"time"
)

// lang ok
func (b *Bot) EventText(in models.InMessage) (text string, numE int) {
	//проверяем, есть ли активный ивент
	numberevent := b.storage.Event.NumActiveEvent(in.Config.CorpName)
	if numberevent == 0 { //ивент не активен
		return "", 0
	} else if numberevent > 0 { //активный ивент
		numE = b.storage.Event.NumberQueueEvents(in.Config.CorpName) //номер кз number FROM rsevent
		text = fmt.Sprintf("\nID %d %s\nㅤ\nㅤ", numE, b.getText(in, "for_event"))
		return text, numE
	}
	return text, numE
}
func (b *Bot) EventStart(in models.InMessage) {
	b.iftipdelete(in)
	//проверяем, есть ли активный ивент
	event1 := b.storage.Event.NumActiveEvent(in.Config.CorpName)
	text := b.getText(in, "info_event_started")
	if event1 > 0 {
		b.ifTipSendTextDelSecond(in, b.getText(in, "event_mode_enabled"), 10)
	} else {
		if in.Tip == ds && (in.Username == "Mentalisit" || b.client.Ds.CheckAdmin(in.UserId, in.Config.DsChannel)) {
			b.storage.Event.EventStartInsert(in.Config.CorpName)
			if in.Config.TgChannel != "" {
				b.client.Tg.SendChannel(in.Config.TgChannel, text)
				b.client.Ds.Send(in.Config.DsChannel, text)
			} else {
				b.client.Ds.Send(in.Config.DsChannel, text)
			}
		} else if in.Tip == tg {
			adminTg, err := b.client.Tg.CheckAdminTg(in.Config.TgChannel, in.Username)
			if err != nil {
				b.log.ErrorErr(err)
			}
			if adminTg || in.Username == "Mentalisit" {
				b.storage.Event.EventStartInsert(in.Config.CorpName)
				if in.Config.DsChannel != "" {
					b.client.Ds.Send(in.Config.DsChannel, text)
					b.client.Tg.SendChannel(in.Config.TgChannel, text)
				} else {
					b.client.Tg.SendChannel(in.Config.TgChannel, text)
				}
			}
		} else {
			text = b.getText(in, "info_event_starting")
			b.ifTipSendTextDelSecond(in, text, 60)
		}
	}
}
func (b *Bot) EventStop(in models.InMessage) {
	b.iftipdelete(in)
	event1 := b.storage.Event.NumActiveEvent(in.Config.CorpName)
	eventStop := b.getText(in, "event_stopped")
	eventNull := b.getText(in, "info_event_not_active")
	if in.Tip == "ds" && (in.Username == "Mentalisit" || b.client.Ds.CheckAdmin(in.UserId, in.Config.DsChannel)) {
		if event1 > 0 {
			b.storage.Event.UpdateActiveEvent0(in.Config.CorpName, event1)
			go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, eventStop, 60)
		} else {
			go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, eventNull, 10)
		}
	} else if in.Tip == tg {
		adminTg, err := b.client.Tg.CheckAdminTg(in.Config.TgChannel, in.Username)
		if err != nil {
			b.log.ErrorErr(err)
		}
		if in.Username == "Mentalisit" || adminTg {
			if event1 > 0 {
				b.storage.Event.UpdateActiveEvent0(in.Config.CorpName, event1)
				go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, eventStop, 60)
			} else {
				go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, eventNull, 10)
			}
		}
	} else {
		text := b.getText(in, "info_event_starting")
		b.ifTipSendTextDelSecond(in, text, 20)
	}
}
func (b *Bot) EventPoints(in models.InMessage, numKZ, points int) {
	b.iftipdelete(in)
	// проверяем активен ли ивент
	event1 := b.storage.Event.NumActiveEvent(in.Config.CorpName)
	message := ""
	if event1 > 0 {
		CountEventNames := b.storage.Event.CountEventNames(in.Config.CorpName, in.NameMention, numKZ, event1)
		if CountEventNames == 0 {
			CountEventNames = b.storage.Event.CountEventNames(in.Config.CorpName, "$"+in.NameMention, numKZ, event1)
		}
		admin := b.checkAdmin(in)
		if CountEventNames > 0 || admin {
			pointsGood := b.storage.Event.CountEventsPoints(in.Config.CorpName, numKZ, event1)
			if pointsGood > 0 && !admin {
				message = b.getText(in, "rs_data_entered")
			} else if pointsGood == 0 || admin {
				countEvent := b.storage.Event.UpdatePoints(in.Config.CorpName, numKZ, points, event1) //if error
				if countEvent == 0 {
					b.ifTipSendTextDelSecond(in, "error Count User = 0", 20)
				}
				message = fmt.Sprintf("%s %d %s", in.Username, points, b.getText(in, "points_added_to_database"))
				b.changeMessageEvent(in, points, countEvent, numKZ, event1)
			}
		} else {
			message = fmt.Sprintf("%s  %s %d", in.NameMention, b.getText(in, "info_points_cannot_be_added"), numKZ)
		}

	} else {
		message = b.getText(in, "event_not_started")
	}
	b.ifTipSendTextDelSecond(in, message, 20)
}
func (b *Bot) changeMessageEvent(in models.InMessage, points, countEvent, numberkz, numberEvent int) {
	nd, nt, t := b.storage.Event.ReadNamesMessage(in.Config.CorpName, numberkz, numberEvent)
	mes1 := fmt.Sprintf("🔴 %s №%d (%s)\n", b.getText(in, "event_game"), t.Numberkz, t.Lvlkz)
	mesOld := fmt.Sprintf("🎉 %s %s %d\nㅤ\nㅤ", b.getText(in, "contributed"), in.Username, points)
	if countEvent == 1 {
		if in.Config.DsChannel != "" {
			text := fmt.Sprintf("%s %s \n%s", mes1, nd.Name1, mesOld)
			b.client.Ds.EditWebhook(text, in.Username, in.Config.DsChannel, t.Dsmesid, in.Ds.Avatar)
		}
		if in.Config.TgChannel != "" {
			b.client.Tg.EditTextParse(in.Config.TgChannel, strconv.Itoa(t.Tgmesid), fmt.Sprintf("%s %s \n%s", mes1, nt.Name1, mesOld), "")
		}
	} else if countEvent == 2 {
		if in.Config.DsChannel != "" {
			text := fmt.Sprintf("%s %s\n %s\n %s", mes1, nd.Name1, nd.Name2, mesOld)
			b.client.Ds.EditWebhook(text, in.Username, in.Config.DsChannel, t.Dsmesid, in.Ds.Avatar)
		}
		if in.Config.TgChannel != "" {
			text := fmt.Sprintf("%s %s\n %s\n %s", mes1, nt.Name1, nt.Name2, mesOld)
			b.client.Tg.EditTextParse(in.Config.TgChannel, strconv.Itoa(t.Tgmesid), text, "")
		}
	} else if countEvent == 3 {
		if in.Config.DsChannel != "" {
			text := fmt.Sprintf("%s %s\n %s\n %s\n %s", mes1, nd.Name1, nd.Name2, nd.Name3, mesOld)
			b.client.Ds.EditWebhook(text, in.Username, in.Config.DsChannel, t.Dsmesid, in.Ds.Avatar)
		}
		if in.Config.TgChannel != "" {
			text := fmt.Sprintf("%s %s\n %s\n %s\n %s", mes1, nt.Name1, nt.Name2, nt.Name3, mesOld)
			b.client.Tg.EditTextParse(in.Config.TgChannel, strconv.Itoa(t.Tgmesid), text, "")
		}
	} else if countEvent == 4 {
		if in.Config.DsChannel != "" {
			text := fmt.Sprintf("%s %s\n %s\n %s\n %s\n %s", mes1, nd.Name1, nd.Name2, nd.Name3, nd.Name4, mesOld)
			b.client.Ds.EditWebhook(text, in.Username, in.Config.DsChannel, t.Dsmesid, in.Ds.Avatar)
		}
		if in.Config.TgChannel != "" {
			text := fmt.Sprintf("%s %s\n %s\n %s\n %s\n %s", mes1, nt.Name1, nt.Name2, nt.Name3, nt.Name4, mesOld)
			b.client.Tg.EditTextParse(in.Config.TgChannel, strconv.Itoa(t.Tgmesid), text, "")
		}
	}
}

// new

func (b *Bot) EventPreStart(in models.InMessage) {
	b.storage.Event.EventInsertPreStart(in.Config.CorpName, -1)
}

func (b *Bot) EventAutoStart() {
	date := time.Now().UTC().Format(time.DateOnly)
	nextDateEventStart, nextDateEventStop := b.storage.Event.ReadEventSchedule()

	send := func(config models.CorporationConfig, text string) {
		if config.TgChannel != "" {
			b.client.Tg.SendChannel(config.TgChannel, text)
		}
		if config.DsChannel != "" {
			b.client.Ds.Send(config.DsChannel, text)
		}
	}
	sendHelp := func(config models.CorporationConfig) {
		if config.DsChannel != "" {
			config = b.sendHelpDs(config, true)
		}
		if config.TgChannel != "" {
			config = b.sendHelpTg(config, true)
		}
		b.storage.ConfigRs.UpdateConfigRs(config)
	}

	if date == nextDateEventStart {
		b.log.Info("Event Starting all")
		corps := b.storage.Event.ReadRsEvent(-1)
		for _, event := range corps {
			ok, config := b.CheckCorpNameConfig(event.CorpName)
			if ok {
				text := b.getLanguageText(config.Country, "info_event_started")
				//Устанавливаем 1 как активный ивент
				b.storage.Event.UpdateActiveEvent(1, event.CorpName, event.NumEvent)

				send(config, text)
				sendHelp(config)

				time.Sleep(1 * time.Second)
			}
		}
	}
	if date == nextDateEventStop {
		b.log.Info("Event Stopping all")
		top := func(conf models.CorporationConfig) {
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

			numEvent := b.storage.Event.NumActiveEvent(conf.CorpName)

			message = fmt.Sprintf("\xF0\x9F\x93\x96 %s %s:\n",
				b.getLanguageText(conf.Country, "top_participants"),
				b.getLanguageText(conf.Country, "event"),
			)
			resultsTop = b.storage.Top.TopAllEventNew(conf.CorpName, numEvent)
			resultsTop = mergeAndSumTops(resultsTop)

			if len(resultsTop) > 0 {
				for _, top := range resultsTop {
					message2 = message2 + format(top)
					number++
				}
			} else {
				return
			}
			if allpoints != 0 {
				message2 = fmt.Sprintf("%s\nTotal: %d", message2, allpoints)

			}

			if conf.DsChannel != "" {
				b.client.Ds.SendEmbedText(conf.DsChannel, message, message2)
			}
			if conf.TgChannel != "" {
				text := message + message2
				if conf.Guildid != "" {
					text = b.client.Ds.ReplaceTextMessage(text, conf.Guildid)
				}
				text = strings.ReplaceAll(text, "@", "")
				b.client.Tg.SendChannel(conf.TgChannel, text)
			}
		}

		corps := b.storage.Event.ReadRsEvent(1)
		for _, event := range corps {
			ok, config := b.CheckCorpNameConfig(event.CorpName)
			if ok {
				//отправляем топ если он не пуст
				top(config)

				//Устанавливаем 0 как закрытый ивент
				b.storage.Event.UpdateActiveEvent(0, event.CorpName, event.NumEvent)

				text := b.getLanguageText(config.Country, "event_stopped")
				send(config, text)
				sendHelp(config)
			}
		}
		b.OptimizationSborkz()
	}
}
