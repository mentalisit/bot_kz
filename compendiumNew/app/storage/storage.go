package storage

import (
	"compendium/storage/postgres"
	postgresv2 "compendium/storage/postgres/postgresV2"

	"github.com/jmoiron/sqlx"
	"github.com/mentalisit/conf/logger"
	"go.uber.org/zap"
)

type Storage struct {
	log   *zap.Logger
	debug bool
	DB    *postgres.Db
	V2    *postgresv2.Db
}

func NewStorage(log *logger.Logger, db *sqlx.DB) *Storage {
	s := &Storage{
		DB: postgres.NewDb(log, db),
		V2: postgresv2.NewDb(log, db),
	}

	//go s.loadDbArray()

	return s
}
