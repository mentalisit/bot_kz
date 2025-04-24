package postgres

import (
	"compendium/models"
	"errors"
	"github.com/jackc/pgx/v5"
	"sort"
)

func (d *Db) WsKillInsert(w models.WsKill) error {
	ctx, cancel := d.getContext()
	defer cancel()
	insert := `INSERT INTO hs_compendium.wskill(guildid, chatid, username, mention, shipname, timestampend,language) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := d.db.Exec(ctx, insert, w.GuildId, w.ChatId, w.UserName, w.Mention, w.ShipName, w.TimestampEnd, w.Language)
	return err
}
func (d *Db) WsKillDelete(w models.WsKill) error {
	ctx, cancel := d.getContext()
	defer cancel()
	deleteRole := `DELETE FROM hs_compendium.wskill WHERE guildid = $1 AND username = $2 AND shipname = $3`
	_, err := d.db.Exec(ctx, deleteRole, w.GuildId, w.UserName, w.ShipName)
	return err
}
func (d *Db) WsKillReadByGuildId(guildid string) ([]models.WsKill, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	selectWsKill := `SELECT * FROM hs_compendium.wskill WHERE guildid = $1`
	var wskill []models.WsKill
	rows, err := d.db.Query(ctx, selectWsKill, guildid)
	defer rows.Close()
	if err != nil {
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
	sort.Slice(wskill, func(i, j int) bool {
		return wskill[i].TimestampEnd > wskill[j].TimestampEnd
	})
	return wskill, nil
}

func (d *Db) WsKillReadAll() ([]models.WsKill, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	selectWsKill := `SELECT * FROM hs_compendium.wskill`
	var wskill []models.WsKill
	rows, err := d.db.Query(ctx, selectWsKill)
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
