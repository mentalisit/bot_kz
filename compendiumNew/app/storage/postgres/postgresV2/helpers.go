package postgresv2

import (
	"fmt"

	"github.com/google/uuid"
)

func (d *Db) UpdateNickname(uid uuid.UUID, oldNick, nickName string) error {
	// Update nickname in multi_accounts
	const query = `
		UPDATE my_compendium.multi_accounts
		SET nickname = $1
		WHERE uuid = $2`

	_, err := d.db.Exec(query, nickName, uid)
	if err != nil {
		d.log.ErrorErr(fmt.Errorf("failed to update nickname in multi_accounts: %w", err))
		return err
	}

	// Update username in technologies
	const query2 = `
		UPDATE my_compendium.technologies
		SET username = $1
		WHERE uid = $2 AND username = $3`

	_, err = d.db.Exec(query2, nickName, uid, oldNick)
	if err != nil {
		d.log.ErrorErr(fmt.Errorf("failed to update username in technologies: %w", err))
		return err
	}

	return nil
}
