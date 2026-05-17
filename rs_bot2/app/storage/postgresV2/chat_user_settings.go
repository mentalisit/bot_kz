package postgresV2

import (
	"database/sql"
	"encoding/json"
	"errors"
)

func (d *Db) SaveChatUserSettings(uuid, gid string, settings map[string]any) error {

	if gid == "" {
		gid = "00000000-0000-0000-0000-000000000000"
	}

	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO my_chat.user_settings (uuid, gid, settings, updated_at)
		VALUES ($1, $2, $3::jsonb, NOW())
		ON CONFLICT (uuid, gid) DO UPDATE SET
			settings = EXCLUDED.settings,
			updated_at = NOW()
	`
	_, err = d.db.Exec(query, uuid, gid, string(settingsJSON))
	return err
}

func (d *Db) GetChatUserSettings(uuid, gid string) (map[string]any, error) {

	if gid == "" {
		gid = "00000000-0000-0000-0000-000000000000"
	}

	var settingsRaw []byte
	err := d.db.QueryRow(`SELECT settings FROM my_chat.user_settings WHERE uuid = $1 AND gid = $2`, uuid, gid).Scan(&settingsRaw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return map[string]any{}, nil
		}
		d.log.ErrorErr(err)
		return map[string]any{}, nil
	}

	var settings map[string]any
	if err = json.Unmarshal(settingsRaw, &settings); err != nil {
		return map[string]any{}, nil
	}

	if settings == nil {
		settings = map[string]any{}
	}
	return settings, nil
}
