package postgres

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"rs/models"
	"time"
)

//func (d *Db) InsertUpdateCorpLevel(l models.LevelCorps) {
//	ctx, cancel := d.GetContext()
//	defer cancel()
//	_, err := d.ReadCorpLevel(l.CorpName)
//	if err != nil {
//		switch {
//		case errors.Is(err, pgx.ErrNoRows):
//			// Если запись не найдена, вставляем новую запись
//			insert := `INSERT INTO kzbot.corpslevel(corpname, level, enddate, hcorp, percent, last_update, relic) VALUES ($1,$2,$3,$4,$5,$6,$7)`
//			_, err = d.db.Exec(ctx, insert, l.CorpName, l.Level, l.EndDate, l.HCorp, l.Percent, l.LastUpdate, l.Relic)
//			if err != nil {
//				d.log.ErrorErr(err)
//			}
//			return
//		case err != nil:
//			d.log.ErrorErr(err)
//			return
//		}
//	}
//
//	upd := `update kzbot.corpslevel set level = $1,enddate = $2,hcorp = $3,percent = $4 where corpname = $5`
//	_, err = d.db.Exec(ctx, upd, l.Level, l.EndDate, l.HCorp, l.Percent, l.CorpName)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//}
//
//func (d *Db) ReadCorpLevel(CorpName string) (models.LevelCorps, error) {
//	ctx, cancel := d.GetContext()
//	defer cancel()
//	var n models.LevelCorps
//	err := d.db.QueryRow(ctx, "SELECT * FROM kzbot.corpslevel WHERE corpname = $1", CorpName).Scan(
//		&n.CorpName, &n.Level, &n.EndDate, &n.HCorp, &n.Percent, &n.LastUpdate, &n.Relic)
//	if err != nil {
//		return models.LevelCorps{}, err
//	}
//	return n, nil
//}

func (d *Db) InsertUpdateCorpLevel(l models.LevelCorps) {
	ctx, cancel := d.GetContext()
	defer cancel()

	_, err := d.ReadCorpLevelByCorpConf(l.CorpName)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			// Если записи нет, вставляем новую
			insert := `INSERT INTO ws.corpslevel (corpname, level, enddate, hcorp, percent, last_update, relic) 
					   VALUES ($1, $2, $3, $4, $5, $6, $7)`
			_, err = d.db.Exec(ctx, insert,
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
	_, err = d.db.Exec(ctx, upd,
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
func (d *Db) ReadCorpLevelByCorpConf(Corp string) (models.LevelCorps, error) {
	ctx, cancel := d.GetContext()
	defer cancel()

	var n models.LevelCorps
	var endDateStr, lastUpdateStr string // Временные переменные для строковых дат

	err := d.db.QueryRow(ctx, "SELECT corpname, level, enddate, hcorp, percent, last_update, relic FROM ws.corpslevel WHERE corpname = $1", Corp).
		Scan(&n.CorpName, &n.Level, &endDateStr, &n.HCorp, &n.Percent, &lastUpdateStr, &n.Relic)
	if err != nil {
		return models.LevelCorps{}, err
	}

	// Парсим строки в time.Time
	n.EndDate, err = time.Parse(time.RFC3339Nano, endDateStr)
	if err != nil {
		//return models.LevelCorps{}, fmt.Errorf("не удалось распарсить EndDate: %w", err)
	}

	n.LastUpdate, err = time.Parse(time.RFC3339Nano, lastUpdateStr)
	if err != nil {
		//return models.LevelCorps{}, fmt.Errorf("не удалось распарсить LastUpdate: %w", err)
	}

	return n, nil
}

func (d *Db) ReadCorpLevelAll() ([]models.LevelCorps, error) {
	ctx, cancel := d.GetContext()
	defer cancel()

	var nn []models.LevelCorps
	sel := "SELECT corpname, level, enddate, hcorp, percent, last_update, relic FROM ws.corpslevel"
	results, err := d.db.Query(ctx, sel)
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}
	defer results.Close()

	for results.Next() {
		var n models.LevelCorps
		var endDateStr, lastUpdateStr string // Временные переменные для строковых дат

		err = results.Scan(&n.CorpName, &n.Level, &endDateStr, &n.HCorp, &n.Percent, &lastUpdateStr, &n.Relic)
		if err != nil {
			d.log.ErrorErr(err)
			continue
		}

		// Парсим строки в time.Time
		n.EndDate, err = time.Parse(time.RFC3339Nano, endDateStr)
		if err != nil {
			d.log.ErrorErr(fmt.Errorf("не удалось распарсить EndDate: %w", err))
			continue
		}

		n.LastUpdate, err = time.Parse(time.RFC3339Nano, lastUpdateStr)
		if err != nil {
			d.log.ErrorErr(fmt.Errorf("не удалось распарсить LastUpdate: %w", err))
			continue
		}

		nn = append(nn, n)
	}

	return nn, nil
}
func (d *Db) ReadCorpsLevelAllOld() ([]models.LevelCorps, error) {
	ctx, cancel := d.GetContext()
	defer cancel()
	var nn []models.LevelCorps
	sel := "SELECT * FROM kzbot.corpslevel"
	results, err := d.db.Query(ctx, sel)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}
	for results.Next() {
		var n models.LevelCorps
		err = results.Scan(&n.CorpName, &n.Level, &n.EndDate, &n.HCorp, &n.Percent, &n.LastUpdate, &n.Relic)
		if err != nil {
			d.log.ErrorErr(err)
		}
		nn = append(nn, n)
	}
	return nn, nil
}
