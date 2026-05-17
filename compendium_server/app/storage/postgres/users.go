package postgres

import (
	"compendium_s/models"

	"github.com/lib/pq"
)

func (d *Db) UsersGetByUserId(userid string) (*models.User, error) {
	var u models.User
	var id int
	selectUser := "SELECT id, userid, username, avatarurl, alts, gamename FROM hs_compendium.users WHERE userid = $1 "
	err := d.db.QueryRow(selectUser, userid).Scan(&id, &u.ID, &u.Username, &u.AvatarURL, pq.Array(&u.Alts), &u.GameName)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
