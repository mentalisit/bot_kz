package postgres

import (
	"discord/models"
)

func (d *Db) ScoreboardInsertParam(p models.ScoreboardParams) {
	ctx, cancel := d.getContext()
	defer cancel()
	insert := `INSERT INTO rs_bot.scoreboard(name,webhookchannel,scorechannel,lastmessage) VALUES ($1,$2,$3,$4)`
	_, err := d.db.Exec(ctx, insert, p.Name, p.ChannelWebhook, p.ChannelScoreboard, p.LastMessageID)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) ScoreboardUpdateParamLastMessageId(p models.ScoreboardParams) {
	ctx, cancel := d.getContext()
	defer cancel()
	update := `UPDATE rs_bot.scoreboard SET lastmessage = $1 where name = $2`
	_, err := d.db.Exec(ctx, update, p.LastMessageID, p.Name)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) ScoreboardReadWebhookChannel(webhookChannel string) *models.ScoreboardParams {
	ctx, cancel := d.getContext()
	defer cancel()
	selectScoreboard := "SELECT name, webhookchannel, scorechannel,lastmessage FROM rs_bot.scoreboard WHERE webhookchannel = $1"
	results, err := d.db.Query(ctx, selectScoreboard, webhookChannel)
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
	ctx, cancel := d.getContext()
	defer cancel()
	selectScoreboard := "SELECT name, webhookchannel, scorechannel,lastmessage FROM rs_bot.scoreboard"
	rows, err := d.db.Query(ctx, selectScoreboard)
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
	ctx, cancel := d.getContext()
	defer cancel()

	sel := "SELECT datestart,datestop,message FROM kzbot.event ORDER BY id DESC LIMIT 1"
	err := d.db.QueryRow(ctx, sel).Scan(&nextDateStart, &nextDateStop, &message)
	if err != nil {
		d.log.ErrorErr(err)
		return "", "", ""
	}
	return nextDateStart, nextDateStop, message
}

func (d *Db) InsertWebhook(ts int64, corp, message string) {
	ctx, cancel := d.getContext()
	defer cancel()
	insert := `INSERT INTO rs_bot.webhooks(tsunix,corp,message) VALUES ($1,$2,$3)`
	_, err := d.db.Exec(ctx, insert, ts, corp, message)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) InsertWebhookType(ts int64, corpName, eventType, message string) {
	ctx, cancel := d.getContext()
	defer cancel()
	insert := `INSERT INTO rs_bot.webhook_type(tsUnix,corpName,eventType,message) VALUES ($1,$2,$3,$4)`
	_, err := d.db.Exec(ctx, insert, ts, corpName, eventType, message)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) LoadNameAliases() (map[string]string, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	rows, err := d.db.Query(ctx, "SELECT alias, canonical_name FROM rs_bot.name_aliases")
	defer rows.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	m := make(map[string]string)
	for rows.Next() {
		var alias, canonical string
		if err := rows.Scan(&alias, &canonical); err != nil {
			return nil, err
		}
		m[alias] = canonical
	}

	return m, nil
}
