package multi

import (
	"compendium_s/models"
	"github.com/google/uuid"
	"time"
)

func (d *Db) CorpMemberInsert(cm models.MultiAccountCorpMember) error {
	ctx, cancel := d.getContext()
	defer cancel()
	var count int
	sel := "SELECT count(*) as count FROM compendium.corpmember WHERE uid = $1 "
	err := d.db.QueryRow(ctx, sel, cm.Uid).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		insert := `INSERT INTO compendium.corpmember(uid,guildids, timezona, zonaoffset, afkfor) VALUES ($1,$2,$3,$4,$5)`
		_, err = d.db.Exec(ctx, insert, cm.Uid, cm.GuildIds, cm.TimeZona, cm.ZonaOffset, cm.AfkFor)
		if err != nil {
			return err
		}
	} else {
		err = d.CorpMemberUpdateGuildIds(cm)
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *Db) CorpMemberUpdateGuildIds(cm models.MultiAccountCorpMember) error {
	ctx, cancel := d.getContext()
	defer cancel()
	sqlUpd := `update compendium.corpmember set guildids = $1 where uid = $2 `
	_, err := d.db.Exec(ctx, sqlUpd, cm.GuildIds, cm.Uid)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) CorpMembersRead(GID uuid.UUID) ([]models.CorpMember, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	sel := "SELECT * FROM compendium.corpmember WHERE $1 = ANY(guildids)"
	results, err := d.db.Query(ctx, sel, GID)
	defer results.Close()
	if err != nil {
		return nil, err
	}
	var mm []models.CorpMember
	for results.Next() {
		var t models.MultiAccountCorpMember
		err = results.Scan(&t.Uid, &t.GuildIds, &t.TimeZona, &t.ZonaOffset, &t.AfkFor)
		m := models.CorpMember{
			TimeZone:    t.TimeZona,
			ZoneOffset:  t.ZonaOffset,
			AfkFor:      t.AfkFor,
			TypeAccount: "ma",
		}
		multiMember, err := d.FindMultiAccountUUID(t.Uid)
		if err != nil {
			d.log.ErrorErr(err)
			continue
		}
		m.Multi = multiMember
		m.Name = multiMember.Nickname
		m.UserId = multiMember.UUID.String()
		m.AvatarUrl = multiMember.AvatarURL
		m.GuildId = GID.String()

		if m.TimeZone != "" {
			t12, t24 := getTimeStrings(m.ZoneOffset)
			m.LocalTime = t12
			m.LocalTime24 = t24
		}

		membersUser, err := d.TechnologiesGetAllCorpMember(m, t.Uid)
		if err != nil {
			return nil, err
		}

		mm = append(mm, membersUser...)
	}
	return mm, nil
}

func (d *Db) CorpMembersApiRead(uid uuid.UUID) (*models.MultiAccountCorpMember, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	sel := "SELECT * FROM compendium.corpmember WHERE uid = $1 "
	results := d.db.QueryRow(ctx, sel, uid)
	var t models.MultiAccountCorpMember
	err := results.Scan(&t.Uid, &t.GuildIds, &t.TimeZona, &t.ZonaOffset, &t.AfkFor)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (d *Db) CorpMemberTZUpdate(uid uuid.UUID, timeZone string, offset int) error {
	ctx, cancel := d.getContext()
	defer cancel()
	sqlUpd := `update compendium.corpmember set zonaoffset = $1, timezona = $2 where uid = $3 `
	_, err := d.db.Exec(ctx, sqlUpd, offset, timeZone, uid)
	if err != nil {
		return err
	}
	return nil
}

func (d *Db) CorpMemberByUId(uid uuid.UUID) (*models.MultiAccountCorpMember, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var u models.MultiAccountCorpMember
	selectUser := "SELECT * FROM compendium.corpmember WHERE uid = $1 "
	err := d.db.QueryRow(ctx, selectUser, uid).Scan(&u.Uid, &u.GuildIds, &u.TimeZona, &u.ZonaOffset, &u.AfkFor)
	if err != nil {
		return nil, err
	}
	return &u, nil
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
