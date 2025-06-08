package storage

import "discord/models"

type Battles interface {
	BattlesInsert(b models.Battles) error
	BattlesGetAll(corpName string, event int) ([]models.PlayerStats, error)
	BattlesTopInsert(b models.BattlesTop) error
	BattlesTopGetAll(corpName string) ([]models.BattlesTop, error)
	LoadNameAliases() (map[string]string, error)
}
