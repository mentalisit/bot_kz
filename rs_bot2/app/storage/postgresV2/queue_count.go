package postgresV2

import (
	"database/sql"
	"errors"
)

// IncrementQueueCount - вставить 1 или обновить значение на +1 для корпорации и уровня
func (d *Db) IncrementQueueCount(corporation, level string) error {

	// Используем UPSERT: вставляем новую запись с count=1 или обновляем существующую, увеличивая count на 1
	query := `
		INSERT INTO rs_bot.queue_count(count, corporation, level) 
		VALUES (1, $1, $2) 
		ON CONFLICT (corporation, level) 
		DO UPDATE SET count = rs_bot.queue_count.count + 1`

	_, err := d.db.Exec(query, corporation, level)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}

	return nil
}

// ReadQueueCount - чтение по корпорации и уровню
func (d *Db) ReadQueueCount(corporation, level string) (int, error) {

	query := "SELECT count FROM rs_bot.queue_count WHERE corporation = $1 AND level = $2"
	row := d.db.QueryRow(query, corporation, level)

	var qc int
	err := row.Scan(&qc)
	if err != nil {
		// If no rows found, return 0 (default value)
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return qc, nil
}
