package storage

import (
	"telegram/config"
	"telegram/storage/postgres"

	"github.com/mentalisit/logger"
)

type Storage struct {
	Db   *postgres.Db
	Conf *config.ConfigBot
}

func NewStorage(log *logger.Logger, cfg *config.ConfigBot) *Storage {
	local := postgres.NewDb(log, cfg)

	s := &Storage{
		Db:   local,
		Conf: cfg,
	}

	//go s.loadDbArray()

	return s
}
