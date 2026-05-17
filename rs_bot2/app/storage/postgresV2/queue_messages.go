package postgresV2

import (
	"database/sql"
	"encoding/json"
	"errors"
	"rs/models"
)

// SaveQueueMessages - вставить или обновить сообщения для корпорации и уровня
func (d *Db) SaveQueueMessages(corporation, level string, messages map[string]models.QueueMessages) error {

	b, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO rs_bot.queue_active_messages (corporation, level, messages)
		VALUES ($1, $2, $3)
		ON CONFLICT (corporation, level)
		DO UPDATE SET messages = EXCLUDED.messages`

	_, err = d.db.Exec(query, corporation, level, b)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}
	return nil
}

// ReadQueueMessages - чтение сообщений для корпорации и уровня
func (d *Db) ReadQueueMessages(corporation, level string) (map[string]models.QueueMessages, error) {

	query := `SELECT messages FROM rs_bot.queue_active_messages WHERE corporation = $1 AND level = $2`

	var b []byte
	err := d.db.QueryRow(query, corporation, level).Scan(&b)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make(map[string]models.QueueMessages), nil
		}
		return nil, err
	}

	var messages map[string]models.QueueMessages
	err = json.Unmarshal(b, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// DeleteQueueMessages - удаление сообщений для корпорации и уровня
func (d *Db) DeleteQueueMessages(corporation, level string) error {

	query := `DELETE FROM rs_bot.queue_active_messages WHERE corporation = $1 AND level = $2`
	_, err := d.db.Exec(query, corporation, level)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}
	return nil
}

// UpdateQueueMessages - обновление сообщений для корпорации и уровня
func (d *Db) UpdateQueueMessages(corporation, level string, messages map[string]models.QueueMessages) error {
	return d.SaveQueueMessages(corporation, level, messages)
}
