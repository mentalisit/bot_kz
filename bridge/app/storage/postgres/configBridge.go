package postgres

import (
	"bridge/models"
	"encoding/json"
	"github.com/lib/pq"
)

func (d *Db) DBReadBridgeConfig() []models.BridgeConfig {
	var cc []models.BridgeConfig
	rows, err := d.db.Query(`SELECT * FROM kzbot.bridge_config`)
	if err != nil {
		d.log.ErrorErr(err)
		return cc
	}
	defer rows.Close()

	for rows.Next() {
		var config models.BridgeConfig
		var channelDs, channelTg []byte
		if err = rows.Scan(&config.Id, &config.NameRelay, &config.HostRelay, pq.Array(&config.Role), &channelDs, &channelTg, pq.Array(&config.ForbiddenPrefixes)); err != nil {
			d.log.ErrorErr(err)
			return cc
		}

		if err = json.Unmarshal(channelDs, &config.ChannelDs); err != nil {
			d.log.ErrorErr(err)
		}

		if err = json.Unmarshal(channelTg, &config.ChannelTg); err != nil {
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

func (d *Db) UpdateBridgeChat(br models.BridgeConfig) {
	channelDs, err := json.Marshal(br.ChannelDs)
	if err != nil {
		d.log.ErrorErr(err)
	}

	channelTg, err := json.Marshal(br.ChannelTg)
	if err != nil {
		d.log.ErrorErr(err)
	}

	upd := `UPDATE kzbot.bridge_config SET role = $1, channel_ds = $2, channel_tg = $3, forbidden_prefixes = $4 WHERE name_relay = $5`
	_, err = d.db.Exec(upd, pq.Array(br.Role), channelDs, channelTg, pq.Array(br.ForbiddenPrefixes), br.NameRelay)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) InsertBridgeChat(br models.BridgeConfig) {
	channelDs, err := json.Marshal(br.ChannelDs)
	if err != nil {
		d.log.ErrorErr(err)
	}

	channelTg, err := json.Marshal(br.ChannelTg)
	if err != nil {
		d.log.ErrorErr(err)
	}

	_, err = d.db.Exec(
		`INSERT INTO kzbot.bridge_config (name_relay, host_relay, role, channel_ds, channel_tg, forbidden_prefixes)
         VALUES ($1, $2, $3, $4, $5, $6)`,
		br.NameRelay, br.HostRelay, pq.Array(br.Role), channelDs, channelTg, pq.Array(br.ForbiddenPrefixes))
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) DeleteBridgeChat(br models.BridgeConfig) {
	del := "DELETE FROM kzbot.bridge_config WHERE name_relay = $1"
	_, err := d.db.Exec(del, br.NameRelay)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

//func (d *Db) DBReadBridgeConfig() []models.BridgeConfig {
//	ctx, cancel := d.GetContext()
//	defer cancel()
//	var cc []models.BridgeConfig
//	rows, err := d.db.Query(ctx, `SELECT * FROM kzbot.bridge_config`)
//	if err != nil {
//		d.log.ErrorErr(err)
//		return cc
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var config models.BridgeConfig
//		var channelDs, channelTg []byte
//		if err = rows.Scan(&config.Id, &config.NameRelay, &config.HostRelay, pq.Array(&config.Role), &channelDs, &channelTg, pq.Array(&config.ForbiddenPrefixes)); err != nil {
//			d.log.ErrorErr(err)
//			return cc
//		}
//
//		if err = json.Unmarshal(channelDs, &config.ChannelDs); err != nil {
//			d.log.ErrorErr(err)
//		}
//
//		if err = json.Unmarshal(channelTg, &config.ChannelTg); err != nil {
//			d.log.ErrorErr(err)
//		}
//
//		cc = append(cc, config)
//	}
//	if err = rows.Err(); err != nil {
//		d.log.ErrorErr(err)
//		return cc
//	}
//	return cc
//}
//func (d *Db) UpdateBridgeChat(br models.BridgeConfig) {
//	ctx, cancel := d.GetContext()
//	defer cancel()
//	channelDs, err := json.Marshal(br.ChannelDs)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//
//	channelTg, err := json.Marshal(br.ChannelTg)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//	upd := `update kzbot.bridge_config set role = $1,channel_ds = $2,channel_tg = $3,forbidden_prefixes = $4 where name_relay = $5`
//	_, err = d.db.Exec(ctx, upd, pq.Array(br.Role), channelDs, channelTg, pq.Array(br.ForbiddenPrefixes), br.NameRelay)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//}
//func (d *Db) InsertBridgeChat(br models.BridgeConfig) {
//	ctx, cancel := d.GetContext()
//	defer cancel()
//	channelDs, err := json.Marshal(br.ChannelDs)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//
//	channelTg, err := json.Marshal(br.ChannelTg)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//
//	_, err = d.db.Exec(ctx,
//		`INSERT INTO kzbot.bridge_config (name_relay, host_relay, role, channel_ds, channel_tg, forbidden_prefixes)
//        VALUES ($1, $2, $3, $4, $5, $6)`,
//		br.NameRelay, br.HostRelay, pq.Array(br.Role), channelDs, channelTg, pq.Array(br.ForbiddenPrefixes))
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//}
//func (d *Db) DeleteBridgeChat(br models.BridgeConfig) {
//	ctx, cancel := d.GetContext()
//	defer cancel()
//	del := "delete from kzbot.bridge_config where name_relay = $1"
//	_, err := d.db.Exec(ctx, del, br.NameRelay)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//}
