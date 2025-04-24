package postgres

import (
	"rs/models"
)

func (d *Db) ScoreboardInsertParam(p models.ScoreboardParams) {
	ctx, cancel := d.getContext()
	defer cancel()
	insert := `INSERT INTO rs_bot.scoreboard(name,webhookchannel,scorechannel,lastmessage) VALUES ($1,$2,$3,$4)`
	_, err := d.db.Exec(ctx, insert, p.Name, p.ChannelWebhook, p.ChannelScoreboardOrMap, p.LastMessageID)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) ScoreboardUpdateParam(p models.ScoreboardParams) {
	ctx, cancel := d.getContext()
	defer cancel()
	update := `UPDATE rs_bot.scoreboard SET webhookchannel = $1,scorechannel = $2 where name = $3`
	_, err := d.db.Exec(ctx, update, p.ChannelWebhook, p.ChannelScoreboardOrMap, p.Name)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) ScoreboardUpdateParamScoreChannels(p models.ScoreboardParams) {
	ctx, cancel := d.getContext()
	defer cancel()
	update := `UPDATE rs_bot.scoreboard SET scorechannel = $1 where name = $2`
	_, err := d.db.Exec(ctx, update, p.ChannelScoreboardOrMap, p.Name)
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
		err = results.Scan(&s.Name, &s.ChannelWebhook, &s.ChannelScoreboardOrMap, &s.LastMessageID)
		if err != nil {
			d.log.ErrorErr(err)
		}
	}
	if s.Name == "" {
		return nil
	}
	return &s
}
func (d *Db) ScoreboardReadName(name string) *models.ScoreboardParams {
	ctx, cancel := d.getContext()
	defer cancel()
	selectScoreboard := "SELECT name, webhookchannel, scorechannel,lastmessage FROM rs_bot.scoreboard WHERE name = $1"
	results, err := d.db.Query(ctx, selectScoreboard, name)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	var s models.ScoreboardParams
	for results.Next() {
		err = results.Scan(&s.Name, &s.ChannelWebhook, &s.ChannelScoreboardOrMap, &s.LastMessageID)
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
		err = rows.Scan(&s.Name, &s.ChannelWebhook, &s.ChannelScoreboardOrMap, &s.LastMessageID)
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
