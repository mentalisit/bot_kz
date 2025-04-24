package storage

import "rs/models"

type Battles interface {
	BattlesInsert(b models.Battles) error
	BattlesGetAll(corpName string, event int) ([]models.PlayerStats, error)
	BattlesTopGetAll(corpName string) ([]models.BattlesTop, error)
	ReadEventScheduleAndMessage() (nextDateStart, nextDateStop, message string)
}
