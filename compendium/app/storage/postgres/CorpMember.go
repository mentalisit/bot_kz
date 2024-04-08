package postgres

import (
	"compendium/models"
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v4"
)

func (d *Db) corpMemberInsert(ctx context.Context, guildid string, u models.CorpMember) {
	// Проверяем существует ли запись для данного игрока в таблице
	var existingGuildID string
	err := d.db.QueryRow(ctx, "SELECT guildid FROM compendium.corpmember WHERE userid = $1", u.UserId).Scan(&existingGuildID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		// Если запись не найдена, вставляем новую запись
		Tech, err := json.Marshal(u.Tech)
		if err != nil {
			d.log.Info(err.Error())
		}

		insert := `INSERT INTO compendium.corpmember(guildid, name, userid, clientuserid, avatar, tech, avatarurl, timezona, zonaoffset, afkfor) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
		_, err = d.db.Exec(ctx, insert, guildid, u.Name, u.UserId, u.ClientUserId, u.Avatar, Tech, u.AvatarUrl, u.TimeZone, u.ZoneOffset, u.AfkFor)
		if err != nil {
			d.log.ErrorErr(err)
		}
		return
	case err != nil:
		d.log.ErrorErr(err)
		return
	default:
		// Если запись найдена, проверяем совпадает ли guildID
		if existingGuildID != guildid {
			Tech, err := json.Marshal(u.Tech)
			if err != nil {
				d.log.Info(err.Error())
			}

			insert := `INSERT INTO compendium.corpmember(guildid, name, userid, clientuserid, avatar, tech, avatarurl, timezona, zonaoffset, afkfor) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
			_, err = d.db.Exec(ctx, insert, guildid, u.Name, u.UserId, u.ClientUserId, u.Avatar, Tech, u.AvatarUrl, u.TimeZone, u.ZoneOffset, u.AfkFor)
			if err != nil {
				d.log.ErrorErr(err)
			}
			return
		}
	}
}
func (d *Db) CorpMemberReadAllByGuildId(ctx context.Context, guildid string) []models.CorpMemberint {
	sel := "SELECT * FROM compendium.corpmember WHERE guildid = $1"
	results, err := d.db.Query(ctx, sel, guildid)
	if err != nil {
		d.log.ErrorErr(err)
	}
	var tt []models.CorpMemberint
	for results.Next() {
		var t models.CorpMemberint
		var TechData []byte
		ttt := make(map[int]models.TechLevel)
		err = results.Scan(&guildid, &t.Name, &t.UserId, &t.ClientUserId, &t.Avatar, &TechData, &t.AvatarUrl, &t.TimeZone, &t.ZoneOffset, &t.AfkFor)
		err = json.Unmarshal(TechData, &ttt)
		t.Tech = make(map[int][2]int)
		for i, level := range ttt {
			t.Tech[i] = [2]int{level.Level}
			//fmt.Println("t ", i, level)
		}

		if err != nil {
			d.log.Info(err.Error())
		}
		tt = append(tt, t)
	}
	return tt
}
func (d *Db) CorpMemberReadByUserId(ctx context.Context, userId string) models.CorpMember {
	sel := "SELECT * FROM compendium.corpmember WHERE userid = $1"
	results, err := d.db.Query(ctx, sel, userId)
	if err != nil {
		d.log.ErrorErr(err)
	}
	var guildid string
	var t models.CorpMember
	for results.Next() {
		var TechData []byte
		err = results.Scan(&guildid, &t.Name, &t.UserId, &t.ClientUserId, &t.Avatar, &TechData, &t.AvatarUrl, &t.TimeZone, &t.ZoneOffset, &t.AfkFor)
		if err != nil {
			d.log.Info(err.Error())
		}
		err = json.Unmarshal(TechData, &t.Tech)
		if err != nil {
			d.log.Info(err.Error())
		}
	}
	return t
}
func (d *Db) CorpMemberTechUpdate(ctx context.Context, userid string, tech models.TechLevels) {
	Tech, err := json.Marshal(tech)
	if err != nil {
		d.log.Info(err.Error())
	}
	sqlUpd := `update compendium.corpmember set tech = $1 where userid = $2`
	_, err = d.db.Exec(ctx, sqlUpd, Tech, userid)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) CorpMemberReadByUserIdByGuildId(ctx context.Context, userId, guildId string) models.CorpMemberint {
	sel := "SELECT * FROM compendium.corpmember WHERE userid = $1 AND guildid = $2"
	results, err := d.db.Query(ctx, sel, userId, guildId)
	if err != nil {
		d.log.ErrorErr(err)
	}
	var guildid string
	var t models.CorpMemberint
	for results.Next() {
		var TechData []byte
		ttt := make(map[int]models.TechLevel)
		err = results.Scan(&guildid, &t.Name, &t.UserId, &t.ClientUserId, &t.Avatar, &TechData, &t.AvatarUrl, &t.TimeZone, &t.ZoneOffset, &t.AfkFor)
		if err != nil {
			d.log.Info(err.Error())
		}
		err = json.Unmarshal(TechData, &ttt)
		if err != nil {
			d.log.Info(err.Error())
		}
		t.Tech = make(map[int][2]int)
		for i, level := range ttt {
			t.Tech[i] = [2]int{level.Level}
		}
	}

	return t
}

func (d *Db) CorpMemberReadByNameByGuildId(ctx context.Context, Name, guildId string) models.CorpMemberint {
	sel := "SELECT * FROM compendium.corpmember WHERE name = $1 AND guildid = $2"
	results, err := d.db.Query(ctx, sel, Name, guildId)
	if err != nil {
		d.log.ErrorErr(err)
	}
	var guildid string
	var t models.CorpMemberint
	for results.Next() {
		var TechData []byte
		ttt := make(map[int]models.TechLevel)
		err = results.Scan(&guildid, &t.Name, &t.UserId, &t.ClientUserId, &t.Avatar, &TechData, &t.AvatarUrl, &t.TimeZone, &t.ZoneOffset, &t.AfkFor)
		if err != nil {
			d.log.Info(err.Error())
		}
		err = json.Unmarshal(TechData, &ttt)
		if err != nil {
			d.log.Info(err.Error())
		}
		t.Tech = make(map[int][2]int)
		for i, level := range ttt {
			t.Tech[i] = [2]int{level.Level}
		}
	}

	return t
}
