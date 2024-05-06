package postgres

import (
	"compendium/models"
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v4"
)

func (d *Db) TechInsert(username, userid, guildid string, tech []byte) error {
	var count int
	sel := "SELECT count(*) as count FROM hs_compendium.tech WHERE guildid = $1 AND userid = $2 AND username = $3"
	err := d.db.QueryRow(context.Background(), sel, guildid, userid, username).Scan(&count)
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
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}
	return tech, nil
}
func (d *Db) TechGetAll(cm models.CorpMember) ([]models.CorpMember, error) {
	var acm []models.CorpMember
	sel := "SELECT username,tech FROM hs_compendium.tech WHERE userid = $1 AND guildid = $2"
	q, err := d.db.Query(context.Background(), sel, cm.UserId, cm.GuildId)
	if err != nil {
		return acm, err
	}
	for q.Next() {
		var ncm models.CorpMember
		ncm = cm
		var tech []byte
		err = q.Scan(&ncm.Name, &tech)
		if err != nil {
			return nil, err
		}
		var techl models.TechLevels
		err = json.Unmarshal(tech, &techl)
		if err != nil {
			return nil, err
		}
		if len(techl) > 0 {
			m := make(map[int][2]int)
			for i, level := range techl {
				m[i] = [2]int{level.Level}
			}
			ncm.Tech = m
		}
		acm = append(acm, ncm)
	}
	return acm, nil
}
func (d *Db) TechUpdate(username, userid, guildid string, tech []byte) error {
	upd := `update hs_compendium.tech set tech = $1 where username = $2 and userid = $3 and guildid = $4`
	updresult, err := d.db.Exec(context.Background(), upd, tech, username, userid, guildid)
	if err != nil {
		return err
	}
	if updresult.RowsAffected() == 0 {
		err = d.TechInsert(username, userid, guildid, tech)
		if err != nil {
			d.log.ErrorErr(err)
			return err
		}
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
func (d *Db) TechGetCount(userid, guildid string) (int, error) {
	var count int
	sel := "SELECT count(*) as count FROM hs_compendium.tech WHERE guildid = $1 AND userid = $2"
	err := d.db.QueryRow(context.Background(), sel, guildid, userid).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
