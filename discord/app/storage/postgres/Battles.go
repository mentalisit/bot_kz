package postgres

import (
	"discord/models"
	"fmt"
)

/*
`
		CREATE TABLE IF NOT EXISTS rs_bot.battles
	(
		id     bigserial        primary key,
		eventId integer NOT NULL DEFAULT 0,
		corporation text NOT NULL DEFAULT '',
		name text NOT NULL DEFAULT '',
		level    integer NOT NULL DEFAULT 0,
		points   integer NOT NULL DEFAULT 0
	);`
*/

func (d *Db) BattlesInsert(b models.Battles) error {
	ctx, cancel := d.getContext()
	defer cancel()
	insert := `INSERT INTO rs_bot.battles(eventid,corporation,name,level,points) VALUES ($1,$2,$3,$4,$5)`
	_, err := d.db.Exec(ctx, insert, b.EventId, b.CorpName, b.Name, b.Level, b.Points)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) BattlesGetAll(corpName string, event int) ([]models.PlayerStats, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	query := `
		SELECT name,
		       SUM(points) AS total_points, 
		       COUNT(*) AS runs,
		       MAX(level) AS max_level		
		FROM rs_bot.battles 
		where eventid=$1 AND corporation=$2
		GROUP BY name;
	`

	rows, err := d.db.Query(ctx, query, event, corpName)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	var stats []models.PlayerStats
	for rows.Next() {
		var ps models.PlayerStats
		if err = rows.Scan(&ps.Player, &ps.Points, &ps.Runs, &ps.Level); err != nil {
			return nil, fmt.Errorf("ошибка чтения данных: %v", err)
		}
		stats = append(stats, ps)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка после чтения строк: %v", err)
	}

	return stats, nil
}
func (d *Db) BattlesUpdate(b models.Battles) error {
	ctx, cancel := d.getContext()
	defer cancel()
	sqlUpd := "update rs_bot.battles set points = $1 where name = $2 AND level = $3 AND eventid = $4 AND corporation = $5"
	_, err := d.db.Exec(ctx, sqlUpd, b.Points, b.Name, b.Level, b.EventId, b.CorpName)
	if err != nil {
		return err
	}
	return nil
}

func (d *Db) BattlesTopInsert(b models.BattlesTop) error {
	topGet, _ := d.BattlesTopGet(b)
	if topGet.Name == b.Name {
		topGet.Count++
		if b.Level > topGet.Level {
			topGet.Level = b.Level
		}
		err := d.BattlesTopUpdate(topGet)
		if err != nil {
			return err
		}
		return nil
	}

	ctx, cancel := d.getContext()
	defer cancel()
	insert := `INSERT INTO rs_bot.battlestop(corporation,name,level,count) VALUES ($1,$2,$3,$4)`
	_, err := d.db.Exec(ctx, insert, b.CorpName, b.Name, b.Level, 1)
	if err != nil {
		return err
	}
	return nil
}

func (d *Db) BattlesTopUpdate(b models.BattlesTop) error {
	ctx, cancel := d.getContext()
	defer cancel()
	sqlUpd := "update rs_bot.battlestop set count = $1, level = $2 where id = $3 "
	_, err := d.db.Exec(ctx, sqlUpd, b.Count, b.Level, b.Id)
	if err != nil {
		return err
	}
	return nil
}

func (d *Db) BattlesTopGetAll(corpName string) ([]models.BattlesTop, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	query := `SELECT * FROM rs_bot.battlestop where corporation=$1 `
	rows, err := d.db.Query(ctx, query, corpName)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	var stats []models.BattlesTop
	for rows.Next() {
		var ps models.BattlesTop
		if err = rows.Scan(&ps.Id, &ps.CorpName, &ps.Name, &ps.Level, &ps.Count); err != nil {
			return nil, fmt.Errorf("ошибка чтения данных: %v", err)
		}
		stats = append(stats, ps)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка после чтения строк: %v", err)
	}

	return stats, nil
}

func (d *Db) BattlesTopGet(b models.BattlesTop) (models.BattlesTop, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	query := `SELECT * FROM rs_bot.battlestop where name=$1 AND corporation=$2`
	row := d.db.QueryRow(ctx, query, b.Name, b.CorpName)

	var ps models.BattlesTop
	err := row.Scan(&ps.Id, &ps.CorpName, &ps.Name, &ps.Level, &ps.Count)
	if err != nil {
		return models.BattlesTop{}, fmt.Errorf("ошибка чтения данных: %v", err)
	}

	return ps, nil
}
