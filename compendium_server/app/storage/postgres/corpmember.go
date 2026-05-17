package postgres

import (
	"compendium_s/models"
	"time"
)

func (d *Db) CorpMembersRead(guildid string) ([]models.CorpMember, error) {
	sel := "SELECT * FROM hs_compendium.corpmember WHERE guildid = $1"
	results, err := d.db.Query(sel, guildid)
	defer results.Close()
	if err != nil {
		return nil, err
	}
	var mm []models.CorpMember
	for results.Next() {
		var t models.CorpMember
		var id int
		err = results.Scan(&id, &t.Name, &t.UserId, &t.GuildId, &t.Avatar, &t.AvatarUrl, &t.TimeZone, &t.ZoneOffset, &t.AfkFor)

		getAll, errGet := d.TechGetAll(t)
		if errGet != nil {
			d.log.InfoStruct("TechGetAll(t)", t)
			return nil, errGet
		}
		user, erru := d.UsersGetByUserId(t.UserId)
		if erru != nil {
			d.log.InfoStruct("UsersGetByUserId", t)
			d.log.ErrorErr(erru)
		}

		for _, member := range getAll {
			if user != nil && member.Name == user.Username && user.GameName != "" {
				member.Name = user.GameName
			}
			// Используем функцию с поддержкой DST
			if member.TimeZone != "" {
				member.LocalTime, member.LocalTime24 = getTimeStringsWithDST(member.TimeZone, member.ZoneOffset)
			}
			member.UserId = t.UserId + "/" + member.Name
			mm = append(mm, member)
		}
	}
	return mm, nil
}

func (d *Db) CorpMemberRead(userid string) ([]models.CorpMember, error) {
	sel := "SELECT * FROM hs_compendium.corpmember WHERE userid = $1"
	results, err := d.db.Query(sel, userid)
	defer results.Close()
	if err != nil {
		return nil, err
	}
	var mm []models.CorpMember
	for results.Next() {
		var t models.CorpMember
		var id int
		err = results.Scan(&id, &t.Name, &t.UserId, &t.GuildId, &t.Avatar, &t.AvatarUrl, &t.TimeZone, &t.ZoneOffset, &t.AfkFor)

		getAll, errGet := d.TechGetAll(t)
		if errGet != nil {
			d.log.InfoStruct("TechGetAll(t)", t)
			return nil, errGet
		}
		user, erru := d.UsersGetByUserId(t.UserId)
		if erru != nil {
			d.log.InfoStruct("UsersGetByUserId", t)
			d.log.ErrorErr(erru)
		}

		for _, member := range getAll {
			if user != nil && member.Name == user.Username && user.GameName != "" {
				member.Name = user.GameName
			}
			mm = append(mm, member)
		}
	}
	return mm, nil
}

func getTimeStrings(offset int) (string, string) {
	// Получаем текущее время в UTC
	now := time.Now().UTC()

	// Применяем смещение к текущему времени в UTC
	offsetDuration := time.Duration(offset) * time.Minute
	timeWithOffset := now.Add(offsetDuration)

	// Форматируем время в 12-часовом формате с AM/PM
	time12HourFormat := timeWithOffset.Format("03:04 PM")

	// Форматируем время в 24-часовом формате
	time24HourFormat := timeWithOffset.Format("15:04")

	return time12HourFormat, time24HourFormat
}

// getTimeStringsWithDST возвращает время с учетом DST (летнего/зимнего времени)
// Если timezone - это название локации (например, "America/New_York"),
// то смещение вычисляется динамически с учетом текущего DST
func getTimeStringsWithDST(timezone string, fallbackOffsetMinutes int) (string, string) {
	now := time.Now()

	// Пытаемся загрузить локацию для динамического расчета DST
	if timezone != "" {
		loc, err := time.LoadLocation(timezone)
		if err == nil {
			timeInLocation := now.In(loc)
			return timeInLocation.Format("03:04 PM"), timeInLocation.Format("15:04")
		}
	}

	// Fallback на фиксированное смещение
	return getTimeStrings(fallbackOffsetMinutes)
}
