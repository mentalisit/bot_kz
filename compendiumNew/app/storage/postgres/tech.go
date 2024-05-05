package postgres

import "context"

func (d *Db) TechInsert(username, userid, guildid string, tech []byte) error {
	var count int
	sel := "SELECT count(*) as count FROM hs_compendium.tech WHERE guildid = $1 AND userid = $2"
	err := d.db.QueryRow(context.Background(), sel, guildid, userid).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		insert := `INSERT INTO hs_compendium.tech(username, userid, guildid, tech) VALUES ($1,$2,$3,$4)`
		_, err = d.db.Exec(context.Background(), insert, username, userid, guildid, tech)
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *Db) TechGet(username, userid, guildid string) ([]byte, error) {
	var tech []byte
	sel := "SELECT tech FROM hs_compendium.tech WHERE userid = $1 AND guildid = $2 AND username = $3"
	err := d.db.QueryRow(context.Background(), sel, userid, guildid, username).Scan(&tech)
	if err != nil {
		return nil, err
	}
	return tech, nil
}
func (d *Db) TechUpdate(username, userid, guildid string, tech []byte) error {
	upd := `update hs_compendium.tech set tech = $1 where username = $2 and userid = $3 and guildid = $4`
	_, err := d.db.Exec(context.Background(), upd, tech, username, userid, guildid)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) TechDelete(username, userid, guildid string) error {
	var count int
	sel := "SELECT count(*) as count FROM hs_compendium.tech WHERE guildid = $1 AND userid = $2 AND username = $3"
	err := d.db.QueryRow(context.Background(), sel, guildid, userid, username).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		del := "delete from hs_compendium.tech where username = $1 and userid = $2 and guildid = $3"
		_, err = d.db.Exec(context.Background(), del, username, userid, guildid)
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *Db) Unsubscribe(ctx context.Context, name, lvlkz string, TgChannel string, tipPing int) {
	del := "delete from kzbot.subscribe where name = $1 AND lvlkz = $2 AND chatid = $3 AND tip = $4"
	_, err := d.db.Exec(ctx, del, name, lvlkz, TgChannel, tipPing)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
