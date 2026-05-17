package postgres

import (
	"database/sql"
	"time"

	"errors"

	"github.com/mentalisit/restapi/models"
)

// CreateCorpInfo создает новую запись о корпорации
func (d *Db) CreateCorpInfo(corp models.CorpInfo) (int64, error) {

	query := `INSERT INTO ws.corps_info(corp_name, corp_id, level, xp, webhook, last_win, date_ended, last_update) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	var id int64
	err := d.db.QueryRow(query,
		corp.CorpName,
		corp.CorpID,
		corp.Level,
		corp.XP,
		corp.Webhook,
		corp.LastWin,
		corp.DateEnded,
		time.Now().UTC()).Scan(&id)

	if err != nil {
		d.log.ErrorErr(err)
		return 0, err
	}

	return id, nil
}

// ReadCorpInfoByCorpID читает запись о корпорации по corp_id
func (d *Db) ReadCorpInfoByCorpID(corpID string) (*models.CorpInfo, error) {

	query := `SELECT id, corp_name, corp_id, level, xp, webhook, last_win, date_ended, last_update 
			  FROM ws.corps_info WHERE corp_id = $1`

	var corp models.CorpInfo
	err := d.db.QueryRow(query, corpID).Scan(
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

// UpdateCorpInfo обновляет запись о корпорации
func (d *Db) UpdateCorpInfo(corp models.CorpInfo) error {

	query := `UPDATE ws.corps_info 
			  SET corp_name = $2, corp_id = $3, level = $4, xp = $5, webhook = $6, last_win = $7, date_ended = $8, last_update = $9 
			  WHERE id = $1`

	_, err := d.db.Exec(query,
		corp.ID,
		corp.CorpName,
		corp.CorpID,
		corp.Level,
		corp.XP,
		corp.Webhook,
		corp.LastWin,
		corp.DateEnded,
		time.Now().UTC(),
	)

	if err != nil {
		d.log.ErrorErr(err)
		return err
	}

	return nil
}
