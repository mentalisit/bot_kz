package storage

import "kz_bot/models"

type LevelCorp interface {
	InsertUpdateCorpLevel(l models.LevelCorp)
	ReadCorpLevel(CorpName string) (models.LevelCorp, error)
}
