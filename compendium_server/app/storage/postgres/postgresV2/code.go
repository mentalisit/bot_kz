package postgresv2

import (
	"compendium_s/models"
	"encoding/json"
)

func (d *Db) CodeInsert(u models.Code) error {
	insert := `INSERT INTO my_compendium.codes(code,time,identity) VALUES ($1,$2,$3)`
	bytes, err := json.Marshal(u.Identity)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}

	_, err = d.db.Exec(insert, u.Code, u.Timestamp, bytes)
	if err != nil {
		return err
	}
	return nil
}

func (d *Db) CodeGet(code string) (*models.Code, error) {
	var u models.Code
	var id int
	var bytes []byte
	selectCode := "SELECT * FROM my_compendium.codes WHERE code = $1"
	err := d.db.QueryRow(selectCode, code).Scan(&id, &u.Code, &u.Timestamp, &bytes)
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
func (d *Db) CodeDelete(code string) {
	del := "delete from my_compendium.codes where code = $1"
	_, err := d.db.Exec(del, code)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
