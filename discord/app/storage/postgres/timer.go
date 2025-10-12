package postgres

import (
	"discord/models"
	"time"
)

func (d *Db) TimerInsert(c models.Timer) {
	ctx, cancel := d.getContext()
	defer cancel()
	insert := `INSERT INTO rs_bot.timer(tip, chatid, mesid, timed) 
				VALUES ($1,$2,$3,$4)`
	_, err := d.db.Exec(ctx, insert, c.Tip, c.ChatId, c.MesId, c.Timed)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) TimerReadMessage(tip string) []models.Timer {
	ctx, cancel := d.getContext()
	defer cancel()

	tu := int(time.Now().UTC().Unix())

	query := `
        SELECT chatid, mesid, timed
        FROM rs_bot.timer
        WHERE tip = $1 AND chatid IS NOT NULL AND mesid <> '' AND $2 > timed;`

	// Выполнение запроса
	rows, _ := d.db.Query(ctx, query, tip, tu)

	defer rows.Close()
	var tt []models.Timer
	for rows.Next() {
		var t models.Timer
		_ = rows.Scan(&t.ChatId, &t.MesId, &t.Timed)
		if t.MesId != "" {
			tt = append(tt, t)
		}
	}
	return tt
}
func (d *Db) TimerDeleteMessage(t models.Timer) {
	ctx, cancel := d.getContext()
	defer cancel()

	query := `DELETE FROM rs_bot.timer WHERE mesid = $1 `
	_, _ = d.db.Exec(ctx, query, t.MesId)
}
