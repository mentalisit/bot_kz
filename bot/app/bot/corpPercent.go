package bot

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"kz_bot/models"
	"time"
)

func (b *Bot) SendPercent(Config models.CorporationConfig) {
	currentCorp, err := b.storage.LevelCorp.ReadCorpLevel(Config.CorpName)
	if err != nil {
		b.log.ErrorErr(err)
	}
	untilTime := currentCorp.EndDate.AddDate(0, 0, 7).Unix()
	if time.Now().UTC().Unix() < untilTime {
		return
	}

	all, err := b.storage.LevelCorp.ReadCorpLevelAll()
	if err != nil {
		b.log.ErrorErr(err)
		return
	}
	var preperingText []string
	for _, corp := range all {
		if corp.HCorp != "" && corp.Percent != 0 {
			untilTime = corp.EndDate.AddDate(0, 0, 7).Unix()
			if time.Now().UTC().Unix() < untilTime {
				preperingText = append(preperingText,
					fmt.Sprintf("%d%% %s %+v\n", percent(corp.Level), corp.HCorp, formatTime(untilTime)))
			}
		}
	}
	sortText := sortByFirstTwoDigits(preperingText)

	text := ""

	for _, s := range sortText {
		text += s
	}

	if Config.DsChannel != "" {
		go b.client.Ds.SendChannelDelSecond(Config.DsChannel, text, 180)
	}
	if Config.TgChannel != "" {
		go b.client.Tg.SendChannelDelSecond(Config.TgChannel, text, 180)
	}
}

func (b *Bot) GetTextPercent(Config models.CorporationConfig, dark bool) string {
	currentCorp, err := b.storage.LevelCorp.ReadCorpLevel(Config.CorpName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			b.storage.LevelCorp.InsertUpdateCorpLevel(models.LevelCorps{
				CorpName:   Config.CorpName,
				Level:      0,
				EndDate:    time.Time{},
				HCorp:      "",
				Percent:    0,
				LastUpdate: time.Time{},
				Relic:      0,
			})
			return ""
		} else {
			b.log.ErrorErr(err)
			return ""
		}
	}
	untilTime := currentCorp.EndDate.AddDate(0, 0, 7).Unix()
	if time.Now().UTC().Unix() < untilTime {
		per := percent(currentCorp.Level)
		if dark {
			return fmt.Sprintf(" %d%%", per+100)
		} else {
			return fmt.Sprintf(" %d%%", per)
		}
	}
	return ""
}
