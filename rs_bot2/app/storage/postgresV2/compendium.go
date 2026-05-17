package postgresV2

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"rs/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

func (d *Db) TechnologiesGet(uid uuid.UUID, name string) (*models.TechLevels, error) {

	var m []byte
	query := `SELECT tech FROM my_compendium.technologies WHERE uid = $1 AND username = $2`

	err := d.db.QueryRow(query, uid, name).Scan(&m)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var tech models.TechLevels
	if err := json.Unmarshal(m, &tech); err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}

	return &tech, nil
}

func (d *Db) TechnologyInsertUpdate(req models.CompendiumTechReq) (*models.TechLevels, error) {

	patch := map[string]models.TechLevel{
		req.Id: {
			Level: req.Level,
			Ts:    time.Now().UTC().UnixMilli(),
		},
	}

	patchRaw, err := json.Marshal(patch)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tech patch: %w", err)
	}

	query := `
        INSERT INTO my_compendium.technologies (uid, username, tech)
        VALUES ($1, $2, $3::jsonb)
        ON CONFLICT (uid, username) 
        DO UPDATE SET tech = technologies.tech || EXCLUDED.tech
        RETURNING tech;`

	var updatedTechRaw []byte
	err = d.db.QueryRow(query, req.Uuid, req.Name, patchRaw).Scan(&updatedTechRaw)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	var allTechs models.TechLevels
	if err := json.Unmarshal(updatedTechRaw, &allTechs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal updated tech: %w", err)
	}

	return &allTechs, nil
}

func (d *Db) CorpMemberByUId(ctx context.Context, uid uuid.UUID) (*models.MultiAccountCorpMember, error) {
	var m models.MultiAccountCorpMember

	query := `
        SELECT uid, guildids, timezona, zonaoffset, afkfor 
        FROM my_compendium.corpmember 
        WHERE uid = $1`

	err := d.db.QueryRow(query, uid).Scan(
		&m.Uid,
		&m.GuildIds,
		&m.TimeZona,
		&m.ZonaOffset,
		&m.AfkFor,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get corp member: %w", err)
	}

	return &m, nil
}

func (d *Db) GuildSave(ctx context.Context, g models.MultiAccountGuildV2) (*models.MultiAccountGuildV2, error) {
	// Если GId не передан, генерируем его сами
	if g.GId == uuid.Nil {
		g.GId = uuid.New()
	}

	query := `
        INSERT INTO my_compendium.guilds (gid, guildname, channels, avatarurl)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (gid) DO UPDATE SET
            guildname = EXCLUDED.guildname,
            channels = EXCLUDED.channels,
            avatarurl = EXCLUDED.avatarurl
        RETURNING gid, guildname, channels, avatarurl`

	var res models.MultiAccountGuildV2

	err := d.db.QueryRow(query,
		g.GId,
		g.GuildName,
		g.Channels,
		g.AvatarUrl,
	).Scan(
		&res.GId,
		&res.GuildName,
		&res.Channels,
		&res.AvatarUrl,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to save guild: %w", err)
	}

	return &res, nil
}

func (d *Db) CorpMembersReadMulti(gid *uuid.UUID) ([]models.CorpMember, error) {

	query := `
       SELECT ma.uuid, ma.nickname, ma.telegram_id, ma.telegram_username, ma.discord_id, ma.discord_username, ma.whatsapp_id, ma.whatsapp_username, ma.avatarurl, ma.alts, ma.created_at, ma.active_account, ma.data, cm.timezona, cm.zonaoffset, cm.afkfor
       FROM my_compendium.corpmember cm
       JOIN my_compendium.multi_accounts ma ON cm.uid = ma.uuid
       WHERE $1 = ANY(cm.guildids)
       ORDER BY ma.nickname`

	rows, err := d.db.Query(query, gid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.CorpMember

	for rows.Next() {
		var mAcc models.MultiAccount
		var timeZone, afkFor string
		var zoneOffset int

		var telegramID, discordID, whatsappID sql.NullString

		err := rows.Scan(
			&mAcc.UUID, &mAcc.Nickname, &telegramID, &mAcc.TelegramUsername,
			&discordID, &mAcc.DiscordUsername, &whatsappID, &mAcc.WhatsappUsername,
			&mAcc.AvatarURL, pq.Array(&mAcc.Alts), &mAcc.CreatedAt, &mAcc.ActiveAccount, &mAcc.Data, &timeZone, &zoneOffset, &afkFor,
		)
		if err != nil {
			return nil, err
		}

		if telegramID.Valid {
			mAcc.TelegramID = telegramID.String
		}
		if discordID.Valid {
			mAcc.DiscordID = discordID.String
		}
		if whatsappID.Valid {
			mAcc.WhatsappID = whatsappID.String
		}

		member := models.CorpMember{
			Multi:      &mAcc,
			TimeZone:   timeZone,
			ZoneOffset: zoneOffset,
			AfkFor:     afkFor,
			Name:       mAcc.Nickname,
			AvatarURL:  mAcc.AvatarURL,
		}

		if member.TimeZone != "" {
			member.LocalTime, member.LocalTime24 = getTimeStringsWithDST(member.TimeZone, member.ZoneOffset)
		}

		memberTech := d.TechnologiesGetUser(mAcc.UUID)
		for _, tech := range memberTech {
			mt := member
			mt.Name = tech.Name
			mt.Tech = tech.Tech
			mt.UserID = mAcc.UUID.String() + "/" + tech.Name

			if mt.Multi.Nickname != mt.Name {
				mt.TypeAccount = "ALT"
			}
			members = append(members, mt)
		}
	}

	return members, nil
}

func (d *Db) TechnologiesGetUser(uid uuid.UUID) []models.Technology {

	query := `SELECT uid, username, tech FROM my_compendium.technologies WHERE uid = $1`
	rows, err := d.db.Query(query, uid)
	if err != nil {
		d.log.ErrorErr(err)
		return nil
	}
	defer rows.Close()

	var techs []models.Technology
	for rows.Next() {
		var t models.Technology
		if err := rows.Scan(&t.Uid, &t.Name, &t.Tech); err != nil {
			d.log.ErrorErr(err)
			continue
		}
		techs = append(techs, t)
	}
	return techs
}

func (d *Db) GetChatsRoles(chatID int64) ([]models.CorpRole, error) {

	query := `SELECT id, name, chat_id FROM telegram.roles WHERE chat_id = $1`
	rows, err := d.db.Query(query, chatID)
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}
	defer rows.Close()

	var roles []models.CorpRole
	for rows.Next() {
		var r models.CorpRole
		if err := rows.Scan(&r.ID, &r.Name, &r.ChatID); err != nil {
			d.log.ErrorErr(err)
			continue
		}
		roles = append(roles, r)
	}

	return roles, nil
}

func getTimeStrings(offset int) (string, string) {
	now := time.Now().UTC()
	offsetDuration := time.Duration(offset) * time.Minute
	timeWithOffset := now.Add(offsetDuration)
	return timeWithOffset.Format("03:04 PM"), timeWithOffset.Format("15:04")
}

func getTimeStringsWithDST(timezone string, fallbackOffsetMinutes int) (string, string) {
	now := time.Now()
	if timezone != "" {
		loc, err := time.LoadLocation(timezone)
		if err == nil {
			timeInLocation := now.In(loc)
			return timeInLocation.Format("03:04 PM"), timeInLocation.Format("15:04")
		}
	}
	return getTimeStrings(fallbackOffsetMinutes)
}

func (d *Db) CorpMemberInsert(ctx context.Context, m models.MultiAccountCorpMember) error {
	query := `
       INSERT INTO my_compendium.corpmember (uid, guildids, timezona, zonaoffset, afkfor)
       VALUES ($1, $2, $3, $4, $5)
       ON CONFLICT (uid) DO UPDATE SET
          guildids = EXCLUDED.guildids,
          timezona = EXCLUDED.timezona,
          zonaoffset = EXCLUDED.zonaoffset,
          afkfor = EXCLUDED.afkfor`

	// Используем Exec для операций, не возвращающих данные
	_, err := d.db.Exec(query,
		m.Uid,
		m.GuildIds, // Передается как массив UUID
		m.TimeZona,
		m.ZonaOffset,
		m.AfkFor,
	)

	if err != nil {
		err = fmt.Errorf("failed to upsert corp member: %w", err)
		d.log.ErrorErr(err)
		return err
	}

	return nil
}

func (d *Db) CorpMemberUpdate(ctx context.Context, m models.MultiAccountCorpMember) error {
	query := `
       UPDATE my_compendium.corpmember
       SET guildids = $1, 
           timezona = $2, 
           zonaoffset = $3, 
           afkfor = $4
       WHERE uid = $5`

	// Выполняем запрос. Порядок аргументов должен строго соответствовать $1...$5
	result, err := d.db.Exec(query,
		m.GuildIds,   // $1
		m.TimeZona,   // $2
		m.ZonaOffset, // $3
		m.AfkFor,     // $4
		m.Uid,        // $5
	)

	if err != nil {
		err = fmt.Errorf("failed to update corp member: %w", err)
		d.log.ErrorErr(err)
		return err
	}

	// Опционально: проверка, была ли обновлена хоть одна строка
	r, _ := result.RowsAffected()
	if r == 0 {
		// Можно либо логировать, либо возвращать специфичную ошибку,
		// если юзер не найден
		d.log.Warn("no rows updated for uid: " + m.Uid.String())
	}

	return nil
}
