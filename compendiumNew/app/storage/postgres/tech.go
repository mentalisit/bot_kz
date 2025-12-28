package postgres

import (
	"compendium/models"
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
		return nil, err
	}
	return tech, nil
}

func (d *Db) TechGetName(username string) ([]byte, string, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var tech []byte
	var userid string
	sel := "SELECT userid,tech FROM hs_compendium.tech WHERE username = $1"
	err := d.db.QueryRow(ctx, sel, username).Scan(&userid, &tech)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, "", err
		}
	}
	return tech, userid, nil
}

func (d *Db) TechGetAll(cm models.CorpMember) ([]models.CorpMember, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var acm []models.CorpMember
	acmMap := make(map[string]models.CorpMember)
	sel := "SELECT username,tech FROM hs_compendium.tech WHERE userid = $1"
	q, err := d.db.Query(ctx, sel, cm.UserId)
	defer q.Close()
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
			acmMap[ncm.Name] = compareCM(acmMap[ncm.Name], ncm)
		}
	}
	for _, member := range acmMap {
		acm = append(acm, member)
	}
	return acm, nil
}

func (d *Db) TechDelete(username, userid string) error {
	ctx, cancel := d.getContext()
	defer cancel()
	var count int
	sel := "SELECT count(*) as count FROM hs_compendium.tech WHERE userid = $1 AND username = $2"
	err := d.db.QueryRow(ctx, sel, userid, username).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		del := "delete from hs_compendium.tech where username = $1 and userid = $2"
		_, err = d.db.Exec(ctx, del, username, userid)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Db) TechGetCount(userid string) (int, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var count int
	sel := "SELECT count(*) as count FROM hs_compendium.tech WHERE userid = $1"
	err := d.db.QueryRow(ctx, sel, userid).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (d *Db) TechGetAllUserId(userid string) ([]models.TechTable, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var acm []models.TechTable
	sel := "SELECT * FROM hs_compendium.tech where userid = $1"
	q, err := d.db.Query(ctx, sel, userid)
	defer q.Close()
	if err != nil {
		return nil, err
	}

	for q.Next() {
		var ncm models.TechTable
		err = q.Scan(&ncm.Id, &ncm.Name, &ncm.NameId, &ncm.GuildId, &ncm.Tech)
		if err != nil {
			return nil, err
		}

		acm = append(acm, ncm)
	}
	if err = q.Err(); err != nil { // Проверка ошибок после завершения итерации
		return nil, err
	}
	return acm, nil
}

func compareCM(cm1, cm2 models.CorpMember) models.CorpMember {
	var cm models.CorpMember
	cm = cm1
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
