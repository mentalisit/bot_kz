package postgres

import (
	"rs/models"
)

func (d *Db) ScoreboardInsertParam(p models.ScoreboardParams) {
	insert := `INSERT INTO rs_bot.scoreboard(name,webhookchannel,scorechannel,lastmessage) VALUES ($1,$2,$3,$4)`
	_, err := d.db.Exec(insert, p.Name, p.ChannelWebhook, p.ChannelScoreboardOrMap, p.LastMessageID)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) ScoreboardUpdateParam(p models.ScoreboardParams) {
	update := `UPDATE rs_bot.scoreboard SET webhookchannel = $1,scorechannel = $2 where name = $3`
	_, err := d.db.Exec(update, p.ChannelWebhook, p.ChannelScoreboardOrMap, p.Name)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) ScoreboardUpdateParamScoreChannels(p models.ScoreboardParams) {
	update := `UPDATE rs_bot.scoreboard SET scorechannel = $1 where name = $2`
	_, err := d.db.Exec(update, p.ChannelScoreboardOrMap, p.Name)
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
	selectScoreboard := "SELECT name, webhookchannel, scorechannel,lastmessage FROM rs_bot.scoreboard WHERE name = $1"
	results, err := d.db.Query(selectScoreboard, name)
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
	selectScoreboard := "SELECT name, webhookchannel, scorechannel,lastmessage FROM rs_bot.scoreboard"
	rows, err := d.db.Query(selectScoreboard)
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
	sel := "SELECT datestart,datestop,message FROM kzbot.event ORDER BY id DESC LIMIT 1"
	err := d.db.QueryRow(sel).Scan(&nextDateStart, &nextDateStop, &message)
	if err != nil {
		d.log.ErrorErr(err)
		return "", "", ""
	}
	return nextDateStart, nextDateStop, message
}
