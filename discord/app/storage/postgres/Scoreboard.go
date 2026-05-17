package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"

	"fmt"

	"github.com/mentalisit/restapi/models"
)

func (d *Db) ScoreboardInsertParam(p models.ScoreboardParams) {

	insert := `INSERT INTO rs_bot.scoreboard(name,webhookchannel,scorechannel,lastmessage) VALUES ($1,$2,$3,$4)`
	_, err := d.db.Exec(insert, p.Name, p.ChannelWebhook, p.ChannelScoreboard, p.LastMessageID)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) ScoreboardUpdateParamLastMessageId(p models.ScoreboardParams) {

	update := `UPDATE rs_bot.scoreboard SET lastmessage = $1 where name = $2`
	_, err := d.db.Exec(update, p.LastMessageID, p.Name)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) ScoreboardReadWebhookChannel(webhookChannel string) *models.ScoreboardParams {

	selectScoreboard := "SELECT name, webhookchannel, scorechannel,lastmessage FROM rs_bot.scoreboard WHERE webhookchannel = $1"
	results, err := d.db.Query(selectScoreboard, webhookChannel)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	var s models.ScoreboardParams
	for results.Next() {
		err = results.Scan(&s.Name, &s.ChannelWebhook, &s.ChannelScoreboard, &s.LastMessageID)
		if err != nil {
			d.log.ErrorErr(err)
		}
	}
	if s.Name == "" {
		return nil
	}
	return &s
}
func (d *Db) ScoreboardReadAll() []models.ScoreboardParams {

	selectScoreboard := "SELECT name, webhookchannel, scorechannel,lastmessage FROM rs_bot.scoreboard"
	rows, err := d.db.Query(selectScoreboard)
	defer rows.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	var ss []models.ScoreboardParams
	for rows.Next() {
		var s models.ScoreboardParams
		err = rows.Scan(&s.Name, &s.ChannelWebhook, &s.ChannelScoreboard, &s.LastMessageID)
		if err != nil {
			d.log.ErrorErr(err)
		}
		ss = append(ss, s)
	}
	return ss
}

func (d *Db) ReadEventScheduleAndMessage() (nextDateStart, nextDateStop, message string) {

	sel := "SELECT datestart,datestop,message FROM kzbot.event ORDER BY id DESC LIMIT 1"
	err := d.db.QueryRow(sel).Scan(&nextDateStart, &nextDateStop, &message)
	if err != nil {
		d.log.ErrorErr(err)
		return "", "", ""
	}
	return nextDateStart, nextDateStop, message
}

func (d *Db) InsertWebhook(ts int64, corp, message string) {

	insert := `INSERT INTO rs_bot.webhooks(tsunix,corp,message) VALUES ($1,$2,$3)`
	_, err := d.db.Exec(insert, ts, corp, message)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) InsertWebhookType(ts int64, corpName, eventType, message string) {

	insert := `INSERT INTO rs_bot.webhook_type(tsUnix,corpName,eventType,message) VALUES ($1,$2,$3,$4)`
	_, err := d.db.Exec(insert, ts, corpName, eventType, message)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) GetNickNameByMergedID(targetID string) (string, error) {

	// Формируем JSON-паттерн для поиска в массиве
	searchPattern, err := json.Marshal([]map[string]string{
		{"id": targetID},
	})
	if err != nil {
		return "", fmt.Errorf("ошибка маршалинга: %w", err)
	}

	query := `
		SELECT nickname 
		FROM my_compendium.multi_accounts 
		WHERE data->'Merged' @> $1::jsonb 
		LIMIT 1
	`

	var nickName string
	// Используем QueryRow, так как нам нужно одно значение
	err = d.db.QueryRow(query, string(searchPattern)).Scan(&nickName)

	if err != nil {
		// Если запись не найдена, pgx возвращает sql.ErrNoRows или pgx.ErrNoRows
		// В этом случае возвращаем пустую строку без ошибки
		if !errors.Is(err, sql.ErrNoRows) {
			d.log.ErrorErr(err)
		}
		return "", err
	}

	return nickName, nil
}
