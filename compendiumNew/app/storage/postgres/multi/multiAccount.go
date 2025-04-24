package multi

import (
	"compendium/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const returningMultiAccount = `
		RETURNING uuid, nickname,
		          telegram_id, telegram_username,
		          discord_id, discord_username,
		          whatsapp_id, whatsapp_username,
		          created_at,
				  avatarUrl, alts`

func scanMultiAccount(row pgx.Row) (*models.MultiAccount, error) {
	var acc models.MultiAccount

	var telegramID, discordID, whatsappID sql.NullString

	err := row.Scan(
		&acc.UUID, &acc.Nickname,
		&telegramID, &acc.TelegramUsername,
		&discordID, &acc.DiscordUsername,
		&whatsappID, &acc.WhatsappUsername,
		&acc.CreatedAt,
		&acc.AvatarURL, &acc.Alts,
	)
	if err != nil {
		return nil, err
	}

	if telegramID.Valid {
		acc.TelegramID = telegramID.String
	}
	if discordID.Valid {
		acc.DiscordID = discordID.String
	}
	if whatsappID.Valid {
		acc.WhatsappID = whatsappID.String
	}

	return &acc, nil
}

func (d *Db) CreateMultiAccountWithPlatform(id, nickname, platform, username string) (*models.MultiAccount, error) {
	ctx, cancel := d.getContext()
	defer cancel()

	var query string

	switch platform {
	case "tg":
		query = `
			INSERT INTO compendium.multi_accounts (nickname, telegram_id, telegram_username)
			VALUES ($1, $2, $3)` + returningMultiAccount
	case "ds":
		query = `
			INSERT INTO compendium.multi_accounts (nickname, discord_id, discord_username)
			VALUES ($1, $2, $3)` + returningMultiAccount
	case "wa":
		query = `
			INSERT INTO compendium.multi_accounts (nickname, whatsapp_id, whatsapp_username)
			VALUES ($1, $2, $3)` + returningMultiAccount
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}

	row := d.db.QueryRow(ctx, query, nickname, id, username)

	acc, err := scanMultiAccount(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		d.log.ErrorErr(err)
		return nil, err
	}
	return acc, nil
}

func (d *Db) FindMultiAccountByUserId(userId string) (*models.MultiAccount, error) {
	ctx, cancel := d.getContext()
	defer cancel()

	var selectQuery = `
		SELECT uuid, nickname,
		       telegram_id, telegram_username,
		       discord_id, discord_username,
		       whatsapp_id, whatsapp_username,
		       created_at,
			   avatarUrl, alts
		FROM compendium.multi_accounts
		WHERE telegram_id = $1 or discord_id = $1 or whatsapp_id = $1`

	row := d.db.QueryRow(ctx, selectQuery, userId)

	acc, err := scanMultiAccount(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		d.log.ErrorErr(err)
		return nil, err
	}
	return acc, nil
}
func (d *Db) FindMultiAccountByUsername(userName string) (*models.MultiAccount, error) {
	ctx, cancel := d.getContext()
	defer cancel()

	var selectQuery = `
		SELECT uuid, nickname,
		       telegram_id, telegram_username,
		       discord_id, discord_username,
		       whatsapp_id, whatsapp_username,
		       created_at,
			   avatarUrl, alts
		FROM compendium.multi_accounts
		WHERE telegram_username = $1 or discord_username = $1 or whatsapp_username = $1 OR $1 = ANY(alts);`

	row := d.db.QueryRow(ctx, selectQuery, userName)

	acc, err := scanMultiAccount(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		d.log.ErrorErr(err)
		return nil, err
	}
	return acc, nil
}

func (d *Db) FindMultiAccountUUID(uid uuid.UUID) (*models.MultiAccount, error) {
	ctx, cancel := d.getContext()
	defer cancel()

	var selectQuery = `
		SELECT uuid, nickname,
		       telegram_id, telegram_username,
		       discord_id, discord_username,
		       whatsapp_id, whatsapp_username,
		       created_at,
			   avatarUrl, alts
		FROM compendium.multi_accounts
		WHERE uuid = $1`

	row := d.db.QueryRow(ctx, selectQuery, uid)

	acc, err := scanMultiAccount(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // не найдено — не ошибка
		}
		return nil, err
	}

	return acc, nil
}

func (d *Db) UpdateMultiAccountInfo(uid uuid.UUID, platform, id, username string) (*models.MultiAccount, error) {
	ctx, cancel := d.getContext()
	defer cancel()

	var query string

	switch platform {
	case "tg":
		query = `
			UPDATE compendium.multi_accounts
			SET telegram_id = $1, telegram_username = $2
			WHERE uuid = $3` + returningMultiAccount
	case "ds":
		query = `
			UPDATE compendium.multi_accounts
			SET discord_id = $1, discord_username = $2
			WHERE uuid = $3` + returningMultiAccount
	case "wa":
		query = `
			UPDATE compendium.multi_accounts
			SET whatsapp_id = $1, whatsapp_username = $2
			WHERE uuid = $3` + returningMultiAccount
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}

	row := d.db.QueryRow(ctx, query, id, username, uid)

	acc, err := scanMultiAccount(row)

	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}

	return acc, nil
}

func (d *Db) UpdateMultiAccountNickname(m models.MultiAccount) (*models.MultiAccount, error) {
	ctx, cancel := d.getContext()
	defer cancel()

	const query = `
		UPDATE compendium.multi_accounts
		SET nickname = $1
		WHERE uuid = $2` + returningMultiAccount

	row := d.db.QueryRow(ctx, query, m.Nickname, m.UUID)

	acc, err := scanMultiAccount(row)

	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (d *Db) UpdateMultiAccountAlts(m models.MultiAccount) (*models.MultiAccount, error) {
	ctx, cancel := d.getContext()
	defer cancel()

	const query = `
		UPDATE compendium.multi_accounts
		SET alts = $1
		WHERE uuid = $2` + returningMultiAccount

	row := d.db.QueryRow(ctx, query, m.Alts, m.UUID)

	acc, err := scanMultiAccount(row)

	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (d *Db) UpdateMultiAccountAvatarUrl(m models.MultiAccount) (*models.MultiAccount, error) {
	ctx, cancel := d.getContext()
	defer cancel()

	const query = `
		UPDATE compendium.multi_accounts
		SET avatarUrl = $1
		WHERE uuid = $2` + returningMultiAccount

	row := d.db.QueryRow(ctx, query, m.AvatarURL, m.UUID)

	acc, err := scanMultiAccount(row)

	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (d *Db) RemoveUserIdAllTable(ma *models.MultiAccount) {
	deleteUser := func(userid string) {
		ctx, cancel := d.getContext()
		defer cancel()
		deleteuser := `DELETE FROM hs_compendium.corpmember WHERE userid = $1 `
		_, _ = d.db.Exec(ctx, deleteuser, userid)
		//deleteuser = `DELETE FROM hs_compendium.list_users WHERE userid = $1 `
		//_, _ = d.db.Exec(ctx, deleteuser, userid)
		deleteuser = `DELETE FROM hs_compendium.tech WHERE userid = $1 `
		_, _ = d.db.Exec(ctx, deleteuser, userid)
		deleteuser = `DELETE FROM hs_compendium.users WHERE userid = $1 `
		_, _ = d.db.Exec(ctx, deleteuser, userid)
	}
	if ma != nil {
		if ma.DiscordID != "" {
			deleteUser(ma.DiscordID)
		}
		if ma.TelegramID != "" {
			deleteUser(ma.TelegramID)
		}
	}

}
