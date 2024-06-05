package postgres

import (
	"context"
	"fmt"
	"kz_bot/models"
)

func (d *Db) TopLevel(ctx context.Context, CorpName, lvlkz string) bool {
	var good = false
	sel := "SELECT name FROM kzbot.sborkz WHERE corpname=$1 AND active=1  AND (lvlkz = $2 OR lvlkz = $3) GROUP BY name LIMIT 50"
	results, err := d.db.Query(ctx, sel, CorpName, lvlkz, "d"+lvlkz)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
		return false
	}

	var name string
	for results.Next() {
		err = results.Scan(&name)
		if len(name) > 0 {
			good = true
			countNames, err1 := d.CountNumberNameActive1(ctx, lvlkz, CorpName, name)
			if err1 != nil {
				return false
			}
			countNamesd, _ := d.CountNumberNameActive1(ctx, "d"+lvlkz, CorpName, name)
			countNames = countNames + countNamesd

			insertTempTopEvent := `INSERT INTO kzbot.temptopevent(name,numkz,points) VALUES ($1,$2,$3)`
			_, err = d.db.Exec(ctx, insertTempTopEvent, name, countNames, 0)
			if err != nil {
				d.log.ErrorErr(err)
				return false
			}
		}
	}
	return good
}
func (d *Db) TopEventLevel(ctx context.Context, CorpName, lvlkz string, numEvent int) bool {
	var good = false
	sel := "SELECT name FROM kzbot.sborkz WHERE corpname=$1 AND active=1  AND lvlkz = $2 AND numberevent = $3 GROUP BY name LIMIT 50"
	results, err := d.db.Query(ctx, sel, CorpName, lvlkz, numEvent)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
		return false
	}
	var name string
	for results.Next() {
		err = results.Scan(&name)
		if len(name) > 0 {
			good = true
			var countNames int
			selC := "SELECT  COUNT(*) as count FROM kzbot.sborkz WHERE name = $1 AND corpname = $2 AND active = 1 AND numberevent = $3 AND lvlkz = $4"
			row := d.db.QueryRow(ctx, selC, name, CorpName, numEvent, lvlkz)
			err = row.Scan(&countNames)
			if err != nil {
				d.log.ErrorErr(err)
				return false
			}

			var points int
			selS := "SELECT  SUM(eventpoints) FROM kzbot.sborkz WHERE name = $1 AND corpname = $2 AND active = 1 AND numberevent = $3 AND lvlkz = $4"
			row4 := d.db.QueryRow(ctx, selS, name, CorpName, numEvent, lvlkz)
			err4 := row4.Scan(&points)
			if err4 != nil {
				d.log.ErrorErr(err)
				return false
			}

			insertTempTopEvent := `INSERT INTO kzbot.temptopevent(name,numkz,points) VALUES ($1,$2,$3)`
			_, err = d.db.Exec(ctx, insertTempTopEvent, name, countNames, points)
			if err != nil {
				d.log.ErrorErr(err)
				return false
			}
		}
	}
	return good
}
func (d *Db) TopTemp(ctx context.Context) string {
	number := 1
	var (
		name, message2    string
		numkz, id, points int
	)

	sel := "SELECT * FROM kzbot.temptopevent ORDER BY numkz DESC"
	results, err := d.db.Query(ctx, sel)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
		return ""
	}

	for results.Next() {
		results.Scan(&id, &name, &numkz, &points)
		message2 = message2 + fmt.Sprintf("%d. %s - %d \n", number, name, numkz)
		number = number + 1
	}

	_, err = d.db.Exec(ctx, "DELETE FROM kzbot.temptopevent")
	if err != nil {
		d.log.ErrorErr(err)
		return ""
	}
	return message2
}
func (d *Db) TopTempEvent(ctx context.Context) string {
	number := 1
	var (
		name, message2               string
		numkz, id, points, allpoints int
	)

	sel := "SELECT * FROM kzbot.temptopevent ORDER BY points DESC"
	results, err := d.db.Query(ctx, sel)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
		return ""
	}

	for results.Next() {
		results.Scan(&id, &name, &numkz, &points)
		message2 = message2 + fmt.Sprintf("%d. %s - %d (%d)\n", number, name, numkz, points)
		number = number + 1
		allpoints += points
	}

	_, err = d.db.Exec(ctx, "DELETE FROM kzbot.temptopevent")
	if err != nil {
		d.log.ErrorErr(err)
		return ""
	}
	return fmt.Sprintf("%s\nTotal: %d", message2, allpoints)
}
func (d *Db) TopAll(ctx context.Context, CorpName string) bool {
	good := false
	sel := "SELECT name FROM kzbot.sborkz WHERE corpname=$1 AND active>0 GROUP BY name LIMIT 50"
	results, err := d.db.Query(ctx, sel, CorpName)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
		return false
	}
	var name string
	for results.Next() {
		err = results.Scan(&name)
		if len(name) > 0 {
			good = true
			var countNames int
			selC := "SELECT COALESCE(SUM(active),0) FROM kzbot.sborkz WHERE corpname = $1 AND name = $2 AND active>0"
			//selC := "SELECT  COUNT(*) as count FROM kzbot.sborkz WHERE name = $1 AND corpname = $2 AND active = 1"
			row := d.db.QueryRow(ctx, selC, CorpName, name)
			err = row.Scan(&countNames)
			if err != nil {
				d.log.ErrorErr(err)
				return false
			}

			insertTempTopEvent := `INSERT INTO kzbot.temptopevent(name,numkz,points) VALUES ($1,$2,$3)`
			_, err = d.db.Exec(ctx, insertTempTopEvent, name, countNames, 0)
			if err != nil {
				d.log.ErrorErr(err)
				return false
			}
		}
	}
	return good
}
func (d *Db) TopAllEvent(ctx context.Context, CorpName string, numberevent int) bool {
	good := false

	//синтаксична помилка в або поблизу \"ASC\"
	sel := "SELECT name FROM kzbot.sborkz WHERE corpname=$1 AND numberevent = $2 AND active=1 GROUP BY name LIMIT 50"
	results, err := d.db.Query(ctx, sel, CorpName, numberevent)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
		return false
	}

	var name string
	for results.Next() {
		err = results.Scan(&name)
		if len(name) > 0 {
			good = true
			var countNames int
			selC := "SELECT  COUNT(*) as count FROM kzbot.sborkz WHERE name = $1 AND corpname = $2 AND active = 1 AND numberevent = $3"
			row := d.db.QueryRow(ctx, selC, name, CorpName, numberevent)
			err = row.Scan(&countNames)
			if err != nil {
				d.log.ErrorErr(err)
				continue
			}
			var points int
			selS := "SELECT  SUM(eventpoints) FROM kzbot.sborkz WHERE name = $1 AND corpname = $2 AND active = 1 AND numberevent = $3"
			row4 := d.db.QueryRow(ctx, selS, name, CorpName, numberevent)
			err4 := row4.Scan(&points)
			if err4 != nil {
				d.log.ErrorErr(err)
				continue
			}

			insertTempTopEvent := `INSERT INTO kzbot.temptopevent(name,numkz,points) VALUES ($1,$2,$3)`
			_, err = d.db.Exec(ctx, insertTempTopEvent, name, countNames, points)
			if err != nil {
				d.log.ErrorErr(err)
				continue
			}
		}
	}
	return good
}
func (d *Db) TopAllEventNew(ctx context.Context, CorpName string, numberevent int) []models.Top {
	sel := "SELECT name FROM kzbot.sborkz WHERE corpname=$1 AND numberevent = $2 AND active=1 GROUP BY name LIMIT 50"
	results, err := d.db.Query(ctx, sel, CorpName, numberevent)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
		return []models.Top{}
	}

	var topslice []models.Top
	var name string
	for results.Next() {
		err = results.Scan(&name)
		if len(name) > 0 {
			var countNames int
			selC := "SELECT  COUNT(*) as count FROM kzbot.sborkz WHERE name = $1 AND corpname = $2 AND active = 1 AND numberevent = $3"
			row := d.db.QueryRow(ctx, selC, name, CorpName, numberevent)
			err = row.Scan(&countNames)
			if err != nil {
				d.log.ErrorErr(err)
				return topslice
			}
			var points int
			selS := "SELECT  SUM(eventpoints) FROM kzbot.sborkz WHERE name = $1 AND corpname = $2 AND active = 1 AND numberevent = $3"
			row4 := d.db.QueryRow(ctx, selS, name, CorpName, numberevent)
			err4 := row4.Scan(&points)
			if err4 != nil {
				d.log.ErrorErr(err)
				return topslice
			}
			top := models.Top{
				Name:   name,
				Numkz:  countNames,
				Points: points,
			}

			topslice = append(topslice, top)
		}
	}
	return topslice
}

func (d *Db) TopAllPerMonth(ctx context.Context, CorpName string) bool {
	good := false
	sel := "SELECT name FROM kzbot.sborkz WHERE corpname=$1 AND active>0 AND to_timestamp(date,'YYYY-MM-DD') >= CURRENT_DATE - INTERVAL '30 days' GROUP BY name LIMIT 50"
	results, err := d.db.Query(ctx, sel, CorpName)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
		return false
	}
	var name string
	for results.Next() {
		err = results.Scan(&name)
		if len(name) > 0 {
			good = true
			var countNames int
			selC := "SELECT COALESCE(SUM(active),0) FROM kzbot.sborkz WHERE corpname = $1 AND name = $2 AND active>0 AND to_timestamp(date,'YYYY-MM-DD') >= CURRENT_DATE - INTERVAL '30 days'"
			row := d.db.QueryRow(ctx, selC, CorpName, name)
			err = row.Scan(&countNames)
			if err != nil {
				d.log.ErrorErr(err)
				return false
			}

			insertTempTopEvent := `INSERT INTO kzbot.temptopevent(name,numkz,points) VALUES ($1,$2,$3)`
			_, err = d.db.Exec(ctx, insertTempTopEvent, name, countNames, 0)
			if err != nil {
				d.log.ErrorErr(err)
				return false
			}
		}
	}
	return good
}
func (d *Db) TopLevelPerMonth(ctx context.Context, CorpName, lvlkz string) bool {
	var good = false
	sel := "SELECT name FROM kzbot.sborkz WHERE corpname=$1 AND active=1  AND (lvlkz = $2 OR lvlkz = $3) AND to_timestamp(date,'YYYY-MM-DD') >= CURRENT_DATE - INTERVAL '30 days' GROUP BY name LIMIT 50"
	results, err := d.db.Query(ctx, sel, CorpName, lvlkz, "d"+lvlkz)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
		return false
	}

	var name string
	for results.Next() {
		err = results.Scan(&name)
		if len(name) > 0 {
			good = true
			var countNumberNameActive1 int
			sel = "SELECT COALESCE(SUM(active),0) FROM kzbot.sborkz WHERE (lvlkz = $1 OR lvlkz = $2) AND corpname = $3 AND name = $4 AND to_timestamp(date,'YYYY-MM-DD') >= CURRENT_DATE - INTERVAL '30 days'"
			row := d.db.QueryRow(ctx, sel, lvlkz, "d"+lvlkz, CorpName, name)
			err = row.Scan(&countNumberNameActive1)
			if err != nil {
				d.log.ErrorErr(err)
			}
			countNames := countNumberNameActive1
			insertTempTopEvent := `INSERT INTO kzbot.temptopevent(name,numkz,points) VALUES ($1,$2,$3)`
			_, err = d.db.Exec(ctx, insertTempTopEvent, name, countNames, 0)
			if err != nil {
				d.log.ErrorErr(err)
				return false
			}
		}
	}
	return good
}
