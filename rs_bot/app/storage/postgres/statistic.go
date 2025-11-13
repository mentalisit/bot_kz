package postgres

import (
	"fmt"
	"rs/models"
)

// StatisticGetName возвращает статистику игрока по каждой паре (event_id, level).
func (d *Db) StatisticGetName(name string) ([]models.Statistic, error) {
	var stat []models.Statistic

	ctx, cancel := d.getContext()
	defer cancel()

	query := `
        SELECT
            eventid,
            max(level),
            SUM(points) AS total_points,
            COUNT(*) AS total_runs
        FROM
            rs_bot.battles
        WHERE
            name = $1
        GROUP BY
            eventid
        ORDER BY
            eventid;
    `

	rows, err := d.db.Query(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	// 2. Обработка результатов
	for rows.Next() {
		var s models.Statistic

		// Порядок сканирования: EventId, Level, Points, Runs
		err = rows.Scan(
			&s.EventId,
			&s.Level,
			&s.Points,
			&s.Runs,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		stat = append(stat, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка во время итерации по результатам: %w", err)
	}

	return stat, nil
}

func (d *Db) GetBattleStats(corporation string, minRecords int) []*models.BattleStats {
	ctx, cancel := d.getContext()
	defer cancel()

	query := `
    WITH max_event AS (
        SELECT MAX(eventid) as max_eventid 
        FROM rs_bot.battles
    ),
    level_avg AS (
        SELECT 
            level,
            AVG(CAST(points AS NUMERIC)) as avg_points_for_level
        FROM rs_bot.battles
        GROUP BY level
    )
    SELECT 
        b.name,
        b.level,
        SUM(b.points) as Points_sum,
        COUNT(*) as records_count,
        ROUND(AVG(CAST(b.points AS NUMERIC)), 0) as average_points,
        ROUND(AVG(CAST(b.points AS NUMERIC)) / la.avg_points_for_level, 1) as Quality
    FROM rs_bot.battles b
    CROSS JOIN max_event me
    JOIN level_avg la ON b.level = la.level
    WHERE b.eventid = me.max_eventid 
        AND b.corporation = $1
    GROUP BY b.name, b.level, la.avg_points_for_level
    HAVING COUNT(*) >= $2
    ORDER BY Quality DESC;`

	results, err := d.db.Query(ctx, query, corporation, minRecords)
	if err != nil {
		d.log.ErrorErr(err)
		return nil
	}
	defer results.Close()

	var stats []*models.BattleStats

	for results.Next() {
		var s models.BattleStats
		err = results.Scan(
			&s.Name,
			&s.Level,
			&s.PointsSum,
			&s.RecordsCount,
			&s.AveragePoints,
			&s.Quality,
		)
		if err != nil {
			d.log.ErrorErr(err)
			continue
		}
		stats = append(stats, &s)
	}

	if len(stats) == 0 {
		return nil
	}

	return stats
}
