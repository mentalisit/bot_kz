package postgres

import (
	"errors"
	"github.com/jackc/pgx/v4"
	"kz_bot/models"
)

func (d *Db) InsertUpdateCorpLevel(l models.LevelCorps) {
	ctx, cancel := d.GetContext()
	defer cancel()
	_, err := d.ReadCorpLevel(l.CorpName)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			// Если запись не найдена, вставляем новую запись
			insert := `INSERT INTO kzbot.corpslevel(corpname, level, enddate, hcorp, percent, last_update, relic) VALUES ($1,$2,$3,$4,$5,$6,$7)`
			_, err = d.db.Exec(ctx, insert, l.CorpName, l.Level, l.EndDate, l.HCorp, l.Percent, l.LastUpdate, l.Relic)
			if err != nil {
				d.log.ErrorErr(err)
			}
			return
		case err != nil:
			d.log.ErrorErr(err)
			return
		}
	}

	upd := `update kzbot.corpslevel set level = $1,enddate = $2,hcorp = $3,percent = $4 where corpname = $5`
	_, err = d.db.Exec(ctx, upd, l.Level, l.EndDate, l.HCorp, l.Percent, l.CorpName)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) ReadCorpLevel(CorpName string) (models.LevelCorps, error) {
	ctx, cancel := d.GetContext()
	defer cancel()
	var n models.LevelCorps
	err := d.db.QueryRow(ctx, "SELECT * FROM kzbot.corpslevel WHERE corpname = $1", CorpName).Scan(
		&n.CorpName, &n.Level, &n.EndDate, &n.HCorp, &n.Percent, &n.LastUpdate, &n.Relic)
	if err != nil {
		return models.LevelCorps{}, err
	}
	return n, nil
}

func (d *Db) ReadCorpLevelAll() ([]models.LevelCorps, error) {
	ctx, cancel := d.GetContext()
	defer cancel()
	var n []models.LevelCorps
	sel := "SELECT * FROM kzbot.corpslevel"
	results, err := d.db.Query(ctx, sel)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}
	for results.Next() {
		var t models.LevelCorps
		err = results.Scan(&t.CorpName, &t.Level, &t.EndDate, &t.HCorp, &t.Percent, &t.LastUpdate, &t.Relic)
		if err != nil {
			d.log.ErrorErr(err)
		}
		n = append(n, t)
	}
	return n, nil
}
