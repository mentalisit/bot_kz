package storage

import (
	"compendium/config"
	"compendium/storage/postgres"
	"compendium/storage/postgres/multi"
	postgresv2 "compendium/storage/postgres/postgresV2"

	"github.com/mentalisit/logger"
	"go.uber.org/zap"
)

type Storage struct {
	log   *zap.Logger
	debug bool
	DB    *postgres.Db
	Multi *multi.Db
	DBv2  *postgresv2.Db
}

func NewStorage(log *logger.Logger, cfg *config.ConfigBot) *Storage {
	local := postgres.NewDb(log, cfg)

	s := &Storage{
		DB:    local,
		Multi: local.Multi,
		DBv2:  postgresv2.NewDb(log, cfg),
	}

	//go s.loadDbArray()

	return s
}
