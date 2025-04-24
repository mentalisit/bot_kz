package dbpostgres

import (
	"fmt"
	"time"
	"ws/models"
)

func (d *Db) ReadCorpsLevelAll() ([]models.LevelCorps, error) {
	ctx, cancel := d.getContext()
	defer cancel()

	var nn []models.LevelCorps
	sel := "SELECT corpname, level, enddate, hcorp, percent, last_update, relic FROM ws.corpslevel"
	results, err := d.pool.Query(ctx, sel)
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

func (d *Db) ReadCorpsLevel(hCorp string) (models.LevelCorps, error) {
	ctx, cancel := d.getContext()
	defer cancel()

	var n models.LevelCorps
	var endDateStr, lastUpdateStr string // Временные переменные для строковых дат

	err := d.pool.QueryRow(ctx, "SELECT corpname, level, enddate, hcorp, percent, last_update, relic FROM ws.corpslevel WHERE hcorp = $1", hCorp).
		Scan(&n.CorpName, &n.Level, &endDateStr, &n.HCorp, &n.Percent, &lastUpdateStr, &n.Relic)
	if err != nil {
		return models.LevelCorps{}, err
	}

	// Парсим строки в time.Time
	n.EndDate, err = time.Parse(time.RFC3339Nano, endDateStr)
	if err != nil {
		return models.LevelCorps{}, fmt.Errorf("не удалось распарсить EndDate: %w", err)
	}

	n.LastUpdate, err = time.Parse(time.RFC3339Nano, lastUpdateStr)
	if err != nil {
		return models.LevelCorps{}, fmt.Errorf("не удалось распарсить LastUpdate: %w", err)
	}

	return n, nil
}
