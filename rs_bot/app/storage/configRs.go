package storage

import "rs/models"

type ConfigRs interface {
	InsertConfigRs(c models.CorporationConfig)
	ReadConfigRs() []models.CorporationConfig
	DeleteConfigRs(c models.CorporationConfig)
	UpdateConfigRs(c models.CorporationConfig)
	ReadConfigForDsChannel(dsChannel string) (conf models.CorporationConfig)
	ReadConfigForTgChannel(tgChannel string) (conf models.CorporationConfig)
	ReadConfigForCorpName(corpName string) (conf models.CorporationConfig)
}

//func (s *Storage) DeleteConfigRs(c models.CorporationConfig) {
//	s.ConfigRs.DeleteConfigRs(c)
//	var a map[string]models.CorporationConfig
//	a = make(map[string]models.CorporationConfig)
//	b := s.ConfigRs.ReadConfigRs()
//	for _, config := range b {
//		a[config.CorpName] = config
//	}
//	s.CorpConfigRS = a
//}
