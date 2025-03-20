package storage

import (
	"discord/config"
	"discord/storage/postgres"
	"github.com/mentalisit/logger"
)

type Storage struct {
	Db         *postgres.Db
	Emoji      Emoji
	Scoreboard Scoreboard
	Battles    Battles
}

func NewStorage(log *logger.Logger, cfg *config.ConfigBot) *Storage {
	local := postgres.NewDb(log, cfg)

	s := &Storage{
		Db:         local,
		Emoji:      local,
		Scoreboard: local,
		Battles:    local,
	}

	//go s.loadDbArray()

	return s
}
