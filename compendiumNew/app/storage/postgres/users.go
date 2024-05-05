package postgres

import (
	"compendium/models"
	"context"
	"github.com/lib/pq"
)

func (d *Db) UsersInsert(u models.User) error {
	count, err := d.UserGetCountByUserId(u.ID)
	if err != nil {
		return err
	}
	if count > 0 {
		err = d.UsersUpdate(u)
		if err != nil {
			return err
		}
	} else {
		insert := `INSERT INTO hs_compendium.users(userid, username, discriminator, avatar, avatarurl, alts) VALUES ($1,$2,$3,$4,$5,$6)`
		_, err = d.db.Exec(context.Background(), insert, u.ID, u.Username, u.Discriminator, u.Avatar, u.AvatarURL, pq.Array(u.Alts))
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *Db) UsersGetByUserId(userid string) (*models.User, error) {
	var u models.User
	var id int
	selectUser := "SELECT * FROM hs_compendium.users WHERE id = $1 "
	err := d.db.QueryRow(context.Background(), selectUser, userid).Scan(&id, &u.ID, &u.Username, &u.Discriminator, &u.Avatar, &u.AvatarURL, pq.Array(&u.Alts))
	if err != nil {
		return nil, err
	}
	return &u, nil
}
func (d *Db) UsersGetByUserName(username string) (*models.User, error) {
	var u models.User
	var id int
	selectUser := "SELECT * FROM hs_compendium.users WHERE username = $1 "
	err := d.db.QueryRow(context.Background(), selectUser, username).Scan(&id, &u.ID, &u.Username, &u.Discriminator, &u.Avatar, &u.AvatarURL, pq.Array(&u.Alts))
	if err != nil {
		return nil, err
	}
	return &u, nil
}
func (d *Db) UserGetCountByUserId(userid string) (int, error) {
	var count int
	sel := "SELECT count(*) as count FROM hs_compendium.users WHERE userid = $1"
	err := d.db.QueryRow(context.Background(), sel, userid).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (d *Db) UsersUpdate(u models.User) error {
	upd := `update hs_compendium.users set avatarurl = $1, alts = $2 where userid = $3 AND username = $4`
	_, err := d.db.Exec(context.Background(), upd, u.AvatarURL, pq.Array(u.Alts), u.ID, u.Username)
	if err != nil {
		return err
	}
	return nil
}
