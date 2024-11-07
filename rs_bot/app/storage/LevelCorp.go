package storage

import "rs/models"

type LevelCorp interface {
	InsertUpdateCorpLevel(l models.LevelCorps)
	ReadCorpLevel(CorpName string) (models.LevelCorps, error)
	ReadCorpLevelAll() ([]models.LevelCorps, error)
}
