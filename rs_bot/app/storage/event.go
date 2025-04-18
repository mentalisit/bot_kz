package storage

import "rs/models"

type Event interface {
	NumActiveEvent(CorpName string) (event1 int)    //номер активного ивента
	NumDeactivEvent(CorpName string) (event0 int)   //номер предыдущего ивента
	UpdateActiveEvent0(CorpName string, event1 int) //отключение активного ивента
	EventStartInsert(CorpName string)               //включение ивента
	CountEventNames(CorpName, mention string, numberkz, numEvent int) (countEventNames int)
	CountEventsPoints(CorpName string, numberkz, numberEvent int) int
	UpdatePoints(CorpName string, numberkz, points, event1 int) int
	ReadNamesMessage(CorpName string, numberkz, numberEvent int) (nd, nt models.Names, t models.Sborkz)
	NumberQueueEvents(CorpName string) int
	//ReadEventSchedule() (start string, stop string)
	ReadEventScheduleAndMessage() (nextDateStart, nextDateStop, message string)
	EventInsertPreStart(CorpName string, activeevent int)
	ReadRsEvent(activeEvent int) []models.RsEvent
	UpdateActiveEvent(activeEvent int, CorpName string, numEvent int)
}
