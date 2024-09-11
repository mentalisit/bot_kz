package kzbotdb

import (
	"context"
	"queue/models"
)

func (d *Db) SelectSborkzActive() []models.Sborkz {
	sel := "SELECT * FROM kzbot.sborkz WHERE active = 0"
	results, err := d.db.Query(context.Background(), sel)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	var tt []models.Sborkz
	for results.Next() {
		var t models.Sborkz
		err = results.Scan(&t.Id, &t.Corpname, &t.Name, &t.Mention, &t.Tip, &t.Dsmesid, &t.Tgmesid, &t.Wamesid, &t.Time, &t.Date, &t.Lvlkz, &t.Numkzn, &t.Numberkz, &t.Numberevent, &t.Eventpoints, &t.Active, &t.Timedown, &t.UserId)
		tt = append(tt, t)

	}
	return tt
}

func (d *Db) SelectSborkzActiveLevel(level string) []models.Sborkz {
	sel := "SELECT * FROM kzbot.sborkz WHERE lvlkz = $1 AND active = 0"
	results, err := d.db.Query(context.Background(), sel, level)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	var tt []models.Sborkz
	for results.Next() {
		var t models.Sborkz
		err = results.Scan(&t.Id, &t.Corpname, &t.Name, &t.Mention, &t.Tip, &t.Dsmesid, &t.Tgmesid, &t.Wamesid, &t.Time, &t.Date, &t.Lvlkz, &t.Numkzn, &t.Numberkz, &t.Numberevent, &t.Eventpoints, &t.Active, &t.Timedown, &t.UserId)
		tt = append(tt, t)

	}
	return tt
}
