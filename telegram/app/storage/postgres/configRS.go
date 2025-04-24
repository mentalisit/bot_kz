package postgres

import (
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"telegram/models"
)

func (d *Db) ReadConfigRs() []models.CorporationConfig {
	ctx, cancel := d.getContext()
	defer cancel()

	var tt []models.CorporationConfig
	results, err := d.db.Query(ctx, "SELECT * FROM kzbot.config")
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
		return tt
	}
	for results.Next() {
		var t models.CorporationConfig
		err = results.Scan(&t.Type, &t.CorpName, &t.DsChannel, &t.TgChannel, &t.MesidDsHelp, &t.MesidTgHelp,
			&t.Country, &t.DelMesComplite, &t.Guildid, &t.Forward)
		tt = append(tt, t)
	}
	return tt
}
func (d *Db) ReadConfigForTgChannel(tgChannel string) (conf models.CorporationConfig) {
	ctx, cancel := d.getContext()
	defer cancel()
	sel := "SELECT * FROM kzbot.config WHERE tgchannel = $1"
	results, err := d.db.Query(ctx, sel, tgChannel)
	defer results.Close()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return
		} else {
			d.log.ErrorErr(err)
		}
	}
	for results.Next() {
		err = results.Scan(&conf.Type, &conf.CorpName, &conf.DsChannel, &conf.TgChannel, &conf.MesidDsHelp, &conf.MesidTgHelp,
			&conf.Country, &conf.DelMesComplite, &conf.Guildid, &conf.Forward)
	}
	return conf
}
func (d *Db) ReadConfigForCorpName(corpName string) (conf models.CorporationConfig) {
	ctx, cancel := d.getContext()
	defer cancel()
	sel := "SELECT * FROM kzbot.config WHERE corpname = $1"
	results, err := d.db.Query(ctx, sel, corpName)
	defer results.Close()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return
		} else {
			d.log.ErrorErr(err)
		}
	}
	for results.Next() {
		err = results.Scan(&conf.Type, &conf.CorpName, &conf.DsChannel, &conf.TgChannel, &conf.MesidDsHelp, &conf.MesidTgHelp,
			&conf.Country, &conf.DelMesComplite, &conf.Guildid, &conf.Forward)
	}
	return conf
}
func (d *Db) DBReadBridgeConfig() []models.BridgeConfig {
	ctx, cancel := d.getContext()
	defer cancel()
	var cc []models.BridgeConfig
	rows, err := d.db.Query(ctx, `SELECT * FROM kzbot.bridge_config`)
	if err != nil {
		d.log.ErrorErr(err)
		return cc
	}
	defer rows.Close()

	for rows.Next() {
		var config models.BridgeConfig
		var channelDs, channelTg []byte
		if err = rows.Scan(&config.Id, &config.NameRelay, &config.HostRelay, &config.Role, &channelDs, &channelTg, &config.ForbiddenPrefixes); err != nil {
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
