package postgres

import "time"

func (d *Db) SaveEventDate(message string) {
	// Получаем текущее время в UTC
	now := time.Now().UTC()

	// Определяем количество дней до ближайшей субботы
	daysUntilSaturday := (time.Saturday - now.Weekday() + 7) % 7
	nextSaturday := now.AddDate(0, 0, int(daysUntilSaturday))
	nextSaturdayMidnight := time.Date(nextSaturday.Year(), nextSaturday.Month(), nextSaturday.Day(), 0, 0, 0, 0, time.UTC)

	dateStart := nextSaturdayMidnight.Format(time.DateOnly)
	dateStop := nextSaturdayMidnight.Add(48 * time.Hour).Format(time.DateOnly)

	ctx, cancel := d.getContext()
	defer cancel()
	insert := `INSERT INTO kzbot.event(dateStart,dateStop,message) 
				VALUES ($1,$2,$3)`
	_, err := d.db.Exec(ctx, insert, dateStart, dateStop, message)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
