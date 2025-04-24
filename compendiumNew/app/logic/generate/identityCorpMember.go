package generate

import "compendium/models"

func GenerateIdentity(m models.IncomingMessage) (models.Identity, models.CorpMember) {
	identity := models.Identity{
		User: models.User{
			ID:        m.NameId,
			Username:  m.Name,
			AvatarURL: m.Avatar,
			Alts:      []string{},
		},
		Guild: models.Guild{
			URL:  m.GuildAvatar,
			ID:   m.GuildId,
			Name: m.GuildName,
			Type: m.Type,
		},
		Token: m.Type + m.GuildId + "." + m.NameId + GenerateToken(174),
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
