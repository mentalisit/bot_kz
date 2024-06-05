package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"kz_bot/models"
)

func (d *Db) InsertUpdateCorpLevel(l models.LevelCorp) {
	ctx := context.Background()
	_, err := d.ReadCorpLevel(l.CorpName)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			// Если запись не найдена, вставляем новую запись
			insert := `INSERT INTO kzbot.corplevel(corpname, level, enddate, hcorp, percent) VALUES ($1,$2,$3,$4,$5)`
			_, err = d.db.Exec(ctx, insert, l.CorpName, l.Level, l.EndDate, l.HCorp, l.Percent)
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
	_, err = d.db.Exec(ctx, upd, l.Level, l.EndDate, l.HCorp, l.Percent, l.CorpName)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) ReadCorpLevel(CorpName string) (models.LevelCorp, error) {
	var n models.LevelCorp
	err := d.db.QueryRow(context.Background(), "SELECT * FROM kzbot.corplevel WHERE corpname = $1", CorpName).Scan(
		&n.CorpName, &n.Level, &n.EndDate, &n.HCorp, &n.Percent)
	if err != nil {
		return models.LevelCorp{}, err
	}
	return n, nil
}

func (d *Db) ReadCorpLevelAll() ([]models.LevelCorp, error) {
	var n []models.LevelCorp
	sel := "SELECT * FROM kzbot.corplevel"
	results, err := d.db.Query(context.Background(), sel)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}
	for results.Next() {
		var t models.LevelCorp
		err = results.Scan(&t.CorpName, &t.Level, &t.EndDate, &t.HCorp, &t.Percent)
		if err != nil {
			d.log.ErrorErr(err)
		}
		n = append(n, t)
	}
	return n, nil
}
