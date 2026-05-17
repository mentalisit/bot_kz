package postgres

import (
	"bridge/models"
	"encoding/json"
	"fmt"
	"strconv"
)

func (d *Db) CreatePoll(data models.Poll2Struct) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	query := `INSERT INTO rs_bot2.poll (id, data) VALUES ($1, $2)`
	_, err = d.db.Exec(query, data.CreateTime, dataJSON)
	return err
}

func (d *Db) GetPollById(id string) (models.Poll2Struct, error) {
	var poll models.Poll2Struct

	// 1. Конвертируем строковый ID в int64 для соответствия типу bigint в БД
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return poll, fmt.Errorf("invalid poll id format: %w", err)
	}

	// 2. Временная структура для сканирования данных из БД
	var row struct {
		Data  json.RawMessage `db:"data"`
		Votes json.RawMessage `db:"votes"`
	}

	// 3. Выполняем запрос
	query := `SELECT data, votes FROM rs_bot2.poll WHERE id = $1`
	err = d.db.Get(&row, query, intID)
	if err != nil {
		return poll, fmt.Errorf("failed to get poll from db: %w", err)
	}

	// 4. Анмаршалим основное тело опроса
	if err := json.Unmarshal(row.Data, &poll); err != nil {
		return poll, fmt.Errorf("failed to unmarshal poll data: %w", err)
	}

	// 5. Опционально: если нужно вернуть голоса вместе со структурой,
	// убедитесь, что в Poll2Struct есть соответствующее поле.
	// Если голоса хранятся отдельно, их можно распарсить в map[string]models.Votes2

	return poll, nil
}

func (d *Db) UpsertVote(pollID int64, vote models.Votes2) error {
	voteJSON, err := json.Marshal(vote)
	if err != nil {
		return err
	}

	// Добавляем явное приведение типов ::text и ::jsonb
	query := `
       UPDATE rs_bot2.poll 
       SET votes = votes || jsonb_build_object($2::text, $3::jsonb)
       WHERE id = $1`

	// Используем d.db.Exec (убедитесь, что передаете аргументы в правильном порядке)
	_, err = d.db.Exec(query, pollID, vote.Uid, voteJSON)
	return err
}

func (d *Db) GetVotes(pollID int64) (map[string]models.Votes2, error) {
	var rawVotes json.RawMessage
	query := `SELECT votes FROM rs_bot2.poll WHERE id = $1`

	err := d.db.Get(&rawVotes, query, pollID)
	if err != nil {
		return nil, err
	}

	votes := make(map[string]models.Votes2)
	if len(rawVotes) > 0 && string(rawVotes) != "{}" {
		if err = json.Unmarshal(rawVotes, &votes); err != nil {
			return nil, err
		}
	}
	return votes, nil
}
