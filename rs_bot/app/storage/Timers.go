package storage

import (
	"rs/models"
)

type Timers interface {
	UpdateMitutsQueue(userid, CorpName string) models.Sborkz
	MinusMin() []models.Sborkz
}
