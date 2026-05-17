package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"rs/models"
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
func (d *Db) InsertConfigRs(c models.CorporationConfig) {
	insert := `INSERT INTO kzbot.config(corpname, dschannel, tgchannel, mesiddshelp, mesidtghelp, country, delmescomplite, guildid, forvard) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := d.db.Exec(insert, c.CorpName, c.DsChannel, c.TgChannel, c.MesidDsHelp, c.MesidTgHelp, c.Country, c.DelMesComplite, c.Guildid, c.Forward)
	if err != nil {
		d.log.ErrorErr(err)
	}
	fmt.Printf("insert %+v\n", c)
}
func (d *Db) DeleteConfigRs(c models.CorporationConfig) {
	del := "delete from kzbot.config where corpname = $1"
	_, err := d.db.Exec(del, c.CorpName)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) UpdateConfigRs(c models.CorporationConfig) {
	upd := `update kzbot.config set mesiddshelp = $1,mesidtghelp = $2,country = $3 where corpname = $4`
	_, err := d.db.Exec(upd, c.MesidDsHelp, c.MesidTgHelp, c.Country, c.CorpName)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) ReadConfigForDsChannel(dsChannel string) (conf models.CorporationConfig) {
	sel := "SELECT * FROM kzbot.config WHERE dschannel = $1"
	results, err := d.db.Query(sel, dsChannel)
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
