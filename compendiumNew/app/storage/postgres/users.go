package postgres

import (
	"compendium/models"
	"github.com/jackc/pgx/v5"
)

func (d *Db) UsersInsert(u models.User) error {
	ctx, cancel := d.getContext()
	defer cancel()
	count, err := d.UserGetCountByUserId(u.ID)
	if err != nil {
		return err
	}
	if count > 0 {
		user, _ := d.UsersGetByUserId(u.ID)
		if len(user.Alts) > 0 {
			u.Alts = user.Alts
		}
		if user.GameName != "" {
			u.GameName = user.GameName
		}
		err = d.UsersUpdate(u)
		if err != nil {
			return err
		}
	} else {
		insert := `INSERT INTO hs_compendium.users(userid, username, discriminator, avatar, avatarurl, alts,gamename) VALUES ($1,$2,$3,$4,$5,$6,$7)`
		_, err = d.db.Exec(ctx, insert, u.ID, u.Username, u.Discriminator, u.Avatar, u.AvatarURL, u.Alts, u.GameName)
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *Db) UsersGetByUserId(userid string) (*models.User, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var u models.User
	var id int
	selectUser := "SELECT * FROM hs_compendium.users WHERE userid = $1 "
	err := d.db.QueryRow(ctx, selectUser, userid).Scan(&id, &u.ID, &u.Username, &u.Discriminator, &u.Avatar, &u.AvatarURL, &u.Alts, &u.GameName)
	if err != nil {
		return nil, pgx.ErrNoRows
	}
	return &u, nil
}
func (d *Db) UsersGetByUserName(username string) (*models.User, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var u models.User
	var id int
	selectUser := "SELECT * FROM hs_compendium.users WHERE username = $1 "
	err := d.db.QueryRow(ctx, selectUser, username).Scan(&id, &u.ID, &u.Username, &u.Discriminator, &u.Avatar, &u.AvatarURL, &u.Alts, &u.GameName)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
func (d *Db) UsersFindByGameName(gameName string) ([]models.User, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var users []models.User
	selectUsers := "SELECT * FROM hs_compendium.users WHERE gamename = $1 "
	results, err := d.db.Query(ctx, selectUsers, gameName)
	defer results.Close()
	if err != nil {
		return nil, err
	}
	for results.Next() {
		var t models.User
		var id int
		err = results.Scan(&id, &t.ID, &t.Username, &t.Discriminator, &t.Avatar, &t.AvatarURL, &t.Alts, &t.GameName)

		if err != nil {
			return nil, err
		}
		users = append(users, t)
	}

	return users, nil
}

func (d *Db) UserGetCountByUserId(userid string) (int, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var count int
	sel := "SELECT count(*) as count FROM hs_compendium.users WHERE userid = $1"
	err := d.db.QueryRow(ctx, sel, userid).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (d *Db) UsersUpdate(u models.User) error {
	ctx, cancel := d.getContext()
	defer cancel()
	upd := `update hs_compendium.users set avatarurl = $1, alts = $2, gamename = $3, username = $4 where userid = $5`
	_, err := d.db.Exec(ctx, upd, u.AvatarURL, u.Alts, u.GameName, u.Username, u.ID)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) UsersGetAll() ([]models.User, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var users []models.User
	selectUsers := "SELECT * FROM hs_compendium.users "
	results, err := d.db.Query(ctx, selectUsers)
	defer results.Close()
	if err != nil {
		return nil, err
	}
	for results.Next() {
		var t models.User
		var id int
		err = results.Scan(&id, &t.ID, &t.Username, &t.Discriminator, &t.Avatar, &t.AvatarURL, &t.Alts, &t.GameName)

		if err != nil {
			return nil, err
		}
		users = append(users, t)
	}
	return users, nil
}
