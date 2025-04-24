package postgres

import (
	"errors"
	"time"
)

func (d *Db) MesidTgUpdate(mesidtg int, lvlkz string, corpname string) error {
	if mesidtg == 0 {
		return errors.New("mesId == null")
	}
	ctx, cancel := d.getContext()
	defer cancel()
	upd := `update kzbot.sborkz set tgmesid = $1 where lvlkz = $2 AND corpname = $3 AND active = 0`
	_, err := d.db.Exec(ctx, upd, mesidtg, lvlkz, corpname)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) MesidDsUpdate(mesidds, lvlkz, corpname string) error {
	if mesidds == "" {
		return errors.New("mesId == null")
	}
	ctx, cancel := d.getContext()
	defer cancel()
	upd := `update kzbot.sborkz set dsmesid = $1 where lvlkz = $2 AND corpname = $3 AND active = 0`
	_, err := d.db.Exec(ctx, upd, mesidds, lvlkz, corpname)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) UpdateCompliteRS(lvlkz string, dsmesid string, tgmesid int, wamesid string, numberkz int, numberevent int, corpname string) error {
	ctx, cancel := d.getContext()
	defer cancel()
	tm := time.Now().UTC()
	mdate := (tm.Format("2006-01-02"))
	mtime := (tm.Format("15:04"))
	upd := `update kzbot.sborkz set active = 1,dsmesid = $1,tgmesid = $2,wamesid = $3,numberkz = $4,numberevent = $5,date = $6,time = $7 
				where lvlkz = $8 AND corpname = $9 AND active = 0`
	_, err := d.db.Exec(ctx, upd, dsmesid, tgmesid, wamesid, numberkz, numberevent, mdate, mtime, lvlkz, corpname)
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
	ctx, cancel := d.getContext()
	defer cancel()
	tm := time.Now().UTC()
	mdate := (tm.Format("2006-01-02"))
	mtime := (tm.Format("15:04"))
	upd := `update kzbot.sborkz set active = 1,numberkz = $1,numberevent = $2,date = $3,time = $4 
				where lvlkz = $5 AND corpname = $6 AND dsmesid = $7 AND tgmesid = $8 AND active = 0`
	_, err := d.db.Exec(ctx, upd, numberkz, numberevent, mdate, mtime, lvlkz, corpname, dsmesid, tgmesid)
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
