package postgres

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"rs/models"
	"time"
)

func (d *Db) OptimizationSborkz() {
	ctx, cancel := d.GetContext()
	defer cancel()
	// Подсчет активных записей и сортировка по имени
	query := `SELECT mention,corpname,lvlkz, SUM(active) AS active_sum FROM kzbot.sborkz GROUP BY corpname, mention,lvlkz ORDER BY mention`
	rows, err := d.db.Query(ctx, query)
	defer rows.Close()
	if err != nil {
		d.log.Info(err.Error())
		return
	}
	for rows.Next() {
		var mention string
		var activeCount int
		var corpname string
		var level string
		if err := rows.Scan(&mention, &corpname, &level, &activeCount); err != nil {
			d.log.Info(err.Error())
			return
		}
		var countNames int
		sel := "SELECT  COUNT(*) as count FROM kzbot.sborkz WHERE mention = $1 AND lvlkz = $2 AND corpname = $3 AND active > 0"
		row := d.db.QueryRow(ctx, sel, mention, level, corpname)
		err := row.Scan(&countNames)
		if err != nil {
			d.log.Info(err.Error())
			return
		}
		if countNames > 5 {
			sel := "SELECT * FROM kzbot.sborkz WHERE lvlkz = $1 AND corpname = $2 AND mention = $3"
			results, err := d.db.Query(ctx, sel, level, corpname, mention)
			defer results.Close()
			if err != nil {
				d.log.ErrorErr(err)
			}
			var t models.Sborkz
			for results.Next() {

				err = results.Scan(&t.Id, &t.Corpname, &t.Name, &t.Mention, &t.Tip, &t.Dsmesid,
					&t.Tgmesid, &t.Wamesid, &t.Time, &t.Date, &t.Lvlkz, &t.Numkzn, &t.Numberkz,
					&t.Numberevent, &t.Eventpoints, &t.Active, &t.Timedown, &t.UserId)
			}
			del := "delete from kzbot.sborkz where mention = $1 and corpname = $2 and lvlkz = $3"
			_, err = d.db.Exec(ctx, del, mention, corpname, level)
			if err != nil {
				d.log.ErrorErr(err)
			}
			tm := time.Now()
			mdate := (tm.Format("2006-01-02"))
			mtime := (tm.Format("15:04"))
			insertSborkztg1 := `INSERT INTO kzbot.sborkz(corpname,name,mention,tip,dsmesid,tgmesid,wamesid,time,date,lvlkz,
		          numkzn,numberkz,numberevent,eventpoints,active,timedown)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`
			_, err = d.db.Exec(ctx, insertSborkztg1, t.Corpname, t.Name, t.Mention, t.Tip, t.Dsmesid, t.Tgmesid,
				t.Wamesid, mtime, mdate, t.Lvlkz, t.Numkzn, t.Numberkz, t.Numberevent, t.Eventpoints, activeCount, t.Timedown)
			if err != nil {
				d.log.ErrorErr(err)
			}
			d.log.Info(fmt.Sprintf("Выполнено сжатие данных игрока %s в корпорации %s кз%s изза %d записей", t.Name, t.Corpname, level, activeCount))
			time.Sleep(1 * time.Second)
		}
	}

	if err := rows.Err(); err != nil {
		d.log.Info(err.Error())
	}
}
func (d *Db) СountName(userid, lvlkz, corpName string) (int, error) {
	ctx, cancel := d.GetContext()
	defer cancel()

	var countNames int
	sel := "SELECT  COUNT(*) as count FROM kzbot.sborkz WHERE userid = $1 AND lvlkz = $2 AND corpname = $3 AND active = 0"
	row := d.db.QueryRow(ctx, sel, userid, lvlkz, corpName)
	err := row.Scan(&countNames)
	if err != nil {
		d.log.ErrorErr(err)
		return 0, err
	}

	return countNames, nil
}
func (d *Db) CountQueue(lvlkz, CorpName string) (int, error) { //проверка сколько игровок в очереди
	ctx, cancel := d.GetContext()
	defer cancel()

	var count int
	sel := "SELECT  COUNT(*) as count FROM kzbot.sborkz WHERE lvlkz = $1 AND corpname = $2 AND active = 0"
	row := d.db.QueryRow(ctx, sel, lvlkz, CorpName)
	err := row.Scan(&count)
	if err != nil {
		d.log.ErrorErr(err)
		return 0, err
	}

	return count, nil
}
func (d *Db) CountNumberNameActive1(lvlkz, CorpName, userid string) (int, error) { // выковыриваем из базы значение количества походов на кз
	ctx, cancel := d.GetContext()
	defer cancel()

	var countNumberNameActive1 int
	sel := "SELECT COALESCE(SUM(active),0) FROM kzbot.sborkz WHERE lvlkz = $1 AND corpname = $2 AND userid = $3"
	//COALESCE(SUM(value), 0)
	row := d.db.QueryRow(ctx, sel, lvlkz, CorpName, userid)
	err := row.Scan(&countNumberNameActive1)
	if err != nil {
		d.log.ErrorErr(err)
		return 0, err
	}
	return countNumberNameActive1, nil
}

func (d *Db) CountNameQueue(userid string) (countNames int) { //проверяем есть ли игрок в других очередях
	ctx, cancel := d.GetContext()
	defer cancel()
	sel := "SELECT  COUNT(*) as count FROM kzbot.sborkz WHERE userid = $1 AND active = 0"
	row := d.db.QueryRow(ctx, sel, userid)
	err := row.Scan(&countNames)
	if err != nil {
		d.log.ErrorErr(err)
	}

	return countNames
}
func (d *Db) CountNameQueueCorp(userid, corp string) (countNames int) { //проверяем есть ли игрок в других очередях
	ctx, cancel := d.GetContext()
	defer cancel()
	sel := "SELECT  COUNT(*) as count FROM kzbot.sborkz WHERE userid = $1 AND corpname = $2 AND active = 0"
	row := d.db.QueryRow(ctx, sel, userid, corp)
	err := row.Scan(&countNames)
	if err != nil {
		d.log.ErrorErr(err)
		return 0
	}

	return countNames
}
func (d *Db) ReadTop5Level(corpname string) []string {
	ctx, cancel := d.GetContext()
	defer cancel()
	query := `
        SELECT lvlkz, COUNT(*) AS lvlkz_count
        FROM kzbot.sborkz
        WHERE corpname=$1
          AND date::timestamp >= CURRENT_DATE - INTERVAL '5 days'
        GROUP BY lvlkz
        ORDER BY lvlkz_count DESC
        LIMIT 5;
    `

	// Выполнение запроса
	rows, err := d.db.Query(ctx, query, corpname)
	defer rows.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}

	var levels []string

	// Итерация по результатам запроса
	for rows.Next() {
		var lvlkz string
		var lvlkzCount int
		if err = rows.Scan(&lvlkz, &lvlkzCount); err != nil {
			d.log.ErrorErr(err)
		}
		levels = append(levels, lvlkz)
	}
	if err = rows.Err(); err != nil {
		d.log.ErrorErr(err)
	}
	return levels
}

func (d *Db) CountQueueNumberNameActive1QueueLvl(lvlkz, CorpName, userid string) (countQueue, countNumberName, NumRsLevel int, errorsAll error) {
	ctx, cancel := d.GetContext()
	defer cancel()

	sel := "SELECT  COUNT(*) as count FROM kzbot.sborkz WHERE lvlkz = $1 AND corpname = $2 AND active = 0"
	row := d.db.QueryRow(ctx, sel, lvlkz, CorpName)
	err := row.Scan(&countQueue)
	if err != nil {
		errorsAll = fmt.Errorf("%+v\n", err)
	}

	ctx, cancel = d.GetContext()
	defer cancel()

	sel = "SELECT COALESCE(SUM(active),0) FROM kzbot.sborkz WHERE lvlkz = $1 AND corpname = $2 AND userid = $3"
	row = d.db.QueryRow(ctx, sel, lvlkz, CorpName, userid)
	err = row.Scan(&countNumberName)
	if err != nil {
		errorsAll = fmt.Errorf("%+v\n%+v\n", errorsAll, err)
	}

	ctx, cancel = d.GetContext()
	defer cancel()

	sel = "SELECT  number FROM kzbot.numkz WHERE lvlkz = $1 AND corpname = $2"
	row = d.db.QueryRow(ctx, sel, lvlkz, CorpName)
	err = row.Scan(&NumRsLevel)
	if err != nil {
		NumRsLevel = 0
		if errors.Is(err, pgx.ErrNoRows) {
			insertSmt := "INSERT INTO kzbot.numkz(lvlkz, number,corpname) VALUES ($1,$2,$3)"
			_, err = d.db.Exec(ctx, insertSmt, lvlkz, NumRsLevel, CorpName)
			if err != nil {
				errorsAll = fmt.Errorf("%+v\n%+v\n", errorsAll, err)
			}
			NumRsLevel = NumRsLevel + 1
		} else {
			errorsAll = fmt.Errorf("%+v\n%+v\n", errorsAll, err)
		}
	} else {
		NumRsLevel++
	}
	return countQueue, countNumberName, NumRsLevel, errorsAll
}
