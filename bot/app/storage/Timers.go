package storage

import (
	"context"
	"kz_bot/models"
)

type Timers interface {
	UpdateMitutsQueue(ctx context.Context, userid, CorpName string) models.Sborkz
	MinusMin(ctx context.Context) []models.Sborkz
}
type TimeDeleteMessage interface {
	TimerDeleteMessage() []models.Timer
	TimerInsert(c models.Timer)
}
