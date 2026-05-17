package postgres

import (
	"database/sql"
	"errors"
	"rs/models"
)

func (d *Db) UpdatePoints(CorpName string, numberkz, points, event1 int) int {
	// считаем количество участников КЗ опр уровня
	var countEvent int
	selec := "SELECT  COUNT(*) as count FROM kzbot.sborkz WHERE numberevent = $1 AND corpname=$2 AND numberkz=$3  AND active=1"
	row := d.db.QueryRow(selec, event1, CorpName, numberkz)
	err := row.Scan(&countEvent)
	if err != nil {
		d.log.ErrorErr(err)
	}
	if countEvent == 0 {
		return 0
	}
	pointsq := points / countEvent
	//вносим очки
	upd := `update kzbot.sborkz set eventpoints=$1 WHERE numberevent = $2 AND corpname =$3 AND numberkz=$4 AND active=1`
	_, err = d.db.Exec(upd, pointsq, event1, CorpName, numberkz)
	if err != nil {
		d.log.ErrorErr(err)
	}
	return countEvent
}
func (d *Db) ReadNamesMessage(CorpName string, numberkz, numberEvent int) (nd, nt models.Names, t models.Sborkz) {
	var name string
	sel := "SELECT * FROM kzbot.sborkz WHERE corpname=$1 AND numberkz=$2 AND numberevent = $3 AND active=1"
	results, err := d.db.Query(sel, CorpName, numberkz, numberEvent)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}

	num := 1
	for results.Next() {
		err = results.Scan(&t.Id, &t.Corpname, &t.Name, &t.Mention, &t.Tip, &t.Dsmesid, &t.Tgmesid, &t.Wamesid, &t.Time, &t.Date, &t.Lvlkz, &t.Numkzn, &t.Numberkz, &t.Numberevent, &t.Eventpoints, &t.Active, &t.Timedown, &t.UserId)
		if t.Tip == "ds" {
			name = t.Mention
		} else {
			name = t.Name
		}
		if num == 1 {
			nd.Name1 = name
		} else if num == 2 {
			nd.Name2 = name
		} else if num == 3 {
			nd.Name3 = name
		} else if num == 4 {
			nd.Name4 = name
		}
		if t.Tip == "tg" {
			name = t.Mention
		} else {
			name = t.Name
		}
		if num == 1 {
			nt.Name1 = name
		} else if num == 2 {
			nt.Name2 = name
		} else if num == 3 {
			nt.Name3 = name
		} else if num == 4 {
			nt.Name4 = name
		}
		num = num + 1
	}
	return nd, nt, t
}
func (d *Db) CountEventNames(CorpName, mention string, numberkz, numEvent int) (countEventNames int) {
	sel := "SELECT  COUNT(*) as count FROM kzbot.sborkz WHERE corpname = $1 AND numberkz=$2  AND active=1 AND mention=$3 AND numberevent = $4"
	row := d.db.QueryRow(sel, CorpName, numberkz, mention, numEvent)
	err := row.Scan(&countEventNames)
	if err != nil {
		d.log.ErrorErr(err)
	}
	return countEventNames
}
func (d *Db) CountEventsPoints(CorpName string, numberkz, numberEvent int) int {
	var countEventPoints int
	sel := "SELECT  COUNT(*) as count FROM kzbot.sborkz WHERE corpname=$1 AND numberkz=$2 AND numberevent = $3 AND active=1 AND eventpoints > 0"
	row := d.db.QueryRow(sel, CorpName, numberkz, numberEvent)
	err := row.Scan(&countEventPoints)
	if err != nil {
		d.log.ErrorErr(err)
	}
	return countEventPoints
}
func (d *Db) NumActiveEvent(CorpName string) (event1 int) { //запрос номера ивента
	sel := "SELECT numevent FROM kzbot.rsevent WHERE corpname=$1 AND activeevent=1 ORDER BY numevent DESC LIMIT 1"
	row := d.db.QueryRow(sel, CorpName)
	err := row.Scan(&event1)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			event1 = 0
		} else {
			d.log.ErrorErr(err)
		}
	}
	return event1
}
func (d *Db) NumDeactivEvent(CorpName string) (event0 int) { //запрос номера последнего ивента
	sel := "SELECT max(numevent) FROM kzbot.rsevent WHERE corpname=$1 AND activeevent=0"
	row := d.db.QueryRow(sel, CorpName)
	err := row.Scan(&event0)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			d.log.ErrorErr(err)
		}
	}
	return event0
}
func (d *Db) UpdateActiveEvent0(CorpName string, event1 int) {
	upd := "UPDATE kzbot.rsevent SET activeevent=0 WHERE corpname=$1 AND numevent=$2"
	_, err := d.db.Exec(upd, CorpName, event1)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) EventStartInsert(CorpName string) {
	event0 := d.NumDeactivEvent(CorpName)
	insertEvent := `INSERT INTO kzbot.rsevent (corpname,numevent,activeevent,number) VALUES ($1,$2,$3,$4)`
	if event0 > 0 {
		numberevent := event0 + 1
		_, err := d.db.Exec(insertEvent, CorpName, numberevent, 1, 1)
		if err != nil {
			d.log.ErrorErr(err)
		}
	} else {
		_, err := d.db.Exec(insertEvent, CorpName, 1, 1, 1)
		if err != nil {
			d.log.ErrorErr(err)
		}
	}
}
func (d *Db) NumberQueueEvents(CorpName string) int {
	var number int
	sel := "SELECT  number FROM kzbot.rsevent WHERE activeevent = 1 AND corpname = $1 "
	row := d.db.QueryRow(sel, CorpName)
	err := row.Scan(&number)
	if err != nil {
		d.log.ErrorErr(err)
	}
	return number
}

// new

// activeevent int -1 prepare, 0 stop , 1 run
func (d *Db) EventInsertPreStart(CorpName string, activeevent int) {
	event0 := d.NumDeactivEvent(CorpName)
	insertEvent := `INSERT INTO kzbot.rsevent (corpname,numevent,activeevent,number) VALUES ($1,$2,$3,$4)`
	if event0 > 0 {
		numberevent := event0 + 1
		_, err := d.db.Exec(insertEvent, CorpName, numberevent, activeevent, 1)
		if err != nil {
			d.log.ErrorErr(err)
		}
	} else {
		_, err := d.db.Exec(insertEvent, CorpName, 1, activeevent, 1)
		if err != nil {
			d.log.ErrorErr(err)
		}
	}
}

// ReadRsEvent activeEvent int -1 prepare, 0 stop , 1 run
func (d *Db) ReadRsEvent(activeEvent int) []models.RsEvent {
	sel := "SELECT * FROM kzbot.rsevent WHERE activeevent=$1"
	results, err := d.db.Query(sel, activeEvent)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	var eventsCorps []models.RsEvent

	for results.Next() {
		var c models.RsEvent
		err = results.Scan(&c.Id, &c.CorpName, &c.NumEvent, &c.ActiveEvent, &c.Number)
		if err != nil {
			d.log.ErrorErr(err)
		}
		eventsCorps = append(eventsCorps, c)
	}

	return eventsCorps
}
func (d *Db) UpdateActiveEvent(activeEvent int, CorpName string, numEvent int) {
	upd := "UPDATE kzbot.rsevent SET activeevent=$1 WHERE corpname=$2 AND numevent=$3"
	_, err := d.db.Exec(upd, activeEvent, CorpName, numEvent)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
