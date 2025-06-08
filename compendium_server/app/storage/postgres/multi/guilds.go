package multi

import (
	"compendium_s/models"
	"github.com/google/uuid"
)

func (d *Db) GuildInsert(u models.MultiAccountGuild) error {
	ctx, cancel := d.getContext()
	defer cancel()
	insert := `INSERT INTO compendium.guilds(guildName,channels,avatarUrl) VALUES ($1,$2,$3)`
	_, err := d.db.Exec(ctx, insert, u.GuildName, u.Channels, u.AvatarUrl)
	if err != nil {
		return err
	}
	return nil
}

func (d *Db) GuildGet(guid *uuid.UUID) (*models.MultiAccountGuild, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var u models.MultiAccountGuild
	selectGuild := "SELECT * FROM compendium.guilds WHERE gid = $1"
	err := d.db.QueryRow(ctx, selectGuild, guid).Scan(&u.GId, &u.GuildName, &u.Channels, &u.AvatarUrl)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (d *Db) GuildGetById(guildId string) (*models.MultiAccountGuild, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	var u models.MultiAccountGuild
	selectGuild := "SELECT * FROM compendium.guilds WHERE $1 = ANY(channels)"
	err := d.db.QueryRow(ctx, selectGuild, guildId).Scan(&u.GId, &u.GuildName, &u.Channels, &u.AvatarUrl)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (d *Db) GuildUpdateAvatar(u models.MultiAccountGuild) error {
	ctx, cancel := d.getContext()
	defer cancel()
	upd := `update compendium.guilds set avatarUrl = $1 where gid = $2`
	_, err := d.db.Exec(ctx, upd, u.AvatarUrl, u.GId)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) GuildUpdateGuildName(u models.MultiAccountGuild) error {
	ctx, cancel := d.getContext()
	defer cancel()
	upd := `update compendium.guilds set guildName = $1 where gid = $2`
	_, err := d.db.Exec(ctx, upd, u.GuildName, u.GId)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) GuildUpdateChannels(u models.MultiAccountGuild) error {
	ctx, cancel := d.getContext()
	defer cancel()
	upd := `update compendium.guilds set channels = $1 where gid = $2`
	_, err := d.db.Exec(ctx, upd, u.Channels, u.GId)
	if err != nil {
		return err
	}
	return nil
}
