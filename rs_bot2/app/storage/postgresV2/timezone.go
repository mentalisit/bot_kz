package postgresV2

import "github.com/google/uuid"

// UpdateCorpMemberTimezone обновляет IANA timezone и offset в таблице corpmember.
// Если записи нет — ничего не делает (не создаёт).
func (d *Db) UpdateCorpMemberTimezone(uid uuid.UUID, timezone string, zoneOffset int) error {

	query := `UPDATE my_compendium.corpmember SET timezona = $1, zonaoffset = $2 WHERE uid = $3`

	_, err := d.db.Exec(query, timezone, zoneOffset, uid)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}
	return nil
}

// UpdateMultiAccountTimezone сохраняет IANA timezone в JSONB поле data таблицы multi_accounts.
func (d *Db) UpdateMultiAccountTimezone(uid uuid.UUID, timezone string) error {

	query := `
        UPDATE my_compendium.multi_accounts 
        SET data = jsonb_set(COALESCE(data, '{}'), '{timezone}', to_jsonb($2::text)) 
        WHERE uuid = $1`

	_, err := d.db.Exec(query, uid, timezone)
	return err
}
