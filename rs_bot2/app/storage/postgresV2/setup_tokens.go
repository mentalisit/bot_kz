package postgresV2

import (
	"database/sql"
	"encoding/json"
	"errors"
	"rs/models"

	"github.com/google/uuid"
)

func (d *Db) InsertOther(p models.Other) uuid.UUID {

	query := `
        INSERT INTO rs_bot.other(uuid, data_type, data, read) 
        VALUES ($1, $2, $3, $4) 
        ON CONFLICT (uuid) DO UPDATE 
        SET data_type = EXCLUDED.data_type, 
            data = EXCLUDED.data,
            read = EXCLUDED.read
        RETURNING uuid
    `

	var uid uuid.UUID
	err := d.db.QueryRow(query, p.Uuid, p.DataType, p.JsonMarshalDataWeb(), false).Scan(&uid)
	if err != nil {
		d.log.ErrorErr(err)
		return uuid.Nil
	}

	return uid
}

func (d *Db) GetOtherByUUID(uid uuid.UUID) (*models.Other, error) {

	query := `SELECT uuid, data_type, data, read  FROM rs_bot.other WHERE uuid = $1`
	var other models.Other
	var dataJSON []byte

	err := d.db.QueryRow(query, uid).Scan(&other.Uuid, &other.DataType, &dataJSON, &other.Read)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		d.log.ErrorErr(err)
		return nil, err
	}

	if len(dataJSON) > 0 {
		err = json.Unmarshal(dataJSON, &other.Data)
		if err != nil {
			d.log.ErrorErr(err)
			return nil, err
		}
	}

	return &other, nil
}

func (d *Db) DeleteOtherByUUID(uid string) {

	del := "delete from rs_bot.other where uuid = $1"
	_, err := d.db.Exec(del, uid)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) UpdateReadOtherByUUID(uid string) {

	upd := `update rs_bot.other set read = $1 where uuid = $2`
	_, err := d.db.Exec(upd, true, uid)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
