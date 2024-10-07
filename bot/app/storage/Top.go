package storage

import (
	"kz_bot/models"
)

type Top interface {
	TopEventLevelNew(CorpName, lvlkz string, numEvent int) []models.Top
	TopAllEventNew(CorpName string, numberevent int) (top []models.Top)
	TopAllPerMonthNew(CorpName string) (top []models.Top)
	TopLevelPerMonthNew(CorpName, lvlkz string) []models.Top
}
