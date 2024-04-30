package storage

import (
	"fmt"
	"github.com/mentalisit/logger"
	"kz_bot/config"
	"kz_bot/models"
	"kz_bot/storage/dictionary"
	"kz_bot/storage/mongo"
	"kz_bot/storage/postgres"
	"kz_bot/storage/reststorage"
)

type Storage struct {
	log               *logger.Logger
	debug             bool
	BridgeConfig      *reststorage.Db
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
	BridgeConfigs     map[string]models.BridgeConfig
	CorpConfigRS      map[string]models.CorporationConfig
	mongo             *mongo.DB
	postgres          *postgres.Db
}

func NewStorage(log *logger.Logger, cfg *config.ConfigBot) *Storage {

	//Initializing a repository from a cloud configuration
	mongoDB := mongo.InitMongoDB(log)

	//REST API
	rdb := reststorage.InitRestApiStorage(log)

	//add language packages
	d := dictionary.NewDictionary(log)

	//Initializing a local repository
	local := postgres.NewDb(log, cfg)

	s := &Storage{
		log:               log,
		BridgeConfig:      rdb,
		TimeDeleteMessage: local,
		ConfigRs:          mongoDB,
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
		BridgeConfigs:     make(map[string]models.BridgeConfig),
		CorpConfigRS:      make(map[string]models.CorporationConfig),
		mongo:             mongoDB,
		postgres:          local,
	}

	go s.loadDbArray()
	return s
}
func (s *Storage) loadDbArray() {
	var bridgeCounter = 0
	var bridge string
	bc := s.BridgeConfig.DBReadBridgeConfig()
	for _, configBridge := range bc {
		s.BridgeConfigs[configBridge.NameRelay] = configBridge
		bridgeCounter++
		bridge = bridge + fmt.Sprintf("%s, ", configBridge.HostRelay)
	}
	fmt.Printf("Загружено конфиг мостов %d : %s\n", bridgeCounter, bridge)

	var c = 0
	var rslist string
	rs := s.ConfigRs.ReadConfigRs()
	for _, r := range rs {
		s.CorpConfigRS[r.CorpName] = r
		c++
		rslist = rslist + fmt.Sprintf("%s, ", r.CorpName)
	}
	fmt.Printf("Загружено конфиг RsBot %d : %s\n", c, rslist)
}
func (s *Storage) ReloadDbArray() {
	CorpConfigRS := make(map[string]models.CorporationConfig)
	BridgeConfigs := make(map[string]models.BridgeConfig)

	s.CorpConfigRS = CorpConfigRS
	s.BridgeConfigs = BridgeConfigs

	bridgeConfig := s.BridgeConfig.DBReadBridgeConfig()
	for _, configBridge := range bridgeConfig {
		s.BridgeConfigs[configBridge.NameRelay] = configBridge
	}
	rs := s.ConfigRs.ReadConfigRs()
	for _, r := range rs {
		s.CorpConfigRS[r.CorpName] = r
	}
}
func (s *Storage) Shutdown() {
	s.mongo.Shutdown()
	s.postgres.Shutdown()
}
