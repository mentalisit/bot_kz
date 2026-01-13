package postgresv2

import (
	"compendium/models"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

func (d *Db) GuildSave(g models.MultiAccountGuildV2) (*models.MultiAccountGuildV2, error) {
	// Если GId не передан, генерируем его сами (или доверьте это БД через DEFAULT)
	if g.GId == uuid.Nil {
		g.GId = uuid.New()
	}

	query := `
        INSERT INTO my_compendium.guilds (gid, guildname, channels, avatarurl)
        VALUES (:gid, :guildname, :channels, :avatarurl)
        ON CONFLICT (gid) DO UPDATE SET
            guildname = EXCLUDED.guildname,
            channels = EXCLUDED.channels,
            avatarurl = EXCLUDED.avatarurl
        RETURNING *`

	var res models.MultiAccountGuildV2
	// NamedQuery подставит всё сам, включая JSON из мапы Channels
	rows, err := d.db.NamedQuery(query, g)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&res)
	}
	return &res, err
}

func (d *Db) GuildUpdateAvatar(u models.MultiAccountGuildV2) error {
	query := `UPDATE my_compendium.guilds SET avatarurl = $1 WHERE gid = $2`
	_, err := d.db.Exec(query, u.AvatarUrl, u.GId)
	return err
}

func (d *Db) GuildUpdateGuildName(u models.MultiAccountGuildV2) error {
	query := `UPDATE my_compendium.guilds SET guildname = $1 WHERE gid = $2`
	_, err := d.db.Exec(query, u.GuildName, u.GId)
	return err
}

func (d *Db) GuildUpdateChannels(u models.MultiAccountGuildV2) error {
	// Просто передаем мапу u.Channels. sqlx сам вызовет .Value() и сделает JSON
	query := `UPDATE my_compendium.guilds SET channels = $1 WHERE gid = $2`
	_, err := d.db.Exec(query, u.Channels, u.GId)
	return err
}

func (d *Db) GuildGet(gid uuid.UUID) (*models.MultiAccountGuildV2, error) {
	var guild models.MultiAccountGuildV2

	// sqlx автоматически вызовет метод Scan у поля Channels (тип GuildChannels)
	query := `SELECT gid, guildname, channels, avatarurl FROM my_compendium.guilds WHERE gid = $1`

	err := d.db.Get(&guild, query, gid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Или ваша логика обработки пустого результата
		}
		return nil, err
	}

	return &guild, nil
}

func (d *Db) GuildGetById(guildid string) (*models.MultiAccountGuildV2, error) {
	gid, err := uuid.Parse(guildid)
	if err == nil {
		return d.GuildGet(gid)
	}

	return d.GuildGetChannel(guildid)
}

func (d *Db) GuildGetChannel(channelID string) (*models.MultiAccountGuildV2, error) {
	var guild models.MultiAccountGuildV2

	// Но самый надежный и быстрый способ для вашего случая (JSONB "key": ["id1", "id2"]):
	optimizedQuery := `
        SELECT gid, guildname, channels, avatarurl 
        FROM my_compendium.guilds 
        WHERE EXISTS (
            SELECT 1 FROM jsonb_each(channels) WHERE value ? $1
        ) LIMIT 1`

	err := d.db.Get(&guild, optimizedQuery, channelID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &guild, nil
}

func (d *Db) GetChatsRoles(chatID int64) ([]models.CorpRole, error) {
	var roles []models.CorpRole

	// Выбираем id и name. Поле ChatId мы заполним автоматически из БД,
	// если добавим его в SELECT, либо вручную, если нужно.
	query := `SELECT id, name, chat_id FROM telegram.roles WHERE chat_id = $1`

	// Select сразу сканирует все строки в слайс структур
	err := d.db.Select(&roles, query, chatID)
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}

	return roles, nil
}

// IsUserSubscribedToRole проверяет, подписан ли пользователь на указанную роль
func (d *Db) IsUserSubscribedToRole(userID, roleID int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM telegram.user_roles WHERE user_id = $1 AND role_id = $2)`

	// sqlx.Get может работать с простыми типами, не только со структурами
	err := d.db.Get(&exists, query, userID, roleID)
	if err != nil {
		return false, err
	}

	return exists, nil
}
