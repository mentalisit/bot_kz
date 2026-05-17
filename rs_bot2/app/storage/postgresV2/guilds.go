package postgresV2

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"rs/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (d *Db) UserCorporationsGet(mAcc *models.MultiAccount) ([]models.MultiAccountGuildV2, error) {

	// Добавляем поле channels в выборку
	query := `
        SELECT g.gid, g.guildname, g.channels, g.avatarurl 
        FROM my_compendium.guilds g
        WHERE g.gid = ANY (
            SELECT UNNEST(guildids) 
            FROM my_compendium.corpmember 
            WHERE uid = $1
        )`

	rows, err := d.db.Query(query, mAcc.UUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user corporations: %w", err)
	}
	defer rows.Close()

	var corporations []models.MultiAccountGuildV2

	for rows.Next() {
		var guild models.MultiAccountGuildV2

		err = rows.Scan(
			&guild.GId,
			&guild.GuildName,
			&guild.Channels,
			&guild.AvatarUrl,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		corporations = append(corporations, guild)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return corporations, nil
}

func (d *Db) GuildGet(gid uuid.UUID) (*models.MultiAccountGuildV2, error) {

	var guild models.MultiAccountGuildV2

	query := `SELECT gid, guildname, channels, avatarurl FROM my_compendium.guilds WHERE gid = $1`

	err := d.db.QueryRow(query, gid).Scan(&guild.GId, &guild.GuildName, &guild.Channels, &guild.AvatarUrl)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &guild, nil
}

func (d *Db) GuildGetChannel(ctx context.Context, channelID string) (*models.MultiAccountGuildV2, error) {
	var guild models.MultiAccountGuildV2

	// Оставляем ваш оптимизированный запрос
	optimizedQuery := `
        SELECT gid, guildname, channels, avatarurl 
        FROM my_compendium.guilds 
        WHERE EXISTS (
            SELECT 1 FROM jsonb_each(channels) WHERE value ? $1
        ) LIMIT 1`

	err := d.db.QueryRow(optimizedQuery, channelID).Scan(
		&guild.GId,
		&guild.GuildName,
		&guild.Channels,
		&guild.AvatarUrl,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get guild by channel: %w", err)
	}

	return &guild, nil
}

func (d *Db) GuildGetFull(gid uuid.UUID) (*models.MultiAccountGuildV2, error) {
	var guild models.MultiAccountGuildV2

	query := `SELECT gid, guildname, channels, avatarurl, data FROM my_compendium.guilds WHERE gid = $1`

	err := d.db.QueryRow(query, gid).Scan(&guild.GId, &guild.GuildName, &guild.Channels, &guild.AvatarUrl, &guild.Data)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &guild, nil
}

func (d *Db) UpdateGuildData(gid uuid.UUID, data models.DataGuild) error {

	query := `UPDATE my_compendium.guilds SET data = $1 WHERE gid = $2`
	_, err := d.db.Exec(query, data, gid)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}

	return nil
}
