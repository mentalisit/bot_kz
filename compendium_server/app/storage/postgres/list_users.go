package postgres

import (
	"strings"
)

func (d *Db) ListUserGetUserIdAndGuildId(token string) (userid string, guildid string, err error) {
	selectUser := "SELECT userid,guildid FROM hs_compendium.list_users WHERE token = $1"
	err = d.db.QueryRow(selectUser, token).Scan(&userid, &guildid)
	if err != nil {
		return "", "", err
	}
	return userid, guildid, nil
}
func (d *Db) ListUserGetByMatch(ttoken string) string {
	selectUser := "SELECT token FROM hs_compendium.list_users"
	results, err := d.db.Query(selectUser)
	defer results.Close()
	if err != nil {
		return ""
	}
	var tokens []string
	for results.Next() {
		var t string
		err = results.Scan(&t)
		if err != nil {
			return ""
		}
		tokens = append(tokens, t)
	}

	for _, token := range tokens {
		if strings.Contains(token, ttoken) {
			return token
		}
	}
	return ""
}
func (d *Db) ListUserUpdateToken(tokenOld, tokenNew string) error {
	upd := `update hs_compendium.list_users set token = $1 where token = $2`
	_, err := d.db.Exec(upd, tokenNew, tokenOld)
	if err != nil {
		return err
	}
	return nil
}
