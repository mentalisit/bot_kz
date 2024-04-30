package postgres

import (
	"compendium/models"
	"context"
)

func (d *Db) guildInsert(ctx context.Context, token string, u models.Guild) {
	insert := `INSERT INTO compendium.guild(token, url,id,name,icon) VALUES ($1,$2,$3,$4,$5)`
	_, err := d.db.Exec(ctx, insert, token, u.URL, u.ID, u.Name, u.Icon)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) guildRead(ctx context.Context, token string) models.Guild {
	selec := "SELECT * FROM compendium.guild WHERE token = $1"
	results, err := d.db.Query(ctx, selec, token)
	if err != nil {
		d.log.ErrorErr(err)
	}
	var t models.Guild
	for results.Next() {
		err = results.Scan(&token, &t.URL, &t.ID, &t.Name, &t.Icon)
		if err != nil {
			d.log.ErrorErr(err)
		}
	}
	return t
}
func (d *Db) guildUpdate(ctx context.Context, oldToken, newToken string) {
	sqlUpd := `update compendium.guild set token = $1 where token = $2`
	_, err := d.db.Exec(ctx, sqlUpd, newToken, oldToken)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
