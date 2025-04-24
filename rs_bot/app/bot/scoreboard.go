package bot

import (
	"fmt"
	"time"
)

func (b *Bot) ReadAndSendPic() {
	nextDateStart, nextDateStop, message := b.storage.Event.ReadEventScheduleAndMessage()
	date1 := time.Now().UTC().Format("02-01-2006")
	date2 := time.Now().UTC().Add(24 * time.Hour).Format("02-01-2006")
	eventId := 0
	title := ""

	paramsReadAll := b.storage.Scoreboard.ScoreboardReadAll()

	//for _, params := range paramsReadAll {
	//	if params.Name == "IX_Легион" {
	//		_, str := params.GetMapOrString()
	//		if str != "" {
	//			mc := make(map[string]string)
	//			mc["ds"] = str
	//			mc["tg"] = "-1002298028181/71634"
	//			marshal, _ := json.Marshal(mc)
	//			params.ChannelScoreboardOrMap = string(marshal)
	//			b.storage.Battles.ScoreboardUpdateParam(params)
	//		}
	//	}
	//}
	seconds := 1790

	for _, wh := range paramsReadAll {

		if wh.ChannelScoreboardOrMap == "" {
			continue
		}

		if date1 == nextDateStart || date2 == nextDateStop {
			eventId = getSeasonNumber(message)
			title = fmt.Sprintf("Сезон №%d %s - %s", eventId, nextDateStart, nextDateStop)
		} else {
			title = ""
		}
		if eventId == 0 && time.Now().UTC().Hour()%6 != 0 {
			break
		} else if eventId == 0 && time.Now().UTC().Hour()%6 == 0 {
			seconds = 1800 * 2 * 6
		}

		m, str := wh.GetMapOrString()

		filename := fmt.Sprintf("Scoreboard_for_%s.png", str)
		if m != nil {
			filename = fmt.Sprintf("Scoreboard_for_%s.png", m["ds"])
		}

		scoreboard := b.helpers.CreateScoreboard(filename, wh.Name, eventId)
		if scoreboard != "" {
			if title == "" {
				title = wh.Name
			}
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
						if eventId == 0 || time.Now().Minute() == 29 {
							b.client.Tg.DelMessageSecond(channel, mid, seconds)
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
