package NewLogic

import (
	"compendium/logic/dictionary"
	"compendium/logic/ds"
	"compendium/logic/tg"
	"compendium/logic/wa"
	"compendium/models"
	"compendium/storage"

	"github.com/mentalisit/logger"
)

type HsLogic struct {
	log   *logger.Logger
	db    *storage.Storage
	Dict  *dictionary.Dictionary
	moron map[models.IncomingMessage]int
	ds    *ds.Client
	tg    *tg.Client
	wa    *wa.Client
}

func NewCompendiumLogic(log *logger.Logger, m chan models.IncomingMessage, db *storage.Storage) *HsLogic {
	c := &HsLogic{
		log:   log,
		db:    db,
		Dict:  dictionary.NewDictionary(log),
		moron: map[models.IncomingMessage]int{},
		ds:    ds.NewClient(log),
		tg:    tg.NewClient(log),
		wa:    wa.NewClient(log),
	}
	//go c.inbox(m)
	return c
}

//func (c *HsLogic) inbox(m chan models.IncomingMessage) {
//	go c.wsKillTimer()
//	for {
//		select {
//		case in := <-m:
//			go c.logic(in)
//
//		}
//	}
//}
