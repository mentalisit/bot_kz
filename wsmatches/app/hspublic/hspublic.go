package hspublic

import (
	"ws/dbpostgres"

	"github.com/mentalisit/logger"
)

type HS struct {
	log *logger.Logger
	//r   *dbredis.Db
	p  *dbpostgres.Db
	Db *dbpostgres.Db
}

func NewHS(log *logger.Logger, passDb string) *HS {

	db := dbpostgres.NewDb(log, passDb)

	return &HS{
		log: log,
		//r:   dbredis.NewDb(log),
		p:  db,
		Db: db,
	}
}
