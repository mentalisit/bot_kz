package postgres

import (
	"context"
	"fmt"
	"kz_bot/models"
)

func (d *Db) ReadConfigRs() []models.CorporationConfig {
	var tt []models.CorporationConfig
	results, err := d.db.Query(context.Background(), "SELECT * FROM kzbot.config")
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
	_, err := d.db.Exec(context.Background(), insert, c.CorpName, c.DsChannel, c.TgChannel, c.MesidDsHelp, c.MesidTgHelp, c.Country, c.DelMesComplite, c.Guildid, c.Forward)
	if err != nil {
		d.log.ErrorErr(err)
	}
	fmt.Printf("insert %+v\n", c)
}
func (d *Db) DeleteConfigRs(c models.CorporationConfig) {
	del := "delete from kzbot.config where corpname = $1"
	_, err := d.db.Exec(context.Background(), del, c.CorpName)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) UpdateConfigRs(c models.CorporationConfig) {
	upd := `update kzbot.config set mesiddshelp = $1,mesidtghelp = $2,country = $3 where corpname = $4`
	_, err := d.db.Exec(context.Background(), upd, c.MesidDsHelp, c.MesidTgHelp, c.Country, c.CorpName)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
