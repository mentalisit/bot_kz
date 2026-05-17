package dbpostgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"ws/models"

	"github.com/jackc/pgx/v5"
)

func (d *Db) GetPollById(id string) (models.PollStruct, []models.Votes, error) {
	var poll models.PollStruct
	var votesList []models.Votes

	// 1. Парсим ID
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return poll, nil, fmt.Errorf("invalid poll id: %w", err)
	}

	// 2. Объявляем переменные для сырых JSON данных
	var rawData []byte
	var rawVotes []byte

	// 3. Выполняем запрос через QueryRow
	// Мы используем контекст (обычно передается сверху, здесь для примера Background)
	query := `SELECT data, votes FROM rs_bot2.poll WHERE id = $1`
	err = d.pool.QueryRow(context.Background(), query, intID).Scan(&rawData, &rawVotes)

	if err != nil {
		if err == pgx.ErrNoRows {
			return poll, nil, nil // Опрос не найден
		}
		return poll, nil, fmt.Errorf("scan poll error: %w", err)
	}

	// 4. Анмаршалим основную дату
	if err := json.Unmarshal(rawData, &poll); err != nil {
		return poll, nil, fmt.Errorf("unmarshal poll data: %w", err)
	}

	// 5. Анмаршалим голоса.
	// Так как мы решили хранить их как объект {"uid": {vote}},
	// сначала читаем в map, а потом превращаем в слайс.
	var votesMap map[string]models.Votes
	if len(rawVotes) > 0 {
		if err := json.Unmarshal(rawVotes, &votesMap); err != nil {
			return poll, nil, fmt.Errorf("unmarshal votes: %w", err)
		}
	}

	// Конвертируем карту в слайс для возврата
	for _, v := range votesMap {
		votesList = append(votesList, v)
	}

	return poll, votesList, nil
}
