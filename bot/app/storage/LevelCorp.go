package storage

import "kz_bot/models"

type LevelCorp interface {
	InsertUpdateCorpLevel(l models.LevelCorps)
	ReadCorpLevel(CorpName string) (models.LevelCorps, error)
	ReadCorpLevelAll() ([]models.LevelCorps, error)
}
