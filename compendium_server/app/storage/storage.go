package storage

import (
	"compendium_s/config"
	"compendium_s/storage/postgres"
	"github.com/mentalisit/logger"
	"go.uber.org/zap"
)

type Storage struct {
	log   *zap.Logger
	debug bool
	DB    *postgres.Db
}

func NewStorage(log *logger.Logger, cfg *config.ConfigBot) *Storage {
	local := postgres.NewDb(log, cfg)

	s := &Storage{

		DB: local,
	}

	//go s.loadDbArray()

	return s
}
