package storage

import (
	"rs/models"
)

type Timers interface {
	UpdateMitutsQueue(userid, CorpName string) models.Sborkz
	MinusMin() []models.Sborkz
}
type TimeDeleteMessage interface {
	TimerDeleteMessage() []models.Timer
	TimerInsert(c models.Timer)
}
