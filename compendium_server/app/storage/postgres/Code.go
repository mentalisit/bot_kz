package postgres

import (
	"compendium_s/models"
	"encoding/json"
)

//func (d *Db) CodeInsert(u models.Code) error {
//	insert := `INSERT INTO hs_compendium.codes(code,timestamp,identity) VALUES ($1,$2,$3)`
//	bytes, err := json.Marshal(u.Identity)
//	if err != nil {
//		d.log.ErrorErr(err)
//		return err
//	}
//
//	_, err = d.db.Exec(context.Background(), insert, u.Code, u.Timestamp, bytes)
//	if err != nil {
//		return err
//	}
//	return nil
//}

func (d *Db) CodeGet(code string) (*models.Code, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var u models.Code
	var id int
	var bytes []byte
	selectCode := "SELECT * FROM hs_compendium.codes WHERE code = $1"
	err := d.db.QueryRow(ctx, selectCode, code).Scan(&id, &u.Code, &bytes, &u.Timestamp)
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
	ctx, cancel := d.getContext()
	defer cancel()
	selectCodes := "SELECT * FROM hs_compendium.codes"

	results, err := d.db.Query(ctx, selectCodes)
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
	ctx, cancel := d.getContext()
	defer cancel()
	del := "delete from hs_compendium.codes where code = $1"
	_, err := d.db.Exec(ctx, del, code)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
