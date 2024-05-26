package postgres

import (
	"context"
)

//func (d *Db) ListUserInsert(token, userid, guildid string) error {
//	count, err := d.ListUserGetCountByGuildIdByUserId(guildid, userid)
//	if err != nil {
//		return err
//	}
//	if count > 0 {
//		err = d.ListUserUpdate(token, userid, guildid)
//		if err != nil {
//			return err
//		}
//	} else {
//		insert := `INSERT INTO hs_compendium.list_users(token, userid, guildid) VALUES ($1,$2,$3)`
//		_, err = d.db.Exec(context.Background(), insert, token, userid, guildid)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

//func (d *Db) ListUserGetCountByGuildIdByUserId(guildid, userid string) (int, error) {
//	var count int
//	sel := "SELECT count(*) as count FROM hs_compendium.list_users WHERE guildid = $1 AND userid = $2"
//	err := d.db.QueryRow(context.Background(), sel, guildid, userid).Scan(&count)
//	if err != nil {
//		return 0, err
//	}
//	return count, nil
//}
//
//func (d *Db) ListUserUpdate(token, userid, guildid string) error {
//	upd := `update hs_compendium.list_users set token = $1 where guildid = $2 AND userid = $3`
//	_, err := d.db.Exec(context.Background(), upd, token, guildid, userid)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (d *Db) ListUserGetToken(userid, guildid string) (string, error) {
//	var token string
//	selectUser := "SELECT token FROM hs_compendium.list_users WHERE userid = $1 AND guildid = $2"
//	err := d.db.QueryRow(context.Background(), selectUser, userid, guildid).Scan(&token)
//	if err != nil {
//		return "", err
//	}
//	return token, nil
//}

func (d *Db) ListUserGetUserIdAndGuildId(token string) (userid string, guildid string, err error) {
	selectUser := "SELECT userid,guildid FROM hs_compendium.list_users WHERE token = $1"
	err = d.db.QueryRow(context.Background(), selectUser, token).Scan(&userid, &guildid)
	if err != nil {
		return "", "", err
	}
	return userid, guildid, nil
}
