package postgres

import (
	"compendium/models"
	"encoding/json"
)

func (d *Db) CodeInsert(u models.Code) error {
	ctx, cancel := d.getContext()
	defer cancel()
	insert := `INSERT INTO hs_compendium.codes(code,timestamp,identity) VALUES ($1,$2,$3)`
	bytes, err := json.Marshal(u.Identity)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}

	_, err = d.db.Exec(ctx, insert, u.Code, u.Timestamp, bytes)
	if err != nil {
		return err
	}
	return nil
}

func (d *Db) CodeGet(code string) (*models.Code, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var u models.Code
	var id int
	var bytes []byte
	selectCode := "SELECT * FROM hs_compendium.codes WHERE code = $1"
	err := d.db.QueryRow(ctx, selectCode, code).Scan(&id, &u.Code, &u.Timestamp, &bytes)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &u.Identity)
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}
	return &u, nil
}
