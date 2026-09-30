package storage

import (
	"compendium_s/storage/postgres"
	postgresv2 "compendium_s/storage/postgres/postgresV2"

	"github.com/jmoiron/sqlx"
	"github.com/mentalisit/conf/logger"
	"go.uber.org/zap"
)

type Storage struct {
	log   *zap.Logger
	debug bool
	DB    *postgres.Db
	DBv2  *postgresv2.Db
}

func NewStorage(log *logger.Logger, db *sqlx.DB) *Storage {
	s := &Storage{
		DB:   postgres.NewDb(log, db),
		DBv2: postgresv2.NewDb(log, db),
	}

	return s
}
