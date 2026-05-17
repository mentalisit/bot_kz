package postgres

import (
	"context"
	"fmt"
)

type GameAccountData struct {
	PlayerID   string `json:"id"`
	PlayerName string `json:"name"`
	OwnerUuid  string `json:"ownerUuid"`
}

func (d *Db) UpsertGamePlayer(ctx context.Context, id, playerName string) error {
	// Используем ON CONFLICT для обновления имени, если ID уже существует.
	// Поле owner оставляем пустым по умолчанию или не трогаем при обновлении.
	query := `
        INSERT INTO rs_bot2.game_players (id, player, owner)
        VALUES ($1, $2, '')
        ON CONFLICT (id) 
        DO UPDATE SET player = EXCLUDED.player;`

	_, err := d.db.Exec(query, id, playerName)
	if err != nil {
		return fmt.Errorf("failed to upsert player %s: %w", id, err)
	}

	return nil
}

func (d *Db) GetMergeGameAccount(name string) []GameAccountData {

	// Выбираем и ID, и Имя игрока (player).
	// Используем TRIM, чтобы игнорировать лишние пробелы при поиске
	query := `SELECT id, player,owner FROM rs_bot2.game_players WHERE TRIM(player) = TRIM($1);`

	rows, err := d.db.Query(query, name)
	if err != nil {
		d.log.ErrorErr(err)
		return nil
	}
	defer rows.Close()

	var acc []GameAccountData
	for rows.Next() {
		var gacc GameAccountData
		// Сканируем два поля в структуру
		err = rows.Scan(&gacc.PlayerID, &gacc.PlayerName, &gacc.OwnerUuid)
		if err != nil {
			d.log.ErrorErr(err)
			continue
		}
		acc = append(acc, gacc)
	}
	return acc
}

func (d *Db) GetMergeGameAccountAll() map[string]string {

	query := `SELECT id, player FROM rs_bot2.game_players;`

	rows, err := d.db.Query(query)
	if err != nil {
		d.log.ErrorErr(err)
		return nil
	}
	defer rows.Close()

	ga := make(map[string]string)

	for rows.Next() {
		var id, player string
		// Сканируем два поля в структуру
		err = rows.Scan(&id, &player)
		if err != nil {
			d.log.ErrorErr(err)
			continue
		}
		ga[id] = player
	}
	return ga
}
