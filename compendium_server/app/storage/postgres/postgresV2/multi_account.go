package postgresv2

import (
	"compendium_s/models"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

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

func (d *Db) CreateMultiAccountFull(m models.MultiAccount) (*models.MultiAccount, error) {
	if m.UUID == uuid.Nil {
		m.UUID = uuid.New()
	}

	// Логируем если никнейм пустой
	if m.Nickname == "" {
		d.log.Warn("CreateMultiAccountFull called with empty nickname",
			zap.String("uuid", m.UUID.String()),
			zap.String("telegram_id", m.TelegramID),
			zap.String("discord_id", m.DiscordID),
			zap.String("whatsapp_id", m.WhatsappID))
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

// Вспомогательная функция для поиска "хоть по чему-нибудь"
//func (d *Db) findAnyBySocialIDs(tg, ds, wa string) (*models.MultiAccount, error) {
//	query := `
//        SELECT uuid, nickname, telegram_id, telegram_username, discord_id,
//               discord_username, whatsapp_id, whatsapp_username, created_at, avatarUrl, alts
//        FROM my_compendium.multi_accounts
//        WHERE (telegram_id = $1 AND telegram_id != '')
//           OR (discord_id = $2 AND discord_id != '')
//           OR (whatsapp_id = $3 AND whatsapp_id != '')
//        LIMIT 1`
//
//	row := d.db.QueryRow(query, tg, ds, wa)
//	return scanMultiAccount(row)
//}
//const returningMultiAccount = `
//		RETURNING uuid, nickname,
//		          telegram_id, telegram_username,
//		          discord_id, discord_username,
//		          whatsapp_id, whatsapp_username,
//		          created_at,
//				  avatarUrl, alts`
//
//func scanMultiAccount(row *sql.Row) (*models.MultiAccount, error) {
//	var acc models.MultiAccount
//	var telegramID, discordID, whatsappID sql.NullString
//
//	err := row.Scan(
//		&acc.UUID, &acc.Nickname,
//		&telegramID, &acc.TelegramUsername,
//		&discordID, &acc.DiscordUsername,
//		&whatsappID, &acc.WhatsappUsername,
//		&acc.CreatedAt,
//		&acc.AvatarURL, pq.Array(&acc.Alts),
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	if telegramID.Valid {
//		acc.TelegramID = telegramID.String
//	}
//	if discordID.Valid {
//		acc.DiscordID = discordID.String
//	}
//	if whatsappID.Valid {
//		acc.WhatsappID = whatsappID.String
//	}
//	if acc.UUID.String() == "00000000-0000-0000-0000-000000000000" {
//		return nil, errors.New("invalid UUID")
//	}
//
//	return &acc, nil
//}
//
//// Multi-account methods
//func (d *Db) FindMultiAccountUUID(uid uuid.UUID) (*models.MultiAccount, error) {
//	var acc models.MultiAccount
//	var telegramID, discordID, whatsappID sql.NullString
//
//	query := `SELECT uuid, nickname, telegram_id, telegram_username, discord_id, discord_username,
//			  whatsapp_id, whatsapp_username, created_at, avatarUrl, alts
//			  FROM my_compendium.multi_accounts WHERE uuid = $1`
//
//	err := d.db.QueryRow(query, uid).Scan(
//		&acc.UUID, &acc.Nickname,
//		&telegramID, &acc.TelegramUsername,
//		&discordID, &acc.DiscordUsername,
//		&whatsappID, &acc.WhatsappUsername,
//		&acc.CreatedAt,
//		&acc.AvatarURL, pq.Array(&acc.Alts),
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	if telegramID.Valid {
//		acc.TelegramID = telegramID.String
//	}
//	if discordID.Valid {
//		acc.DiscordID = discordID.String
//	}
//	if whatsappID.Valid {
//		acc.WhatsappID = whatsappID.String
//	}
//
//	return &acc, nil
//}
//func (d *Db) FindMultiAccountByUserId(userid string) (*models.MultiAccount, error) {
//	var acc models.MultiAccount
//	var telegramID, discordID, whatsappID sql.NullString
//
//	query := `SELECT uuid, nickname, telegram_id, telegram_username, discord_id, discord_username,
//			  whatsapp_id, whatsapp_username, created_at, avatarUrl, alts
//			  FROM my_compendium.multi_accounts WHERE discord_id = $1 OR telegram_id = $1 OR whatsapp_id = $1`
//
//	err := d.db.QueryRow(query, userid).Scan(
//		&acc.UUID, &acc.Nickname,
//		&telegramID, &acc.TelegramUsername,
//		&discordID, &acc.DiscordUsername,
//		&whatsappID, &acc.WhatsappUsername,
//		&acc.CreatedAt,
//		&acc.AvatarURL, pq.Array(&acc.Alts),
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	if telegramID.Valid {
//		acc.TelegramID = telegramID.String
//	}
//	if discordID.Valid {
//		acc.DiscordID = discordID.String
//	}
//	if whatsappID.Valid {
//		acc.WhatsappID = whatsappID.String
//	}
//
//	return &acc, nil
//}
//// CreateMultiAccountFull создает новый аккаунт или обновляет существующий со всеми данными
//func (d *Db) CreateMultiAccountFull(m models.MultiAccount) (*models.MultiAccount, error) {
//	if m.UUID.String() == "00000000-0000-0000-0000-000000000000" {
//		m.UUID = uuid.New()
//	}
//	query := `
//		INSERT INTO my_compendium.multi_accounts (uuid,
//			nickname, telegram_id, telegram_username,
//			discord_id, discord_username, whatsapp_id, whatsapp_username,
//			avatarUrl, alts
//		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
//		ON CONFLICT (uuid) DO UPDATE SET
//			nickname = EXCLUDED.nickname,
//			telegram_id = EXCLUDED.telegram_id,
//			telegram_username = EXCLUDED.telegram_username,
//			discord_id = EXCLUDED.discord_id,
//			discord_username = EXCLUDED.discord_username,
//			whatsapp_id = EXCLUDED.whatsapp_id,
//			whatsapp_username = EXCLUDED.whatsapp_username,
//			avatarUrl = EXCLUDED.avatarUrl,
//			alts = EXCLUDED.alts` + returningMultiAccount
//
//	row := d.db.QueryRow(query, m.UUID,
//		m.Nickname, m.TelegramID, m.TelegramUsername,
//		m.DiscordID, m.DiscordUsername, m.WhatsappID, m.WhatsappUsername,
//		m.AvatarURL, pq.Array(m.Alts), // Конвертируем []string в PostgreSQL массив
//	)
//
//	acc, err := scanMultiAccount(row)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return nil, nil
//		}
//		d.log.ErrorErr(err)
//		return nil, err
//	}
//
//	return acc, nil
//}
//func (d *Db) CreateMultiAccountFull(m models.MultiAccount) (*models.MultiAccount, error) {
//	// 1. Подготовка данных (защита от NULL в массивах)
//	alts := m.Alts
//	if alts == nil {
//		alts = []string{}
//	}
//
//	if m.UUID == uuid.Nil {
//		m.UUID = uuid.New()
//	}
//
//	// 2. Основной запрос (Upsert по UUID)
//	query := `
//       INSERT INTO my_compendium.multi_accounts (uuid,
//          nickname, telegram_id, telegram_username,
//          discord_id, discord_username, whatsapp_id, whatsapp_username,
//          avatarUrl, alts
//       ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
//       ON CONFLICT (uuid) DO UPDATE SET
//         nickname = EXCLUDED.nickname,
//         telegram_id = EXCLUDED.telegram_id,
//         telegram_username = EXCLUDED.telegram_username,
//         discord_id = EXCLUDED.discord_id,
//         discord_username = EXCLUDED.discord_username,
//         whatsapp_id = EXCLUDED.whatsapp_id,
//         whatsapp_username = EXCLUDED.whatsapp_username,
//         avatarUrl = EXCLUDED.avatarUrl,
//         alts = EXCLUDED.alts` + returningMultiAccount
//
//	row := d.db.QueryRow(query, m.UUID,
//		m.Nickname, m.TelegramID, m.TelegramUsername,
//		m.DiscordID, m.DiscordUsername, m.WhatsappID, m.WhatsappUsername,
//		m.AvatarURL, pq.Array(alts))
//
//	acc, err := scanMultiAccount(row)
//
//	// 3. ОБРАБОТКА ОШИБКИ ДУБЛИКАТА
//	if err != nil {
//		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
//			// Если мы здесь, значит uuid новый, но telegram_id (или другой ID) уже занят.
//			// Просто ищем существующую запись по социальным ID.
//			d.log.Info("Конфликт уникальности, возвращаем существующую запись")
//			return d.findAnyBySocialIDs(m.TelegramID, m.DiscordID, m.WhatsappID)
//		}
//		return nil, err
//	}
//
//	return acc, nil
//}
