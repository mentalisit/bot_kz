package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/mentalisit/restapi/models"
)

func (d *Db) ReadConfigRs() []models.CorporationConfig {

	var tt []models.CorporationConfig
	results, err := d.db.Query("SELECT * FROM kzbot.config")
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

	sel := "SELECT * FROM kzbot.config WHERE tgchannel = $1"
	results, err := d.db.Query(sel, tgChannel)
	defer results.Close()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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

	sel := "SELECT * FROM kzbot.config WHERE corpname = $1"
	results, err := d.db.Query(sel, corpName)
	defer results.Close()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
func (d *Db) DeleteConfigRs(c models.CorporationConfig) {

	del := "delete from kzbot.config where corpname = $1"
	_, err := d.db.Exec(del, c.CorpName)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
