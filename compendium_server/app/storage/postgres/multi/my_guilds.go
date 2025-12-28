package multi

import (
	"compendium_s/models"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func (d *Db) GuildGetV2(gid *uuid.UUID) (*models.MultiAccountGuildV2, error) {
	var guild models.MultiAccountGuildV2
	var channelsData []byte

	query := `SELECT gid, GuildName, Channels, AvatarUrl FROM my_compendium.guilds WHERE gid = $1`
	err := d.db.QueryRow(context.Background(), query, gid).Scan(
		&guild.GId, &guild.GuildName, &channelsData, &guild.AvatarUrl,
	)
	if err != nil {
		return nil, err
	}

	// Convert JSONB to map[string]string
	guild.Channels = make(map[string][]string)
	if len(channelsData) > 0 {
		err = json.Unmarshal(channelsData, &guild.Channels)
		if err != nil {
			d.log.ErrorErr(err)
			return nil, err
		}
	}

	return &guild, nil
}

func (d *Db) GuildGetByIdV2(guildid string) (*models.MultiAccountGuildV2, error) {
	if gid, err := uuid.Parse(guildid); err == nil {
		return d.GuildGetV2(&gid)
	}
	return nil, nil
}

// GetChatRoles возвращает роли определенного чата
func (d *Db) GetChatsRoles(chatID int64) ([]models.CorpRole, error) {
	query := `SELECT id, name FROM telegram.roles WHERE chat_id = $1`

	rows, err := d.db.Query(context.Background(), query, chatID)
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
	err := d.db.QueryRow(context.Background(), query, userID, roleID).Scan(&isSubscribed)
	if err != nil {
		return false, fmt.Errorf("failed to check user role subscription: %w", err)
	}

	return isSubscribed, nil
}
