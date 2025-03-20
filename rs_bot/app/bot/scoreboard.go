package bot

import (
	"fmt"
	"time"
)

func (b *Bot) ReadAndSendPic() {
	nextDateStart, nextDateStop, message := b.storage.Event.ReadEventScheduleAndMessage()
	date1 := time.Now().UTC().Format("02-01-2006")
	date2 := time.Now().UTC().Add(24 * time.Hour).Format("02-01-2006")
	eventId := time.Now().UTC().Day() - 10
	//todo need event read
	currentDay := time.Now().UTC().Format("02-01-2006")
	title := fmt.Sprintf("Сезон №%d  %s - %s", eventId, currentDay, currentDay)
	if date1 == nextDateStart || date2 == nextDateStop {
		eventId = getSeasonNumber(message)
		title = fmt.Sprintf("Сезон №%d %s - %s", eventId, nextDateStart, nextDateStop)
	}

	paramsReadAll := b.storage.Battles.ScoreboardParamsReadAll()

	for _, wh := range paramsReadAll {
		filename := fmt.Sprintf("Scoreboard_for_%s.png", wh.ChannelScoreboard)

		scoreboard := b.helpers.CreateScoreboard(filename, wh.Name, eventId)
		if scoreboard != "" {
			err := b.client.Ds.SendOrEditEmbedImageScoreboard(wh.ChannelScoreboard, title, filename)
			if err != nil {
				b.log.ErrorErr(err)
			}
		}
	}
}
