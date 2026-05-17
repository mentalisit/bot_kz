package bot2

import (
	"fmt"
	"time"
)

func (b *Bot) ReadAndSendPic() {
	tn := time.Now().UTC()
	nextDateStart, nextDateStop, message := b.storage.ReadEventScheduleAndMessage()
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

		paramsReadAll := b.storage.ScoreboardReadAll()

		for _, wh := range paramsReadAll {

			if len(wh.Channels) == 0 {
				continue
			}

			filename := fmt.Sprintf("Scoreboard_for_%s.png", wh.Uid)

			scoreboard := b.helpers.CreateScoreboard(filename, wh.Uid, eventId)
			if scoreboard != "" {
				for _, ch := range wh.Channels {
					if ch.TypeMessenger == "ds" {
						err := b.client.Ds.SendOrEditEmbedImageScoreboard(ch.ChannelId, title, filename)
						if err != nil {
							b.log.ErrorErr(err)
						}
					} else if ch.TypeMessenger == "tg" {
						mid, err := b.client.Tg.SendPicScoreboard(ch.ChannelId, title, filename)
						if err != nil {
							b.log.ErrorErr(err)
						}
						if tn.Weekday() == time.Sunday && tn.Hour() == 23 && tn.Minute() == 59 {
							b.log.Info("final scoreboard event")
						} else {
							b.client.Tg.DelMessageSecond(ch.ChannelId, mid, seconds)
						}
					} else if ch.TypeMessenger == "wa" {
						if tn.Minute() == 59 && (tn.Hour() == 5 || tn.Hour() == 11 || tn.Hour() == 17 || tn.Hour() == 23) {
							_, err := b.client.Wa.SendPicScoreboard(ch.ChannelId, title, filename)
							if err != nil {
								b.log.ErrorErr(err)
							}
						}
					}
				}

				time.Sleep(1 * time.Second)
			}
		}
	}
}
