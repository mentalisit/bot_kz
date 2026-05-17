package postgres

import (
	"encoding/json"

	"github.com/mentalisit/restapi/models"
)

func (d *Db) ReadConfigV2() []models.CorporationConfigV2 {

	var tt []models.CorporationConfigV2
	results, err := d.db.Query("SELECT * FROM rs_bot.config_rs")
	if err != nil {
		d.log.ErrorErr(err)
		return tt
	}
	defer results.Close()
	for results.Next() {
		var t models.CorporationConfigV2
		var channelsJSON, bonusesJSON, helpMessageJSON []byte
		err = results.Scan(&t.Uid, &channelsJSON, &bonusesJSON, &helpMessageJSON)
		if err != nil {
			d.log.ErrorErr(err)
		}
		err = json.Unmarshal(channelsJSON, &t.Channels)
		if err != nil {
			d.log.ErrorErr(err)
		}
		err = json.Unmarshal(bonusesJSON, &t.Bonuses)
		if err != nil {
			d.log.ErrorErr(err)
		}
		err = json.Unmarshal(helpMessageJSON, &t.HelpMessage)
		if err != nil {
			d.log.ErrorErr(err)
		}
		tt = append(tt, t)
	}
	return tt
}
