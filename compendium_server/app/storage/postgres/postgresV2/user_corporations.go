package postgresv2

import (
	"compendium_s/models"
)

// UserCorporationsGet возвращает список корпораций, в которых состоит пользователь
// Теперь работает ТОЛЬКО с мульти-гильдиями (multi-guilds)
func (d *Db) UserCorporationsGet(identity *models.Identity) ([]models.Guild, error) {
	var corporations []models.Guild

	// Проверяем, является ли пользователь multi-аккаунтом
	if identity.MAccount != nil {
		// Это multi-аккаунт
		corpMember, err := d.CorpMemberByUId(identity.MAccount.UUID)
		if err == nil && corpMember != nil {
			for _, guildId := range corpMember.GuildIds {
				guild, err := d.GuildGet(guildId)
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
	}
	return corporations, nil
}
