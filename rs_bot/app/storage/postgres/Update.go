package postgres

func (d *Db) MesidTgUpdate(mesidtg int, lvlkz string, corpname string) error {
	ctx, cancel := d.GetContext()
	defer cancel()
	upd := `update kzbot.sborkz set tgmesid = $1 where lvlkz = $2 AND corpname = $3 AND active = 0`
	_, err := d.db.Exec(ctx, upd, mesidtg, lvlkz, corpname)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) MesidDsUpdate(mesidds, lvlkz, corpname string) error {
	ctx, cancel := d.GetContext()
	defer cancel()
	upd := `update kzbot.sborkz set dsmesid = $1 where lvlkz = $2 AND corpname = $3 AND active = 0`
	_, err := d.db.Exec(ctx, upd, mesidds, lvlkz, corpname)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) UpdateCompliteRS(lvlkz string, dsmesid string, tgmesid int, wamesid string, numberkz int, numberevent int, corpname string) error {
	ctx, cancel := d.GetContext()
	defer cancel()
	upd := `update kzbot.sborkz set active = 1,dsmesid = $1,tgmesid = $2,wamesid = $3,numberkz = $4,numberevent = $5 
				where lvlkz = $6 AND corpname = $7 AND active = 0`
	_, err := d.db.Exec(ctx, upd, dsmesid, tgmesid, wamesid, numberkz, numberevent, lvlkz, corpname)
	if err != nil {
		return err
	}

	updN := `update kzbot.numkz set number=number+1 where lvlkz = $1 AND corpname = $2`
	_, err = d.db.Exec(ctx, updN, lvlkz, corpname)
	if err != nil {
		return err
	}
	if numberevent > 0 {
		updE := `update kzbot.rsevent set number = number+1  where corpname = $1 AND activeevent = 1`
		_, err = d.db.Exec(ctx, updE, corpname)
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *Db) UpdateCompliteSolo(lvlkz string, dsmesid string, tgmesid int, numberkz int, numberevent int, corpname string) error {
	ctx, cancel := d.GetContext()
	defer cancel()
	upd := `update kzbot.sborkz set active = 1,numberkz = $1,numberevent = $2 
				where lvlkz = $3 AND corpname = $4 AND dsmesid = $5 AND tgmesid = $6 AND active = 0`
	_, err := d.db.Exec(ctx, upd, numberkz, numberevent, lvlkz, corpname, dsmesid, tgmesid)
	if err != nil {
		return err
	}

	updN := `update kzbot.numkz set number=number+1 where lvlkz = $1 AND corpname = $2`
	_, err = d.db.Exec(ctx, updN, lvlkz, corpname)
	if err != nil {
		return err
	}
	if numberevent > 0 {
		updE := `update kzbot.rsevent set number = number+1  where corpname = $1 AND activeevent = 1`
		_, err = d.db.Exec(ctx, updE, corpname)
		if err != nil {
			return err
		}
	}
	return nil
}
