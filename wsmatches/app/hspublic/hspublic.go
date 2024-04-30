package hspublic

import (
	"github.com/mentalisit/logger"
	"ws/dbpostgres"
	"ws/dbredis"
)

type HS struct {
	log *logger.Logger
	r   *dbredis.Db
	p   *dbpostgres.Db
}

func NewHS(log *logger.Logger) *HS {
	return &HS{
		log: log,
		r:   dbredis.NewDb(log),
		p:   dbpostgres.NewDb(log),
	}
}
