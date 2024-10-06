package storage

import (
	"github.com/mentalisit/logger"
	"telegram/config"
	"telegram/storage/postgres"
)

type Storage struct {
	Db *postgres.Db
}

func NewStorage(log *logger.Logger, cfg *config.ConfigBot) *Storage {
	local := postgres.NewDb(log, cfg)

	s := &Storage{
		Db: local,
	}

	//go s.loadDbArray()

	return s
}
