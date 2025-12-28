package multi

import (
	"compendium_s/models"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (d *Db) TechnologiesInsert(uid uuid.UUID, username string, tech []byte) error {
	ctx, cancel := d.getContext()
	defer cancel()
	var count int
	sel := "SELECT count(*) as count FROM compendium.technologies WHERE uid = $1 AND username = $2"
	err := d.db.QueryRow(ctx, sel, uid, username).Scan(&count)
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
		insert := `INSERT INTO compendium.technologies(uid, username, tech) VALUES ($1,$2,$3)`
		_, err = d.db.Exec(ctx, insert, uid, username, tech)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Db) TechnologiesGet(uid uuid.UUID, username string) ([]byte, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var tech []byte
	sel := "SELECT tech FROM compendium.technologies WHERE uid = $1 AND username = $2"
	err := d.db.QueryRow(ctx, sel, uid, username).Scan(&tech)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}
	return tech, nil
}

func (d *Db) TechnologiesGetName(username string) ([]byte, uuid.UUID, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var tech []byte
	var uid uuid.UUID
	sel := "SELECT uid,tech FROM compendium.technologies WHERE username = $1"
	err := d.db.QueryRow(ctx, sel, username).Scan(&uid, &tech)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, uuid.Nil, err
		}
	}
	return tech, uid, nil
}

func (d *Db) TechnologiesGetAllCorpMember(cm models.CorpMember, UId uuid.UUID) ([]models.CorpMember, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var acm []models.CorpMember
	sel := "SELECT username,tech FROM compendium.technologies WHERE uid = $1"
	q, err := d.db.Query(ctx, sel, UId)
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
			return acm, err
		}

		var techl models.TechLevels
		err = json.Unmarshal(tech, &techl)
		if err != nil {
			return acm, err
		}
		if len(techl) > 0 {
			//ncm.Tech = make(map[int][2]int)
			//for i, level := range techl {
			//	ncm.Tech[i] = [2]int{level.Level, int(level.Ts)}
			//}
			ncm.Tech = make(models.TechLevels)
			for i, level := range techl {
				ncm.Tech[i] = level
			}
			ncm.UserId = ncm.UserId + "/" + ncm.Name
		}
		acm = append(acm, ncm)
	}
	if err = q.Err(); err != nil { // Проверка ошибок после завершения итерации
		return nil, err
	}
	return acm, nil
}

func (d *Db) TechnologiesGetMember(UId uuid.UUID) []models.Technology {
	ctx, cancel := d.getContext()
	defer cancel()
	var tt []models.Technology
	sel := "SELECT username,tech FROM compendium.technologies WHERE uid = $1"
	q, err := d.db.Query(ctx, sel, UId)
	defer q.Close()
	if err != nil {
		d.log.ErrorErr(err)
		return nil
	}

	for q.Next() {
		var t models.Technology
		var tech []byte
		err = q.Scan(&t.Name, &tech)
		if err != nil {
			d.log.ErrorErr(err)
			return nil
		}

		err = json.Unmarshal(tech, &t.Tech)
		if err != nil {
			d.log.ErrorErr(err)
			return nil
		}
		tt = append(tt, t)
	}
	return tt
}

func (d *Db) TechnologiesUpdate(uid uuid.UUID, username string, tech []byte) error {
	ctx, cancel := d.getContext()
	defer cancel()
	upd := `update compendium.technologies set tech = $1 where username = $2 and uid = $3`
	updresult, err := d.db.Exec(ctx, upd, tech, username, uid)
	if err != nil {
		return err
	}
	if updresult.RowsAffected() == 0 {
		err = d.TechnologiesInsert(uid, username, tech)
		if err != nil {
			d.log.ErrorErr(err)
			return err
		}
	}
	return nil
}
func (d *Db) TechnologiesUpdateUsername(uid uuid.UUID, oldUsername string, NewUsername string) error {
	ctx, cancel := d.getContext()
	defer cancel()
	upd := `update compendium.technologies set username = $1 where username = $2 and uid = $3`
	_, err := d.db.Exec(ctx, upd, NewUsername, oldUsername, uid)
	if err != nil {
		return err
	}
	return nil
}

func (d *Db) TechnologiesDelete(uid uuid.UUID, username string) error {
	ctx, cancel := d.getContext()
	defer cancel()
	var count int
	sel := "SELECT count(*) as count FROM compendium.technologies WHERE uid = $1 AND username = $2"
	err := d.db.QueryRow(ctx, sel, uid, username).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		del := "delete from compendium.technologies where username = $1 and uid = $2"
		_, err = d.db.Exec(ctx, del, username, uid)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Db) TechnologiesGetCount(uid uuid.UUID) (int, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var count int
	sel := "SELECT count(*) as count FROM compendium.technologies WHERE uid = $1"
	err := d.db.QueryRow(ctx, sel, uid).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
