package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

func (d *Db) SaveBridgeMap(msgMap map[string]string) error {

	jsonData, err := json.Marshal(msgMap)
	if err != nil {
		return fmt.Errorf("failed to marshal message map: %w", err)
	}

	// Вставляем новую запись с автоинкрементным ID
	query := `
        INSERT INTO rs_bot.message_maps (message_ids)
        VALUES ($1);
    `

	_, err = d.db.ExecContext(context.Background(), query, jsonData)
	if err != nil {
		return fmt.Errorf("failed to save bridge map: %w", err)
	}
	return nil
}

// GetMapByLinkedID ищет карту, содержащую заданную пару ChatID и MessageID.
func (d *Db) GetMapByLinkedID(msg map[string]string) (map[string]string, error) {

	fragmentData, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search fragment: %w", err)
	}

	query := `
        SELECT message_ids FROM rs_bot.message_maps 
        WHERE message_ids @> $1
        LIMIT 1 -- Предполагаем, что найдена будет только одна карта
    `

	var jsonData []byte

	// Выполняем быстрый запрос благодаря GIN-индексу
	err = d.db.QueryRow(query, string(fragmentData)).Scan(&jsonData)

	if err == sql.ErrNoRows {
		return nil, nil // Карта не найдена
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query bridge map by linked ID: %w", err)
	}

	// Десериализация
	var msgMap map[string]string
	if err := json.Unmarshal(jsonData, &msgMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}

	return msgMap, nil
}
