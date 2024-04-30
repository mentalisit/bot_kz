package storage

import (
	"compendium/config"
	"compendium/storage/postgres"
	"github.com/mentalisit/logger"
	"go.uber.org/zap"
)

type Storage struct {
	log   *zap.Logger
	debug bool
	Temp  *postgres.Db
}

func NewStorage(log *logger.Logger, cfg *config.ConfigBot) *Storage {
	local := postgres.NewDb(log, cfg)

	s := &Storage{

		Temp: local,
	}

	//go s.loadDbArray()

	return s
}
