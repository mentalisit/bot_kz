package models

import "time"

//type LevelCorp struct {
//	CorpName string
//	Level    int
//	EndDate  time.Time
//	HCorp    string
//	Percent  int
//}

type LevelCorps struct {
	CorpName   string
	Level      int
	EndDate    time.Time
	HCorp      string
	Percent    int
	LastUpdate time.Time
	Relic      int
}
