package postgresV2

import (
	"database/sql"
	"errors"
	"rs/models"
)

// InsertActiveQueue inserts a new active queue entry with JSONB data
func (d *Db) InsertActiveQueue(p models.QueueActive) error {

	insertQuery := `INSERT INTO rs_bot.queue_active(data, remaining_time) VALUES ($1, $2)`
	_, err := d.db.Exec(insertQuery, p.JsonMarshalData(), p.RemainingTime)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}

	return nil
}

func (d *Db) GetActiveQueue() ([]models.QueueActive, error) {

	query := `SELECT id, data, remaining_time FROM rs_bot.queue_active`
	rows, err := d.db.Query(query)
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}
	defer rows.Close()

	var result []models.QueueActive
	for rows.Next() {
		var qa models.QueueActive
		var dataJSON []byte

		err = rows.Scan(&qa.ID, &dataJSON, &qa.RemainingTime)
		if err != nil {
			d.log.ErrorErr(err)
			continue
		}
		qa.JsonUnmarshalData(dataJSON)
		result = append(result, qa)
	}

	return result, nil
}

// GetActiveQueueByCorpAndLevel retrieves active queue entries for a specific corporation and level
func (d *Db) GetActiveQueueByCorpAndLevel(corpName, lvlRS string) ([]models.QueueActive, error) {

	query := `SELECT id, data, remaining_time FROM rs_bot.queue_active 
              WHERE data->>'corporation_uuid' = $1 AND data->>'lvl_rs' = $2`
	rows, err := d.db.Query(query, corpName, lvlRS)
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}
	defer rows.Close()

	var result []models.QueueActive
	for rows.Next() {
		var qa models.QueueActive
		var dataJSON []byte

		err := rows.Scan(&qa.ID, &dataJSON, &qa.RemainingTime)
		if err != nil {
			d.log.ErrorErr(err)
			continue
		}

		qa.JsonUnmarshalData(dataJSON)

		result = append(result, qa)
	}

	return result, nil
}

// UpdateActiveQueueRemainingTime updates the remaining time for a queue entry
func (d *Db) UpdateActiveQueueRemainingTime(id int64, remainingTime int64) error {

	query := `UPDATE rs_bot.queue_active SET remaining_time = $1 WHERE id = $2`
	_, err := d.db.Exec(query, remainingTime, id)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}

	return nil
}

// CountActiveQueueByUser counts active queue entries for a user by corporation
func (d *Db) CountActiveQueueByUser(userId, corpUuid string) (int, error) {

	query := `SELECT COUNT(*) FROM rs_bot.queue_active WHERE data->>'user_id' = $1 AND data->>'corporation_uuid' = $2`
	var count int
	err := d.db.QueryRow(query, userId, corpUuid).Scan(&count)
	if err != nil {
		d.log.ErrorErr(err)
		return 0, err
	}

	return count, nil
}

// GetActiveQueueByUserAndCorp gets active queue entries for a user by corporation
func (d *Db) GetActiveQueueByUserAndCorp(userId, corpUuid string) ([]models.QueueActive, error) {

	query := `SELECT id, data, remaining_time FROM rs_bot.queue_active 
			  WHERE data->>'user_id' = $1 AND data->>'corporation_uuid' = $2`
	rows, err := d.db.Query(query, userId, corpUuid)
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}
	defer rows.Close()

	var result []models.QueueActive
	for rows.Next() {
		var qa models.QueueActive
		var dataJSON []byte

		err = rows.Scan(&qa.ID, &dataJSON, &qa.RemainingTime)
		if err != nil {
			d.log.ErrorErr(err)
			continue
		}
		qa.JsonUnmarshalData(dataJSON)
		result = append(result, qa)
	}

	return result, nil
}

// InsertCompleteQueue inserts a completed queue entry
func (d *Db) InsertCompleteQueue(qa models.QueueActive) error {

	query := `INSERT INTO rs_bot.queue_complete (data) VALUES ($1)`
	_, err := d.db.Exec(query, qa.JsonMarshalData())
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}

	return nil
}

// MinusMin decrements remaining_time by 1 for all active queues and returns all records
func (d *Db) MinusMin() ([]models.QueueActive, error) {

	// First, decrement remaining_time by 1 for all records where remaining_time > 0
	updateQuery := `UPDATE rs_bot.queue_active SET remaining_time = GREATEST(remaining_time - 1, 0) WHERE remaining_time > 0`
	_, err := d.db.Exec(updateQuery)
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}

	// Then, get all records
	selectQuery := `SELECT id, data, remaining_time FROM rs_bot.queue_active`
	rows, err := d.db.Query(selectQuery)
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}
	defer rows.Close()

	var result []models.QueueActive
	for rows.Next() {
		var qa models.QueueActive
		var dataJSON []byte

		err = rows.Scan(&qa.ID, &dataJSON, &qa.RemainingTime)
		if err != nil {
			d.log.ErrorErr(err)
			continue
		}
		qa.JsonUnmarshalData(dataJSON)
		result = append(result, qa)
	}

	return result, nil
}

// DeleteActiveQueue deletes an active queue entry by ID
func (d *Db) DeleteActiveQueue(id int64) error {

	query := `DELETE FROM rs_bot.queue_active WHERE id = $1`
	_, err := d.db.Exec(query, id)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}

	return nil
}

// DeleteActiveQueueByUserAndLevel deletes active queue entries for a user in a specific level
func (d *Db) DeleteActiveQueueByUserAndLevel(userID, corpUuid, lvlRS string) error {

	query := `DELETE FROM rs_bot.queue_active 
              WHERE data->>'user_id' = $1 AND data->>'corporation_uuid' = $2 AND data->>'lvl_rs' = $3`
	_, err := d.db.Exec(query, userID, corpUuid, lvlRS)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}

	return nil
}

// MoveToCompleteQueue moves an active queue entry to the complete queue
func (d *Db) MoveToCompleteQueue(rs *models.Rs) error {

	for _, a := range rs.U {
		fullBytes := a.GetFullMap(rs.QueueMessages)
		// Insert into complete queue
		insertQuery := `INSERT INTO rs_bot.queue_complete(data) VALUES ($1)`
		_, err := d.db.Exec(insertQuery, fullBytes)
		if err != nil {
			d.log.ErrorErr(err)
			return err
		}
		if a.ID != 0 {
			err = d.DeleteActiveQueue(a.ID)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

// GetActiveQueueByUser retrieves active queue entries for a specific user
func (d *Db) GetActiveQueueByUser(userID string) ([]models.QueueActive, error) {

	query := `SELECT id, data, remaining_time FROM rs_bot.queue_active 
              WHERE data->>'user_id' = $1`
	rows, err := d.db.Query(query, userID)
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}
	defer rows.Close()

	var result []models.QueueActive
	for rows.Next() {
		var qa models.QueueActive
		var dataJSON []byte

		err = rows.Scan(&qa.ID, &dataJSON, &qa.RemainingTime)
		if err != nil {
			d.log.ErrorErr(err)
			continue
		}

		qa.JsonUnmarshalData(dataJSON)

		result = append(result, qa)
	}

	return result, nil
}

// CountActiveQueueByCorpAndLevel counts active queue entries for a specific corporation and level
func (d *Db) CountActiveQueueByCorpAndLevel(corpUuid, lvlRS string) (int, error) {

	var count int
	query := `SELECT COUNT(*) FROM rs_bot.queue_active 
              WHERE data->>'corporation_uuid' = $1 AND data->>'lvl_rs' = $2`
	err := d.db.QueryRow(query, corpUuid, lvlRS).Scan(&count)
	if err != nil {
		d.log.ErrorErr(err)
		return 0, err
	}

	return count, nil
}

// CheckUserInActiveQueue - проверка пользователя в активной таблице queue_active
func (d *Db) CheckUserInActiveQueue(userid, lvlkz, corpUuid string) (bool, error) {

	var count int
	sel := "SELECT COUNT(*) FROM rs_bot.queue_active WHERE data->>'user_id' = $1 AND data->>'lvl_rs' = $2 AND data->>'corporation_uuid' = $3"
	row := d.db.QueryRow(sel, userid, lvlkz, corpUuid)
	err := row.Scan(&count)
	if err != nil {
		d.log.ErrorErr(err)
		return false, err
	}

	return count > 0, nil
}

// GetQueueState returns active status, queue count, completed count for user, and total queue count for the level in one query
func (d *Db) GetQueueState(userID, lvlRS, corpUuid string) (active bool, countQueue, numberName, numberLevel int, err error) {

	query := `
		SELECT 
			(SELECT COUNT(*) > 0 FROM rs_bot.queue_active WHERE data->>'user_id' = $1 AND data->>'lvl_rs' = $2 AND data->>'corporation_uuid' = $3) as active,
			(SELECT COUNT(*) FROM rs_bot.queue_active WHERE data->>'corporation_uuid' = $3 AND data->>'lvl_rs' = $2) as count_queue,
			(SELECT COUNT(*) FROM rs_bot.queue_complete WHERE data->>'corporation_uuid' = $3 AND data->>'lvl_rs' = $2 AND data->>'user_id' = $1) as number_name,
			COALESCE((SELECT count FROM rs_bot.queue_count WHERE corporation = $3 AND level = $2), 0) as number_level
	`

	err = d.db.QueryRow(query, userID, lvlRS, corpUuid).Scan(&active, &countQueue, &numberName, &numberLevel)
	if err != nil {
		d.log.ErrorErr(err)
		return false, 0, 0, 0, err
	}

	return active, countQueue, numberName, numberLevel, nil
}

// GetActiveQueueLevelsByCorp retrieves all active queue levels for a corporation
func (d *Db) GetActiveQueueLevelsByCorp(corpUuid string) ([]string, error) {

	query := `SELECT DISTINCT data->>'lvl_rs' FROM rs_bot.queue_active 
              WHERE data->>'corporation_uuid' = $1`
	rows, err := d.db.Query(query, corpUuid)
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}
	defer rows.Close()

	var levels []string
	for rows.Next() {
		var level string
		err := rows.Scan(&level)
		if err != nil {
			d.log.ErrorErr(err)
			continue
		}
		levels = append(levels, level)
	}

	return levels, nil
}

// CountCompletedQueueByCorpLevelUser - получение количества завершенных очередей для корпорации, уровня и пользователя
func (d *Db) CountCompletedQueueByCorpLevelUser(corpUuid, lvlRS, userID string) (int, error) {

	var count int
	sel := "SELECT COUNT(*) FROM rs_bot.queue_complete WHERE data->>'corporation_uuid' = $1 AND data->>'lvl_rs' = $2 AND data->>'user_id' = $3"
	row := d.db.QueryRow(sel, corpUuid, lvlRS, userID)
	err := row.Scan(&count)
	if err != nil {
		d.log.ErrorErr(err)
		return 0, err
	}

	return count, nil
}

func (d *Db) CountUserIdQueue(userid string) int {
	//проверяем есть ли игрок в других очередях

	sel := "SELECT COUNT(*) FROM rs_bot.queue_active WHERE data->>'user_id' = $1"
	var count int
	row := d.db.QueryRow(sel, userid)
	err := row.Scan(&count)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			d.log.ErrorErr(err)
		}
		return 0
	}

	return count
}
