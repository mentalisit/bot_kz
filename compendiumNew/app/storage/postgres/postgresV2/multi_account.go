package postgresv2

import (
	"compendium/models"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Multi-account methods
func (d *Db) FindMultiAccountUUID(uid uuid.UUID) (*models.MultiAccount, error) {
	var acc models.MultiAccount

	// sqlx сам сопоставит колонки с тегами db:"..." в структуре
	query := `SELECT * FROM my_compendium.multi_accounts WHERE uuid = $1`

	err := d.db.Get(&acc, query, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // не найдено — не ошибка
		}
		return nil, err
	}

	return &acc, nil
}

func (d *Db) FindMultiAccountByUserId(userid string) (*models.MultiAccount, error) {
	var acc models.MultiAccount

	// Используем именованный запрос для красоты или обычный
	query := `SELECT * FROM my_compendium.multi_accounts 
              WHERE discord_id = $1 OR telegram_id = $1 OR whatsapp_id = $1 
              LIMIT 1`

	// sqlx.Get сам заполнит все поля, включая слайс Alts и указатели ID
	err := d.db.Get(&acc, query, userid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &acc, nil
}

func (d *Db) FindMultiAccountByUserName(userName string) (*models.MultiAccount, error) {
	var acc models.MultiAccount

	// Используем SELECT *, так как sqlx сам сопоставит колонки с полями структуры
	query := `
        SELECT * FROM my_compendium.multi_accounts 
        WHERE nickname = $1 
           OR $1 = ANY(alts) 
           OR telegram_username = $1 
           OR discord_username = $1 
           OR whatsapp_username = $1 
        LIMIT 1`

	// d.db.Get — это обертка sqlx, которая делает QueryRow + Scan автоматически
	err := d.db.Get(&acc, query, userName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		d.log.ErrorErr(err)
		return nil, err
	}

	return &acc, nil
}

func (d *Db) CreateMultiAccountWithPlatform(id, nickname, platform, username string) (*models.MultiAccount, error) {
	// Определяем имя колонки ID и Username в зависимости от платформы
	var ma models.MultiAccount

	switch platform {
	case "tg":
		ma.TelegramUsername = nickname
		ma.TelegramID = id
	case "ds":
		ma.DiscordUsername = nickname
		ma.DiscordID = id
	case "wa":
		ma.WhatsappUsername = nickname
		ma.WhatsappID = id
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}

	return d.CreateMultiAccountFull(ma)
}

func (d *Db) UpdateMultiAccountAlts(m models.MultiAccount) (*models.MultiAccount, error) {
	// В sqlx мы просто используем RETURNING * для заполнения всей структуры
	const query = `
       UPDATE my_compendium.multi_accounts
       SET alts = $1
       WHERE uuid = $2
       RETURNING *`

	var acc models.MultiAccount

	// sqlx.Get выполнит запрос и вызовет наш метод Scan для поля alts
	err := d.db.Get(&acc, query, m.Alts, m.UUID)
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}

	return &acc, nil
}

func (d *Db) CreateMultiAccountFull(m models.MultiAccount) (*models.MultiAccount, error) {
	if m.UUID == uuid.Nil {
		m.UUID = uuid.New()
	}

	query := `
       INSERT INTO my_compendium.multi_accounts (
          uuid, nickname, telegram_id, telegram_username, 
          discord_id, discord_username, whatsapp_id, whatsapp_username, 
          avatarurl, alts
       ) VALUES (:uuid, :nickname, :telegram_id, :telegram_username, 
                 :discord_id, :discord_username, :whatsapp_id, :whatsapp_username, 
                 :avatarurl, :alts)
       ON CONFLICT (uuid) DO UPDATE SET 
          nickname=EXCLUDED.nickname, telegram_id=EXCLUDED.telegram_id, 
          avatarurl=EXCLUDED.avatarurl, alts=EXCLUDED.alts
       RETURNING *`

	// Выполняем запрос, передавая структуру m целиком
	rows, err := d.db.NamedQuery(query, m)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return d.findAnyBySocialIDs(m.TelegramID, m.DiscordID, m.WhatsappID)
		}
		return nil, err
	}
	defer rows.Close()

	var acc models.MultiAccount
	if rows.Next() {
		err = rows.StructScan(&acc)
	}
	return &acc, err
}

func (d *Db) findAnyBySocialIDs(tg, ds, wa string) (*models.MultiAccount, error) {
	var acc models.MultiAccount
	query := `
        SELECT * FROM my_compendium.multi_accounts 
        WHERE (telegram_id = $1 AND telegram_id != '')
           OR (discord_id = $2 AND discord_id != '')
           OR (whatsapp_id = $3 AND whatsapp_id != '')
        LIMIT 1`

	err := d.db.Get(&acc, query, tg, ds, wa)
	return &acc, err
}

func (d *Db) UpdateMultiAccountNickname(m models.MultiAccount) (*models.MultiAccount, error) {
	const query = `
       UPDATE my_compendium.multi_accounts
       SET nickname = $1
       WHERE uuid = $2
       RETURNING *`

	var acc models.MultiAccount
	err := d.db.Get(&acc, query, m.AvatarURL, m.UUID)
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}

	return &acc, nil
}

func (d *Db) UpdateMultiAccount(m models.MultiAccount) (*models.MultiAccount, error) {
	query := `
       UPDATE my_compendium.multi_accounts
       SET nickname = :nickname,
           telegram_id = :telegram_id, telegram_username = :telegram_username,
           discord_id = :discord_id, discord_username = :discord_username,
           whatsapp_id = :whatsapp_id, whatsapp_username = :whatsapp_username,
           avatarurl = :avatarurl, alts = :alts
       WHERE uuid = :uuid
       RETURNING *`

	// Используем NamedQuery, чтобы sqlx сам вытащил данные из структуры m
	rows, err := d.db.NamedQuery(query, m)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var acc models.MultiAccount
	if rows.Next() {
		err = rows.StructScan(&acc)
	}
	return &acc, err
}

func (d *Db) UpdateMultiAccountAvatarUrl(m models.MultiAccount) (*models.MultiAccount, error) {
	const query = `
       UPDATE my_compendium.multi_accounts
       SET avatarurl = $1
       WHERE uuid = $2
       RETURNING *`

	var acc models.MultiAccount
	err := d.db.Get(&acc, query, m.AvatarURL, m.UUID)
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}

	return &acc, nil
}
