package storage

import (
	"github.com/mentalisit/logger"
	"rs/config"
	"rs/storage/dictionary"
	"rs/storage/postgres"
)

type Storage struct {
	log               *logger.Logger
	ConfigRs          ConfigRs
	TimeDeleteMessage TimeDeleteMessage
	Dictionary        *dictionary.Dictionary
	Subscribe         Subscribe
	Emoji             Emoji
	Count             Count
	Top               Top
	Update            Update
	Timers            Timers
	DbFunc            DbFunc
	Event             Event
	LevelCorp         LevelCorp
	Postgres          *postgres.Db
	Scoreboard        Scoreboard
	Battles           Battles
	postgres          *postgres.Db
}

func NewStorage(log *logger.Logger, cfg *config.ConfigBot) *Storage {
	//add language packages
	d := dictionary.NewDictionary(log)

	//Initializing a local repository
	local := postgres.NewDb(log, cfg)

	s := &Storage{
		log:               log,
		TimeDeleteMessage: local,
		ConfigRs:          local,
		Dictionary:        d,
		Subscribe:         local,
		Emoji:             local,
		Count:             local,
		Top:               local,
		Update:            local,
		Timers:            local,
		DbFunc:            local,
		Event:             local,
		LevelCorp:         local,
		postgres:          local,
		Scoreboard:        local,
		Battles:           local,
		Postgres:          local,
	}

	//go s.loadDbArray()
	return s
}

//func (s *Storage) loadDbArray() {
//	s.BridgeConfigs = restapi.ReadBridgeConfig()
//
//	//var c = 0
//	//var rslist string
//	//rs := s.ConfigRs.ReadConfigRs()
//	//for _, r := range rs {
//	//	s.CorpConfigRS[r.CorpName] = r
//	//	c++
//	//	rslist = rslist + fmt.Sprintf("%s, ", r.CorpName)
//	//}
//	//fmt.Printf("Загружено конфиг RsBot %d : %s\n", c, rslist)
//}
//func (s *Storage) ReloadDbArray() {
//	s.BridgeConfigs = restapi.ReadBridgeConfig()
//
//	//CorpConfigRS := make(map[string]models.CorporationConfig)
//
//	//s.CorpConfigRS = CorpConfigRS
//	//rs := s.ConfigRs.ReadConfigRs()
//	//for _, r := range rs {
//	//	s.CorpConfigRS[r.CorpName] = r
//	//}
//}

func (s *Storage) Shutdown() {
	s.postgres.Shutdown()
}

//type EventNumber interface {
//	InsertEventNumber(number, event int, status bool) error
//	GetEventNumber() (events []models.Events, err error)
//	UpdateEventStatus(id int64, newStatus bool) error
//}
