package postgresv2

import (
	"compendium/models"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const returningMultiAccount = `
		RETURNING uuid, nickname,
		          telegram_id, telegram_username,
		          discord_id, discord_username,
		          whatsapp_id, whatsapp_username,
		          created_at,
				  avatarUrl, alts`

func scanMultiAccount(row *sql.Row) (*models.MultiAccount, error) {
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

// Multi-account methods
func (d *Db) FindMultiAccountUUID(uid uuid.UUID) (*models.MultiAccount, error) {
	query := `SELECT uuid, nickname, telegram_id, telegram_username, discord_id, discord_username,
			  whatsapp_id, whatsapp_username, created_at, avatarUrl, alts
			  FROM my_compendium.multi_accounts WHERE uuid = $1`

	row := d.db.QueryRow(query, uid)
	acc, err := scanMultiAccount(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		d.log.ErrorErr(err)
		return nil, err
	}
	return acc, nil
}

func (d *Db) FindMultiAccountByUserId(userid string) (*models.MultiAccount, error) {
	query := `SELECT uuid, nickname, telegram_id, telegram_username, discord_id, discord_username,
			  whatsapp_id, whatsapp_username, created_at, avatarUrl, alts
			  FROM my_compendium.multi_accounts WHERE discord_id = $1 OR telegram_id = $1 OR whatsapp_id = $1`

	row := d.db.QueryRow(query, userid)
	acc, err := scanMultiAccount(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		d.log.ErrorErr(err)
		return nil, err
	}
	return acc, nil
}

func (d *Db) CreateMultiAccountWithPlatform(id, nickname, platform, username string) (*models.MultiAccount, error) {
	var query string
	switch platform {
	case "tg":
		query = `
			INSERT INTO my_compendium.multi_accounts (nickname, telegram_id, telegram_username)
			VALUES ($1, $2, $3)` + returningMultiAccount
	case "ds":
		query = `
			INSERT INTO my_compendium.multi_accounts (nickname, discord_id, discord_username)
			VALUES ($1, $2, $3)` + returningMultiAccount
	case "wa":
		query = `
			INSERT INTO my_compendium.multi_accounts (nickname, whatsapp_id, whatsapp_username)
			VALUES ($1, $2, $3)` + returningMultiAccount
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
	row := d.db.QueryRow(query, nickname, id, username)
	acc, err := scanMultiAccount(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		d.log.ErrorErr(err)
		return nil, err
	}
	return acc, nil
}

func (d *Db) UpdateMultiAccountInfo(uid uuid.UUID, platform, id, username string) (*models.MultiAccount, error) {
	var query string

	switch platform {
	case "tg":
		query = `
			UPDATE my_compendium.multi_accounts
			SET telegram_id = $1, telegram_username = $2
			WHERE uuid = $3` + returningMultiAccount
	case "ds":
		query = `
			UPDATE my_compendium.multi_accounts
			SET discord_id = $1, discord_username = $2
			WHERE uuid = $3` + returningMultiAccount
	case "wa":
		query = `
			UPDATE my_compendium.multi_accounts
			SET whatsapp_id = $1, whatsapp_username = $2
			WHERE uuid = $3` + returningMultiAccount
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}

	row := d.db.QueryRow(query, id, username, uid)

	acc, err := scanMultiAccount(row)

	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}

	return acc, nil
}

func (d *Db) UpdateMultiAccountNickname(m models.MultiAccount) (*models.MultiAccount, error) {
	const query = `
		UPDATE my_compendium.multi_accounts
		SET nickname = $1
		WHERE uuid = $2` + returningMultiAccount

	row := d.db.QueryRow(query, m.Nickname, m.UUID)

	acc, err := scanMultiAccount(row)

	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (d *Db) UpdateMultiAccountAlts(m models.MultiAccount) (*models.MultiAccount, error) {
	const query = `
		UPDATE my_compendium.multi_accounts
		SET alts = $1
		WHERE uuid = $2` + returningMultiAccount

	row := d.db.QueryRow(query, m.Alts, m.UUID)

	acc, err := scanMultiAccount(row)

	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (d *Db) UpdateMultiAccountAvatarUrl(m models.MultiAccount) (*models.MultiAccount, error) {
	const query = `
		UPDATE my_compendium.multi_accounts
		SET avatarUrl = $1
		WHERE uuid = $2` + returningMultiAccount

	row := d.db.QueryRow(query, m.AvatarURL, m.UUID)

	acc, err := scanMultiAccount(row)

	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (d *Db) CreateMultiAccountFull(m models.MultiAccount) (*models.MultiAccount, error) {
	query := `
		INSERT INTO my_compendium.multi_accounts (
			nickname, telegram_id, telegram_username,
			discord_id, discord_username, whatsapp_id, whatsapp_username,
			avatarUrl, alts
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)` + returningMultiAccount

	row := d.db.QueryRow(query,
		m.Nickname, m.TelegramID, m.TelegramUsername,
		m.DiscordID, m.DiscordUsername, m.WhatsappID, m.WhatsappUsername,
		m.AvatarURL, pq.Array(m.Alts),
	)

	acc, err := scanMultiAccount(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		d.log.ErrorErr(err)
		return nil, err
	}

	return acc, nil
}

// FindMultiAccountByDiscordID ищет аккаунт по Discord ID
func (d *Db) FindMultiAccountByDiscordID(discordID string) (*models.MultiAccount, error) {
	query := `SELECT uuid, nickname, telegram_id, telegram_username, discord_id, discord_username,
			  whatsapp_id, whatsapp_username, created_at, avatarUrl, alts
			  FROM my_compendium.multi_accounts WHERE discord_id = $1`

	row := d.db.QueryRow(query, discordID)
	acc, err := scanMultiAccount(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Аккаунт не найден
		}
		d.log.ErrorErr(err)
		return nil, err
	}
	return acc, nil
}

// FindMultiAccountByTelegramID ищет аккаунт по Telegram ID
func (d *Db) FindMultiAccountByTelegramID(telegramID string) (*models.MultiAccount, error) {
	query := `SELECT uuid, nickname, telegram_id, telegram_username, discord_id, discord_username,
			  whatsapp_id, whatsapp_username, created_at, avatarUrl, alts
			  FROM my_compendium.multi_accounts WHERE telegram_id = $1`

	row := d.db.QueryRow(query, telegramID)
	acc, err := scanMultiAccount(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Аккаунт не найден
		}
		d.log.ErrorErr(err)
		return nil, err
	}
	return acc, nil
}

// LinkDiscordToAccount добавляет Discord ID к существующему аккаунту
func (d *Db) LinkDiscordToAccount(uuid string, discordID, discordUsername string) error {
	query := `UPDATE my_compendium.multi_accounts
			  SET discord_id = $1, discord_username = $2
			  WHERE uuid = $3`

	_, err := d.db.Exec(query, discordID, discordUsername, uuid)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}
	return nil
}

// LinkTelegramToAccount добавляет Telegram ID к существующему аккаунту
func (d *Db) LinkTelegramToAccount(uuid string, telegramID, telegramUsername string) error {
	query := `UPDATE my_compendium.multi_accounts
			  SET telegram_id = $1, telegram_username = $2
			  WHERE uuid = $3`

	_, err := d.db.Exec(query, telegramID, telegramUsername, uuid)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}
	return nil
}

// MergeAccounts объединяет два аккаунта в один
func (d *Db) MergeAccounts(keepUUID, mergeUUID string) error {
	// Конвертируем строки в UUID
	keepUUIDParsed, err := uuid.Parse(keepUUID)
	if err != nil {
		return fmt.Errorf("invalid keep UUID: %w", err)
	}
	mergeUUIDParsed, err := uuid.Parse(mergeUUID)
	if err != nil {
		return fmt.Errorf("invalid merge UUID: %w", err)
	}

	// Получаем данные аккаунта, который будет объединен
	mergeAcc, err := d.FindMultiAccountUUID(mergeUUIDParsed)
	if err != nil {
		return fmt.Errorf("failed to find merge account: %w", err)
	}
	if mergeAcc == nil {
		return fmt.Errorf("merge account not found")
	}

	// Обновляем аккаунт, который сохраняем
	updateQuery := `UPDATE my_compendium.multi_accounts SET
		telegram_id = COALESCE(telegram_id, $1),
		telegram_username = COALESCE(telegram_username, $2),
		discord_id = COALESCE(discord_id, $3),
		discord_username = COALESCE(discord_username, $4),
		whatsapp_id = COALESCE(whatsapp_id, $5),
		whatsapp_username = COALESCE(whatsapp_username, $6),
		avatarUrl = COALESCE(avatarUrl, $7),
		alts = CASE WHEN alts IS NULL OR array_length(alts, 1) IS NULL THEN $8
					ELSE alts || $8 END
		WHERE uuid = $9`

	_, err = d.db.Exec(updateQuery,
		mergeAcc.TelegramID, mergeAcc.TelegramUsername,
		mergeAcc.DiscordID, mergeAcc.DiscordUsername,
		mergeAcc.WhatsappID, mergeAcc.WhatsappUsername,
		mergeAcc.AvatarURL, pq.Array(mergeAcc.Alts), keepUUIDParsed)
	if err != nil {
		d.log.ErrorErr(err)
		return fmt.Errorf("failed to update account: %w", err)
	}

	// Удаляем объединенный аккаунт
	deleteQuery := `DELETE FROM my_compendium.multi_accounts WHERE uuid = $1`
	_, err = d.db.Exec(deleteQuery, mergeUUID)
	if err != nil {
		d.log.ErrorErr(err)
		return fmt.Errorf("failed to delete merged account: %w", err)
	}

	return nil
}
