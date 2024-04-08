package Compendium

import "C"
import (
	"compendium/models"
	"compendium/storage"
	"github.com/mentalisit/logger"
)

type Compendium struct {
	log *logger.Logger
	in  models.IncomingMessage
	db  *storage.Storage
}

func NewCompendium(log *logger.Logger, m chan models.IncomingMessage, db *storage.Storage) *Compendium {
	c := &Compendium{
		log: log,
		db:  db,
	}
	go c.inbox(m)
	return c
}

func (c *Compendium) inbox(m chan models.IncomingMessage) {
	for {
		select {
		case in := <-m:
			go c.logic(in)

		}
	}
}
