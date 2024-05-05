package postgresOLD

import (
	"compendium/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
)

func (d *Db) corpMemberInsert(ctx context.Context, guildid string, u models.CorpMember) {
	var existingGuildID string
	err := d.db.QueryRow(ctx, "SELECT guildid FROM compendium.corpmember WHERE userid = $1 AND guildid = $2 AND name = $3", u.UserId, guildid, u.Name).Scan(&existingGuildID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			fmt.Println("pgx.ErrNoRows")
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
		}
	}
	fmt.Println("существует в бд")
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
		var id int
		ttt := make(map[int]models.TechLevel)
		err = results.Scan(&id, &guildid, &t.Name, &t.UserId, &t.ClientUserId, &t.Avatar, &TechData, &t.AvatarUrl, &t.TimeZone, &t.ZoneOffset, &t.AfkFor)
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
func (d *Db) CorpMemberReadByUserId(ctx context.Context, userId, guildid string) models.CorpMember {
	sel := "SELECT * FROM compendium.corpmember WHERE userid = $1 AND guildid = $2"
	results, err := d.db.Query(ctx, sel, userId, guildid)
	if err != nil {
		d.log.ErrorErr(err)
	}
	var t models.CorpMember
	for results.Next() {
		var TechData []byte
		var id int
		err = results.Scan(&id, &guildid, &t.Name, &t.UserId, &t.ClientUserId, &t.Avatar, &TechData, &t.AvatarUrl, &t.TimeZone, &t.ZoneOffset, &t.AfkFor)
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
func (d *Db) CorpMemberReadByUserIdByName(ctx context.Context, userId, guildid, name string) models.CorpMember {
	sel := "SELECT * FROM compendium.corpmember WHERE userid = $1 AND guildid = $2 AND name = $3"
	results, err := d.db.Query(ctx, sel, userId, guildid, name)
	if err != nil {
		d.log.ErrorErr(err)
	}
	var t models.CorpMember
	for results.Next() {
		var TechData []byte
		var id int
		err = results.Scan(&id, &guildid, &t.Name, &t.UserId, &t.ClientUserId, &t.Avatar, &TechData, &t.AvatarUrl, &t.TimeZone, &t.ZoneOffset, &t.AfkFor)
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

func (d *Db) CorpMemberTechUpdate(ctx context.Context, userid, guildid, name string, tech models.TechLevels) {
	Tech, err := json.Marshal(tech)
	if err != nil {
		d.log.Info(err.Error())
	}
	sqlUpd := `update compendium.corpmember set tech = $1 where userid = $2 AND guildid = $3 AND name = $4`
	upd, err := d.db.Exec(ctx, sqlUpd, Tech, userid, guildid, name)
	ErrNoRows := false
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ErrNoRows = true
		} else {
			d.log.ErrorErr(err)
		}
	}
	if upd.RowsAffected() == 0 || ErrNoRows {
		member := d.CorpMemberReadByUserId(ctx, userid, guildid)
		member.Name = name
		member.Tech = tech
		d.corpMemberInsert(ctx, guildid, member)
	}
}
func (d *Db) CorpMemberReadByUserIdByGuildIdByName(ctx context.Context, userId, guildId, name string) models.CorpMemberint {
	sel := "SELECT * FROM compendium.corpmember WHERE userid = $1 AND guildid = $2 AND name = $3"
	results, err := d.db.Query(ctx, sel, userId, guildId, name)
	if err != nil {
		d.log.ErrorErr(err)
	}
	var guildid string
	var t models.CorpMemberint
	for results.Next() {
		var TechData []byte
		var id int
		ttt := make(map[int]models.TechLevel)
		err = results.Scan(&id, &guildid, &t.Name, &t.UserId, &t.ClientUserId, &t.Avatar, &TechData, &t.AvatarUrl, &t.TimeZone, &t.ZoneOffset, &t.AfkFor)
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
		var id int
		ttt := make(map[int]models.TechLevel)
		err = results.Scan(&id, &guildid, &t.Name, &t.UserId, &t.ClientUserId, &t.Avatar, &TechData, &t.AvatarUrl, &t.TimeZone, &t.ZoneOffset, &t.AfkFor)
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

func (d *Db) CorpMemberReadByNameByGuildId(ctx context.Context, Name, guildid string) models.CorpMemberint {
	sel := "SELECT * FROM compendium.corpmember WHERE name = $1 AND guildid = $2"
	results, err := d.db.Query(ctx, sel, Name, guildid)
	if err != nil {
		d.log.ErrorErr(err)
	}
	var t models.CorpMemberint
	for results.Next() {
		var TechData []byte
		var id int
		ttt := make(map[int]models.TechLevel)
		err = results.Scan(&id, &guildid, &t.Name, &t.UserId, &t.ClientUserId, &t.Avatar, &TechData, &t.AvatarUrl, &t.TimeZone, &t.ZoneOffset, &t.AfkFor)
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
