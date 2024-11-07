package logic

import "C"
import (
	"compendium/logic/dictionary"
	"compendium/logic/ds"
	"compendium/logic/tg"
	"compendium/models"
	"compendium/storage"
	"github.com/mentalisit/logger"
)

type Hs struct {
	log *logger.Logger
	//in         models.IncomingMessage
	db         *storage.Storage
	corpMember CorpMember
	tech       Tech
	guilds     Guilds
	listUser   ListUser
	users      Users
	guildsRole GuildRoles
	Dict       *dictionary.Dictionary
	moron      map[models.IncomingMessage]int
	ds         *ds.Client
	tg         *tg.Client
}

func NewCompendium(log *logger.Logger, m chan models.IncomingMessage, db *storage.Storage) *Hs {
	c := &Hs{
		log:        log,
		db:         db,
		corpMember: db.DB,
		tech:       db.DB,
		guilds:     db.DB,
		listUser:   db.DB,
		users:      db.DB,
		guildsRole: db.DB,
		Dict:       dictionary.NewDictionary(log),
		moron:      map[models.IncomingMessage]int{},
		ds:         ds.NewClient(log),
		tg:         tg.NewClient(log),
	}
	go c.inbox(m)
	return c
}

func (c *Hs) inbox(m chan models.IncomingMessage) {
	go c.wsKillTimer()
	for {
		select {
		case in := <-m:
			go c.logic(in)

		}
	}
}
