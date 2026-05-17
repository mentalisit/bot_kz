package postgresV2

import (
	"encoding/json"
	"regexp"
	"rs/models"
	"sort"
	"strconv"

	"github.com/google/uuid"
)

// ScoreboardInsertParam вставляет новую запись scoreboard_config
func (d *Db) ScoreboardInsertParam(p models.ScoreboardParamsV2) {

	channelsJSON, _ := json.Marshal(p.Channels)
	gameJSON, _ := json.Marshal(p.Game)

	// SQL запрос с логикой UPSERT (Update or Insert)
	query := `
        INSERT INTO rs_bot.scoreboard_config (uid, game, channels) 
        VALUES ($1, $2, $3)
        ON CONFLICT (uid) 
        DO UPDATE SET 
            game = EXCLUDED.game, 
            channels = EXCLUDED.channels`

	_, err := d.db.Exec(query, p.Uid, gameJSON, channelsJSON)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

// ScoreboardUpdateParamChannels обновляет только channels по uid
func (d *Db) ScoreboardUpdateParamChannels(p models.ScoreboardParamsV2) {

	channelsJSON, _ := json.Marshal(p.Channels)
	update := `UPDATE rs_bot.scoreboard_config SET channels = $1 WHERE uid = $2`
	_, err := d.db.Exec(update, channelsJSON, p.Uid)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) ScoreboardUpdateParamGame(p models.ScoreboardParamsV2) {

	gameJSON, _ := json.Marshal(p.Game)
	update := `UPDATE rs_bot.scoreboard_config SET game = $1 WHERE uid = $2`
	_, err := d.db.Exec(update, gameJSON, p.Uid)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

// ScoreboardReadByChannelId ищет запись по ChannelId в JSON поле channels
func (d *Db) ScoreboardReadByChannelId(channelId string) *models.ScoreboardParamsV2 {

	// Ищем channelId в JSON массиве channels используя LATERAL
	// Структура: [{"channel_id": "...", ...}, ...]
	query := `
		SELECT DISTINCT sc.uid, sc.game, sc.channels
		FROM rs_bot.scoreboard_config sc,
		LATERAL jsonb_array_elements(sc.channels) AS elem
		WHERE elem->>'ChannelId' = $1`

	results, err := d.db.Query(query, channelId)
	if err != nil {
		d.log.ErrorErr(err)
		return nil
	}
	defer results.Close()

	var s models.ScoreboardParamsV2
	var uid uuid.UUID
	var channelsJSON []byte
	var gameJSON []byte
	if results.Next() {
		err = results.Scan(&uid, &gameJSON, &channelsJSON)
		if err != nil {
			d.log.ErrorErr(err)
			return nil
		}

		s.Uid = uid.String()
		_ = json.Unmarshal(channelsJSON, &s.Channels)
		_ = json.Unmarshal(gameJSON, &s.Game)
	}
	if s.Uid == "" {
		return nil
	}
	return &s
}

// ScoreboardReadAll читает все записи
func (d *Db) ScoreboardReadAll() []models.ScoreboardParamsV2 {

	selectScoreboard := "SELECT uid, game, channels FROM rs_bot.scoreboard_config"
	rows, err := d.db.Query(selectScoreboard)
	if err != nil {
		d.log.ErrorErr(err)
		return nil
	}
	defer rows.Close()

	var ss []models.ScoreboardParamsV2
	for rows.Next() {
		var s models.ScoreboardParamsV2
		var uid uuid.UUID
		var channelsJSON []byte
		var gameJSON []byte
		err = rows.Scan(&uid, &gameJSON, &channelsJSON)
		if err != nil {
			d.log.ErrorErr(err)
			continue
		}
		s.Uid = uid.String()
		_ = json.Unmarshal(channelsJSON, &s.Channels)
		_ = json.Unmarshal(gameJSON, &s.Game)
		ss = append(ss, s)
	}
	return ss
}

// ScoreboardReadByUid ищет запись по uid
func (d *Db) ScoreboardReadByUid(uidStr string) *models.ScoreboardParamsV2 {

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		d.log.ErrorErr(err)
		return nil
	}

	selectScoreboard := "SELECT uid, game, channels FROM rs_bot.scoreboard_config WHERE uid = $1"
	results, err := d.db.Query(selectScoreboard, uid)
	if err != nil {
		d.log.ErrorErr(err)
		return nil
	}
	defer results.Close()

	var s models.ScoreboardParamsV2
	var uidRes uuid.UUID
	var channelsJSON []byte
	var gameJSON []byte
	for results.Next() {
		err = results.Scan(&uidRes, &gameJSON, &channelsJSON)
		if err != nil {
			d.log.ErrorErr(err)
			continue
		}
		s.Uid = uidRes.String()
		_ = json.Unmarshal(gameJSON, &s.Game)
		_ = json.Unmarshal(channelsJSON, &s.Channels)
	}
	if s.Uid == "" {
		return nil
	}
	return &s
}

// ScoreboardDeleteByUid удаляет запись по uid
func (d *Db) ScoreboardDeleteByUid(uid string) {

	deleteQuery := `DELETE FROM rs_bot.scoreboard_config WHERE uid = $1`
	_, err := d.db.Exec(deleteQuery, uid)
	if err != nil {
		d.log.ErrorErr(err)
	}
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

func (d *Db) ReadEventScheduleAll() []models.ScheduleEvents {

	getSeasonNumber := func(text string) int {
		re := regexp.MustCompile(`Season (\d+)`)
		matches := re.FindStringSubmatch(text)

		if len(matches) < 2 {
			return 0
		}
		seasonNumber, _ := strconv.Atoi(matches[1])
		return seasonNumber
	}

	sel := "SELECT datestart,datestop,message FROM kzbot.event"
	rows, err := d.db.Query(sel)
	if err != nil {
		d.log.ErrorErr(err)
		return nil
	}
	defer rows.Close()

	var ss []models.ScheduleEvents
	for rows.Next() {
		var s models.ScheduleEvents
		var message string
		err = rows.Scan(&s.NextDateStart, &s.NextDateStop, &message)
		if err != nil {
			d.log.ErrorErr(err)
			continue
		}
		s.Season = getSeasonNumber(message)

		ss = append(ss, s)
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Season > ss[j].Season
	})

	return ss
}
