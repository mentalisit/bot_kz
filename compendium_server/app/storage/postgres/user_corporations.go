package postgres

import (
	"compendium_s/models"

	"github.com/google/uuid"
)

// UserCorporationsGet возвращает список корпораций, в которых состоит пользователь
func (d *Db) UserCorporationsGet(identity *models.Identity) ([]models.Guild, error) {
	ctx, cancel := d.getContext()
	defer cancel()

	var corporations []models.Guild

	// Сначала проверяем, является ли пользователь multi-аккаунтом
	if identity.MultiAccount != nil {
		// Это multi-аккаунт
		corpMember, err := d.Multi.CorpMemberByUId(identity.MultiAccount.UUID)
		if err == nil && corpMember != nil {
			for _, guildId := range corpMember.GuildIds {
				guild, err := d.Multi.GuildGet(&guildId)
				if err == nil && guild != nil {
					corporations = append(corporations, models.Guild{
						ID:   guild.GId.String(),
						Name: guild.GuildName,
						URL:  guild.AvatarUrl,
						Type: "mg",
					})
				}
			}
		}
		return corporations, nil
	}

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
		// Сначала пытаемся найти в multi гильдиях
		if gid, err := uuid.Parse(guildId); err == nil {
			guild, err := d.Multi.GuildGet(&gid)
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

		// Если не нашли в multi, ищем в обычных гильдиях
		guild, err := d.GuildGet(guildId)
		if err == nil && guild != nil {
			corporations = append(corporations, *guild)
		}
	}

	return corporations, nil
}
