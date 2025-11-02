package bot

import (
	"fmt"
	"time"
)

func (b *Bot) ReadAndSendPic(tn time.Time) {
	nextDateStart, nextDateStop, message := b.storage.Event.ReadEventScheduleAndMessage()
	date1 := tn.Format("02-01-2006")
	date2 := tn.Add(24 * time.Hour).Format("02-01-2006")
	if date1 == nextDateStart || date2 == nextDateStop {
		seconds := 600

		if tn.Hour() < 9 {
			if tn.Minute() == 59 || tn.Minute() == 29 {
				seconds = 1800
			} else {
				return
			}
		}

		eventId := getSeasonNumber(message)
		title := fmt.Sprintf("Сезон №%d %s - %s", eventId, nextDateStart, nextDateStop)

		paramsReadAll := b.storage.Scoreboard.ScoreboardReadAll()

		for _, wh := range paramsReadAll {

			if wh.ChannelScoreboardOrMap == "" {
				continue
			}

			m, str := wh.GetMapOrString()

			filename := fmt.Sprintf("Scoreboard_for_%s.png", wh.ChannelWebhook)

			scoreboard := b.helpers.CreateScoreboard(filename, wh.Name, eventId)
			if scoreboard != "" {
				if str == "" {
					for s, channel := range m {
						if s == "ds" {
							err := b.client.Ds.SendOrEditEmbedImageScoreboard(channel, title, filename)
							if err != nil {
								b.log.ErrorErr(err)
							}
						} else if s == "tg" {
							mid, err := b.client.Tg.SendPicScoreboard(channel, title, filename)
							if err != nil {
								b.log.ErrorErr(err)
							}
							if tn.Weekday() == time.Sunday && tn.Hour() == 23 && tn.Minute() == 59 {
								b.log.Info("final scoreboard event")
							} else {
								b.client.Tg.DelMessageSecond(channel, mid, seconds)
							}
						} else if s == "wa" {
							if tn.Minute() == 59 && (tn.Hour() == 5 || tn.Hour() == 11 || tn.Hour() == 17 || tn.Hour() == 23) {
								_, err := b.client.Wa.SendPicScoreboard(channel, title, filename)
								if err != nil {
									b.log.ErrorErr(err)
								}
							}
						}
					}
				} else {
					err := b.client.Ds.SendOrEditEmbedImageScoreboard(wh.ChannelScoreboardOrMap, title, filename)
					if err != nil {
						b.log.ErrorErr(err)
					}
				}

				time.Sleep(1 * time.Second)
			}
		}
	}
}
