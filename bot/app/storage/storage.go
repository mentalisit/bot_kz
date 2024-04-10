package storage

import (
	"fmt"
	"github.com/mentalisit/logger"
	"go.uber.org/zap"
	"kz_bot/config"
	"kz_bot/models"
	"kz_bot/storage/mongo"
	"kz_bot/storage/postgres"
	"kz_bot/storage/reststorage"
	"kz_bot/storage/words"
	"time"
)

type Storage struct {
	log               *zap.Logger
	debug             bool
	BridgeConfig      *reststorage.Db
	ConfigRs          ConfigRs
	TimeDeleteMessage TimeDeleteMessage
	Words             *words.Words
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
}

func NewStorage(log *logger.Logger, cfg *config.ConfigBot) *Storage {

	//инициализируем и читаем репозиторий из облока конфига конфигурации
	mongoDB := mongo.InitMongoDB(log)

	rdb := reststorage.InitRestApiStorage(log)
	//corp := CorpsConfig.NewCorps(log, cfg)

	//подключаю языковой пакет
	w := words.NewWords()

	//инициализируем локальный репозиторий
	local := postgres.NewDb(log, cfg)

	s := &Storage{
		//CorpsConfig:       corp,
		//HadesClient:       mongoDB,
		BridgeConfig:      rdb,
		TimeDeleteMessage: mongoDB,
		ConfigRs:          mongoDB,
		Words:             w,
		Subscribe:         local,
		Emoji:             local,
		Count:             local,
		Top:               local,
		Update:            local,
		Timers:            local,
		DbFunc:            local,
		Event:             local,
		LevelCorp:         local,
		//CorporationHades:  make(map[string]models.CorporationHadesClient),
		BridgeConfigs: make(map[string]models.BridgeConfig),
		CorpConfigRS:  make(map[string]models.CorporationConfig),
	}

	go s.loadDbArray()
	return s
}
func (s *Storage) loadDbArray() {
	var b = 0
	var bridge string
	bc := s.BridgeConfig.DBReadBridgeConfig()
	for _, configBridge := range bc {
		s.BridgeConfigs[configBridge.NameRelay] = configBridge
		b++
		bridge = bridge + fmt.Sprintf("%s, ", configBridge.HostRelay)
	}
	fmt.Printf("Загружено конфиг мостов %d : %s\n", b, bridge)

	var c = 0
	var rslist string
	rs := s.ConfigRs.ReadConfigRs()
	for _, r := range rs {
		s.CorpConfigRS[r.CorpName] = r
		s.LevelCorp.InsertUpdateCorpLevel(models.LevelCorp{
			CorpName: r.CorpName,
			Level:    0,
			EndDate:  time.Time{},
			HCorp:    "",
			Percent:  0,
		})
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

	bc := s.BridgeConfig.DBReadBridgeConfig()
	for _, configBridge := range bc {
		s.BridgeConfigs[configBridge.NameRelay] = configBridge
	}
	rs := s.ConfigRs.ReadConfigRs()
	for _, r := range rs {
		s.CorpConfigRS[r.CorpName] = r
	}
}
