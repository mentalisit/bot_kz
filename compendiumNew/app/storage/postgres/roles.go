package postgres

import (
	"compendium/models"
	"strconv"
)

func (d *Db) GuildRoleCreate(guildid string, RoleName string) error {
	ctx, cancel := d.GetContext()
	defer cancel()
	insert := `INSERT INTO hs_compendium.guildroles(guildid, role) VALUES ($1,$2)`
	_, err := d.db.Exec(ctx, insert, guildid, RoleName)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) GuildRoleExist(guildid string, RoleName string) bool {
	ctx, cancel := d.GetContext()
	defer cancel()
	var count int
	selectid := "SELECT count(*) FROM hs_compendium.guildroles WHERE guildid = $1 AND role = $2"
	_ = d.db.QueryRow(ctx, selectid, guildid, RoleName).Scan(&count)
	if count > 0 {
		return true
	}
	return false
}
func (d *Db) GuildRoleDelete(guildid string, RoleName string) error {
	ctx, cancel := d.GetContext()
	defer cancel()
	deleteRole := `DELETE FROM hs_compendium.guildroles WHERE guildid = $1 AND role = $2`
	_, err := d.db.Exec(ctx, deleteRole, guildid, RoleName)
	if err != nil {
		return err
	}
	return nil
}

func (d *Db) GuildRolesRead(guildid string) ([]models.CorpRole, error) {
	ctx, cancel := d.GetContext()
	defer cancel()
	selectRoLes := `SELECT id,role FROM hs_compendium.guildroles WHERE guildid = $1`
	var roles []models.CorpRole
	rows, err := d.db.Query(ctx, selectRoLes, guildid)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var role models.CorpRole
		var id int
		err = rows.Scan(&id, &role.Name)
		if err != nil {
			return nil, err
		}
		role.Id = strconv.Itoa(id)
		roles = append(roles, role)
	}
	return roles, nil
}

func (d *Db) GuildRolesSubscribe(guildid, RoleName, userName, userid string) error {
	ctx, cancel := d.GetContext()
	defer cancel()
	insert := `INSERT INTO hs_compendium.userroles(guildid, role,username,userid) VALUES ($1,$2,$3,$4)`
	_, err := d.db.Exec(ctx, insert, guildid, RoleName, userName, userid)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) GuildRolesExistSubscribe(guildid, RoleName, userid string) bool {
	ctx, cancel := d.GetContext()
	defer cancel()
	var id int
	selectid := "SELECT id FROM hs_compendium.userroles WHERE guildid = $1 AND role = $2 AND userid = $3"
	err := d.db.QueryRow(ctx, selectid, guildid, RoleName, userid).Scan(&id)
	if err != nil {
		return false
	}
	return true
}
func (d *Db) GuildRolesDeleteSubscribe(guildid, RoleName, userid string) error {
	ctx, cancel := d.GetContext()
	defer cancel()
	deleteSubscribe := `DELETE FROM hs_compendium.userroles WHERE guildid = $1 AND role = $2 AND userid = $3`
	_, err := d.db.Exec(ctx, deleteSubscribe, guildid, RoleName, userid)
	if err != nil {
		return err
	}
	return nil
}

//func (d *Db) GuildRolesReadSubscribeUsers(guildid, RoleName string) ([]string, error) {
//	ctx, cancel := d.GetContext()
//	defer cancel()
//	selectUsers := `SELECT username FROM compendium.guildrole WHERE guildid = $1 AND role = $2`
//	var users []string
//	rows, err := d.db.Query(ctx, selectUsers, guildid, RoleName)
//	defer rows.Close()
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//	for rows.Next() {
//		var user string
//		err = rows.Scan(&user)
//		if err != nil {
//			return nil, err
//		}
//		users = append(users, user)
//	}
//	return users, nil
//}
