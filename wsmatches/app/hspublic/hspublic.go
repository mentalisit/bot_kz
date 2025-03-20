package hspublic

import (
	"github.com/mentalisit/logger"
	"ws/dbpostgres"
)

type HS struct {
	log *logger.Logger
	//r   *dbredis.Db
	p *dbpostgres.Db
}

func NewHS(log *logger.Logger, passDb string) *HS {
	return &HS{
		log: log,
		//r:   dbredis.NewDb(log),
		p: dbpostgres.NewDb(log, passDb),
	}
}
