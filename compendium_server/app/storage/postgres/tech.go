package postgres

import (
	"compendium_s/models"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (d *Db) TechInsert(username, userid string, tech []byte) error {
	ctx, cancel := d.getContext()
	defer cancel()
	var count int
	sel := "SELECT count(*) as count FROM hs_compendium.tech WHERE userid = $1 AND username = $2"
	err := d.db.QueryRow(ctx, sel, userid, username).Scan(&count)
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
		_, err = d.db.Exec(ctx, insert, username, userid, "all", tech)
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *Db) TechGet(username, userid string) ([]byte, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var tech []byte
	sel := "SELECT tech FROM hs_compendium.tech WHERE userid = $1 AND username = $2"
	err := d.db.QueryRow(ctx, sel, userid, username).Scan(&tech)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}
	return tech, nil
}

func (d *Db) TechGetAll(cm models.CorpMember) ([]models.CorpMember, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var acm []models.CorpMember
	acmMap := make(map[string]models.CorpMember)
	sel := "SELECT username,tech FROM hs_compendium.tech WHERE userid = $1"
	q, err := d.db.Query(ctx, sel, cm.UserId)
	if err != nil {
		return acm, err
	}
	defer q.Close()
	for q.Next() {
		var ncm models.CorpMember
		ncm = cm
		var tech []byte
		err = q.Scan(&ncm.Name, &tech)
		if err != nil {
			return nil, err
		}

		if cm.Name != ncm.Name {
			ncm.UserId = ncm.UserId + "/" + ncm.Name
		}

		var techl models.TechLevels
		err = json.Unmarshal(tech, &techl)
		if err != nil {
			return nil, err
		}
		if len(techl) > 0 {
			m := make(models.TechLevels)
			for i, level := range techl {
				m[i] = level
			}
			ncm.Tech = m
		}

		if acmMap[ncm.Name].Name == "" {
			acmMap[ncm.Name] = ncm
		} else {
			acmMap[ncm.Name] = compare(acmMap[ncm.Name], ncm)
		}
	}
	for _, member := range acmMap {
		acm = append(acm, member)
	}
	return acm, nil
}

func (d *Db) TechUpdate(username, userid string, tech []byte) error {
	ctx, cancel := d.getContext()
	defer cancel()
	upd := `update hs_compendium.tech set tech = $1 where username = $2 and userid = $3`
	updresult, err := d.db.Exec(ctx, upd, tech, username, userid)
	if err != nil {
		return err
	}
	if updresult.RowsAffected() == 0 {
		err = d.TechInsert(username, userid, tech)
		if err != nil {
			d.log.ErrorErr(err)
			return err
		}
	}
	return nil
}

func compare(cm1, cm2 models.CorpMember) models.CorpMember {
	var cm models.CorpMember
	cm = cm1
	if cm2.ZoneOffset != 0 {
		cm.ZoneOffset = cm2.ZoneOffset
	}
	if cm2.TimeZone != "" {
		cm.TimeZone = cm2.TimeZone
	}
	if cm2.AfkFor != "" {
		cm.AfkFor = cm2.AfkFor
	}
	cm.Tech = make(models.TechLevels)
	for i, level := range cm1.Tech {
		if cm2.Tech[i].Ts == level.Ts {
			cm.Tech[i] = level
		} else if cm2.Tech[i].Ts > level.Ts {
			cm.Tech[i] = cm2.Tech[i]
		} else if cm2.Tech[i].Ts < level.Ts {
			cm.Tech[i] = level
		}
	}
	return cm
}
