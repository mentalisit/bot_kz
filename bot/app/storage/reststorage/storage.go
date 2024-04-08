package reststorage

import (
	"github.com/mentalisit/logger"
)

type Db struct {
	log *logger.Logger
}

func InitRestApiStorage(log *logger.Logger) *Db {
	d := &Db{log: log}
	return d
}
