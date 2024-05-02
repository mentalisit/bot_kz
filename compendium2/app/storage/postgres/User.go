package postgres

import (
	"compendium/models"
	"context"
	"github.com/lib/pq"
)

func (d *Db) userInsert(ctx context.Context, token string, u models.User) {
	insert := `INSERT INTO compendium.user(token, id, username, discriminator, avatar, avatarurl, alts) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := d.db.Exec(ctx, insert, token, u.ID, u.Username, u.Discriminator, u.Avatar, u.AvatarURL, pq.Array(u.Alts))
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) userRead(ctx context.Context, token string) models.User {
	selec := "SELECT * FROM compendium.user WHERE token = $1"
	results, err := d.db.Query(ctx, selec, token)
	if err != nil {
		d.log.ErrorErr(err)
	}
	var t models.User
	for results.Next() {
		err = results.Scan(&token, &t.ID, &t.Username, &t.Discriminator, &t.Avatar, &t.AvatarURL, pq.Array(&t.Alts))
		if err != nil {
			d.log.ErrorErr(err)
		}
	}
	return t
}
func (d *Db) userUpdate(ctx context.Context, oldToken, newToken string) {
	sqlUpd := `update compendium.user set token = $1 where token = $2`
	_, err := d.db.Exec(ctx, sqlUpd, newToken, oldToken)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) UserUpdateAlts(ctx context.Context, token string, alts []string) {
	sqlUpd := `UPDATE compendium.user SET alts = $2 WHERE token = $1`
	if _, err := d.db.Exec(ctx, sqlUpd, token, pq.Array(alts)); err != nil {
		d.log.ErrorErr(err)
	}
}

//func (d *Db) userFindById(ctx context.Context, userId string) models.User {
//	selec := "SELECT * FROM compendium.user WHERE token = $1"
//	results, err := d.db.Query(ctx, selec, token)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//	var t models.User
//	for results.Next() {
//		err = results.Scan(&token, &t.ID, &t.Username, &t.Discriminator, &t.Avatar, &t.AvatarURL)
//		if err != nil {
//			d.log.ErrorErr(err)
//		}
//	}
//	return t
//}
