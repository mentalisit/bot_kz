package models

import "time"

type LevelCorp struct {
	CorpName string
	Level    int
	EndDate  time.Time
	HCorp    string
	Percent  int
}
