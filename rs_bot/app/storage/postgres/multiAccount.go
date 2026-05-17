package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"rs/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	MAReturn = ` RETURNING uuid, nickname, telegram_id, telegram_username, discord_id, discord_username, whatsapp_id, whatsapp_username, avatarurl, alts, created_at`
	MASelect = `SELECT uuid, nickname, telegram_id, telegram_username, discord_id, discord_username, whatsapp_id, whatsapp_username, avatarurl, alts, created_at `
)

func scanMultiAccount(row *sql.Row) (*models.MultiAccount, error) {
	var acc models.MultiAccount

	var telegramID, discordID, whatsappID sql.NullString
	err := row.Scan(
		&acc.UUID, &acc.Nickname,
		&telegramID, &acc.TelegramUsername,
		&discordID, &acc.DiscordUsername,
		&whatsappID, &acc.WhatsappUsername,
		&acc.AvatarURL,
		pq.Array(&acc.Alts), &acc.CreatedAt,
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
	// Логируем если никнейм пустой
	if nickname == "" {
		d.log.Warn("CreateMultiAccountWithPlatform called with empty nickname",
			zap.String("user_id", id),
			zap.String("platform", platform),
			zap.String("username", username))
	}

	var query string

	switch platform {
	case "tg":
		query = `
			INSERT INTO my_compendium.multi_accounts (nickname, telegram_id, telegram_username)
			VALUES ($1, $2, $3)` + MAReturn
	case "ds":
		query = `
			INSERT INTO my_compendium.multi_accounts (nickname, discord_id, discord_username)
			VALUES ($1, $2, $3)` + MAReturn
	case "wa":
		query = `
			INSERT INTO my_compendium.multi_accounts (nickname, whatsapp_id, whatsapp_username)
			VALUES ($1, $2, $3)` + MAReturn
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

func (d *Db) FindMultiAccountByUserId(userId string) (*models.MultiAccount, error) {
	var selectQuery = MASelect + `FROM my_compendium.multi_accounts
		WHERE telegram_id = $1 or discord_id = $1 or whatsapp_id = $1`

	row := d.db.QueryRow(selectQuery, userId)

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

func (d *Db) TechnologiesGetAll(u uuid.UUID) ([]models.TechUser, error) {
	var results []models.TechUser

	query := `
		SELECT username, tech
		FROM my_compendium.technologies
		WHERE uid = $1
	`

	rows, err := d.db.Query(query, u)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.TechUser
		var techBytes []byte

		err = rows.Scan(&r.Name, &techBytes)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(techBytes, &r.Tech)
		if err != nil {
			return nil, err
		}

		results = append(results, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
