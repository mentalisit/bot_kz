package postgres

import (
	"compendium/models"
	"context"
	"encoding/json"
)

func (d *Db) CorpMemberInsert(cm models.CorpMember) error {
	var count int
	sel := "SELECT count(*) as count FROM hs_compendium.corpmember WHERE guildid = $1 AND userid = $2"
	err := d.db.QueryRow(context.Background(), sel, cm.GuildId, cm.UserId).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		insert := `INSERT INTO hs_compendium.corpmember(username, userid, guildid, avatar, avatarurl, timezona, zonaoffset, afkfor) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
		_, err = d.db.Exec(context.Background(), insert, cm.Name, cm.UserId, cm.GuildId, cm.Avatar, cm.AvatarUrl, cm.TimeZone, cm.ZoneOffset, cm.AfkFor)
		if err != nil {
			return err
		}
	}
	techBytes, err := json.Marshal(cm.Tech)
	if err != nil {
		return err
	}
	err = d.TechInsert(cm.Name, cm.UserId, cm.GuildId, techBytes)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) CorpMembersRead(guildid string) ([]models.CorpMember, error) {
	sel := "SELECT * FROM hs_compendium.corpmember WHERE guildid = $1"
	results, err := d.db.Query(context.Background(), sel, guildid)
	if err != nil {
		return nil, err
	}
	var mm []models.CorpMember
	for results.Next() {
		var t models.CorpMember
		var id int
		ttt := make(map[int]models.TechLevel)
		err = results.Scan(&id, &t.Name, &t.UserId, &t.GuildId, &t.Avatar, &t.AvatarUrl, &t.TimeZone, &t.ZoneOffset, &t.AfkFor)

		bytes, errGet := d.TechGet(t.Name, t.UserId, t.GuildId)
		if errGet != nil {
			return nil, errGet
		}
		err = json.Unmarshal(bytes, &ttt)
		if err != nil {
			return nil, err
		}
		t.Tech = make(map[int][2]int)
		for i, level := range ttt {
			t.Tech[i] = [2]int{level.Level, int(level.Ts)}
			//fmt.Println("t ", i, level)
		}
		mm = append(mm, t)
	}
	return mm, nil
}

//func (d *Db) CorpMemberReadByUserId(ctx context.Context, userId, guildid string) models.CorpMember {
//	sel := "SELECT * FROM compendium.corpmember WHERE userid = $1 AND guildid = $2"
//	results, err := d.db.Query(ctx, sel, userId, guildid)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//	var t models.CorpMember
//	for results.Next() {
//		var TechData []byte
//		var id int
//		err = results.Scan(&id, &guildid, &t.Name, &t.UserId, &t.ClientUserId, &t.Avatar, &TechData, &t.AvatarUrl, &t.TimeZone, &t.ZoneOffset, &t.AfkFor)
//		if err != nil {
//			d.log.Info(err.Error())
//		}
//		err = json.Unmarshal(TechData, &t.Tech)
//		if err != nil {
//			d.log.Info(err.Error())
//		}
//	}
//	return t
//}
//func (d *Db) CorpMemberReadByUserIdByName(ctx context.Context, userId, guildid, name string) models.CorpMember {
//	sel := "SELECT * FROM compendium.corpmember WHERE userid = $1 AND guildid = $2 AND name = $3"
//	results, err := d.db.Query(ctx, sel, userId, guildid, name)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//	var t models.CorpMember
//	for results.Next() {
//		var TechData []byte
//		var id int
//		err = results.Scan(&id, &guildid, &t.Name, &t.UserId, &t.ClientUserId, &t.Avatar, &TechData, &t.AvatarUrl, &t.TimeZone, &t.ZoneOffset, &t.AfkFor)
//		if err != nil {
//			d.log.Info(err.Error())
//		}
//		err = json.Unmarshal(TechData, &t.Tech)
//		if err != nil {
//			d.log.Info(err.Error())
//		}
//	}
//	return t
//}
//
//func (d *Db) CorpMemberTechUpdate(ctx context.Context, userid, guildid, name string, tech models.TechLevels) {
//	Tech, err := json.Marshal(tech)
//	if err != nil {
//		d.log.Info(err.Error())
//	}
//	sqlUpd := `update compendium.corpmember set tech = $1 where userid = $2 AND guildid = $3 AND name = $4`
//	upd, err := d.db.Exec(ctx, sqlUpd, Tech, userid, guildid, name)
//	ErrNoRows := false
//	if err != nil {
//		if errors.Is(err, pgx.ErrNoRows) {
//			ErrNoRows = true
//		} else {
//			d.log.ErrorErr(err)
//		}
//	}
//	if upd.RowsAffected() == 0 || ErrNoRows {
//		member := d.CorpMemberReadByUserId(ctx, userid, guildid)
//		member.Name = name
//		member.Tech = tech
//		d.corpMemberInsert(ctx, guildid, member)
//	}
//}
//func (d *Db) CorpMemberReadByUserIdByGuildIdByName(ctx context.Context, userId, guildId, name string) models.CorpMemberint {
//	sel := "SELECT * FROM compendium.corpmember WHERE userid = $1 AND guildid = $2 AND name = $3"
//	results, err := d.db.Query(ctx, sel, userId, guildId, name)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//	var guildid string
//	var t models.CorpMemberint
//	for results.Next() {
//		var TechData []byte
//		var id int
//		ttt := make(map[int]models.TechLevel)
//		err = results.Scan(&id, &guildid, &t.Name, &t.UserId, &t.ClientUserId, &t.Avatar, &TechData, &t.AvatarUrl, &t.TimeZone, &t.ZoneOffset, &t.AfkFor)
//		if err != nil {
//			d.log.Info(err.Error())
//		}
//		err = json.Unmarshal(TechData, &ttt)
//		if err != nil {
//			d.log.Info(err.Error())
//		}
//		t.Tech = make(map[int][2]int)
//		for i, level := range ttt {
//			t.Tech[i] = [2]int{level.Level}
//		}
//	}
//
//	return t
//}
//func (d *Db) CorpMemberReadByUserIdByGuildId(ctx context.Context, userId, guildId string) models.CorpMemberint {
//	sel := "SELECT * FROM compendium.corpmember WHERE userid = $1 AND guildid = $2"
//	results, err := d.db.Query(ctx, sel, userId, guildId)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//	var guildid string
//	var t models.CorpMemberint
//	for results.Next() {
//		var TechData []byte
//		var id int
//		ttt := make(map[int]models.TechLevel)
//		err = results.Scan(&id, &guildid, &t.Name, &t.UserId, &t.ClientUserId, &t.Avatar, &TechData, &t.AvatarUrl, &t.TimeZone, &t.ZoneOffset, &t.AfkFor)
//		if err != nil {
//			d.log.Info(err.Error())
//		}
//		err = json.Unmarshal(TechData, &ttt)
//		if err != nil {
//			d.log.Info(err.Error())
//		}
//		t.Tech = make(map[int][2]int)
//		for i, level := range ttt {
//			t.Tech[i] = [2]int{level.Level}
//		}
//	}
//
//	return t
//}
//
//func (d *Db) CorpMemberReadByNameByGuildId(ctx context.Context, Name, guildid string) models.CorpMemberint {
//	sel := "SELECT * FROM compendium.corpmember WHERE name = $1 AND guildid = $2"
//	results, err := d.db.Query(ctx, sel, Name, guildid)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//	var t models.CorpMemberint
//	for results.Next() {
//		var TechData []byte
//		var id int
//		ttt := make(map[int]models.TechLevel)
//		err = results.Scan(&id, &guildid, &t.Name, &t.UserId, &t.ClientUserId, &t.Avatar, &TechData, &t.AvatarUrl, &t.TimeZone, &t.ZoneOffset, &t.AfkFor)
//		if err != nil {
//			d.log.Info(err.Error())
//		}
//		err = json.Unmarshal(TechData, &ttt)
//		if err != nil {
//			d.log.Info(err.Error())
//		}
//		t.Tech = make(map[int][2]int)
//		for i, level := range ttt {
//			t.Tech[i] = [2]int{level.Level}
//		}
//	}
//
//	return t
//}
