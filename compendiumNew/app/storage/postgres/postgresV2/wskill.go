package postgresv2

import (
	"compendium/models"
	"errors"
	"sort"

	"github.com/jackc/pgx/v5"
)

func (d *Db) WsKillInsert(w models.WsKill) error {
	insert := `INSERT INTO my_compendium.wskill(guildid, chatid, username, mention, shipname, timestampend,language) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := d.db.Exec(insert, w.GuildId, w.ChatId, w.UserName, w.Mention, w.ShipName, w.TimestampEnd, w.Language)
	return err
}
func (d *Db) WsKillDelete(w models.WsKill) error {
	deleteRole := `DELETE FROM my_compendium.wskill WHERE guildid = $1 AND username = $2 AND shipname = $3`
	_, err := d.db.Exec(deleteRole, w.GuildId, w.UserName, w.ShipName)
	return err
}
func (d *Db) WsKillReadByGuildId(guildid string) ([]models.WsKill, error) {
	selectWsKill := `SELECT * FROM my_compendium.wskill WHERE guildid = $1`
	var wskill []models.WsKill
	rows, err := d.db.Query(selectWsKill, guildid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var w models.WsKill
		var id int
		err = rows.Scan(&id, &w.GuildId, &w.ChatId, &w.UserName, &w.Mention, &w.ShipName, &w.TimestampEnd, &w.Language)
		if err != nil {
			break
		}
		wskill = append(wskill, w)
	}
	sort.Slice(wskill, func(i, j int) bool {
		return wskill[i].TimestampEnd > wskill[j].TimestampEnd
	})
	return wskill, nil
}

func (d *Db) WsKillReadAll() ([]models.WsKill, error) {
	selectWsKill := `SELECT * FROM my_compendium.wskill`
	var wskill []models.WsKill
	rows, err := d.db.Query(selectWsKill)
	defer rows.Close()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	for rows.Next() {
		var w models.WsKill
		var id int
		err = rows.Scan(&id, &w.GuildId, &w.ChatId, &w.UserName, &w.Mention, &w.ShipName, &w.TimestampEnd, &w.Language)
		if err != nil {
			break
		}
		wskill = append(wskill, w)
	}
	return wskill, nil
}
