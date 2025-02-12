package postgres

import (
	"discord/models"
	"github.com/jackc/pgx/v5"
)

func (d *Db) ReadTop5Level(corpname string) []string {
	ctx, cancel := d.getContext()
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

func (d *Db) ReadMesIdDS(mesid string) (string, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	sel := "SELECT lvlkz FROM kzbot.sborkz WHERE dsmesid = $1 AND active = 0"
	results, err := d.db.Query(ctx, sel, mesid)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	a := []string{}
	var dsmesid string
	for results.Next() {
		var t models.Sborkz
		err = results.Scan(&t.Lvlkz)
		a = append(a, t.Lvlkz)
	}
	a = RemoveDuplicates(a)
	if len(a) > 0 {
		dsmesid = a[0]
		return dsmesid, nil
	} else {
		return "", err
	}
}
func RemoveDuplicates[T comparable](a []T) []T {
	result := make([]T, 0, len(a))
	temp := map[T]struct{}{}
	for _, item := range a {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func (d *Db) NumActiveEvent(CorpName string) (event1 int) { //запрос номера ивента
	ctx, cancel := d.getContext()
	defer cancel()
	sel := "SELECT numevent FROM kzbot.rsevent WHERE corpname=$1 AND activeevent=1 ORDER BY numevent DESC LIMIT 1"
	row := d.db.QueryRow(ctx, sel, CorpName)
	err := row.Scan(&event1)
	if err != nil {
		if err == pgx.ErrNoRows {
			event1 = 0
		} else {
			d.log.ErrorErr(err)
		}
	}
	return event1
}
