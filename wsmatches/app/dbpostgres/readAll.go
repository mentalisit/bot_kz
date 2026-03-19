package dbpostgres

import (
	"database/sql"
	"errors"
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

// GetAllCorpInfo возвращает все записи о корпорациях
func (d *Db) GetAllCorpInfo() ([]models.CorpInfo, error) {
	ctx, cancel := d.getContext()
	defer cancel()

	query := `SELECT id, corp_name, corp_id, level, xp, webhook, last_win, date_ended, last_update 
			  FROM ws.corps_info ORDER BY corp_name`

	rows, err := d.pool.Query(ctx, query)
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}
	defer rows.Close()

	var corps []models.CorpInfo
	for rows.Next() {
		var corp models.CorpInfo
		err := rows.Scan(
			&corp.ID,
			&corp.CorpName,
			&corp.CorpID,
			&corp.Level,
			&corp.XP,
			&corp.Webhook,
			&corp.LastWin,
			&corp.DateEnded,
			&corp.LastUpdate,
		)
		if err != nil {
			d.log.ErrorErr(err)
			return nil, err
		}
		corps = append(corps, corp)
	}

	if err = rows.Err(); err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}

	return corps, nil
}

// ReadCorpInfoByCorpID читает запись о корпорации по corp_id
func (d *Db) ReadCorpInfoByCorpID(corpID string) (*models.CorpInfo, error) {
	ctx, cancel := d.getContext()
	defer cancel()

	query := `SELECT id, corp_name, corp_id, level, xp, webhook, last_win, date_ended, last_update 
			  FROM ws.corps_info WHERE corp_id = $1`

	var corp models.CorpInfo
	err := d.pool.QueryRow(ctx, query, corpID).Scan(
		&corp.ID,
		&corp.CorpName,
		&corp.CorpID,
		&corp.Level,
		&corp.XP,
		&corp.Webhook,
		&corp.LastWin,
		&corp.DateEnded,
		&corp.LastUpdate,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		d.log.ErrorErr(err)
		return nil, err
	}

	return &corp, nil
}
