package postgresv2

import (
	"compendium_s/models"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	MAReturn = ` RETURNING uuid, nickname, telegram_id, telegram_username, discord_id, discord_username, whatsapp_id, whatsapp_username, avatarurl, alts, created_at`
	MASelect = `SELECT uuid, nickname, telegram_id, telegram_username, discord_id, discord_username, whatsapp_id, whatsapp_username, avatarurl, alts, created_at `
)

func (d *Db) FindMultiAccountUUID(uid uuid.UUID) (*models.MultiAccount, error) {
	var acc models.MultiAccount

	// sqlx сам сопоставит колонки с тегами db:"..." в структуре
	query := MASelect + `FROM my_compendium.multi_accounts WHERE uuid = $1`

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
	query := MASelect + `FROM my_compendium.multi_accounts 
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
          avatarurl=EXCLUDED.avatarurl, alts=EXCLUDED.alts` + MAReturn

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
	query := MASelect + `FROM my_compendium.multi_accounts 
        WHERE (telegram_id = $1 AND telegram_id != '')
           OR (discord_id = $2 AND discord_id != '')
           OR (whatsapp_id = $3 AND whatsapp_id != '')
        LIMIT 1`

	err := d.db.Get(&acc, query, tg, ds, wa)
	return &acc, err
}
