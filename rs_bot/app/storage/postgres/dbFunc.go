package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"rs/models"
	"rs/pkg/utils"
	"strconv"
	"time"
)

func (d *Db) ReadAll(lvlkz, CorpName string) (users models.Users) {
	ctx, cancel := d.GetContext()
	defer cancel()

	u := models.Users{
		User1: models.Sborkz{},
	}
	user := 1
	sel := "SELECT * FROM kzbot.sborkz WHERE lvlkz = $1 AND corpname = $2 AND active = 0"
	results, err := d.db.Query(ctx, sel, lvlkz, CorpName)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	for results.Next() {
		var t models.Sborkz
		err = results.Scan(&t.Id, &t.Corpname, &t.Name, &t.Mention, &t.Tip, &t.Dsmesid,
			&t.Tgmesid, &t.Wamesid, &t.Time, &t.Date, &t.Lvlkz, &t.Numkzn, &t.Numberkz,
			&t.Numberevent, &t.Eventpoints, &t.Active, &t.Timedown, &t.UserId)
		if user == 1 {
			u.User1 = t
		} else if user == 2 {
			u.User2 = &t
		} else if user == 3 {
			u.User3 = &t
		} else if user == 4 {
			u.User4 = &t
		}
		user = user + 1
	}

	return u
}
func (d *Db) ReadAllActive() (sb []models.Sborkz) {
	//sel := "SELECT * FROM kzbot.sborkz" //" WHERE active > 0"
	sel := "SELECT id, corpname, name, COALESCE(date, '') AS date, COALESCE(lvlkz, '') AS lvlkz, " +
		"numberevent, eventpoints FROM kzbot.sborkz WHERE active > 0"

	results, err := d.db.Query(context.Background(), sel)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	for results.Next() {
		var t models.Sborkz
		err = results.Scan(&t.Id, &t.Corpname, &t.Name, &t.Date, &t.Lvlkz,
			&t.Numberevent, &t.Eventpoints)
		if err != nil {
			fmt.Println(err)
		}
		sb = append(sb, t)
	}

	return sb
}
func (d *Db) DeleteSborkzId(id int) {
	ctx, cancel := d.GetContext()
	defer cancel()

	del := "delete from kzbot.sborkz where id = $1"
	_, err := d.db.Exec(ctx, del, id)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) UpdateSborkz(active string, id int) {
	ctx, cancel := d.GetContext()
	defer cancel()
	upd := `update kzbot.sborkz set active = $1 where id = $2`
	_, err := d.db.Exec(ctx, upd, active, id)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) UpdateSborkzPoints(active string, id int, points int) {
	ctx, cancel := d.GetContext()
	defer cancel()
	upd := `update kzbot.sborkz set active = $1, eventpoints = $2 where id = $3`
	_, err := d.db.Exec(ctx, upd, active, points, id)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) InsertQueue(dsmesid, wamesid, CorpName, name, userid, nameMention, tip, lvlkz, timekz string, tgmesid, numkzN int) {
	ctx, cancel := d.GetContext()
	defer cancel()
	numevent := 0 // d.NumActiveEvent(CorpName)
	tm := time.Now()
	mdate := (tm.Format("2006-01-02"))
	mtime := (tm.Format("15:04"))

	timekzz, _ := strconv.Atoi(timekz)
	if timekzz == 0 {
		timekzz = 1
	}

	insertSborkztg1 := `INSERT INTO kzbot.sborkz(corpname,name,userid,mention,tip,dsmesid,tgmesid,wamesid,time,date,lvlkz,
                   numkzn,numberkz,numberevent,eventpoints,active,timedown) 
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err := d.db.Exec(ctx, insertSborkztg1, CorpName, name, userid, nameMention, tip, dsmesid, tgmesid,
		wamesid, mtime, mdate, lvlkz, numkzN, 0, numevent, 0, 0, timekzz)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) InsertQueueSolo(dsmesid, wamesid, CorpName, name, userid, nameMention, tip, lvlkz string, tgmesid, numevent, numberkz, numkzN, points int) {
	ctx, cancel := d.GetContext()
	defer cancel()
	tm := time.Now()
	mdate := (tm.Format("2006-01-02"))
	mtime := (tm.Format("15:04"))

	timekzz := 1

	insertSborkztg1 := `INSERT INTO kzbot.sborkz(corpname,name,userid,mention,tip,dsmesid,tgmesid,wamesid,time,date,lvlkz,
                   numkzn,numberkz,numberevent,eventpoints,active,timedown) 
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err := d.db.Exec(ctx, insertSborkztg1, CorpName, name, userid, nameMention, tip, dsmesid, tgmesid,
		wamesid, mtime, mdate, lvlkz, numkzN, numberkz, numevent, points, 1, timekzz)
	if err != nil {
		d.log.ErrorErr(err)
	}

	updN := `update kzbot.numkz set number=number+1 where lvlkz = $1 AND corpname = $2`
	_, err = d.db.Exec(ctx, updN, lvlkz, CorpName)
	if err != nil {
		d.log.ErrorErr(err)
	}
	if numevent > 0 {
		updE := `update kzbot.rsevent set number = number+1  where corpname = $1 AND activeevent = 1`
		_, err = d.db.Exec(ctx, updE, CorpName)
		if err != nil {
			d.log.ErrorErr(err)
		}
	}
}

func (d *Db) ElseTrue(userid string) []models.Sborkz {
	ctx, cancel := d.GetContext()
	defer cancel()

	sel := "SELECT * FROM kzbot.sborkz WHERE userid = $1 AND active = 0"
	results, err := d.db.Query(ctx, sel, userid)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	var tt []models.Sborkz
	var t models.Sborkz
	for results.Next() {
		err = results.Scan(&t.Id, &t.Corpname, &t.Name, &t.Mention, &t.Tip, &t.Dsmesid, &t.Tgmesid, &t.Wamesid, &t.Time, &t.Date, &t.Lvlkz, &t.Numkzn, &t.Numberkz, &t.Numberevent, &t.Eventpoints, &t.Active, &t.Timedown, &t.UserId)
		tt = append(tt, t)
	}

	return tt
}
func (d *Db) DeleteQueue(userid, lvlkz, CorpName string) {
	ctx, cancel := d.GetContext()
	defer cancel()

	del := "delete from kzbot.sborkz where userid = $1 AND lvlkz = $2 AND corpname = $3 AND active = 0"
	_, err := d.db.Exec(ctx, del, userid, lvlkz, CorpName)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) ReadMesIdDS(mesid string) (string, error) {
	ctx, cancel := d.GetContext()
	defer cancel()

	sel := "SELECT lvlkz FROM kzbot.sborkz WHERE dsmesid = $1 AND active = 0"
	results, err := d.db.Query(ctx, sel, mesid)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	a := []string{}
	var dsmesid string
	for results.Next() {
		var t models.Sborkz
		err = results.Scan(&t.Lvlkz)
		a = append(a, t.Lvlkz)
	}
	a = utils.RemoveDuplicates(a)

	if len(a) > 0 {
		dsmesid = a[0]
		return dsmesid, nil
	} else {
		return "", err
	}
}

func (d *Db) P30Pl(lvlkz, CorpName, userid string) int {
	ctx, cancel := d.GetContext()
	defer cancel()

	var timedown int
	sel := "SELECT timedown FROM kzbot.sborkz WHERE lvlkz = $1 AND corpname = $2 AND active = 0 AND userid = $3"
	results, err := d.db.Query(ctx, sel, lvlkz, CorpName, userid)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	for results.Next() {
		err = results.Scan(&timedown)
	}

	return timedown
}
func (d *Db) UpdateTimedown(lvlkz, CorpName, userid string) {
	ctx, cancel := d.GetContext()
	defer cancel()

	upd := `update kzbot.sborkz set timedown = timedown+30 where lvlkz = $1 AND corpname = $2 AND userid = $3`
	_, err := d.db.Exec(ctx, upd, lvlkz, CorpName, userid)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) Queue(corpname string) []string {
	ctx, cancel := d.GetContext()
	defer cancel()

	sel := "SELECT lvlkz FROM kzbot.sborkz WHERE corpname = $1 AND active = 0"
	results, err := d.db.Query(ctx, sel, corpname)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	var lvl []string
	for results.Next() {
		var t models.Sborkz
		err = results.Scan(&t.Lvlkz)

		lvl = append(lvl, t.Lvlkz)

	}

	return lvl
}

func (d *Db) OneMinutsTimer() []string {
	ctx, cancel := d.GetContext()
	defer cancel()
	var count int //количество активных игроков
	sel := "SELECT  COUNT(*) as count FROM kzbot.sborkz WHERE active = 0"
	row := d.db.QueryRow(ctx, sel)
	err := row.Scan(&count)
	if err != nil {
		d.log.ErrorErr(err)
	}
	var CorpActive0 []string
	if count > 0 {
		a := []string{}
		aa := []string{}
		selC := "SELECT corpname FROM kzbot.sborkz WHERE active = 0"
		results, err1 := d.db.Query(ctx, selC)
		defer results.Close()
		if err1 != nil {
			d.log.ErrorErr(err)
		}
		var corpname string // ищим корпорации
		for results.Next() {
			err = results.Scan(&corpname)
			a = append(a, corpname)
		}
		a = utils.RemoveDuplicates(a)

		for _, corp := range a {
			skip := false
			for _, u := range aa {
				if corp == u {
					skip = true
					break
				}
			}
			if !skip {
				CorpActive0 = append(CorpActive0, corp)
			}
		}
	}

	return CorpActive0
}
func (d *Db) MessageUpdateMin(corpname string) ([]string, []int) {
	ctx, cancel := d.GetContext()
	defer cancel()

	var countCorp int
	var ds []string
	var tg []int
	sel := "SELECT  COUNT(*) as count FROM kzbot.sborkz WHERE corpname = $1 AND active = 0"
	row := d.db.QueryRow(ctx, sel, corpname)
	err := row.Scan(&countCorp)
	if err != nil {
		d.log.ErrorErr(err)
	}
	if countCorp > 0 {
		selS := "SELECT dsmesid,tgmesid FROM kzbot.sborkz WHERE corpname = $1 AND active = 0"
		results, err1 := d.db.Query(ctx, selS, corpname)
		defer results.Close()
		if err1 != nil {
			d.log.Error(err1.Error())
		}
		for results.Next() {

			var dsmesid string
			var tgmesid int

			err = results.Scan(&dsmesid, &tgmesid)

			if dsmesid != "" {
				ds = append(ds, dsmesid)
			}
			if tgmesid != 0 {
				tg = append(tg, tgmesid)
			}

		}
	}
	ds = utils.RemoveDuplicates(ds)
	tg = utils.RemoveDuplicates(tg)

	return ds, tg
}
func (d *Db) MessageUpdateDS(dsmesid string, config models.CorporationConfig) models.InMessage {
	ctx, cancel := d.GetContext()
	defer cancel()

	sel := "SELECT * FROM kzbot.sborkz WHERE dsmesid = $1 AND active = 0"
	results, err := d.db.Query(ctx, sel, dsmesid)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	var t models.Sborkz
	for results.Next() {
		err = results.Scan(&t.Id, &t.Corpname, &t.Name, &t.Mention, &t.Tip, &t.Dsmesid, &t.Tgmesid, &t.Wamesid, &t.Time, &t.Date, &t.Lvlkz, &t.Numkzn, &t.Numberkz, &t.Numberevent, &t.Eventpoints, &t.Active, &t.Timedown, &t.UserId)
	}
	in := models.InMessage{
		Tip:         "ds",
		Username:    t.Name,
		UserId:      t.UserId,
		NameMention: t.Mention,
		Lvlkz:       t.Lvlkz,
		Timekz:      strconv.Itoa(t.Timedown),
		Ds: struct {
			Mesid   string
			Guildid string
			Avatar  string
		}{
			Mesid:   t.Dsmesid,
			Guildid: config.Guildid,
		},
		Config: config,
		//Option: models.Option{
		//	Edit:   true,
		//	Update: true},
		Opt: []string{models.OptionMessageUpdateDS, models.OptionEdit},
	}
	return in

}
func (d *Db) MessageUpdateTG(tgmesid int, config models.CorporationConfig) models.InMessage {
	ctx, cancel := d.GetContext()
	defer cancel()

	sel := "SELECT * FROM kzbot.sborkz WHERE tgmesid = $1 AND active = 0"
	results, err := d.db.Query(ctx, sel, tgmesid)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	var t models.Sborkz
	for results.Next() {
		err = results.Scan(&t.Id, &t.Corpname, &t.Name, &t.Mention, &t.Tip, &t.Dsmesid, &t.Tgmesid, &t.Wamesid, &t.Time, &t.Date, &t.Lvlkz, &t.Numkzn, &t.Numberkz, &t.Numberevent, &t.Eventpoints, &t.Active, &t.Timedown, &t.UserId)
	}
	in := models.InMessage{
		Tip:         "tg",
		Username:    t.Name,
		NameMention: t.Mention,
		Lvlkz:       t.Lvlkz,
		Timekz:      strconv.Itoa(t.Timedown),
		Tg: struct {
			Mesid int
			//Nameid int64
		}{
			Mesid: t.Tgmesid,
			//Nameid: 0
		},
		Config: config,
		//Option: models.Option{
		//	Edit:   true,
		//	Update: true,},
		Opt: []string{models.OptionMessageUpdateTG, models.OptionEdit},
	}
	return in
}
func (d *Db) NumberQueueLvl(lvlkz, CorpName string) (int, error) {
	ctx, cancel := d.GetContext()
	defer cancel()

	var number int
	sel := "SELECT  number FROM kzbot.numkz WHERE lvlkz = $1 AND corpname = $2"
	row := d.db.QueryRow(ctx, sel, lvlkz, CorpName)
	err := row.Scan(&number)
	if err != nil {
		if err == pgx.ErrNoRows {
			number = 0
			insertSmt := "INSERT INTO kzbot.numkz(lvlkz, number,corpname) VALUES ($1,$2,$3)"
			_, err = d.db.Exec(ctx, insertSmt, lvlkz, number, CorpName)
			if err != nil {
				d.log.ErrorErr(err)
			}
			return number + 1, nil
		} else {
			d.log.ErrorErr(err)
			return 0, err
		}
	}

	return number + 1, nil
}
