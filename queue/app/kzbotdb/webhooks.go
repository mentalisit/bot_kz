package kzbotdb

import (
	"context"
	"fmt"
	"queue/models"
	"time"
)

// RedStarStarted , RedStarEnded ...
func (d *Db) GetWebhooksByEventTypeAndAfterTs(eventType string, afterTs int64) ([]models.Webhook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, tsunix, corp, message
		FROM rs_bot.webhooks
		WHERE tsunix > $1
		  AND message::text LIKE $2
	`
	searchPattern := fmt.Sprintf("%%%q%%", eventType) // "%\"RedStarEnded\"%"

	rows, err := d.db.Query(ctx, query, afterTs, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var results []models.Webhook

	for rows.Next() {
		var w models.Webhook
		err = rows.Scan(&w.ID, &w.TsUnix, &w.Corp, &w.Message)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		results = append(results, w)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return results, nil
}

func (d *Db) GetWebhooksByAfter(afterTs int64) ([]models.Webhook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := d.db.Query(ctx, `
		SELECT id, tsunix, corp, message
		FROM rs_bot.webhooks
		WHERE tsunix > $1
	`, afterTs)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var results []models.Webhook

	for rows.Next() {
		var w models.Webhook
		err = rows.Scan(&w.ID, &w.TsUnix, &w.Corp, &w.Message)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		results = append(results, w)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return results, nil
}

func (d *Db) BattlesGetSeasonAll(season int) ([]models.BattlesEvent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Добавляем COUNT(*) в SELECT
	query := `
        SELECT 
            name, 
            MAX(level), 
            SUM(points), 
            COUNT(*) as entries_count
        FROM rs_bot.battles 
        WHERE eventid = $1 
          AND corporation IN ('IX Легион', 'русь ', 'best')
        GROUP BY name 
        ORDER BY SUM(points) DESC`

	rows, err := d.db.Query(ctx, query, season)
	if err != nil {
		return nil, fmt.Errorf("ошибка SQL: %v", err)
	}
	defer rows.Close()

	var stats []models.BattlesEvent
	for rows.Next() {
		var ps models.BattlesEvent
		// Важно: порядок в Scan должен совпадать с порядком в SELECT
		if err = rows.Scan(&ps.Name, &ps.LevelMax, &ps.Points, &ps.EntriesCount); err != nil {
			return nil, fmt.Errorf("ошибка Scan: %v", err)
		}
		ps.EventId = season
		stats = append(stats, ps)
	}

	return stats, nil
}
