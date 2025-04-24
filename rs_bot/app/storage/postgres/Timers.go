package postgres

import (
	"rs/models"
)

func (d *Db) UpdateMitutsQueue(userid, CorpName string) models.Sborkz {
	ctx, cancel := d.getContext()
	defer cancel()

	sel := "SELECT * FROM kzbot.sborkz WHERE userid = $1 AND corpname = $2 AND active = 0"
	results, err := d.db.Query(ctx, sel, userid, CorpName)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	var t models.Sborkz
	for results.Next() {

		err = results.Scan(&t.Id, &t.Corpname, &t.Name, &t.Mention, &t.Tip, &t.Dsmesid, &t.Tgmesid, &t.Wamesid, &t.Time,
			&t.Date, &t.Lvlkz, &t.Numkzn, &t.Numberkz, &t.Numberevent, &t.Eventpoints, &t.Active, &t.Timedown, &t.UserId)

		if t.UserId == userid && t.Timedown <= 3 {
			upd := "update kzbot.sborkz set timedown = timedown + 30 where active = 0 AND userid = $1 AND corpname = $2"
			_, err = d.db.Exec(ctx, upd, t.UserId, t.Corpname)
			if err != nil {
				d.log.ErrorErr(err)
			}
		}
	}

	return t
}

func (d *Db) MinusMin() []models.Sborkz {
	ctx, cancel := d.getContext()
	defer cancel()
	upd := `update kzbot.sborkz set timedown = timedown - 1 where active = 0`
	_, err := d.db.Exec(ctx, upd)
	if err != nil {
		d.log.ErrorErr(err)
		return []models.Sborkz{}
	}

	sel := "SELECT * FROM kzbot.sborkz WHERE active = 0"
	results, err := d.db.Query(ctx, sel)
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

func (d *Db) TimerInsert(c models.Timer) {
	ctx, cancel := d.getContext()
	defer cancel()
	insert := `INSERT INTO kzbot.timer(dsmesid, dschatid, tgmesid, tgchatid, timed) 
				VALUES ($1,$2,$3,$4,$5)`
	_, err := d.db.Exec(ctx, insert, c.Dsmesid, c.Dschatid, c.Tgmesid, c.Tgchatid, c.Timed)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) TimerDeleteMessage() []models.Timer {
	ctx, cancel := d.getContext()
	defer cancel()
	query := `UPDATE kzbot.timer SET timed = timed - 60 WHERE timed > 60`

	_, _ = d.db.Exec(ctx, query)

	query = `SELECT * FROM kzbot.timer WHERE timed <= 60`

	// Выполнение запроса
	rows, _ := d.db.Query(ctx, query)

	defer rows.Close()
	var tt []models.Timer
	for rows.Next() {
		var id int
		var t models.Timer
		_ = rows.Scan(&id, &t.Dsmesid, &t.Dschatid, &t.Tgmesid, &t.Tgchatid, &t.Timed)
		tt = append(tt, t)
	}
	query = `DELETE FROM kzbot.timer WHERE dsmesid = $1 AND tgmesid = $2`
	for _, t := range tt {
		_, _ = d.db.Exec(ctx, query, t.Dsmesid, t.Tgmesid)
	}
	return tt
}
