package postgresv2

import (
	"compendium/models"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

const returningGuilds = `
		RETURNING gid, guildname, channels, avatarUrl`

func scanGuilds(row *sql.Row) (*models.MultiAccountGuildV2, error) {
	var g models.MultiAccountGuildV2
	var channelsData []byte

	err := row.Scan(&g.GId, &g.GuildName, &channelsData, &g.AvatarUrl)
	if err != nil {
		return nil, err
	}

	// Convert JSONB to map[string][]string
	g.Channels = make(map[string][]string)
	if len(channelsData) > 0 {
		err = json.Unmarshal(channelsData, &g.Channels)
		if err != nil {
			return nil, err
		}
	}
	return &g, nil
}
func (d *Db) GuildInsert(u models.MultiAccountGuildV2) (*models.MultiAccountGuildV2, error) {
	insert := `INSERT INTO my_compendium.guilds(guildName,channels,avatarUrl) VALUES ($1,$2,$3)` + returningGuilds
	row := d.db.QueryRow(insert, u.GuildName, u.ChannelsBytes(), u.AvatarUrl)
	return scanGuilds(row)
}
func (d *Db) GuildInsertFull(u models.MultiAccountGuildV2) (*models.MultiAccountGuildV2, error) {
	insert := `INSERT INTO my_compendium.guilds(gid, guildName,channels,avatarUrl) 
				VALUES ($1,$2,$3,$4)` + returningGuilds
	row := d.db.QueryRow(insert, u.GId, u.GuildName, u.ChannelsBytes(), u.AvatarUrl)
	return scanGuilds(row)
}

func (d *Db) GuildUpdateAvatar(u models.MultiAccountGuildV2) error {
	upd := `update my_compendium.guilds set avatarUrl = $1 where gid = $2`
	_, err := d.db.Exec(upd, u.AvatarUrl, u.GId)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) GuildUpdateGuildName(u models.MultiAccountGuildV2) error {
	upd := `update my_compendium.guilds set guildName = $1 where gid = $2`
	_, err := d.db.Exec(upd, u.GuildName, u.GId)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) GuildUpdateChannels(u models.MultiAccountGuildV2) error {
	upd := `update my_compendium.guilds set channels = $1 where gid = $2`
	_, err := d.db.Exec(upd, u.ChannelsBytes(), u.GId)
	if err != nil {
		return err
	}
	return nil
}

func (d *Db) GuildGet(gid *uuid.UUID) (*models.MultiAccountGuildV2, error) {
	query := `SELECT gid, GuildName, Channels, AvatarUrl FROM my_compendium.guilds WHERE gid = $1`
	row := d.db.QueryRow(query, gid)
	return scanGuilds(row)
}

func (d *Db) GuildGetById(guild string) (*models.MultiAccountGuildV2, error) {
	if gid, err := uuid.Parse(guild); err == nil {
		return d.GuildGet(&gid)
	}
	return nil, nil
}

func (d *Db) GuildGetChatId(ChatId string) (*models.MultiAccountGuildV2, error) {
	query := `SELECT gid, GuildName, Channels, AvatarUrl FROM my_compendium.guilds
		WHERE EXISTS (
		    SELECT 1 FROM jsonb_object_keys(channels) AS k
		    CROSS JOIN jsonb_array_elements_text(channels->k) AS v
		    WHERE v = $1
		)`
	row := d.db.QueryRow(query, ChatId)
	return scanGuilds(row)
}

// GetChatRoles возвращает роли определенного чата
func (d *Db) GetChatsRoles(chatID int64) ([]models.CorpRole, error) {
	query := `SELECT id, name FROM telegram.roles WHERE chat_id = $1`

	rows, err := d.db.Query(query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to query chats roles: %w", err)
	}
	defer rows.Close()

	var roles []models.CorpRole
	for rows.Next() {
		var r models.CorpRole
		r.ChatId = chatID
		if err := rows.Scan(&r.Id, &r.Name); err != nil {
			return nil, fmt.Errorf("failed to scan roles: %w", err)
		}
		roles = append(roles, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return roles, nil
}

// IsUserSubscribedToRole проверяет, подписан ли пользователь на указанную роль
func (d *Db) IsUserSubscribedToRole(userID, roleID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM telegram.user_roles WHERE user_id = $1 AND role_id = $2)`

	var isSubscribed bool
	err := d.db.QueryRow(query, userID, roleID).Scan(&isSubscribed)
	if err != nil {
		return false, fmt.Errorf("failed to check user role subscription: %w", err)
	}

	return isSubscribed, nil
}
