package postgresOLD

import (
	"compendium/models"
	"context"
	"strconv"
)

func (d *Db) CreateRole(guildid string, RoleName string) {
	insert := `INSERT INTO compendium.roles(guildid, role) VALUES ($1,$2)`
	_, err := d.db.Exec(context.Background(), insert, guildid, RoleName)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) ExistRole(guildid string, RoleName string) bool {
	var id int
	selectid := "SELECT id FROM compendium.roles WHERE guildid = $1 AND role = $2"
	err := d.db.QueryRow(context.Background(), selectid, guildid, RoleName).Scan(&id)
	if err != nil {
		return false
	}
	return true
}
func (d *Db) DeleteRole(guildid string, RoleName string) {
	deleteRole := `DELETE FROM compendium.roles WHERE guildid = $1 AND role = $2`
	_, err := d.db.Exec(context.Background(), deleteRole, guildid, RoleName)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) ReadRoles(guildid string) []models.CorpRole {
	selectRoLes := `SELECT * FROM compendium.roles WHERE guildid = $1`
	var roles []models.CorpRole
	rows, err := d.db.Query(context.Background(), selectRoLes, guildid)
	if err != nil {
		d.log.ErrorErr(err)
	}
	defer rows.Close()
	for rows.Next() {
		var role models.CorpRole
		var id int
		err = rows.Scan(&id, &guildid, &role.Name)
		if err != nil {
			d.log.ErrorErr(err)
			return nil
		}
		role.Id = strconv.Itoa(id)
		roles = append(roles, role)
	}
	return roles
}

func (d *Db) RoleSubscribe(guildid, RoleName, userName, userid string) {
	insert := `INSERT INTO compendium.guildrole(guildid, role,username,userid) VALUES ($1,$2,$3,$4)`
	_, err := d.db.Exec(context.Background(), insert, guildid, RoleName, userName, userid)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) ExistSubscribe(guildid, RoleName, userid string) bool {
	var id int
	selectid := "SELECT id FROM compendium.guildrole WHERE guildid = $1 AND role = $2 AND userid = $3"
	err := d.db.QueryRow(context.Background(), selectid, guildid, RoleName, userid).Scan(&id)
	if err != nil {
		return false
	}
	return true
}
func (d *Db) DeleteSubscribe(guildid, RoleName, userid string) {
	deleteSubscribe := `DELETE FROM compendium.guildrole WHERE guildid = $1 AND role = $2 AND userid = $3`
	_, err := d.db.Exec(context.Background(), deleteSubscribe, guildid, RoleName, userid)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) ReadSubscribeUsers(guildid, RoleName string) []string {
	selectUsers := `SELECT username FROM compendium.guildrole WHERE guildid = $1 AND role = $2`
	var users []string
	rows, err := d.db.Query(context.Background(), selectUsers, guildid, RoleName)
	if err != nil {
		d.log.ErrorErr(err)
		return []string{}
	}
	defer rows.Close()
	for rows.Next() {
		var user string
		err = rows.Scan(&user)
		if err != nil {
			d.log.ErrorErr(err)
			return nil
		}
		users = append(users, user)
	}
	return users
}
