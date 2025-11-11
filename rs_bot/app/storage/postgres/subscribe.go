package postgres

import (
	"context"
	"fmt"
	"rs/models"
	"time"
)

func (d *Db) SubscribePing(s models.Subscribe) (subscribes []models.Subscribe) {
	ctx, cancel := d.getContext()
	defer cancel()

	sel := "SELECT subscribe.name,mention,userid FROM kzbot.subscribe WHERE lvlkz = $1 AND chatid = $2 AND tip = $3"
	if rows, err := d.db.Query(ctx, sel, s.Lvlkz, s.ChatId, s.Tip); err == nil {
		for rows.Next() {
			var subc models.Subscribe
			rows.Scan(&subc.Name, &subc.Mention, &subc.UserId)
			if s.UserId == subc.UserId {
				continue
			}
			subscribes = append(subscribes, subc)
		}
		rows.Close()
	}
	return subscribes
}

func (d *Db) CheckSubscribe(s models.Subscribe) int {
	ctx, cancel := d.getContext()
	defer cancel()
	var counts int
	sel := "SELECT  COUNT(*) as count FROM kzbot.subscribe WHERE name = $1 AND lvlkz = $2 AND chatid = $3 AND tip = $4"
	row := d.db.QueryRow(ctx, sel, s.Name, s.Lvlkz, s.ChatId, s.Tip)
	err := row.Scan(&counts)
	if err != nil {
		d.log.ErrorErr(err)
	}
	return counts
}
func (d *Db) Subscribe(s models.Subscribe) {
	if d.CheckSubscribe(s) > 0 {
		return
	}
	ctx, cancel := d.getContext()
	defer cancel()
	insertSubscribe := `INSERT INTO kzbot.subscribe (name, mention, lvlkz, tip, chatid, userid) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := d.db.Exec(ctx, insertSubscribe, s.Name, s.Mention, s.Lvlkz, s.Tip, s.ChatId, s.UserId)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) Unsubscribe(s models.Subscribe) {
	ctx, cancel := d.getContext()
	defer cancel()
	del := "delete from kzbot.subscribe where name = $1 AND lvlkz = $2 AND chatid = $3 AND tip = $4"
	_, err := d.db.Exec(ctx, del, s.Name, s.Lvlkz, s.ChatId, s.Tip)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) UpdateUserIdSubscribe() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sel := "SELECT name,userid FROM kzbot.sborkz WHERE tip = $1"
	results, err := d.db.Query(ctx, sel, "tg")
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	var tt []models.Sborkz
	for results.Next() {
		var t models.Sborkz
		err = results.Scan(&t.Name, &t.UserId)
		if err != nil {
			d.log.ErrorErr(err)
		}
		tt = append(tt, t)
	}
	fmt.Printf("tt len %d\n", len(tt))

	nt := make(map[string]string)
	for _, s := range tt {
		nt[s.Name] = s.UserId
	}
	fmt.Printf("nt clean len %d\n", len(nt))

	for name, userId := range nt {
		fmt.Printf("user %s ", name)
		var counts int
		sel = "SELECT  COUNT(*) as count FROM kzbot.subscribe WHERE name = $1 AND userid = $2"
		row := d.db.QueryRow(ctx, sel, name, userId)
		err = row.Scan(&counts)
		if err != nil {
			d.log.ErrorErr(err)
		}
		if counts == 0 {
			fmt.Println("update")
			upd := `update kzbot.subscribe set userid = $1 where name = $2`
			_, err := d.db.Exec(ctx, upd, userId, name)
			if err != nil {
				d.log.ErrorErr(err)
			}
		}
		fmt.Printf("\n")
	}

}
