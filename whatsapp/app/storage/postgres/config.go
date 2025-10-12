package postgres

import (
	"encoding/json"
	"errors"
	"whatsapp/models"

	"github.com/jackc/pgx/v5"
)

func (d *Db) ReadConfigRs() ([]models.CorporationConfig, error) {
	ctx, cancelFunc := d.getContext()
	defer cancelFunc()
	var tt []models.CorporationConfig
	results, err := d.db.Query(ctx, "SELECT * FROM kzbot.config")
	defer results.Close()
	if err != nil {
		return nil, err
	}
	for results.Next() {
		var t models.CorporationConfig
		err = results.Scan(&t.Type, &t.CorpName, &t.DsChannel, &t.TgChannel, &t.MesidDsHelp, &t.MesidTgHelp,
			&t.Country, &t.DelMesComplite, &t.Guildid, &t.Forward)
		tt = append(tt, t)
	}
	return tt, nil
}
func (d *Db) ReadConfigForDsChannel(dsChannel string) (conf models.CorporationConfig, err error) {
	ctx, cancel := d.getContext()
	defer cancel()
	sel := "SELECT * FROM kzbot.config WHERE dschannel = $1"
	results, err := d.db.Query(ctx, sel, dsChannel)
	defer results.Close()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return
		} else {
			return conf, nil
		}
	}
	for results.Next() {
		err = results.Scan(&conf.Type, &conf.CorpName, &conf.DsChannel, &conf.TgChannel, &conf.MesidDsHelp, &conf.MesidTgHelp,
			&conf.Country, &conf.DelMesComplite, &conf.Guildid, &conf.Forward)
	}
	return conf, nil
}
func (d *Db) ReadConfigForCorpName(corpName string) (conf models.CorporationConfig, err error) {
	ctx, cancel := d.getContext()
	defer cancel()
	sel := "SELECT * FROM kzbot.config WHERE corpname = $1"
	results, err := d.db.Query(ctx, sel, corpName)
	defer results.Close()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return
		} else {
			return conf, err
		}
	}
	for results.Next() {
		err = results.Scan(&conf.Type, &conf.CorpName, &conf.DsChannel, &conf.TgChannel, &conf.MesidDsHelp, &conf.MesidTgHelp,
			&conf.Country, &conf.DelMesComplite, &conf.Guildid, &conf.Forward)
	}
	return conf, nil
}

//func (d *Db) DBReadBridgeConfig() ([]models.BridgeConfig, error) {
//	ctx, cancel := d.getContext()
//	defer cancel()
//	var cc []models.BridgeConfig
//	rows, err := d.db.Query(ctx, `SELECT * FROM kzbot.bridge_config`)
//	if err != nil {
//		return cc, err
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var config models.BridgeConfig
//		var channelDs, channelTg []byte
//		if err = rows.Scan(&config.Id, &config.NameRelay, &config.HostRelay, &config.Role, &channelDs, &channelTg, &config.ForbiddenPrefixes); err != nil {
//			return cc, err
//		}
//
//		if err = json.Unmarshal(channelDs, &config.ChannelDs); err != nil {
//			return cc, err
//		}
//
//		if err = json.Unmarshal(channelTg, &config.ChannelTg); err != nil {
//			return cc, err
//		}
//
//		cc = append(cc, config)
//	}
//	if err = rows.Err(); err != nil {
//		return cc, err
//	}
//	return cc, nil
//}

func (d *Db) DBReadBridgeConfig() []models.Bridge2Config {
	var cc []models.Bridge2Config
	ctx, cancel := d.getContext()
	defer cancel()
	rows, err := d.db.Query(ctx, `SELECT * FROM rs_bot.bridge_config`)
	if err != nil {
		d.log.ErrorErr(err)
		return cc
	}
	defer rows.Close()

	for rows.Next() {
		var config models.Bridge2Config
		var channel []byte
		if err = rows.Scan(&config.Id, &config.NameRelay, &config.HostRelay, &config.Role, &channel, &config.ForbiddenPrefixes); err != nil {
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
