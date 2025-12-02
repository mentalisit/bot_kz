package postgres

import (
	"database/sql"
	"errors"
	"telegram/models"

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

func (d *Db) CreateMultiAccount(acc models.MultiAccount) (*models.MultiAccount, error) {
	ctx, cancel := d.getContext()
	defer cancel()

	query := `INSERT INTO my_compendium.multi_accounts (nickname, 
            telegram_id, telegram_username,
            discord_id, discord_username)
			VALUES ($1, $2, $3, $4 ,$5)` + returningMultiAccount

	row := d.db.QueryRow(ctx, query, acc.Nickname,
		acc.TelegramID, acc.TelegramUsername,
		acc.DiscordID, acc.DiscordUsername)

	accNew, err := scanMultiAccount(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		d.log.ErrorErr(err)
		return nil, err
	}
	return accNew, nil
}

func (d *Db) FindMultiAccountByTelegramID(telegramID string) (*models.MultiAccount, error) {
	ctx, cancel := d.getContext()
	defer cancel()

	var selectQuery = `
		SELECT uuid, nickname,
		       telegram_id, telegram_username,
		       discord_id, discord_username,
		       whatsapp_id, whatsapp_username,
		       created_at,
			   avatarUrl, alts
		FROM my_compendium.multi_accounts
		WHERE telegram_id = $1`

	row := d.db.QueryRow(ctx, selectQuery, telegramID)

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

func (d *Db) UpdateMultiAccount(acc models.MultiAccount) (*models.MultiAccount, error) {
	ctx, cancel := d.getContext()
	defer cancel()

	query := `UPDATE my_compendium.multi_accounts
			SET discord_id = $1, discord_username = $2, nickname =$3
			WHERE uuid = $4` + returningMultiAccount

	row := d.db.QueryRow(ctx, query, acc.DiscordID, acc.DiscordUsername, acc.Nickname, acc.UUID)

	accNew, err := scanMultiAccount(row)

	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}

	return accNew, nil
}
