package postgres

import (
	"compendium_s/models"
	"context"
)

//func (d *Db) GuildInsert(u models.Guild) error {
//	count, err := d.GuildGetCountByGuildId(u.ID)
//	if err != nil {
//		return err
//	}
//	if count > 0 {
//		err = d.GuildUpdate(u)
//		if err != nil {
//			return err
//		}
//	} else {
//		insert := `INSERT INTO hs_compendium.guilds(url,guildid,name,icon,type) VALUES ($1,$2,$3,$4,$5)`
//		_, err = d.db.Exec(context.Background(), insert, u.URL, u.ID, u.Name, u.Icon, u.Type)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

func (d *Db) GuildGet(guildid string) (*models.Guild, error) {
	var u models.Guild
	var id int
	selectGuild := "SELECT * FROM hs_compendium.guilds WHERE guildid = $1 "
	err := d.db.QueryRow(context.Background(), selectGuild, guildid).Scan(&id, &u.URL, &u.ID, &u.Name, &u.Icon, &u.Type)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

//func (d *Db) GuildGetCountByGuildId(guildid string) (int, error) {
//	var count int
//	sel := "SELECT count(*) as count FROM hs_compendium.guilds WHERE guildid = $1"
//	err := d.db.QueryRow(context.Background(), sel, guildid).Scan(&count)
//	if err != nil {
//		return 0, err
//	}
//	return count, nil
//}
//
//func (d *Db) GuildUpdate(u models.Guild) error {
//	upd := `update hs_compendium.guilds set url = $1 where guildid = $2`
//	_, err := d.db.Exec(context.Background(), upd, u.URL, u.ID)
//	if err != nil {
//		return err
//	}
//	return nil
//}
