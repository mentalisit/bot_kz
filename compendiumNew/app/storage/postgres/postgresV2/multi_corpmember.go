package postgresv2

import (
	"compendium/models"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Multi-corp member methods
func (d *Db) CorpMemberByUId(uid uuid.UUID) (*models.MultiAccountCorpMember, error) {
	var corpMember models.MultiAccountCorpMember
	var guildIds []sql.NullString

	query := `SELECT uid, guildIds, timeZona, zonaOffset, afkFor FROM my_compendium.corpMember WHERE uid = $1`
	err := d.db.QueryRow(query, uid).Scan(
		&corpMember.Uid, pq.Array(&guildIds), &corpMember.TimeZona, &corpMember.ZonaOffset, &corpMember.AfkFor,
	)
	if err != nil {
		return nil, err
	}

	// Convert string array to UUID array
	corpMember.GuildIds = make([]uuid.UUID, 0, len(guildIds))
	for _, gidNull := range guildIds {
		if gidNull.Valid {
			if gid, err := uuid.Parse(gidNull.String); err == nil {
				corpMember.GuildIds = append(corpMember.GuildIds, gid)
			}
		}
	}

	return &corpMember, nil
}

func (d *Db) CorpMembersReadMulti(gid *uuid.UUID) ([]models.CorpMemberV2, error) {
	query := `
		SELECT ma.uuid, ma.nickname, ma.discord_id, ma.telegram_id, ma.whatsapp_id, ma.avatarurl,
			   cm.timeZona, cm.zonaOffset, cm.afkFor
		FROM my_compendium.corpMember cm
		JOIN my_compendium.multi_accounts ma ON cm.uid = ma.uuid
		WHERE $1 = ANY(cm.guildIds)
		ORDER BY ma.nickname
	`

	rows, err := d.db.Query(query, gid)
	if err != nil {
		d.log.ErrorErr(fmt.Errorf("failed to query corp members: %w", err))
		return nil, err
	}
	defer rows.Close()

	var members []models.CorpMemberV2
	for rows.Next() {
		var member models.CorpMemberV2
		var discordID, telegramID, whatsappID sql.NullString
		var memberUUID uuid.UUID
		var nickname string

		// Initialize MultiAccount struct
		member.Multi = &models.MultiAccount{}

		err := rows.Scan(
			&memberUUID, &nickname, &discordID, &telegramID, &whatsappID, &member.AvatarUrl,
			&member.TimeZone, &member.ZoneOffset, &member.AfkFor,
		)
		if err != nil {
			d.log.ErrorErr(fmt.Errorf("failed to scan corp member row: %w", err))
			continue
		}

		// Fill MultiAccount data
		member.Multi.UUID = memberUUID
		member.Multi.Nickname = nickname

		if discordID.Valid {
			member.Multi.DiscordID = discordID.String
		}
		if telegramID.Valid {
			member.Multi.TelegramID = telegramID.String
		}
		if whatsappID.Valid {
			member.Multi.WhatsappID = whatsappID.String
		}

		// Set CorpMember fields
		member.Name = nickname
		member.UserUUID = memberUUID.String()

		members = append(members, member)
	}

	// Check for any iteration errors
	if err := rows.Err(); err != nil {
		d.log.ErrorErr(fmt.Errorf("error iterating corp member rows: %w", err))
		return nil, err
	}

	return members, nil
}

func (d *Db) CorpMemberInsert(corpMember models.MultiAccountCorpMember) error {
	query := `
		INSERT INTO my_compendium.corpMember (uid, guildIds, timeZona, zonaOffset, afkFor)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := d.db.Exec(query, corpMember.Uid, pq.Array(corpMember.GuildIds), corpMember.TimeZona, corpMember.ZonaOffset, corpMember.AfkFor)
	if err != nil {
		d.log.ErrorErr(fmt.Errorf("failed to insert corp member: %w", err))
		return err
	}

	return nil
}

func (d *Db) CorpMemberUpdate(corpMember models.MultiAccountCorpMember) error {
	query := `
		UPDATE my_compendium.corpMember
		SET guildIds = $2, timeZona = $3, zonaOffset = $4, afkFor = $5
		WHERE uid = $1`

	_, err := d.db.Exec(query, corpMember.Uid, pq.Array(corpMember.GuildIds), corpMember.TimeZona, corpMember.ZonaOffset, corpMember.AfkFor)
	if err != nil {
		d.log.ErrorErr(fmt.Errorf("failed to update corp member: %w", err))
		return err
	}

	return nil
}
