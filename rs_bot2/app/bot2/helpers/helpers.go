package helpers

import (
	"rs/storage/postgresV2"

	"github.com/mentalisit/logger"
)

const ds = "ds"
const tg = "tg"

type Helpers struct {
	log     *logger.Logger
	storage *postgresV2.Db
}

func NewHelpers(log *logger.Logger, storage *postgresV2.Db) *Helpers {
	return &Helpers{
		log:     log,
		storage: storage,
	}
}
