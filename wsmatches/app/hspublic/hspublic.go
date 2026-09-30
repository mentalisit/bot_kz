package hspublic

import (
	"ws/dbpostgres"

	"github.com/jmoiron/sqlx"
	"github.com/mentalisit/conf/logger"
)

type HS struct {
	log *logger.Logger
	Db  *dbpostgres.Db
}

func NewHS(log *logger.Logger, db *sqlx.DB) *HS {

	st := dbpostgres.NewDb(log, db)

	return &HS{
		log: log,
		Db:  st,
	}
}
