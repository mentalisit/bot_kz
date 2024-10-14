package postgres

import (
	"compendium/models"
)

func (d *Db) GuildInsert(u models.Guild) error {
	ctx, cancel := d.GetContext()
	defer cancel()
	count, err := d.GuildGetCountByGuildId(u.ID)
	if err != nil {
		return err
	}
	if count > 0 {
		err = d.GuildUpdate(u)
		if err != nil {
			return err
		}
	} else {
		insert := `INSERT INTO hs_compendium.guilds(url,guildid,name,icon,type) VALUES ($1,$2,$3,$4,$5)`
		_, err = d.db.Exec(ctx, insert, u.URL, u.ID, u.Name, u.Icon, u.Type)
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *Db) GuildGet(guildid string) (*models.Guild, error) {
	ctx, cancel := d.GetContext()
	defer cancel()
	var u models.Guild
	var id int
	selectGuild := "SELECT * FROM hs_compendium.guilds WHERE guildid = $1 "
	err := d.db.QueryRow(ctx, selectGuild, guildid).Scan(&id, &u.URL, &u.ID, &u.Name, &u.Icon, &u.Type)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
func (d *Db) GuildGetCountByGuildId(guildid string) (int, error) {
	ctx, cancel := d.GetContext()
	defer cancel()
	var count int
	sel := "SELECT count(*) as count FROM hs_compendium.guilds WHERE guildid = $1"
	err := d.db.QueryRow(ctx, sel, guildid).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (d *Db) GuildUpdate(u models.Guild) error {
	ctx, cancel := d.GetContext()
	defer cancel()
	upd := `update hs_compendium.guilds set url = $1 where guildid = $2`
	_, err := d.db.Exec(ctx, upd, u.URL, u.ID)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) GuildGetAll() ([]models.Guild, error) {
	ctx, cancel := d.GetContext()
	defer cancel()
	sel := "SELECT * FROM hs_compendium.guilds"
	results, err := d.db.Query(ctx, sel)
	defer results.Close()
	if err != nil {
		return nil, err
	}
	var mm []models.Guild
	for results.Next() {
		var t models.Guild
		var id int

		err = results.Scan(&id, &t.URL, &t.ID, &t.Name, &t.Icon, &t.Type)

		mm = append(mm, t)
	}
	return mm, nil
}
