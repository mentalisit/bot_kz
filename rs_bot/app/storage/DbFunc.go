package storage

import (
	"rs/models"
)

type DbFunc interface {
	ReadAll(lvlkz, CorpName string) (users models.Users)
	InsertQueue(dsmesid, wamesid, CorpName, name, userid, nameMention, tip, lvlkz, timekz string, tgmesid, numkzN int)
	ElseTrue(userid string) []models.Sborkz
	DeleteQueue(userid, lvlkz, CorpName string)
	ReadMesIdDS(mesid string) (string, error)
	P30Pl(lvlkz, CorpName, userid string) int
	UpdateTimedown(lvlkz, CorpName, userid string)
	Queue(corpname string) []string
	OneMinutsTimer() []string
	MessageUpdateMin(corpname string) ([]string, []int)
	MessageupdateDS(dsmesid string, config models.CorporationConfig) models.InMessage
	MessageupdateTG(tgmesid int, config models.CorporationConfig) models.InMessage
	NumberQueueLvl(lvlkzs, CorpName string) (int, error)
	OptimizationSborkz()
}
