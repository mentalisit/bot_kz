package postgres

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/mentalisit/restapi/models"
)

// CorpMemberByUId получает члена корпорации по UUID
func (d *Db) CorpMemberByUId(uid uuid.UUID) (*models.MultiAccountCorpMember, error) {
	query := `SELECT uid, guildids, timezona, zonaoffset, afkfor FROM my_compendium.corpmember WHERE uid = $1`

	var m models.MultiAccountCorpMember
	err := d.db.QueryRow(query, uid).Scan(&m.Uid, &m.GuildIds, &m.TimeZona, &m.ZonaOffset, &m.AfkFor)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		d.log.ErrorErr(err)
		return nil, err
	}

	return &m, nil
}

// TechnologiesGetAllCorpMember получает все технологии члена корпорации
func (d *Db) TechnologiesGetAll(uid uuid.UUID) ([]models.Tech, error) {

	query := `SELECT uid, username, tech FROM my_compendium.technologies WHERE uid = $1`

	rows, err := d.db.Query(query, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.Tech
	for rows.Next() {
		var tech models.Tech
		err := rows.Scan(&tech.Uid, &tech.Username, &tech.Tech)
		if err != nil {
			return nil, err
		}
		results = append(results, tech)
	}

	return results, nil
}
