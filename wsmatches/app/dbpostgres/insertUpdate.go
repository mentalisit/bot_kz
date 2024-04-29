package dbpostgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"ws/models"
)

func (d *Db) InsertUpdateCorpLevel(l models.LevelCorp) {
	ctx := context.Background()
	_, err := d.ReadCorpLevel(l.CorpName)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			// Если запись не найдена, вставляем новую запись
			insert := `INSERT INTO kzbot.corplevel(corpname, level, enddate, hcorp, percent) VALUES ($1,$2,$3,$4,$5)`
			_, err = d.pool.Exec(ctx, insert, l.CorpName, l.Level, l.EndDate, l.HCorp, l.Percent)
			if err != nil {
				d.log.ErrorErr(err)
			}
			return
		case err != nil:
			d.log.ErrorErr(err)
			return
		}
	}

	upd := `update kzbot.corplevel set level = $1,enddate = $2,hcorp = $3,percent = $4 where corpname = $5`
	_, err = d.pool.Exec(ctx, upd, l.Level, l.EndDate, l.HCorp, l.Percent, l.CorpName)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
