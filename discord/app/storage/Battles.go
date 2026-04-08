package storage

import "github.com/mentalisit/restapi/models"

type Battles interface {
	BattlesInsert(b models.Battles, timestamp string) error
	BattlesGetAll(corpName string, event int) ([]models.PlayerStats, error)
	BattlesTopInsert(b models.BattlesTop) error
	BattlesTopGetAll(corpName string) ([]models.BattlesTop, error)
	LoadNameAliases() (map[string]string, error)
}
