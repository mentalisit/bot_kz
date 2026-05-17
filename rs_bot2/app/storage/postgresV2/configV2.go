package postgresV2

import (
	"encoding/json"
	"errors"
	"fmt"
	"rs/models"

	"github.com/jackc/pgx/v5"
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

func (d *Db) ReadConfigV2Uid(Uid string) *models.CorporationConfigV2 {

	sel := "SELECT * FROM rs_bot.config_rs WHERE uid = $1"
	results, err := d.db.Query(sel, Uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		d.log.ErrorErr(err)

	}
	defer results.Close()

	var t models.CorporationConfigV2

	for results.Next() {
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
	}
	return &t
}

func (d *Db) InsertConfigV2(c models.CorporationConfigV2) {

	marshalChannels, _ := json.Marshal(c.Channels)
	marshalBonuses, _ := json.Marshal(c.Bonuses)
	marshalHelpMessages, _ := json.Marshal(c.HelpMessage)

	insert := `INSERT INTO rs_bot.config_rs(uid, channels, bonuses, help_message) VALUES ($1,$2,$3,$4)`
	_, err := d.db.Exec(insert, c.Uid, marshalChannels, marshalBonuses, marshalHelpMessages)
	if err != nil {
		d.log.ErrorErr(err)
	}
	fmt.Printf("insert %+v\n", c)
}

func (d *Db) DeleteConfigV2(c models.CorporationConfigV2) {

	del := "delete from rs_bot.config_rs where uid = $1"
	_, err := d.db.Exec(del, c.Uid)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) UpdateConfigV2Channels(c models.CorporationConfigV2) {

	marshalChannels, _ := json.Marshal(c.Channels)
	upd := `update rs_bot.config_rs set channels = $1 where uid = $2`
	_, err := d.db.Exec(upd, marshalChannels, c.Uid)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) UpdateConfigV2Bonuses(c models.CorporationConfigV2) {

	marshalBonuses, _ := json.Marshal(c.Bonuses)
	upd := `update rs_bot.config_rs set bonuses = $1 where uid = $2`
	_, err := d.db.Exec(upd, marshalBonuses, c.Uid)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) UpdateConfigV2HelpMessage(c models.CorporationConfigV2) {

	marshalHelpMessages, _ := json.Marshal(c.HelpMessage)
	upd := `update rs_bot.config_rs set help_message = $1 where uid = $2`
	_, err := d.db.Exec(upd, marshalHelpMessages, c.Uid)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
