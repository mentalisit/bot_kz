package logic

import (
	"compendium/logic/dictionary"
	"compendium/logic/ds"
	"compendium/logic/tg"
	"compendium/logic/wa"
	"compendium/models"
	"compendium/storage"
	postgresv2 "compendium/storage/postgres/postgresV2"

	"github.com/mentalisit/logger"
)

type Hs struct {
	log        *logger.Logger
	db         *storage.Storage
	DbV2       *postgresv2.Db
	corpMember CorpMember
	tech       Tech
	users      Users
	guildsRole GuildRoles
	Dict       *dictionary.Dictionary
	moron      map[models.IncomingMessage]int
	ds         *ds.Client
	tg         *tg.Client
	wa         *wa.Client
}

func NewCompendium(log *logger.Logger, m chan models.IncomingMessage, db *storage.Storage) *Hs {
	c := &Hs{
		log:        log,
		db:         db,
		DbV2:       db.V2,
		corpMember: db.DB,
		tech:       db.DB,
		users:      db.DB,
		guildsRole: db.DB,
		Dict:       dictionary.NewDictionary(log),
		moron:      map[models.IncomingMessage]int{},
		ds:         ds.NewClient(log),
		tg:         tg.NewClient(log),
		wa:         wa.NewClient(log),
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
