package postgres

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"rs/models"
)

func (d *Db) ReadConfigRs() []models.CorporationConfig {
	ctx, cancel := d.GetContext()
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
func (d *Db) InsertConfigRs(c models.CorporationConfig) {
	ctx, cancel := d.GetContext()
	defer cancel()
	insert := `INSERT INTO kzbot.config(corpname, dschannel, tgchannel, mesiddshelp, mesidtghelp, country, delmescomplite, guildid, forvard) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := d.db.Exec(ctx, insert, c.CorpName, c.DsChannel, c.TgChannel, c.MesidDsHelp, c.MesidTgHelp, c.Country, c.DelMesComplite, c.Guildid, c.Forward)
	if err != nil {
		d.log.ErrorErr(err)
	}
	fmt.Printf("insert %+v\n", c)
}
func (d *Db) DeleteConfigRs(c models.CorporationConfig) {
	ctx, cancel := d.GetContext()
	defer cancel()
	del := "delete from kzbot.config where corpname = $1"
	_, err := d.db.Exec(ctx, del, c.CorpName)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) UpdateConfigRs(c models.CorporationConfig) {
	ctx, cancel := d.GetContext()
	defer cancel()
	upd := `update kzbot.config set mesiddshelp = $1,mesidtghelp = $2,country = $3 where corpname = $4`
	_, err := d.db.Exec(ctx, upd, c.MesidDsHelp, c.MesidTgHelp, c.Country, c.CorpName)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) ReadConfigForDsChannel(dsChannel string) (conf models.CorporationConfig) {
	ctx, cancel := d.GetContext()
	defer cancel()
	sel := "SELECT * FROM kzbot.config WHERE dschannel = $1"
	results, err := d.db.Query(ctx, sel, dsChannel)
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
func (d *Db) ReadConfigForTgChannel(tgChannel string) (conf models.CorporationConfig) {
	ctx, cancel := d.GetContext()
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
	ctx, cancel := d.GetContext()
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
