package postgres

import (
	"compendium/models"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
)

func (d *Db) CorpMemberInsert(cm models.CorpMember) error {
	ctx, cancel := d.getContext()
	defer cancel()
	var count int
	selCount := "SELECT count(*) as count FROM hs_compendium.corpmember WHERE guildid = $1 AND userid = $2"
	if len(cm.GuildId) < 24 {
		_ = d.db.QueryRow(ctx, selCount, cm.GuildId, cm.UserId).Scan(&count)
		if count != 0 {
			sqlUpd := `update hs_compendium.corpmember set guildid = $1 where userid = $2 AND guildid = $3`
			_, _ = d.db.Exec(ctx, sqlUpd, cm.MGuild.GuildId(), cm.UserId, cm.GuildId)
		}
	} else {
		_ = d.db.QueryRow(ctx, selCount, cm.MGuild.GuildId(), cm.UserId).Scan(&count)
		if count == 0 {
			insert := `INSERT INTO hs_compendium.corpmember(username, userid, guildid, avatar, avatarurl, timezona, zonaoffset, afkfor) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
			_, err := d.db.Exec(ctx, insert, cm.Name, cm.UserId, cm.MGuild.GuildId(), "", cm.AvatarUrl, cm.TimeZone, cm.ZoneOffset, cm.AfkFor)
			if err != nil {
				return err
			}
		}
	}

	techBytes, err := json.Marshal(cm.Tech)
	if err != nil {
		return err
	}

	if len(cm.GuildId) < 24 {
		count, _ = d.TechGetCount(cm.UserId)
		if count != 0 {
			upd := `update hs_compendium.tech set guildid = $1 where userid = $2 and guildid = $3`
			_, _ = d.db.Exec(ctx, upd, cm.MGuild.GuildId(), cm.UserId, cm.GuildId)
			return nil
		}
	}

	err = d.TechInsert(cm.Name, cm.UserId, techBytes)
	if err != nil {
		return err
	}
	return nil
}
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
		var garbich string
		//ttt := make(map[int]models.TechLevel)
		err = results.Scan(&id, &t.Name, &t.UserId, &t.GuildId, &garbich, &t.AvatarUrl, &t.TimeZone, &t.ZoneOffset, &t.AfkFor)

		getAll, errGet := d.TechGetAll(t)
		if errGet != nil {
			return nil, errGet
		}
		user, erru := d.UsersGetByUserId(t.UserId)
		if erru != nil {
			d.UsersInsert(models.User{
				ID:        t.UserId,
				Username:  t.Name,
				AvatarURL: t.AvatarUrl,
				Alts:      []string{},
			})
			d.log.InfoStruct("New UsersGetByUserId", t)
			d.log.ErrorErr(err)
		}

		for _, member := range getAll {
			if user != nil && member.Name == user.Username && user.GameName != "" {
				member.Name = user.GameName
			}
			// Используем функцию с поддержкой DST
			if member.TimeZone != "" {
				t12, t24 := getTimeStringsWithDST(member.TimeZone, member.ZoneOffset)
				member.LocalTime = t12
				member.LocalTime24 = t24
			}
			mm = append(mm, member)
		}

	}
	return mm, nil
}

func (d *Db) CorpMembersApiRead(guildid, userid string) ([]models.CorpMember, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	sel := "SELECT * FROM hs_compendium.corpmember WHERE guildid = $1 AND userid = $2"
	results, err := d.db.Query(ctx, sel, guildid, userid)
	defer results.Close()
	if err != nil {
		return nil, err
	}
	var mm []models.CorpMember
	for results.Next() {
		var t models.CorpMember
		var id int
		var garbich string
		err = results.Scan(&id, &t.Name, &t.UserId, &t.GuildId, &garbich, &t.AvatarUrl, &t.TimeZone, &t.ZoneOffset, &t.AfkFor)

		mm, err = d.TechGetAll(t)
		if err != nil {
			return nil, err
		}
		return mm, nil

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
func (d *Db) CorpMemberTZUpdate(userid, guildid, timeZone string, offset int) error {
	ctx, cancel := d.getContext()
	defer cancel()
	sqlUpd := `update hs_compendium.corpmember set zonaoffset = $1,timezona = $2 where userid = $3 AND guildid = $4`
	row, err := d.db.Exec(ctx, sqlUpd, offset, timeZone, userid, guildid)
	if row.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return err
}

func (d *Db) CorpMemberByUserId(userId string) (*models.CorpMember, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var u models.CorpMember
	var id int
	var garbich string
	selectUser := "SELECT * FROM hs_compendium.corpmember WHERE userid = $1 "
	err := d.db.QueryRow(ctx, selectUser, userId).Scan(&id, &u.Name, &u.UserId, &u.GuildId, &garbich, &u.AvatarUrl, &u.TimeZone, &u.ZoneOffset, &u.AfkFor)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
func (d *Db) CorpMemberAvatarUpdate(userid, guildid, avatarurl string) error {
	ctx, cancel := d.getContext()
	defer cancel()
	sqlUpd := `update hs_compendium.corpmember set avatarurl = $1 where userid = $2 AND guildid = $3`
	row, err := d.db.Exec(ctx, sqlUpd, avatarurl, userid, guildid)
	if row.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	if err != nil {
		return err
	}
	user, err := d.UsersGetByUserId(userid)
	if err != nil {
		return err
	}
	u := *user
	u.AvatarURL = avatarurl

	err = d.UsersUpdate(u)
	if err != nil {
		return err
	}

	return nil
}
func (d *Db) CorpMemberDelete(guildid string, nameId string) error {
	ctx, cancel := d.getContext()
	defer cancel()
	deleteMember := `DELETE FROM hs_compendium.corpmember WHERE guildid = $1 AND userid = $2`
	_, err := d.db.Exec(ctx, deleteMember, guildid, nameId)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) CorpMemberDeleteAlt(guildid string, nameId string, name string) error {
	ctx, cancel := d.getContext()
	defer cancel()
	deleteMember := `DELETE FROM hs_compendium.corpmember WHERE guildid = $1 AND userid = $2 AND username = $3`
	_, err := d.db.Exec(ctx, deleteMember, guildid, nameId, name)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) CorpMembersReadByUserId(UserId string) ([]models.CorpMember, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	sel := "SELECT * FROM hs_compendium.corpmember WHERE userid = $1"
	results, err := d.db.Query(ctx, sel, UserId)
	defer results.Close()
	if err != nil {
		return nil, err
	}
	var mm []models.CorpMember
	var garbich string
	for results.Next() {
		var t models.CorpMember
		var id int
		err = results.Scan(&id, &t.Name, &t.UserId, &t.GuildId, &garbich, &t.AvatarUrl, &t.TimeZone, &t.ZoneOffset, &t.AfkFor)

		mm = append(mm, t)
	}
	return mm, nil
}
