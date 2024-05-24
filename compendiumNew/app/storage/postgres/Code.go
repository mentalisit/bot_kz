package postgres

import (
	"compendium/models"
	"context"
	"encoding/json"
)

func (d *Db) CodeInsert(u models.Code) error {
	insert := `INSERT INTO hs_compendium.codes(code,timestamp,identity) VALUES ($1,$2,$3)`
	bytes, err := json.Marshal(u.Identity)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}

	_, err = d.db.Exec(context.Background(), insert, u.Code, u.Timestamp, bytes)
	if err != nil {
		return err
	}
	return nil
}

func (d *Db) CodeGet(code string) (*models.Code, error) {
	var u models.Code
	var id int
	var bytes []byte
	selectCode := "SELECT * FROM hs_compendium.codes WHERE code = $1"
	err := d.db.QueryRow(context.Background(), selectCode, code).Scan(&id, &u.Code, &u.Timestamp, &bytes)
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
