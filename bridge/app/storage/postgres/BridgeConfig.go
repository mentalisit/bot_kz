package postgres

import (
	"bridge/models"
	"encoding/json"
	"github.com/lib/pq"
)

func (d *Db) DBReadBridgeConfig2() []models.Bridge2Config {
	var cc []models.Bridge2Config
	rows, err := d.db.Query(`SELECT * FROM rs_bot.bridge_config`)
	if err != nil {
		d.log.ErrorErr(err)
		return cc
	}
	defer rows.Close()

	for rows.Next() {
		var config models.Bridge2Config
		var channel []byte
		if err = rows.Scan(&config.Id, &config.NameRelay, &config.HostRelay, pq.Array(&config.Role), &channel, pq.Array(&config.ForbiddenPrefixes)); err != nil {
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

func (d *Db) UpdateBridge2Chat(br models.Bridge2Config) {
	channel, err := json.Marshal(br.Channel)
	if err != nil {
		d.log.ErrorErr(err)
	}

	upd := `UPDATE rs_bot.bridge_config SET role = $1, channel = $2, forbidden_prefixes = $3 WHERE name_relay = $4`
	_, err = d.db.Exec(upd, pq.Array(br.Role), channel, pq.Array(br.ForbiddenPrefixes), br.NameRelay)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) InsertBridge2Chat(br models.Bridge2Config) {
	channel, err := json.Marshal(br.Channel)
	if err != nil {
		d.log.ErrorErr(err)
	}

	_, err = d.db.Exec(
		`INSERT INTO rs_bot.bridge_config (name_relay, host_relay, role, channel, forbidden_prefixes)
         VALUES ($1, $2, $3, $4, $5)`,
		br.NameRelay, br.HostRelay, pq.Array(br.Role), channel, pq.Array(br.ForbiddenPrefixes))
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) DeleteBridge2Chat(br models.Bridge2Config) {
	del := "DELETE FROM rs_bot.bridge_config WHERE name_relay = $1"
	_, err := d.db.Exec(del, br.NameRelay)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
