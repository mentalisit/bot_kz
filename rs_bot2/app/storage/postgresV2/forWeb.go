package postgresV2

import (
	"encoding/json"
	"rs/models"
)

func (d *Db) ReadBridgeConfigByChannelId(channelId string) (models.Bridge2Config, bool) {
	var config models.Bridge2Config

	// Ищем channelId в JSON поле channel используя LATERAL JOIN
	// Структура: {"tg": [{"channel_id": "...", ...}], "discord": [...]}
	query := `
		SELECT DISTINCT bc.id, bc.name_relay, bc.host_relay, bc.role, bc.channel, bc.forbidden_prefixes
		FROM rs_bot.bridge_config bc,
		LATERAL jsonb_each(bc.channel) AS platforms,
		LATERAL jsonb_array_elements(platforms.value) AS elem
		WHERE elem->>'channel_id' = $1`

	rows, err := d.db.Query(query, channelId)
	if err != nil {
		d.log.ErrorErr(err)
		return config, false
	}
	defer rows.Close()

	if rows.Next() {
		var channel []byte
		if err = rows.Scan(&config.Id, &config.NameRelay, &config.HostRelay, &config.Role, &channel, &config.ForbiddenPrefixes); err != nil {
			d.log.ErrorErr(err)
			return config, false
		}

		if err = json.Unmarshal(channel, &config.Channel); err != nil {
			d.log.ErrorErr(err)
		}

		return config, true
	}

	return config, false
}

func (d *Db) InsertBridgeConfig(br models.Bridge2Config) error {

	channel, err := json.Marshal(br.Channel)
	if err != nil {
		return err
	}

	insert := `INSERT INTO rs_bot.bridge_config (name_relay, host_relay, role, channel, forbidden_prefixes)
		VALUES ($1, $2, $3, $4, $5)`
	_, err = d.db.Exec(insert, br.NameRelay, br.HostRelay, br.Role, channel, br.ForbiddenPrefixes)
	if err != nil {
		d.log.ErrorErr(err)
	}
	return err
}

func (d *Db) UpdateBridgeConfig(br models.Bridge2Config) error {

	channel, err := json.Marshal(br.Channel)
	if err != nil {
		return err
	}

	upd := `UPDATE rs_bot.bridge_config 
		SET name_relay = $1, host_relay = $2, role = $3, channel = $4, forbidden_prefixes = $5 
		WHERE id = $6`
	_, err = d.db.Exec(upd, br.NameRelay, br.HostRelay, br.Role, channel, br.ForbiddenPrefixes, br.Id)
	if err != nil {
		d.log.ErrorErr(err)
	}
	return err
}

func (d *Db) UpdateBridgeConfigNameRelay(oldName string, newName string) error {

	upd := `UPDATE rs_bot.bridge_config 
		SET name_relay = $1	WHERE name_relay = $2`
	_, err := d.db.Exec(upd, newName, oldName)
	if err != nil {
		d.log.ErrorErr(err)
	}
	return nil
}

func (d *Db) ReadBridgeConfigByNameRelay(nameRelay string) (models.Bridge2Config, bool) {
	var config models.Bridge2Config

	query := `SELECT * FROM rs_bot.bridge_config WHERE name_relay = $1`
	rows, err := d.db.Query(query, nameRelay)
	if err != nil {
		d.log.ErrorErr(err)
		return config, false
	}
	defer rows.Close()

	if rows.Next() {
		var channel []byte
		if err = rows.Scan(&config.Id, &config.NameRelay, &config.HostRelay, &config.Role, &channel, &config.ForbiddenPrefixes); err != nil {
			d.log.ErrorErr(err)
			return config, false
		}

		if err = json.Unmarshal(channel, &config.Channel); err != nil {
			d.log.ErrorErr(err)
		}

		return config, true
	}

	return config, false
}

func (d *Db) DeleteBridge2Chat(br models.Bridge2Config) {
	del := "DELETE FROM rs_bot.bridge_config WHERE name_relay = $1"
	_, err := d.db.Exec(del, br.NameRelay)

	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) DBReadBBBridgeConfig2() []models.Bridge2Config {
	var cc []models.Bridge2Config
	rows, err := d.db.Query(`SELECT host_relay,channel FROM rs_bot.bridge_config`)
	if err != nil {
		d.log.ErrorErr(err)
		return cc
	}
	defer rows.Close()

	for rows.Next() {
		var config models.Bridge2Config
		var channel []byte
		if err = rows.Scan(&config.HostRelay, &channel); err != nil {
			d.log.ErrorErr(err)
			return cc
		}

		if err = json.Unmarshal(channel, &config.Channel); err != nil {
			d.log.ErrorErr(err)
		}

		cc = append(cc, config)
	}
	if err = rows.Err(); err != nil {
		d.log.ErrorErr(err)
		return cc
	}
	return cc
}
