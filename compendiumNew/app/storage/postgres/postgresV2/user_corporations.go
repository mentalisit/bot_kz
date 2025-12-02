package postgresv2

import (
	"compendium/models"

	"github.com/google/uuid"
)

// UserCorporationsGet возвращает список корпораций, в которых состоит пользователь
// Теперь работает ТОЛЬКО с мульти-гильдиями (multi-guilds)
func (d *Db) UserCorporationsGet(identity *models.IdentityV2) ([]models.Guild, error) {
	var corporations []models.Guild

	// Проверяем, является ли пользователь multi-аккаунтом
	if identity.MultiAccount != nil {
		// Это multi-аккаунт
		corpMember, err := d.CorpMemberByUId(identity.MultiAccount.UUID)
		if err == nil && corpMember != nil {
			for _, guildId := range corpMember.GuildIds {
				guild, err := d.GuildGet(&guildId)
				if err == nil && guild != nil {
					corporations = append(corporations, models.Guild{
						ID:   guild.GId.String(),
						Name: guild.GuildName,
						URL:  guild.AvatarUrl,
						Type: "mg", // multi-guild
					})
				}
			}
		}
		return corporations, nil
	}

	// Для V2 получаем гильдии через multi-account
	// Получаем guildIds из таблицы corpmember по UUID мульти-аккаунта
	if identity.MultiAccount == nil {
		return corporations, nil // нет мульти-аккаунта
	}
	sel := "SELECT DISTINCT guildid FROM my_compendium.corpmember WHERE userid = $1"
	rows, err := d.db.Query(sel, identity.MultiAccount.UUID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var guildIds []string
	for rows.Next() {
		var guildId string
		err = rows.Scan(&guildId)
		if err != nil {
			continue
		}
		guildIds = append(guildIds, guildId)
	}

	// Получаем информацию о каждой мульти-гильдии
	for _, guildId := range guildIds {
		// Все гильдии теперь являются мульти-гильдиями (UUID)
		if gid, err := uuid.Parse(guildId); err == nil {
			guild, err := d.GuildGet(&gid)
			if err == nil && guild != nil {
				corporations = append(corporations, models.Guild{
					ID:   guild.GId.String(),
					Name: guild.GuildName,
					URL:  guild.AvatarUrl,
					Type: "mg", // multi-guild
				})
			}
		}
	}

	return corporations, nil
}
