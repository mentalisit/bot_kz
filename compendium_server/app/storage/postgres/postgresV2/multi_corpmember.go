package postgresv2

import (
	"compendium_s/models"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Multi-corp member methods
func (d *Db) CorpMemberByUId(uid uuid.UUID) (*models.MultiAccountCorpMember, error) {
	var corpMember models.MultiAccountCorpMember

	query := `SELECT uid, guildIds, timeZona, zonaOffset, afkFor FROM my_compendium.corpMember WHERE uid = $1`
	err := d.db.QueryRow(query, uid).Scan(
		&corpMember.Uid, pq.Array(&corpMember.GuildIds), &corpMember.TimeZona, &corpMember.ZonaOffset, &corpMember.AfkFor,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &corpMember, nil
}

func (d *Db) CorpMembersReadMulti(gid *uuid.UUID) ([]models.CorpMember, error) {
	query := `
		SELECT ma.uuid, ma.nickname, ma.discord_id, ma.discord_username,
			   ma.telegram_id, ma.telegram_username, ma.whatsapp_id, ma.whatsapp_username,
			   ma.avatarurl, ma.alts, cm.timeZona, cm.zonaOffset, cm.afkFor
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

	var members []models.CorpMember
	for rows.Next() {
		var member models.CorpMember
		var discordID, discordUsername, telegramID, telegramUsername, whatsappID, whatsappUsername sql.NullString
		var memberUUID uuid.UUID
		var nickname, avatarURL string
		var alts []string

		// Initialize MultiAccount struct
		member.Multi = &models.MultiAccount{}

		err := rows.Scan(
			&memberUUID, &nickname, &discordID, &discordUsername,
			&telegramID, &telegramUsername, &whatsappID, &whatsappUsername,
			&avatarURL, pq.Array(alts), &member.TimeZone, &member.ZoneOffset, &member.AfkFor,
		)
		if err != nil {
			d.log.ErrorErr(fmt.Errorf("failed to scan corp member row: %w", err))
			continue
		}

		// Fill MultiAccount data
		member.Multi.UUID = memberUUID
		member.Multi.Nickname = nickname
		member.Multi.AvatarURL = avatarURL
		member.Multi.Alts = alts

		if discordID.Valid {
			member.Multi.DiscordID = discordID.String
		}
		if discordUsername.Valid {
			member.Multi.DiscordUsername = discordUsername.String
		}
		if telegramID.Valid {
			member.Multi.TelegramID = telegramID.String
		}
		if telegramUsername.Valid {
			member.Multi.TelegramUsername = telegramUsername.String
		}
		if whatsappID.Valid {
			member.Multi.WhatsappID = whatsappID.String
		}
		if whatsappUsername.Valid {
			member.Multi.WhatsappUsername = whatsappUsername.String
		}

		// Set CorpMember fields
		member.Name = nickname
		member.AvatarUrl = avatarURL

		if member.TimeZone != "" {
			t12, t24 := getTimeStrings(member.ZoneOffset)
			member.LocalTime = t12
			member.LocalTime24 = t24
		}

		// Get technologies for this user
		memberTech := d.TechnologiesGetUser(memberUUID)
		for _, tech := range memberTech {
			memberTemp := member
			memberTemp.Name = tech.Name
			memberTemp.Tech = tech.Tech
			memberTemp.UserId = memberUUID.String() + "/" + tech.Name
			if memberTemp.Multi.Nickname != memberTemp.Name {
				memberTemp.TypeAccount = "ALT"
			}
			members = append(members, memberTemp)
		}
	}

	// Check for any iteration errors
	if err := rows.Err(); err != nil {
		d.log.ErrorErr(fmt.Errorf("error iterating corp member rows: %w", err))
		return nil, err
	}

	return members, nil
}

func getTimeStrings(offset int) (string, string) {
	// Получаем текущее время в UTC
	now := time.Now().UTC()

	// Применяем смещение к текущему времени в UTC
	offsetDuration := time.Duration(offset) * time.Minute
	timeWithOffset := now.Add(offsetDuration)

	// Форматируем время в 12-часовом формате с AM/PM
	time12HourFormat := timeWithOffset.Format("03:04 PM")

	// Форматируем время в 24-часовом формате
	time24HourFormat := timeWithOffset.Format("15:04")

	return time12HourFormat, time24HourFormat
}

func (d *Db) CorpMemberInsert(corpMember models.MultiAccountCorpMember) error {
	// Check if multi_accounts record exists, create it if not
	checkQuery := `SELECT 1 FROM my_compendium.multi_accounts WHERE uuid = $1`
	var exists int
	err := d.db.QueryRow(checkQuery, corpMember.Uid).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			// Create minimal multi_accounts record
			createQuery := `
				INSERT INTO my_compendium.multi_accounts (uuid)
				VALUES ($1)
				ON CONFLICT (uuid) DO NOTHING`
			_, err = d.db.Exec(createQuery, corpMember.Uid)
			if err != nil {
				d.log.ErrorErr(fmt.Errorf("failed to create multi account: %w", err))
				return err
			}
		} else {
			d.log.ErrorErr(fmt.Errorf("failed to check multi account existence: %w", err))
			return err
		}
	}

	query := `
		INSERT INTO my_compendium.corpMember (uid, guildIds, timeZona, zonaOffset, afkFor)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (uid) DO UPDATE SET
			guildIds = EXCLUDED.guildIds,
			timeZona = EXCLUDED.timeZona,
			zonaOffset = EXCLUDED.zonaOffset,
			afkFor = EXCLUDED.afkFor`

	_, err = d.db.Exec(query, corpMember.Uid, pq.Array(corpMember.GuildIds), corpMember.TimeZona, corpMember.ZonaOffset, corpMember.AfkFor)
	if err != nil {
		d.log.ErrorErr(fmt.Errorf("failed to upsert corp member: %w", err))
		return err
	}

	return nil
}
