package postgresV2

import (
	"encoding/json"
	"errors"
	"rs/models"

	"github.com/jackc/pgx/v5"
)

// CreateLocalGuild создает новую запись о локальной гильдии
func (d *Db) CreateLocalGuild(l models.LocalGuild) (int, error) {

	marshalData, err := json.Marshal(l.Data)
	if err != nil {
		d.log.ErrorErr(err)
		return 0, err
	}

	var gid int
	query := `INSERT INTO rs_bot2.local_guilds (name, data) VALUES ($1, $2) RETURNING gid`
	err = d.db.QueryRow(query, l.Name, marshalData).Scan(&gid)
	if err != nil {
		d.log.ErrorErr(err)
		return 0, err
	}

	return gid, nil
}

// ReadLocalGuildByName читает запись о локальной гильдии по имени
func (d *Db) ReadLocalGuildByName(name string) (*models.LocalGuild, error) {

	query := `SELECT gid, name, data FROM rs_bot2.local_guilds WHERE name = $1`
	var l models.LocalGuild
	var dataJSON []byte

	err := d.db.QueryRow(query, name).Scan(&l.Gid, &l.Name, &dataJSON)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		d.log.ErrorErr(err)
		return nil, err
	}

	err = json.Unmarshal(dataJSON, &l.Data)
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}

	return &l, nil
}

// UpdateLocalGuildData обновляет поле data для локальной гильдии по gid
func (d *Db) UpdateLocalGuildData(gid int, data models.LocalGuildData) error {

	marshalData, err := json.Marshal(data)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}

	query := `UPDATE rs_bot2.local_guilds SET data = $1 WHERE gid = $2`
	_, err = d.db.Exec(query, marshalData, gid)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}

	return nil
}
