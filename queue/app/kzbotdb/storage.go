package kzbotdb

import (
	"github.com/jmoiron/sqlx"
	"github.com/mentalisit/conf/logger"

	_ "github.com/lib/pq"
)

type Db struct {
	db  *sqlx.DB
	log *logger.Logger
}

func NewDb(log *logger.Logger, db *sqlx.DB) *Db {
	d := &Db{
		db:  db,
		log: log,
	}
	return d
}
