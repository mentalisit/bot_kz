package postgres

import (
	"fmt"
	"rs/models"
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
	ctx, cancel := d.GetContext()
	defer cancel()
	insert := `INSERT INTO rs_bot.battles(eventid,corporation,name,level,points) VALUES ($1,$2,$3,$4,$5)`
	_, err := d.db.Exec(ctx, insert, b.EventId, b.CorpName, b.Name, b.Level, b.Points)
	if err != nil {
		return err
	}
	return nil
}

func (d *Db) BattlesGetAll(corpName string, event int) ([]models.PlayerStats, error) {
	ctx, cancel := d.GetContext()
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
		if err := rows.Scan(&ps.Player, &ps.Points, &ps.Runs, &ps.Level); err != nil {
			return nil, fmt.Errorf("ошибка чтения данных: %v", err)
		}
		stats = append(stats, ps)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка после чтения строк: %v", err)
	}

	return stats, nil
}

func (d *Db) ScoreboardParamsReadAll() []models.ScoreboardParams {
	ctx, cancel := d.GetContext()
	defer cancel()
	query := `SELECT name,webhookchannel,scorechannel FROM rs_bot.scoreboard`
	rows, err := d.db.Query(ctx, query)
	if err != nil {
		d.log.ErrorErr(err)
		return nil
	}
	defer rows.Close()

	var params []models.ScoreboardParams
	for rows.Next() {
		var ps models.ScoreboardParams
		if err = rows.Scan(&ps.Name, &ps.ChannelWebhook, &ps.ChannelScoreboard); err != nil {
			d.log.ErrorErr(err)
			return nil
		}
		params = append(params, ps)
	}

	if err := rows.Err(); err != nil {
		d.log.ErrorErr(err)
		return nil
	}

	return params
}

func (d *Db) BattlesTopGetAll(corpName string) ([]models.BattlesTop, error) {
	ctx, cancel := d.GetContext()
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

func (d *Db) IdentifyGetPoints() (ss []models.Identify, err error) {
	ctx, cancel := d.GetContext()
	defer cancel()

	query := `
	SELECT * 
	FROM rs_bot.identify
	WHERE points > 0 
	AND participants <> '';`

	rows, err := d.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %v", err)
	}
	defer rows.Close()

	var records []models.Identify
	for rows.Next() {
		var rec models.Identify
		if err := rows.Scan(&rec.Id, &rec.MID, &rec.SolarId, &rec.Author, &rec.Count, &rec.Participants, &rec.Points, &rec.StartTime); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании: %v", err)
		}
		records = append(records, rec)
	}

	return records, nil
}
