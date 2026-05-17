package postgresV2

import (
	"fmt"
	"rs/models"
)

func (d *Db) UpdateGamePlayer(p models.GameAccountData) error {

	// Правильный синтаксис: колонки перечисляются через запятую
	query := `UPDATE rs_bot2.game_players SET player = $1, owner = $2 WHERE id = $3;`

	_, err := d.db.Exec(query, p.PlayerName, p.OwnerUuid, p.PlayerID)
	if err != nil {
		return fmt.Errorf("failed to update player %+v\n: %w", p, err)
	}

	return nil
}

func (d *Db) GetMergeGameAccount(name string) []models.GameAccountData {

	// Выбираем и ID, и Имя игрока (player).
	// Используем TRIM, чтобы игнорировать лишние пробелы при поиске
	query := `SELECT id, player,owner FROM rs_bot2.game_players WHERE TRIM(player) = TRIM($1);`

	rows, err := d.db.Query(query, name)
	if err != nil {
		d.log.ErrorErr(err)
		return nil
	}
	defer rows.Close()

	var acc []models.GameAccountData
	for rows.Next() {
		var gacc models.GameAccountData
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

func (d *Db) GetMergeGameAccountAll() map[string]models.GameAccountData {

	query := `SELECT id, player, owner FROM rs_bot2.game_players;`

	rows, err := d.db.Query(query)
	if err != nil {
		d.log.ErrorErr(err)
		return nil
	}
	defer rows.Close()

	ga := make(map[string]models.GameAccountData)

	for rows.Next() {
		var gacc models.GameAccountData
		// Сканируем два поля в структуру
		err = rows.Scan(&gacc.PlayerID, &gacc.PlayerName, &gacc.OwnerUuid)
		if err != nil {
			d.log.ErrorErr(err)
			continue
		}
		ga[gacc.PlayerName] = gacc
	}
	return ga
}
