package postgres

import (
	"compendium/models"
	"context"
	"strconv"
)

func (d *Db) GuildRoleCreate(guildid string, RoleName string) error {
	insert := `INSERT INTO hs_compendium.guildroles(guildid, role) VALUES ($1,$2)`
	_, err := d.db.Exec(context.Background(), insert, guildid, RoleName)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) GuildRoleExist(guildid string, RoleName string) bool {
	var count int
	selectid := "SELECT count(*) FROM hs_compendium.guildroles WHERE guildid = $1 AND role = $2"
	_ = d.db.QueryRow(context.Background(), selectid, guildid, RoleName).Scan(&count)
	if count > 0 {
		return true
	}
	return false
}
func (d *Db) GuildRoleDelete(guildid string, RoleName string) error {
	deleteRole := `DELETE FROM hs_compendium.guildroles WHERE guildid = $1 AND role = $2`
	_, err := d.db.Exec(context.Background(), deleteRole, guildid, RoleName)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) GuildRolesRead(guildid string) ([]models.CorpRole, error) {
	selectRoLes := `SELECT id,role FROM hs_compendium.guildroles WHERE guildid = $1`
	var roles []models.CorpRole
	rows, err := d.db.Query(context.Background(), selectRoLes, guildid)
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
	insert := `INSERT INTO hs_compendium.userroles(guildid, role,username,userid) VALUES ($1,$2,$3,$4)`
	_, err := d.db.Exec(context.Background(), insert, guildid, RoleName, userName, userid)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) GuildRolesExistSubscribe(guildid, RoleName, userid string) bool {
	var id int
	selectid := "SELECT id FROM hs_compendium.userroles WHERE guildid = $1 AND role = $2 AND userid = $3"
	err := d.db.QueryRow(context.Background(), selectid, guildid, RoleName, userid).Scan(&id)
	if err != nil {
		return false
	}
	return true
}
func (d *Db) GuildRolesDeleteSubscribe(guildid, RoleName, userid string) error {
	deleteSubscribe := `DELETE FROM hs_compendium.userroles WHERE guildid = $1 AND role = $2 AND userid = $3`
	_, err := d.db.Exec(context.Background(), deleteSubscribe, guildid, RoleName, userid)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) GuildRolesReadSubscribeUsers(guildid, RoleName string) ([]string, error) {
	selectUsers := `SELECT username FROM compendium.guildrole WHERE guildid = $1 AND role = $2`
	var users []string
	rows, err := d.db.Query(context.Background(), selectUsers, guildid, RoleName)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user string
		err = rows.Scan(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
