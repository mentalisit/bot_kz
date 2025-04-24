package postgres

import "github.com/google/uuid"

func (d *Db) ListUserInsert(token, userid, guildid string) error {
	ctx, cancel := d.getContext()
	defer cancel()
	count, err := d.ListUserGetCountByGuildIdByUserId(guildid, userid)
	if err != nil {
		return err
	}
	if count > 0 {
		err = d.ListUserUpdate(token, userid, guildid)
		if err != nil {
			return err
		}
	} else {
		insert := `INSERT INTO hs_compendium.list_users(token, userid, guildid) VALUES ($1,$2,$3)`
		_, err = d.db.Exec(ctx, insert, token, userid, guildid)
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *Db) ListUserGetCountByGuildIdByUserId(guildid, userid string) (int, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var count int
	sel := "SELECT count(*) as count FROM hs_compendium.list_users WHERE guildid = $1 AND userid = $2"
	err := d.db.QueryRow(ctx, sel, guildid, userid).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (d *Db) ListUserUpdate(token, userid, guildid string) error {
	ctx, cancel := d.getContext()
	defer cancel()
	upd := `update hs_compendium.list_users set token = $1 where guildid = $2 AND userid = $3`
	_, err := d.db.Exec(ctx, upd, token, guildid, userid)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) ListUserGetToken(userid, guildid string) (string, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var token string
	selectUser := "SELECT token FROM hs_compendium.list_users WHERE userid = $1 AND guildid = $2"
	err := d.db.QueryRow(ctx, selectUser, userid, guildid).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}
func (d *Db) ListUserGetUserIdAndGuildId(token string) (userid string, guildid string, err error) {
	ctx, cancel := d.getContext()
	defer cancel()
	selectUser := "SELECT userid,guildid FROM hs_compendium.list_users WHERE token = $1"
	err = d.db.QueryRow(ctx, selectUser, token).Scan(&userid, &guildid)
	if err != nil {
		return "", "", err
	}
	return userid, guildid, nil
}
func (d *Db) ListUserUpdateToken(tokenOld, tokenNew string) error {
	ctx, cancel := d.getContext()
	defer cancel()
	upd := `update hs_compendium.list_users set token = $1 where token = $2`
	_, err := d.db.Exec(ctx, upd, tokenNew, tokenOld)
	if err != nil {
		return err
	}
	return nil
}

func (d *Db) MultiListUserInsert(token, guildId string, uid uuid.UUID) error {
	ctx, cancel := d.getContext()
	defer cancel()
	count, err := d.MultiListUserGetCountByUUID(guildId, uid)
	if err != nil {
		return err
	}
	if count > 0 {
		err = d.MultiListUserUpdate(token, guildId, uid)
		if err != nil {
			return err
		}
	} else {
		insert := `INSERT INTO compendium.list_users(uid, guildid, token) VALUES ($1,$2,$3)`
		_, err = d.db.Exec(ctx, insert, uid, guildId, token)
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *Db) MultiListUserGetCountByUUID(guildId string, uid uuid.UUID) (int, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var count int
	sel := "SELECT count(*) as count FROM compendium.list_users WHERE uid = $1 AND guildid = $2"
	err := d.db.QueryRow(ctx, sel, uid, guildId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (d *Db) MultiListUserUpdate(token, guildId string, uid uuid.UUID) error {
	ctx, cancel := d.getContext()
	defer cancel()
	upd := `update compendium.list_users set token = $1 where guildid = $2 AND uid = $3`
	_, err := d.db.Exec(ctx, upd, token, guildId, uid)
	if err != nil {
		return err
	}
	return nil
}

//todo get token
