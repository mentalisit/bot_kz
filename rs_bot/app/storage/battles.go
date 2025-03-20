package storage

import "rs/models"

type Battles interface {
	BattlesInsert(b models.Battles) error
	BattlesGetAll(corpName string, event int) ([]models.PlayerStats, error)
	ScoreboardParamsReadAll() []models.ScoreboardParams
	BattlesTopGetAll(corpName string) ([]models.BattlesTop, error)
	IdentifyGetPoints() (ss []models.Identify, err error)
}
