package generate

import "compendium/models"

func GenerateIdentity(m models.IncomingMessage) models.Identity {
	i := models.Identity{
		User: models.User{
			ID:        m.NameId,
			Username:  m.Name,
			AvatarURL: m.Avatar,
			Alts:      []string{},
		},
	}

	if m.MultiAccount != nil {
		i.MultiAccount = m.MultiAccount
	}

	if m.MAcc != nil && m.MGuild != nil {
		token, err := JWTGenerateToken(m.MAcc.UUID, m.MGuild.GId)
		if err == nil {
			i.Token = token
		}
		i.MGuild = m.MGuild
		i.MAccount = m.MAcc
		i.User = models.User{
			ID:        m.MAcc.UUID.String(),
			Username:  m.MAcc.Nickname,
			AvatarURL: m.MAcc.AvatarURL,
			Alts:      m.MAcc.Alts,
			GameName:  m.MAcc.Nickname,
		}
		i.Guild = models.Guild{
			URL:  m.MGuild.AvatarUrl,
			ID:   m.MGuild.GId.String(),
			Name: m.MGuild.GuildName,
			Type: m.Type,
		}
	}
	return i
}
