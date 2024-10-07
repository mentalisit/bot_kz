package postgres

import (
	"discord/models"
)

func (d *Db) ReadConfigRs() []models.CorporationConfig {
	ctx, cancelFunc := d.getContext()
	defer cancelFunc()
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
