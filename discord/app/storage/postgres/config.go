package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/lib/pq"
	"github.com/mentalisit/restapi/models"
)

func (d *Db) ReadConfigRs() ([]models.CorporationConfig, error) {
	var tt []models.CorporationConfig
	results, err := d.db.Query("SELECT * FROM kzbot.config")
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

	sel := "SELECT * FROM kzbot.config WHERE dschannel = $1"
	results, err := d.db.Query(sel, dsChannel)
	defer results.Close()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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

	sel := "SELECT * FROM kzbot.config WHERE corpname = $1"
	results, err := d.db.Query(sel, corpName)
	defer results.Close()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
//	cancel := d.getContext()
//	defer cancel()
//	var cc []models.BridgeConfig
//	rows, err := d.db.Query(`SELECT * FROM kzbot.bridge_config`)
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
