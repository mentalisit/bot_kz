package postgresv2

import (
	"compendium/models"
	"database/sql"
	"errors"
)

func (d *Db) WsKillInsert(w models.WsKill) error {
	query := `INSERT INTO my_compendium.wskill (guildid, chatid, username, mention, shipname, timestampend, language) 
              VALUES (:guildid, :chatid, :username, :mention, :shipname, :timestampend, :language)`

	_, err := d.db.NamedExec(query, w)
	return err
}

func (d *Db) WsKillDelete(w models.WsKill) error {
	deleteRole := `DELETE FROM my_compendium.wskill WHERE guildid = $1 AND username = $2 AND shipname = $3`
	_, err := d.db.Exec(deleteRole, w.GuildId, w.UserName, w.ShipName)
	return err
}
func (d *Db) WsKillReadByGuildId(guildID string) ([]models.WsKill, error) {
	var skills []models.WsKill
	// Сортируем прямо в SQL — это быстрее и чище
	query := `SELECT id, guildid, chatid, username, mention, shipname, timestampend, language 
              FROM my_compendium.wskill 
              WHERE guildid = $1 
              ORDER BY timestampend DESC`

	err := d.db.Select(&skills, query, guildID)
	return skills, err
}

func (d *Db) WsKillReadAll() ([]models.WsKill, error) {
	var skills []models.WsKill
	query := `SELECT id, guildid, chatid, username, mention, shipname, timestampend, language 
              FROM my_compendium.wskill`

	err := d.db.Select(&skills, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return skills, nil
}
