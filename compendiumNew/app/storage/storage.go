package storage

import (
	"compendium/config"
	"compendium/storage/postgres"
	"compendium/storage/postgres/multi"
	"github.com/mentalisit/logger"
	"go.uber.org/zap"
)

type Storage struct {
	log   *zap.Logger
	debug bool
	DB    *postgres.Db
	Multi *multi.Db
}

func NewStorage(log *logger.Logger, cfg *config.ConfigBot) *Storage {
	local := postgres.NewDb(log, cfg)

	s := &Storage{
		DB:    local,
		Multi: local.Multi,
	}

	//go s.loadDbArray()

	return s
}
