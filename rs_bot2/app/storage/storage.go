package storage

import (
	"rs/config"
	"rs/storage/dictionary"
	"rs/storage/postgresV2"

	"github.com/mentalisit/logger"
)

type Storage struct {
	log        *logger.Logger
	Dictionary *dictionary.Dictionary
	//Postgres   *postgres.Db
	//postgres   *postgres.Db
	V2 *postgresV2.Db
}

func NewStorage(log *logger.Logger, cfg *config.ConfigBot) *Storage {
	//add language packages
	d := dictionary.NewDictionary(log)

	//Initializing a local repository
	//local := postgres.NewDb(log, cfg)
	V2 := postgresV2.NewDb(log, cfg)

	s := &Storage{
		log:        log,
		Dictionary: d,
		//postgres:   local,
		//Postgres:   local,
		V2: V2,
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
	//s.postgres.Shutdown()
}

//type EventNumber interface {
//	InsertEventNumber(number, event int, status bool) error
//	GetEventNumber() (events []models.Events, err error)
//	UpdateEventStatus(id int64, newStatus bool) error
//}
