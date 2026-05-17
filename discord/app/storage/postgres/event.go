package postgres

import (
	"regexp"
	"strconv"
	"time"
)

func (d *Db) SaveEventDate(msg string, season int) {
	lastEvent := d.readLastEvent()
	if lastEvent == season {
		return
	}
	// Получаем текущее время в UTC
	now := time.Now().UTC()

	// Определяем количество дней до ближайшей субботы
	daysUntilSaturday := (time.Saturday - now.Weekday() + 7) % 7
	nextSaturday := now.AddDate(0, 0, int(daysUntilSaturday))
	nextSaturdayMidnight := time.Date(nextSaturday.Year(), nextSaturday.Month(), nextSaturday.Day(), 0, 0, 0, 0, time.UTC)

	dateStart := nextSaturdayMidnight.Format("02-01-2006")
	dateStop := nextSaturdayMidnight.Add(48 * time.Hour).Format("02-01-2006")

	//old
	insert := `INSERT INTO kzbot.event(dateStart,dateStop,message) 
				VALUES ($1,$2,$3)`
	_, err := d.db.Exec(insert, dateStart, dateStop, msg)
	if err != nil {
		d.log.ErrorErr(err)
	}

	//new
	insert = `INSERT INTO rs_bot2.event_schedule(dateStart,dateStop,season) 
				VALUES ($1,$2,$3)`
	_, err = d.db.Exec(insert, dateStart, dateStop, season)
	if err != nil {
		d.log.ErrorErr(err)
	}

}

func (d *Db) readLastEvent() int {

	selectEvent := "SELECT message FROM kzbot.event ORDER BY id DESC LIMIT 1"

	var message string
	err := d.db.QueryRow(selectEvent).Scan(&message)
	if err != nil {
		d.log.ErrorErr(err)
		return 0
	}

	reRsEvent := regexp.MustCompile(`Season (\d+) of the Corporation Red Star event has just started!`)
	match := reRsEvent.FindStringSubmatch(message)
	if len(match) > 1 {
		season, _ := strconv.Atoi(match[1])
		return season
	}

	return 0
}
