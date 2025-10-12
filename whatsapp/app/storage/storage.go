package storage

import (
	"whatsapp/config"
	"whatsapp/storage/postgres"

	"github.com/mentalisit/logger"
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
