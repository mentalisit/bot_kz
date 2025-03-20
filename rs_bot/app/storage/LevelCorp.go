package storage

import "rs/models"

type LevelCorp interface {
	InsertUpdateCorpLevel(l models.LevelCorps)
	ReadCorpLevelByCorpConf(CorpName string) (models.LevelCorps, error)
	ReadCorpLevelAll() ([]models.LevelCorps, error)

	ReadCorpsLevelAllOld() ([]models.LevelCorps, error)
}
