package postgres

import (
	"compendium/models"
	"context"
)

func (d *Db) IdentityInsert(ctx context.Context, i models.Identity) {
	insert := `INSERT INTO compendium.identity(id) VALUES ($1)`
	_, err := d.db.Exec(ctx, insert, i.Token)
	if err != nil {
		d.log.ErrorErr(err)
	}
	d.userInsert(ctx, i.Token, i.User)
	d.guildInsert(ctx, i.Token, i.Guild)
	cm := models.CorpMember{
		Name:         i.User.Username,
		UserId:       i.User.ID,
		ClientUserId: i.User.ID,
		Avatar:       i.User.Avatar,
		AvatarUrl:    i.User.AvatarURL,
	}
	d.corpMemberInsert(ctx, i.Guild.ID, cm)
}
func (d *Db) IdentityUpdateToken(ctx context.Context, oldToken, newToken string) {
	sqlUpd := `update compendium.identity set id = $1 where id = $2`
	_, err := d.db.Exec(ctx, sqlUpd, newToken, oldToken)
	if err != nil {
		d.log.ErrorErr(err)
	}
	d.userUpdate(ctx, oldToken, newToken)
	d.guildUpdate(ctx, oldToken, newToken)
}
func (d *Db) IdentityRead(ctx context.Context, token string) models.Identity {
	selec := "SELECT * FROM compendium.identity WHERE id = $1"
	results, err := d.db.Query(ctx, selec, token)
	if err != nil {
		d.log.ErrorErr(err)
	}
	var t models.Identity
	for results.Next() {
		err = results.Scan(&t.Token)
		if err != nil {
			d.log.ErrorErr(err)
		}
	}
	t.User = d.userRead(ctx, token)
	t.Guild = d.guildRead(ctx, token)

	return t
}
