package storage

import (
	"telegram/storage/postgres"

	"github.com/jmoiron/sqlx"
	"github.com/mentalisit/conf/logger"
)

type Storage struct {
	Db *postgres.Db
}

func NewStorage(log *logger.Logger, db *sqlx.DB) *Storage {
	local := postgres.NewDb(log, db)

	s := &Storage{
		Db: local,
	}

	//go s.loadDbArray()

	return s
}
