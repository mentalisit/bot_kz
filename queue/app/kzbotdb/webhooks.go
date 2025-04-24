package kzbotdb

import (
	"context"
	"fmt"
	"queue/models"
	"time"
)

// GetWebhooksByCorp возвращает все записи по значению corp
func (d *Db) GetWebhooksByCorp(corp string) ([]models.Webhook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := d.db.Query(ctx, `
		SELECT id, tsunix, corp, message
		FROM rs_bot.webhooks
		WHERE corp = $1
	`, corp)
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

// GetWebhooksByCorpAfter возвращает все записи по corp, у которых tsunix > указанного значения
func (d *Db) GetWebhooksByCorpAfter(corp string, afterTs int64) ([]models.Webhook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := d.db.Query(ctx, `
		SELECT id, tsunix, corp, message
		FROM rs_bot.webhooks
		WHERE corp = $1 AND tsunix > $2
	`, corp, afterTs)
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
