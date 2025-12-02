package storage

import (
	"compendium_s/config"
	"compendium_s/storage/postgres"
	postgresv2 "compendium_s/storage/postgres/postgresV2"

	"github.com/mentalisit/logger"
	"go.uber.org/zap"
)

type Storage struct {
	log   *zap.Logger
	debug bool
	DB    *postgres.Db
	DBv2  *postgresv2.Db
}

func NewStorage(log *logger.Logger, cfg *config.ConfigBot) *Storage {
	s := &Storage{
		DB:   postgres.NewDb(log, cfg),
		DBv2: postgresv2.NewDb(log, cfg),
	}

	return s
}
