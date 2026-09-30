package postgres

import (
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/mentalisit/conf/logger"
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

func (d *Db) Shutdown() {
	d.db.Close()
}
