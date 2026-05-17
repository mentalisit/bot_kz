package postgresV2

import (
	"rs/models"

	"github.com/google/uuid"
)

func (d *Db) UpdateMergedAccounts(uid uuid.UUID, accounts []models.GameAccountData) error {

	query := `
        UPDATE my_compendium.multi_accounts 
        SET data = jsonb_set(COALESCE(data, '{}'), '{Merged}', $2::jsonb) 
        WHERE uuid = $1`

	_, err := d.db.Exec(query, uid, accounts)
	return err
}

func (d *Db) UpdateWsPhone(uid uuid.UUID, wsPhone string) error {

	query := `
        UPDATE my_compendium.multi_accounts 
        SET data = jsonb_set(COALESCE(data, '{}'), '{wsPhone}', to_jsonb($2::text)) 
        WHERE uuid = $1`

	_, err := d.db.Exec(query, uid, wsPhone)
	return err
}

func (d *Db) UpdateNotifyPM(uid uuid.UUID, notifyPM bool) error {

	query := `
        UPDATE my_compendium.multi_accounts 
        SET data = jsonb_set(COALESCE(data, '{}'), '{notifyPM}', to_jsonb($2::boolean)) 
        WHERE uuid = $1`

	_, err := d.db.Exec(query, uid, notifyPM)
	return err
}
