package storage

import (
	"bridge/storage/postgres"

	"github.com/jmoiron/sqlx"
	"github.com/mentalisit/conf/logger"
	"go.uber.org/zap"
)

type Storage struct {
	log   *zap.Logger
	debug bool
	DB    *postgres.Db
}

func NewStorage(log *logger.Logger, db *sqlx.DB) *Storage {
	local := postgres.NewDb(log, db)

	s := &Storage{
		DB: local,
	}

	return s
}
