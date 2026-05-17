package postgres

import (
	"bridge/models"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

func (d *Db) GuildGetChannel(guildId string) (*models.MultiAccountGuildV2, error) {
	var guild models.MultiAccountGuildV2

	// Но самый надежный и быстрый способ для вашего случая (JSONB "key": ["id1", "id2"]):
	optimizedQuery := `
        SELECT gid, guildname, channels, avatarurl 
        FROM my_compendium.guilds 
        WHERE EXISTS (
            SELECT 1 FROM jsonb_each(channels) WHERE value ? $1
        ) LIMIT 1`

	err := d.db.Get(&guild, optimizedQuery, guildId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &guild, nil
}

func (d *Db) FindMultiAccountUidByUserId(userid string) (uuid.UUID, error) {
	var uid uuid.UUID

	// Используем именованный запрос для красоты или обычный
	query := `SELECT uuid FROM my_compendium.multi_accounts 
              WHERE discord_id = $1 OR telegram_id = $1 OR whatsapp_id = $1 
              LIMIT 1`

	// sqlx.Get сам заполнит все поля, включая слайс Alts и указатели ID
	err := d.db.Get(&uid, query, userid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, nil
		}
		return uuid.Nil, err
	}

	return uid, nil
}
