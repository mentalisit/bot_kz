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
			URL:  m.MultiGuild.AvatarUrl,
			ID:   m.MultiGuild.GuildId(),
			Name: m.MultiGuild.GuildName,
			Type: m.Type,
		},
		//Token: m.Type + m.GuildId + "." + m.NameId + GenerateToken(174),
	}
	if m.MultiGuild != nil {
		identity.MultiGuild = m.MultiGuild
	}
	if m.MultiAccount != nil {
		identity.MultiAccount = m.MultiAccount
	}
	//if m.MultiGuild != nil {
	//	identity.Guild = models.Guild{
	//		URL:  m.MultiGuild.AvatarUrl,
	//		ID:   m.MultiGuild.GuildId(),
	//		Name: m.MultiGuild.GuildName,
	//		Type: "mg",
	//	}
	//}
	if m.MultiAccount != nil && m.MultiAccount.DiscordID != "" && m.MultiAccount.TelegramID != "" {
		identity.User = models.User{
			ID:        m.MultiAccount.UUID.String(),
			Username:  m.MultiAccount.Nickname,
			AvatarURL: m.MultiAccount.AvatarURL,
			Alts:      m.MultiAccount.Alts,
		}
	}

	cm := models.CorpMember{
		Name:      identity.User.Username,
		UserId:    identity.User.ID,
		GuildId:   m.MultiGuild.GuildId(),
		Avatar:    "",
		Tech:      map[int][2]int{},
		AvatarUrl: identity.User.AvatarURL,
	}
	if m.MultiAccount != nil && m.MultiAccount.DiscordID != "" && m.MultiAccount.TelegramID != "" {
		cm.Name = m.MultiAccount.Nickname
		cm.AvatarUrl = m.MultiAccount.AvatarURL
		cm.MultiAccount = m.MultiAccount
	}
	if m.MultiGuild != nil {
		cm.MultiGuild = m.MultiGuild
	}
	userJWT, err := JWTGenerateTokenForUser(identity)
	if err == nil {
		identity.Token = userJWT
	}

	return identity, cm
}
