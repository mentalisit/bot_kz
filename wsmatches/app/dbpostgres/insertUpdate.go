package dbpostgres

import (
	"errors"
	"github.com/jackc/pgx/v5"
	"time"
	"ws/models"
)

func (d *Db) InsertUpdateCorpsLevelOld(l models.LevelCorps) {
	ctx, cancel := d.GetContext()
	defer cancel()
	_, err := d.ReadCorpsLevel(l.HCorp)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			// Если запись не найдена, вставляем новую запись
			insert := `INSERT INTO kzbot.corpslevel(corpname, level, enddate, hcorp, percent, last_update, relic) VALUES ($1,$2,$3,$4,$5,$6,$7)`
			_, err = d.pool.Exec(ctx, insert, l.CorpName, l.Level, l.EndDate, l.HCorp, l.Percent, l.LastUpdate, l.Relic)
			if err != nil {
				d.log.ErrorErr(err)
			}
			return
		case err != nil:
			d.log.ErrorErr(err)
			return
		}
	}

	upd := `update kzbot.corpslevel set level = $1,enddate = $2,hcorp = $3,percent = $4,last_update = $5,relic = $6 where corpname = $7`
	_, err = d.pool.Exec(ctx, upd, l.Level, l.EndDate, l.HCorp, l.Percent, l.LastUpdate, l.Relic, l.CorpName)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) InsertUpdateCorpsLevel(l models.LevelCorps) {
	ctx, cancel := d.GetContext()
	defer cancel()

	_, err := d.ReadCorpsLevel(l.HCorp)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			// Если записи нет, вставляем новую
			insert := `INSERT INTO ws.corpslevel (corpname, level, enddate, hcorp, percent, last_update, relic) 
					   VALUES ($1, $2, $3, $4, $5, $6, $7)`
			_, err = d.pool.Exec(ctx, insert,
				l.CorpName,
				l.Level,
				l.EndDate.Format(time.RFC3339Nano), // Преобразуем в строку
				l.HCorp,
				l.Percent,
				l.LastUpdate.Format(time.RFC3339Nano), // Преобразуем в строку
				l.Relic,
			)
			if err != nil {
				d.log.ErrorErr(err)
			}
			return
		default:
			d.log.ErrorErr(err)
			return
		}
	}

	// Если запись есть, обновляем
	upd := `UPDATE ws.corpslevel 
			SET level = $1, enddate = $2, hcorp = $3, percent = $4, last_update = $5, relic = $6 
			WHERE corpname = $7`
	_, err = d.pool.Exec(ctx, upd,
		l.Level,
		l.EndDate.Format(time.RFC3339Nano), // Преобразуем в строку
		l.HCorp,
		l.Percent,
		l.LastUpdate.Format(time.RFC3339Nano), // Преобразуем в строку
		l.Relic,
		l.CorpName,
	)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
