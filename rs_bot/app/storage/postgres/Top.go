package postgres

import (
	"fmt"
	"rs/models"
	"sort"
)

func (d *Db) TopEventLevelNew(CorpName, lvlkz string, numEvent int) []models.Top {
	ctx, cancel := d.GetContext()
	defer cancel()
	var top []models.Top
	sel := "SELECT mention FROM kzbot.sborkz WHERE corpname=$1 AND active=1  AND (lvlkz = $2 OR lvlkz = $3) AND numberevent = $4 GROUP BY mention LIMIT 50"
	results, err := d.db.Query(ctx, sel, CorpName, lvlkz, "d"+lvlkz, numEvent)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
		return top
	}

	for results.Next() {
		var u models.Top
		err = results.Scan(&u.Name)
		if len(u.Name) > 0 {
			selC := "SELECT  COUNT(*) as count FROM kzbot.sborkz WHERE mention = $1 AND corpname = $2 AND active = 1 AND numberevent = $3 AND (lvlkz = $4 OR lvlkz = $5)"
			row := d.db.QueryRow(ctx, selC, u.Name, CorpName, numEvent, lvlkz, "d"+lvlkz)
			err = row.Scan(&u.Numkz)
			if err != nil {
				d.log.ErrorErr(err)
				return top
			}

			selS := "SELECT  SUM(eventpoints) FROM kzbot.sborkz WHERE mention = $1 AND corpname = $2 AND active = 1 AND numberevent = $3 AND (lvlkz = $4 OR lvlkz = $5)"
			row4 := d.db.QueryRow(ctx, selS, u.Name, CorpName, numEvent, lvlkz, "d"+lvlkz)
			err4 := row4.Scan(&u.Points)
			if err4 != nil {
				d.log.ErrorErr(err)
				return top
			}
			top = append(top, u)
		}
	}
	sort.Slice(top, func(i, j int) bool {
		return top[i].Points > top[j].Points
	})
	return top
}

func (d *Db) TopAllEventNew(CorpName string, numberevent int) (top []models.Top) {
	ctx, cancel := d.GetContext()
	defer cancel()
	sel := "SELECT mention FROM kzbot.sborkz WHERE corpname=$1 AND numberevent = $2 AND active=1 GROUP BY mention"
	results, err := d.db.Query(ctx, sel, CorpName, numberevent)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
		return
	}

	for results.Next() {
		var u models.Top
		err = results.Scan(&u.Name)
		if len(u.Name) > 0 {
			selC := "SELECT  COUNT(*) as count FROM kzbot.sborkz WHERE mention = $1 AND corpname = $2 AND active = 1 AND numberevent = $3"
			row := d.db.QueryRow(ctx, selC, u.Name, CorpName, numberevent)
			err = row.Scan(&u.Numkz)
			if err != nil {
				d.log.ErrorErr(err)
				continue
			}
			selS := "SELECT  SUM(eventpoints) FROM kzbot.sborkz WHERE mention = $1 AND corpname = $2 AND active = 1 AND numberevent = $3"
			row4 := d.db.QueryRow(ctx, selS, u.Name, CorpName, numberevent)
			err4 := row4.Scan(&u.Points)
			if err4 != nil {
				d.log.ErrorErr(err)
				continue
			}
			top = append(top, u)
		}
	}
	sort.Slice(top, func(i, j int) bool {
		return top[i].Points > top[j].Points
	})
	return top
}

func (d *Db) TopAllPerMonthNew(CorpName string) (top []models.Top) {
	ctx, cancel := d.GetContext()
	defer cancel()
	sel := "SELECT name FROM kzbot.sborkz WHERE corpname=$1 AND active>0 AND to_timestamp(date,'YYYY-MM-DD') >= CURRENT_DATE - INTERVAL '30 days' GROUP BY name LIMIT 50"
	results, err := d.db.Query(ctx, sel, CorpName)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
	for results.Next() {
		var u models.Top
		err = results.Scan(&u.Name)
		if len(u.Name) > 0 {
			selC := "SELECT COALESCE(SUM(active),0) FROM kzbot.sborkz WHERE corpname = $1 AND name = $2 AND active>0 AND to_timestamp(date,'YYYY-MM-DD') >= CURRENT_DATE - INTERVAL '30 days'"
			row := d.db.QueryRow(ctx, selC, CorpName, u.Name)
			err = row.Scan(&u.Numkz)
			if err != nil {
				d.log.ErrorErr(err)
				return
			}
			top = append(top, u)
		}
	}
	sort.Slice(top, func(i, j int) bool {
		return top[i].Numkz > top[j].Numkz
	})
	return top
}

func (d *Db) TopLevelPerMonthNew(CorpName, lvlkz string) []models.Top {
	ctx, cancel := d.GetContext()
	defer cancel()
	var top []models.Top
	sel := "SELECT name FROM kzbot.sborkz WHERE corpname=$1 AND active=1  AND (lvlkz = $2 OR lvlkz = $3) AND to_timestamp(date,'YYYY-MM-DD') >= CURRENT_DATE - INTERVAL '30 days' GROUP BY name LIMIT 50"
	results, err := d.db.Query(ctx, sel, CorpName, lvlkz, "d"+lvlkz)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
		return top
	}
	for results.Next() {
		var u models.Top
		err = results.Scan(&u.Name)
		if len(u.Name) > 0 {
			sel = "SELECT COALESCE(SUM(active),0) FROM kzbot.sborkz WHERE (lvlkz = $1 OR lvlkz = $2) AND corpname = $3 AND name = $4 AND to_timestamp(date,'YYYY-MM-DD') >= CURRENT_DATE - INTERVAL '30 days'"
			row := d.db.QueryRow(ctx, sel, lvlkz, "d"+lvlkz, CorpName, u.Name)
			err = row.Scan(&u.Numkz)
			if err != nil {
				d.log.ErrorErr(err)
			}
		}
		top = append(top, u)
	}
	sort.Slice(top, func(i, j int) bool {
		return top[i].Numkz > top[j].Numkz
	})
	return top
}

func (d *Db) RedStarFightGetStar() (ss []models.RedStarFight, err error) {
	ctx, cancel := d.GetContext()
	defer cancel()

	query := `
		SELECT * FROM rs_bot.redstarfight`

	rows, err := d.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %v", err)
	}
	defer rows.Close()

	var records []models.RedStarFight
	for rows.Next() {
		var rec models.RedStarFight
		if err = rows.Scan(&rec.Id, &rec.GameMId, &rec.SolarId, &rec.SendId, &rec.Author, &rec.Level, &rec.Count, &rec.Participants,
			&rec.Points, &rec.EventId, &rec.StartTime, &rec.ClientId); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании: %v", err)
		}
		records = append(records, rec)
	}

	return records, nil
}
