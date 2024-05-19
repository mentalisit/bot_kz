package generate

import "compendium/models"

func GenerateIdentity(m models.IncomingMessage) (models.Identity, models.CorpMember) {
	//проверить если есть NameId то предложить соединить для двух корпораций
	identity := models.Identity{
		User: models.User{
			ID:        m.NameId,
			Username:  m.Name,
			Avatar:    m.AvatarF,
			AvatarURL: m.Avatar,
			Alts:      []string{},
		},
		Guild: models.Guild{
			URL:  m.GuildAvatar,
			ID:   m.GuildId,
			Name: m.GuildName,
			Icon: m.GuildAvatarF,
			Type: m.Type,
		},
		Token: GenerateToken(),
	}

	cm := models.CorpMember{
		Name:      m.Name,
		UserId:    m.NameId,
		GuildId:   m.GuildId,
		Avatar:    m.AvatarF,
		Tech:      map[int][2]int{},
		AvatarUrl: m.Avatar,
	}
	return identity, cm
}
