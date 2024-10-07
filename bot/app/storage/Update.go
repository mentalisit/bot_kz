package storage

type Update interface {
	MesidTgUpdate(mesidtg int, lvlkz string, corpname string) error
	MesidDsUpdate(mesidds, lvlkz, corpname string) error

	UpdateCompliteRS(lvlkz string, dsmesid string, tgmesid int, wamesid string, numberkz int, numberevent int, corpname string) error
	UpdateCompliteSolo(lvlkz string, dsmesid string, tgmesid int, numberkz int, numberevent int, corpname string) error
}
