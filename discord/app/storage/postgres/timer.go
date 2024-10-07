package postgres

import (
	"discord/models"
)

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
