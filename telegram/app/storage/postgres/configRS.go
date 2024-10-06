package postgres

import (
	"context"
	"telegram/models"
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
