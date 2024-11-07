package postgres

import (
	"telegram/models"
	"time"
)

func (d *Db) TimerInsert(c models.Timer) {
	ctx, cancel := d.GetContext()
	defer cancel()
	insert := `INSERT INTO kzbot.timer(dsmesid, dschatid, tgmesid, tgchatid, timed) 
				VALUES ($1,$2,$3,$4,$5)`
	_, err := d.db.Exec(ctx, insert, c.Dsmesid, c.Dschatid, c.Tgmesid, c.Tgchatid, c.Timed)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) TimerReadMessage() []models.Timer {
	ctx, cancel := d.GetContext()
	defer cancel()

	tu := int(time.Now().UTC().Unix())

	query := `
        SELECT tgmesid, tgchatid, timed
        FROM kzbot.timer
        WHERE tgmesid IS NOT NULL AND tgmesid <> '' AND $1 > timed;`

	// Выполнение запроса
	rows, _ := d.db.Query(ctx, query, tu)

	defer rows.Close()
	var tt []models.Timer
	for rows.Next() {
		var t models.Timer
		_ = rows.Scan(&t.Tgmesid, &t.Tgchatid, &t.Timed)
		if t.Tgmesid != "" {
			tt = append(tt, t)
		}
	}
	return tt
}
func (d *Db) TimerDeleteMessage(t models.Timer) {
	ctx, cancel := d.GetContext()
	defer cancel()

	query := `DELETE FROM kzbot.timer WHERE tgmesid = $1 `
	_, _ = d.db.Exec(ctx, query, t.Tgmesid)
}
