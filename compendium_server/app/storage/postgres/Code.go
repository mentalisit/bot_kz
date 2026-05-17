package postgres

import (
	"compendium_s/models"
	"encoding/json"
)

func (d *Db) CodeGet(code string) (*models.Code, error) {
	var u models.Code
	var id int
	var bytes []byte
	selectCode := "SELECT * FROM hs_compendium.codes WHERE code = $1"
	err := d.db.QueryRow(selectCode, code).Scan(&id, &u.Code, &bytes, &u.Timestamp)
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
func (d *Db) CodeAllGet() []models.Code {
	selectCodes := "SELECT * FROM hs_compendium.codes"

	results, err := d.db.Query(selectCodes)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	var uu []models.Code

	for results.Next() {
		var u models.Code
		var id int
		var bytes []byte

		_ = results.Scan(&id, &u.Code, &bytes, &u.Timestamp)

		err = json.Unmarshal(bytes, &u.Identity)
		if err != nil {
			d.log.ErrorErr(err)
		}
		uu = append(uu, u)
	}
	return uu
}
func (d *Db) CodeDelete(code string) {
	del := "delete from hs_compendium.codes where code = $1"
	_, err := d.db.Exec(del, code)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
