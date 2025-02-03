package postgres

import (
	"compendium/models"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
)

func (d *Db) TechInsert(username, userid, guildid string, tech []byte) error {
	ctx, cancel := d.GetContext()
	defer cancel()
	var count int
	sel := "SELECT count(*) as count FROM hs_compendium.tech WHERE guildid = $1 AND userid = $2 AND username = $3"
	err := d.db.QueryRow(ctx, sel, guildid, userid, username).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		if len(tech) <= 2 || tech == nil {
			techEmpty := make(map[int]models.TechLevel)
			techEmpty[701] = models.TechLevel{
				Ts:    0,
				Level: 0,
			}
			tech, _ = json.Marshal(techEmpty)
		}
		insert := `INSERT INTO hs_compendium.tech(username, userid, guildid, tech) VALUES ($1,$2,$3,$4)`
		_, err = d.db.Exec(ctx, insert, username, userid, guildid, tech)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Db) TechGet(username, userid, guildid string) ([]byte, error) {
	ctx, cancel := d.GetContext()
	defer cancel()
	var tech []byte
	sel := "SELECT tech FROM hs_compendium.tech WHERE userid = $1 AND guildid = $2 AND username = $3"
	err := d.db.QueryRow(ctx, sel, userid, guildid, username).Scan(&tech)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}
	return tech, nil
}
func (d *Db) TechGetName(username, guildid string) ([]byte, string, error) {
	ctx, cancel := d.GetContext()
	defer cancel()
	var tech []byte
	var userid string
	sel := "SELECT userid,tech FROM hs_compendium.tech WHERE guildid = $1 AND username = $2"
	err := d.db.QueryRow(ctx, sel, guildid, username).Scan(&userid, &tech)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, "", err
		}
	}
	return tech, userid, nil
}

func (d *Db) TechGetAll(cm models.CorpMember) ([]models.CorpMember, error) {
	ctx, cancel := d.GetContext()
	defer cancel()
	var acm []models.CorpMember
	sel := "SELECT username,tech FROM hs_compendium.tech WHERE userid = $1 AND guildid = $2"
	q, err := d.db.Query(ctx, sel, cm.UserId, cm.GuildId)
	defer q.Close()
	if err != nil {
		return nil, err
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
			ncm.Tech = make(map[int][2]int)
			for i, level := range techl {
				ncm.Tech[i] = [2]int{level.Level, int(level.Ts)}
			}
		}
		acm = append(acm, ncm)
	}
	if err = q.Err(); err != nil { // Проверка ошибок после завершения итерации
		return nil, err
	}
	return acm, nil
}

func (d *Db) TechUpdate(username, userid, guildid string, tech []byte) error {
	ctx, cancel := d.GetContext()
	defer cancel()
	upd := `update hs_compendium.tech set tech = $1 where username = $2 and userid = $3 and guildid = $4`
	updresult, err := d.db.Exec(ctx, upd, tech, username, userid, guildid)
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
	ctx, cancel := d.GetContext()
	defer cancel()
	var count int
	sel := "SELECT count(*) as count FROM hs_compendium.tech WHERE guildid = $1 AND userid = $2 AND username = $3"
	err := d.db.QueryRow(ctx, sel, guildid, userid, username).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		del := "delete from hs_compendium.tech where username = $1 and userid = $2 and guildid = $3"
		_, err = d.db.Exec(ctx, del, username, userid, guildid)
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *Db) Unsubscribe(name, lvlkz string, TgChannel string, tipPing int) {
	ctx, cancel := d.GetContext()
	defer cancel()
	del := "delete from kzbot.subscribe where name = $1 AND lvlkz = $2 AND chatid = $3 AND tip = $4"
	_, err := d.db.Exec(ctx, del, name, lvlkz, TgChannel, tipPing)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) TechGetCount(userid, guildid string) (int, error) {
	ctx, cancel := d.GetContext()
	defer cancel()
	var count int
	sel := "SELECT count(*) as count FROM hs_compendium.tech WHERE guildid = $1 AND userid = $2"
	err := d.db.QueryRow(ctx, sel, guildid, userid).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
