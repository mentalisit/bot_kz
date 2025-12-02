package postgresv2

import (
	"compendium_s/models"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
)

// Multi-technologies methods
func (d *Db) TechnologiesGet(uid uuid.UUID, name string) (*models.TechLevels, error) {
	var techData []byte
	query := `SELECT tech FROM my_compendium.technologies WHERE uid = $1 AND username = $2`
	err := d.db.QueryRow(query, uid, name).Scan(&techData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No tech data found
		}
		return nil, err
	}

	var tech models.TechLevels
	err = json.Unmarshal(techData, &tech)
	if err != nil {
		return nil, err
	}

	return &tech, nil
}

func (d *Db) TechnologiesUpdate(uid uuid.UUID, username string, tech models.TechLevels) error {
	techData, err := json.Marshal(tech)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO my_compendium.technologies (uid, username, tech)
		VALUES ($1, $2, $3)
		ON CONFLICT (uid) DO UPDATE SET
			username = EXCLUDED.username,
			tech = EXCLUDED.tech
	`
	_, err = d.db.Exec(query, uid, username, techData)
	return err
}
