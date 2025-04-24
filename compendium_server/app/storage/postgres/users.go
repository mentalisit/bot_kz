package postgres

import (
	"compendium_s/models"
)

//func (d *Db) UsersInsert(u models.User) error {
//	count, err := d.UserGetCountByUserId(u.ID)
//	if err != nil {
//		return err
//	}
//	if count > 0 {
//		user, _ := d.UsersGetByUserId(u.ID)
//		if len(user.Alts) > 0 {
//			u.Alts = user.Alts
//		}
//		if user.GameName != "" {
//			u.GameName = user.GameName
//		}
//		err = d.UsersUpdate(u)
//		if err != nil {
//			return err
//		}
//	} else {
//		insert := `INSERT INTO hs_compendium.users(userid, username, discriminator, avatar, avatarurl, alts,gamename) VALUES ($1,$2,$3,$4,$5,$6,$7)`
//		_, err = d.db.Exec(context.Background(), insert, u.ID, u.Username, u.Discriminator, u.Avatar, u.AvatarURL, pq.Array(u.Alts), u.GameName)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

func (d *Db) UsersGetByUserId(userid string) (*models.User, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var u models.User
	var id int
	selectUser := "SELECT * FROM hs_compendium.users WHERE userid = $1 "
	err := d.db.QueryRow(ctx, selectUser, userid).Scan(&id, &u.ID, &u.Username, &u.Discriminator, &u.Avatar, &u.AvatarURL, &u.Alts, &u.GameName)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

//func (d *Db) UsersGetByUserName(username string) (*models.User, error) {
//	var u models.User
//	var id int
//	selectUser := "SELECT * FROM hs_compendium.users WHERE username = $1 "
//	err := d.db.QueryRow(context.Background(), selectUser, username).Scan(&id, &u.ID, &u.Username, &u.Discriminator, &u.Avatar, &u.AvatarURL, pq.Array(&u.Alts), &u.GameName)
//	if err != nil {
//		return nil, err
//	}
//	return &u, nil
//}
//
//func (d *Db) UserGetCountByUserId(userid string) (int, error) {
//	var count int
//	sel := "SELECT count(*) as count FROM hs_compendium.users WHERE userid = $1"
//	err := d.db.QueryRow(context.Background(), sel, userid).Scan(&count)
//	if err != nil {
//		return 0, err
//	}
//	return count, nil
//}
//func (d *Db) UsersUpdate(u models.User) error {
//	upd := `update hs_compendium.users set avatarurl = $1, alts = $2, gamename = $3, username = $4 where userid = $5`
//	_, err := d.db.Exec(context.Background(), upd, u.AvatarURL, pq.Array(u.Alts), u.GameName, u.Username, u.ID)
//	if err != nil {
//		return err
//	}
//	return nil
//}
