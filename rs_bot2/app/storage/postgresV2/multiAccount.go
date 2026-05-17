package postgresV2

import (
	"database/sql"
	"errors"
	"fmt"
	"rs/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	MAReturn = ` RETURNING uuid, nickname, telegram_id, telegram_username, discord_id, discord_username, whatsapp_id, whatsapp_username, avatarurl, alts, created_at, active_account, data`
	MASelect = `SELECT uuid, nickname, telegram_id, telegram_username, discord_id, discord_username, whatsapp_id, whatsapp_username, avatarurl, alts, created_at, active_account, data `
)

func scanMultiAccount(row pgx.Row) (*models.MultiAccount, error) {
	var acc models.MultiAccount

	var telegramID, discordID, whatsappID sql.NullString
	err := row.Scan(
		&acc.UUID, &acc.Nickname,
		&telegramID, &acc.TelegramUsername,
		&discordID, &acc.DiscordUsername,
		&whatsappID, &acc.WhatsappUsername,
		&acc.AvatarURL,
		pq.Array(&acc.Alts), &acc.CreatedAt,
		&acc.ActiveAccount, &acc.Data,
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
		if errors.Is(err, pgx.ErrNoRows) {
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		d.log.ErrorErr(err)
		return nil, err
	}
	return acc, nil
}

func (d *Db) FindMultiAccountByUUId(uId string) (*models.MultiAccount, error) {

	var selectQuery = MASelect + `FROM my_compendium.multi_accounts
		WHERE uuid = $1`

	row := d.db.QueryRow(selectQuery, uId)

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

// GetUserNicknamesByUUIDs возвращает map[uuid]nickname для списка UUID
func (d *Db) GetUserNicknamesByUUIDs(uuids []string) (map[string]string, error) {
	if len(uuids) == 0 {
		return map[string]string{}, nil
	}

	// Преобразуем строки в UUID
	uuidObjs := make([]uuid.UUID, 0, len(uuids))
	for _, u := range uuids {
		if uid, err := uuid.Parse(u); err == nil {
			uuidObjs = append(uuidObjs, uid)
		}
	}

	query := `SELECT uuid, nickname FROM my_compendium.multi_accounts WHERE uuid = ANY($1)`
	rows, err := d.db.Query(query, uuidObjs)
	if err != nil {
		d.log.ErrorErr(err)
		return map[string]string{}, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var uid uuid.UUID
		var nickname string
		if err := rows.Scan(&uid, &nickname); err == nil {
			result[uid.String()] = nickname
		}
	}

	return result, nil
}

func (d *Db) TechnologiesGetAll(u uuid.UUID) ([]models.TechUser, error) {
	var results []models.TechUser

	query := `SELECT username, tech FROM my_compendium.technologies WHERE uid = $1`

	res, err := d.db.Query(query, u)
	if err != nil {
		return nil, err
	}
	for res.Next() {
		var r models.TechUser
		err = res.Scan(&r.Name, &r.Tech)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

func (d *Db) UpdateMultiAccountNickAltsActive(uid uuid.UUID, nickName string, alts []string, activeAccount string) error {

	query := `UPDATE my_compendium.multi_accounts
			SET nickname = $1, alts = $2, active_account = $3
			WHERE uuid = $4`

	_, err := d.db.Exec(query, nickName, alts, activeAccount, uid.String())
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}
	return nil
}
func (d *Db) UpdateMultiAccountSocial(uid uuid.UUID, platform, socialId, username string) error {

	var clearQuery string
	var query string

	// Если мы не отвязываем аккаунт (socialId != ""), то сначала очищаем его у других
	if socialId != "" {
		switch platform {
		case "discord":
			clearQuery = `UPDATE my_compendium.multi_accounts 
	                 SET discord_id = '', discord_username = '' 
	                 WHERE discord_id = $1 AND uuid != $2`
		case "telegram":
			clearQuery = `UPDATE my_compendium.multi_accounts 
	                 SET telegram_id = '', telegram_username = '' 
	                 WHERE telegram_id = $1 AND uuid != $2`
		}

		if clearQuery != "" {
			_, err := d.db.Exec(clearQuery, socialId, uid.String())
			if err != nil {
				d.log.ErrorErr(err)
				// Не прерываемся, пробуем основной апдейт
			}
		}
	}

	switch platform {
	case "discord":
		query = `UPDATE my_compendium.multi_accounts 
                 SET discord_id = $1, discord_username = $2 
                 WHERE uuid = $3`
	case "telegram":
		query = `UPDATE my_compendium.multi_accounts 
                 SET telegram_id = $1, telegram_username = $2 
                 WHERE uuid = $3`
	default:
		return fmt.Errorf("unsupported platform for social update: %s", platform)
	}

	_, err := d.db.Exec(query, socialId, username, uid.String())
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}
	return nil
}
