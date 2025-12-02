package postgresv2

import (
	"compendium/models"
	"encoding/json"
)

func (d *Db) CodeInsert(u models.CodeV2) error {
	insert := `INSERT INTO hs_compendium.codes(code,timestamp,identity) VALUES ($1,$2,$3)`
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
