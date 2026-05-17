package models

import "time"

type LevelCorps struct {
	CorpName   string
	Level      int
	EndDate    time.Time
	HCorp      string
	Percent    int
	LastUpdate time.Time
	Relic      int
}
