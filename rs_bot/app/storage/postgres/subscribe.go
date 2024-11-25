package postgres

import (
	"sort"
	"strings"
)

func (d *Db) SubscribePing(nameMention, lvlkz string, tipPing int, TgChannel string) string {
	ctx, cancel := d.GetContext()
	defer cancel()
	var name1 string
	var nameIds []string

	sel := "SELECT nameid FROM kzbot.subscribe WHERE lvlkz = $1 AND chatid = $2 AND tip = $3"
	if rows, err := d.db.Query(ctx, sel, lvlkz, TgChannel, tipPing); err == nil {
		for rows.Next() {
			rows.Scan(&name1)
			if nameMention == name1 {
				continue
			}
			nameIds = append(nameIds, name1)
		}
		rows.Close()
	}

	uniqueMap := make(map[string]struct{})
	var result []string

	for _, str := range nameIds {
		if _, exists := uniqueMap[str]; !exists {
			uniqueMap[str] = struct{}{}
			result = append(result, str)
		}
	}

	sort.Strings(result)

	var men string

	for _, sName := range result {
		men = men + sName + ", "
	}

	return strings.TrimSuffix(men, ", ")
}
func (d *Db) CheckSubscribe(name, lvlkz string, TgChannel string, tipPing int) int {
	ctx, cancel := d.GetContext()
	defer cancel()
	var counts int
	sel := "SELECT  COUNT(*) as count FROM kzbot.subscribe WHERE name = $1 AND lvlkz = $2 AND chatid = $3 AND tip = $4"
	row := d.db.QueryRow(ctx, sel, name, lvlkz, TgChannel, tipPing)
	err := row.Scan(&counts)
	if err != nil {
		d.log.ErrorErr(err)
	}
	return counts
}
func (d *Db) Subscribe(name, nameMention, lvlkz string, tipPing int, TgChannel string) {
	if d.CheckSubscribe(name, lvlkz, TgChannel, tipPing) > 0 {
		return
	}
	ctx, cancel := d.GetContext()
	defer cancel()
	insertSubscribe := `INSERT INTO kzbot.subscribe (name, nameid, lvlkz, tip, chatid, timestart, timeend) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := d.db.Exec(ctx, insertSubscribe, name, nameMention, lvlkz, tipPing, TgChannel, "0", "0")
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) Unsubscribe(name, lvlkz string, TgChannel string, tipPing int) {
	ctx, cancel := d.GetContext()
	defer cancel()
	del := "delete from kzbot.subscribe where name = $1 AND lvlkz = $2 AND chatid = $3 AND tip = $4"
	_, err := d.db.Exec(ctx, del, name, lvlkz, TgChannel, tipPing)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
