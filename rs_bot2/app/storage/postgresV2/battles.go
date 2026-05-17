package postgresV2

import (
	"context"
	"fmt"
	"rs/models"
	"time"
)

func (d *Db) BattlesGetAll(corpName string, event int) ([]models.PlayerStats, error) {

	query := `
		SELECT name,
		       SUM(points) AS total_points, 
		       COUNT(*) AS runs,
		       MAX(level) AS max_level		
		FROM rs_bot.battles 
		where eventid=$1 AND corporation=$2
		GROUP BY name;
	`

	rows, err := d.db.Query(query, event, corpName)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	var stats []models.PlayerStats
	for rows.Next() {
		var ps models.PlayerStats
		if err := rows.Scan(&ps.Player, &ps.Points, &ps.Runs, &ps.Level); err != nil {
			return nil, fmt.Errorf("ошибка чтения данных: %v", err)
		}
		stats = append(stats, ps)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка после чтения строк: %v", err)
	}

	return stats, nil
}

func (d *Db) BattlesGetAllId(event int) ([]models.PlayerStats, error) {

	query := `
		SELECT name,
		       corporation,
		       SUM(points) AS total_points, 
		       COUNT(*) AS runs,
		       MAX(level) AS max_level		
		FROM rs_bot.battles 
		WHERE eventid = $1 
		GROUP BY name, corporation;
	`

	rows, err := d.db.Query(query, event)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	var stats []models.PlayerStats

	for rows.Next() {
		var ps models.PlayerStats
		if err := rows.Scan(
			&ps.Player,
			&ps.CorporationName,
			&ps.Points,
			&ps.Runs,
			&ps.Level,
		); err != nil {
			return nil, fmt.Errorf("ошибка чтения данных: %v", err)
		}
		stats = append(stats, ps)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка после чтения строк: %v", err)
	}

	return stats, nil
}

func (d *Db) ScoreboardParamsReadAll() []models.ScoreboardParams {

	query := `SELECT name,webhookchannel,scorechannel FROM rs_bot.scoreboard`
	rows, err := d.db.Query(query)
	if err != nil {
		d.log.ErrorErr(err)
		return nil
	}
	defer rows.Close()

	var params []models.ScoreboardParams
	for rows.Next() {
		var ps models.ScoreboardParams
		if err = rows.Scan(&ps.Name, &ps.ChannelWebhook, &ps.ChannelScoreboardOrMap); err != nil {
			d.log.ErrorErr(err)
			return nil
		}
		params = append(params, ps)
	}

	if err := rows.Err(); err != nil {
		d.log.ErrorErr(err)
		return nil
	}

	return params
}

func (d *Db) BattlesTopGetAll(corpName string) ([]models.BattlesTop, error) {

	query := `SELECT * FROM rs_bot.battlestop where corporation=$1 `
	rows, err := d.db.Query(query, corpName)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	var stats []models.BattlesTop
	for rows.Next() {
		var ps models.BattlesTop
		if err = rows.Scan(&ps.Id, &ps.CorpName, &ps.Name, &ps.Level, &ps.Count); err != nil {
			return nil, fmt.Errorf("ошибка чтения данных: %v", err)
		}
		stats = append(stats, ps)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка после чтения строк: %v", err)
	}

	return stats, nil
}

func (d *Db) DeleteOldWebhooks() {

	// Вычисляем Unix-время 7 дней назад
	// 7 дней * 24 часа * 3600 секунд
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Unix()

	query := `DELETE FROM rs_bot.webhooks WHERE tsunix < $1`

	tag, err := d.db.Exec(query, sevenDaysAgo)
	if err != nil {
		d.log.ErrorErr(fmt.Errorf("failed to delete old webhooks: %w", err))
		return
	}
	count, _ := tag.RowsAffected()
	if count > 0 {
		fmt.Printf("🧹 Удалено старых webhooks: %d\n", count)
	}

	query = `DELETE FROM rs_bot.webhook_type WHERE tsunix < $1`

	tag, err = d.db.Exec(query, sevenDaysAgo)
	if err != nil {
		//d.log.ErrorErr(fmt.Errorf("failed to delete old rs_bot.webhook_type: %w", err))
		return
	}
	count, _ = tag.RowsAffected()
	if count > 0 {
		fmt.Printf("🧹 Удалено старых webhook_type: %d\n", count)
	}

	query = `
        DELETE FROM rs_bot.message_maps 
        WHERE created_at < now() - interval '7 days'`

	tag, err = d.db.Exec(query)
	if err != nil {
		//d.log.ErrorErr(fmt.Errorf("failed to delete old message_maps: %w", err))
		return
	}

	count, _ = tag.RowsAffected()
	if count > 0 {
		fmt.Printf("🧹 Удалено старых карт сообщений: %d\n", count)
	}

}

func (d *Db) BattlesGetForEvent(name string, event int) ([]models.EventGame, error) {

	query := `
		SELECT id, eventid, corporation, name, level, points 
		FROM rs_bot.battles 
		WHERE eventid=$1 AND name=$2
		ORDER BY id DESC; 
	`

	rows, err := d.db.Query(query, event, name)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения BattlesGetForEvent: %w", err)
	}
	defer rows.Close()

	var stats []models.EventGame
	for rows.Next() {
		var ps models.EventGame
		if err := rows.Scan(&ps.ID, &ps.Season, &ps.CorporationName, &ps.Player, &ps.Level, &ps.Points); err != nil {
			return nil, fmt.Errorf("ошибка сканирования: %w", err)
		}
		stats = append(stats, ps)
	}

	// 4. Проверка на ошибки итерации
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка в итераторе строк: %w", err)
	}

	return stats, nil
}

type PlayerInfo struct {
	PlayerID   string `json:"PlayerID"`
	PlayerName string `json:"PlayerName"`
}

// Общая структура для всех типов событий
type WebhookMessage struct {
	EventType             string       `json:"EventType"`
	OurParticipants       []PlayerInfo `json:"OurParticipants"`
	OpponentParticipants  []PlayerInfo `json:"OpponentParticipants"`
	PlayersWhoContributed []PlayerInfo `json:"PlayersWhoContributed"`
	Players               []PlayerInfo `json:"Players"`
}

func (d *Db) GetAllPlayersFromWebhooks(ctx context.Context) (map[string]string, error) {
	query := `SELECT message FROM rs_bot.webhook_type`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	// Используем map для хранения уникальных игроков: ID -> Name
	playersMap := make(map[string]string)

	for rows.Next() {
		var msg WebhookMessage
		// pgx автоматически распарсит jsonb в структуру msg
		err := rows.Scan(&msg)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		// Собираем игроков из всех возможных слайсов
		collectPlayers(playersMap, msg.OurParticipants)
		collectPlayers(playersMap, msg.OpponentParticipants)
		collectPlayers(playersMap, msg.PlayersWhoContributed)
		collectPlayers(playersMap, msg.Players)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return playersMap, nil
}

// Вспомогательная функция для наполнения карты
func collectPlayers(m map[string]string, list []PlayerInfo) {
	for _, p := range list {
		if p.PlayerID != "" {
			m[p.PlayerID] = p.PlayerName
		}
	}
}
