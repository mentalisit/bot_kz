package postgres

import (
	"compendium_s/models"
	"time"
)

//func (d *Db) CorpMemberInsert(cm models.CorpMember) error {
//	var count int
//	sel := "SELECT count(*) as count FROM hs_compendium.corpmember WHERE guildid = $1 AND userid = $2"
//	err := d.db.QueryRow(context.Background(), sel, cm.GuildId, cm.UserId).Scan(&count)
//	if err != nil {
//		return err
//	}
//	if count == 0 {
//		insert := `INSERT INTO hs_compendium.corpmember(username, userid, guildid, avatar, avatarurl, timezona, zonaoffset, afkfor) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
//		_, err = d.db.Exec(context.Background(), insert, cm.Name, cm.UserId, cm.GuildId, cm.Avatar, cm.AvatarUrl, cm.TimeZone, cm.ZoneOffset, cm.AfkFor)
//		if err != nil {
//			return err
//		}
//	}
//	techBytes, err := json.Marshal(cm.Tech)
//	if err != nil {
//		return err
//	}
//	if len(cm.Tech) == 0 {
//		tech := make(map[int]models.TechLevel)
//		tech[701] = models.TechLevel{
//			Ts:    0,
//			Level: 0,
//		}
//		techBytes, _ = json.Marshal(tech)
//	}
//	err = d.TechInsert(cm.Name, cm.UserId, cm.GuildId, techBytes)
//	if err != nil {
//		return err
//	}
//	return nil
//}

func (d *Db) CorpMembersRead(guildid string) ([]models.CorpMember, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	sel := "SELECT * FROM hs_compendium.corpmember WHERE guildid = $1"
	results, err := d.db.Query(ctx, sel, guildid)
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
	ctx, cancel := d.getContext()
	defer cancel()
	sel := "SELECT * FROM hs_compendium.corpmember WHERE userid = $1"
	results, err := d.db.Query(ctx, sel, userid)
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
