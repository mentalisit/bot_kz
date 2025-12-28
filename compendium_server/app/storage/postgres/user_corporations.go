package postgres

import (
	"compendium_s/models"
)

// UserCorporationsGet возвращает список корпораций, в которых состоит пользователь
func (d *Db) UserCorporationsGet(identity *models.Identity) ([]models.Guild, error) {
	ctx, cancel := d.getContext()
	defer cancel()

	var corporations []models.Guild

	// Обычный аккаунт - получаем все гильдии, где пользователь является членом
	sel := "SELECT DISTINCT guildid FROM hs_compendium.corpmember WHERE userid = $1"
	results, err := d.db.Query(ctx, sel, identity.User.ID)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	var guildIds []string
	for results.Next() {
		var guildId string
		err = results.Scan(&guildId)
		if err != nil {
			continue
		}
		guildIds = append(guildIds, guildId)
	}

	// Получаем информацию о каждой гильдии
	for _, guildId := range guildIds {
		guild, err := d.Multi.GuildGetByIdV2(guildId)
		if err == nil && guild != nil {
			corporations = append(corporations, models.Guild{
				ID:   guild.GId.String(),
				Name: guild.GuildName,
				URL:  guild.AvatarUrl,
				Type: "mg",
			})
			continue
		}
	}

	return corporations, nil
}
