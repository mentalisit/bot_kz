package postgresOLD

import (
	"compendium/models"
	"context"
	"github.com/lib/pq"
)

func (d *Db) UserReadByUserIdByUsername(ctx context.Context, userid, username string) (*models.User, error) {
	var uu models.User
	var token string
	selectUser := "SELECT * FROM compendium.user WHERE id = $1 AND username = $2"
	err := d.db.QueryRow(ctx, selectUser, userid, username).Scan(&token, &uu.ID, &uu.Username, &uu.Discriminator, &uu.Avatar, &uu.AvatarURL, pq.Array(&uu.Alts))
	if err != nil {
		return nil, err
	}
	return &uu, nil
}
func (d *Db) userUpdateTokenAvatarUrlAlts(ctx context.Context, token, avatarUrl string, alts []string, userId, username string) error {
	upd := `update compendium.user set token = $1, avatarurl = $2, alts = $3 where id = $4 AND username = $5`
	_, err := d.db.Exec(ctx, upd, token, avatarUrl, pq.Array(alts), userId, username)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) userInsert(ctx context.Context, token string, u models.User) {
	uu, err := d.UserReadByUserIdByUsername(ctx, u.ID, u.Username)
	if err != nil {
		insert := `INSERT INTO compendium.user(token, id, username, discriminator, avatar, avatarurl, alts) VALUES ($1,$2,$3,$4,$5,$6,$7)`
		_, err = d.db.Exec(ctx, insert, token, u.ID, u.Username, u.Discriminator, u.Avatar, u.AvatarURL, pq.Array(u.Alts))
		if err != nil {
			d.log.ErrorErr(err)
		}
	} else {
		if uu.ID == u.ID && uu.Username == u.Username {
			alts := u.Alts
			if len(uu.Alts) > 0 {
				alts = uu.Alts
			}
			err = d.userUpdateTokenAvatarUrlAlts(ctx, token, u.AvatarURL, alts, u.ID, u.Username)
			if err != nil {
				d.log.ErrorErr(err)
			}
		}
	}
}
func (d *Db) userReadWhereToken(ctx context.Context, token string) models.User {
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
func (d *Db) UserUpdateAlts(ctx context.Context, username, userid string, alts []string) {
	sqlUpd := `UPDATE compendium.user SET alts = $1 WHERE username = $2 AND id = $3`
	if _, err := d.db.Exec(ctx, sqlUpd, pq.Array(alts), username, userid); err != nil {
		d.log.ErrorErr(err)
	}
}
