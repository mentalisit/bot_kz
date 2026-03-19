package postgresv2

import (
	"compendium/models"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// Multi-corp member methods
func (d *Db) CorpMemberByUId(uid uuid.UUID) (*models.MultiAccountCorpMember, error) {
	var m models.MultiAccountCorpMember
	query := `SELECT * FROM my_compendium.corpmember WHERE uid = $1`

	err := d.db.Get(&m, query, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}
func (d *Db) CorpMembersReadMulti(gid *uuid.UUID) ([]models.CorpMember, error) {
	// 1. Читаем данные из БД одной командой
	// Используем плоскую структуру для сканирования JOIN
	type tempRow struct {
		models.MultiAccount        // Включает все поля ma.*
		TimeZone            string `db:"timezona"`
		ZoneOffset          int    `db:"zonaoffset"`
		AfkFor              string `db:"afkfor"`
	}

	var rows []tempRow
	query := `
       SELECT ma.*, cm.timezona, cm.zonaoffset, cm.afkfor
       FROM my_compendium.corpmember cm
       JOIN my_compendium.multi_accounts ma ON cm.uid = ma.uuid
       WHERE $1 = ANY(cm.guildids)
       ORDER BY ma.nickname`

	if err := d.db.Select(&rows, query, gid); err != nil {
		return nil, err
	}

	var members []models.CorpMember

	// 2. Постобработка (расчет времени и технологий)
	for _, r := range rows {
		// Создаем базовый объект участника
		mAcc := r.MultiAccount // Копируем данные аккаунта
		member := models.CorpMember{
			MAcc:       &mAcc,
			TimeZone:   r.TimeZone,
			ZoneOffset: r.ZoneOffset,
			AfkFor:     r.AfkFor,
			Name:       r.Nickname,
			AvatarUrl:  r.AvatarURL,
		}

		members = append(members, member)
	}

	return members, nil
}

func (d *Db) CorpMemberInsert(m models.MultiAccountCorpMember) error {
	query := `
       INSERT INTO my_compendium.corpMember (uid, guildids, timezona, zonaoffset, afkfor)
       VALUES (:uid, :guildids, :timezona, :zonaoffset, :afkfor)
       ON CONFLICT (uid) DO UPDATE SET
          guildids = EXCLUDED.guildids,
          timezona = EXCLUDED.timezona,
          zonaoffset = EXCLUDED.zonaoffset,
          afkfor = EXCLUDED.afkfor`

	// Если GuildIds nil, NamedExec корректно обработает наш кастомный тип UUIDArray
	_, err := d.db.NamedExec(query, m)
	if err != nil {
		d.log.ErrorErr(fmt.Errorf("failed to upsert corp member: %w", err))
		return err
	}

	return nil
}

func (d *Db) CorpMemberUpdate(m models.MultiAccountCorpMember) error {
	query := `
       UPDATE my_compendium.corpmember
       SET guildids = :guildids, 
           timezona = :timezona, 
           zonaoffset = :zonaoffset, 
           afkfor = :afkfor
       WHERE uid = :uid`

	// NamedExec автоматически сопоставит поля структуры с именами после двоеточия
	_, err := d.db.NamedExec(query, m)
	if err != nil {
		d.log.ErrorErr(fmt.Errorf("failed to update corp member: %w", err))
		return err
	}

	return nil
}
