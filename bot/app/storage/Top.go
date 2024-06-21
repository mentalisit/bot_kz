package storage

import (
	"kz_bot/models"
)

type Top interface {
	//TopLevel(ctx context.Context, CorpName, lvlkz string) bool
	//TopEventLevel(ctx context.Context, CorpName, lvlkz string, numEvent int) bool
	TopEventLevelNew(CorpName, lvlkz string, numEvent int) []models.Top
	//TopTemp(ctx context.Context) string
	//TopTempEvent(ctx context.Context) string
	//TopAll(ctx context.Context, CorpName string) bool
	//TopAllEvent(ctx context.Context, CorpName string, numberevent int) bool
	TopAllEventNew(CorpName string, numberevent int) (top []models.Top)
	//TopAllPerMonth(ctx context.Context, CorpName string) bool
	TopAllPerMonthNew(CorpName string) (top []models.Top)
	//TopLevelPerMonth(ctx context.Context, CorpName, lvlkz string) bool
	TopLevelPerMonthNew(CorpName, lvlkz string) []models.Top
}
