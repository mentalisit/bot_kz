package serverV2

import (
	"compendium_s/models"
	"strings"
)

func (s *ServerV2) GetTokenIdentity(token string) *models.IdentityV2 {
	var i models.IdentityV2
	if strings.HasPrefix(token, "my_compendium_") {
		i.Token = token
		// удаляем префикс "my_compendium_"
		uid, err := GetTokenData(token[14:])
		if err == nil {
			i.MultiAccount, _ = s.db.FindMultiAccountUUID(uid)
			return &i
		} else {
			s.log.ErrorErr(err)
		}
	}
	return &i
}

//func (s *ServerV2) getRoles(i models.Guild) []models.CorpRole {
//	if len(i.ID) > 24 {
//		parse, _ := uuid.Parse(i.ID)
//		GUILD, err := s.db.GuildGet(&parse)
//		if err != nil {
//			return nil
//		}
//		var rolesAll []models.CorpRole
//		for channel := range GUILD.Channels {
//			if strings.HasPrefix(channel, "-100") {
//				everyone := []models.CorpRole{{
//					Id:   "",
//					Name: "@everyone",
//				}}
//				roles, err := s.db.GuildRolesRead(channel)
//				if err != nil {
//					s.log.ErrorErr(err)
//				}
//				if len(roles) > 0 {
//					everyone = append(everyone, roles...)
//				}
//				rolesAll = append(rolesAll, everyone...)
//			}
//		}
//		return rolesAll
//	}
//
//	if i.Type == "tg" {
//		everyone := []models.CorpRole{{
//			Id:   "",
//			Name: "@everyone",
//		}}
//		roles, err := s.db.GuildRolesRead(i.ID)
//		if err != nil {
//			s.log.ErrorErr(err)
//		}
//		if len(roles) > 0 {
//			everyone = append(everyone, roles...)
//		}
//		return everyone
//	}
//	return nil
//}
//
//func (s *ServerV2) CheckCode(code string) models.Identity {
//	// Code functionality removed in V2
//	return models.Identity{}
//}
//func (s *ServerV2) CleanOldCodes() {
//	// Code functionality removed in V2
//}

func (s *ServerV2) refreshToken(token string) string {
	//if strings.HasPrefix(token, "Multi_") || strings.HasPrefix(token, "identity_") {
	//return token
	//}

	//if len(token) < 60 {
	//	newToken := s.checkPrefixToken(token)
	//	err := s.db.ListUserUpdateToken(token, newToken)
	//	if err != nil {
	//		return token
	//	}
	//	return newToken
	//}
	return token
}

//func (s *ServerV2) checkPrefixToken(token string) string {
//	userid, guildid, err := s.db.ListUserGetUserIdAndGuildId(token)
//	if err != nil || userid == "" || guildid == "" {
//		s.log.Info("err || userid == nil || guildid == nil")
//		return token
//	}
//
//	// Для V2 используем только мульти-гильдии
//	multiGuild, err := s.db.GuildGetById(guildid)
//	if err != nil || multiGuild == nil || multiGuild.GuildName == "" {
//		s.log.Info("err || multiGuild not found or empty name")
//		return token
//	}
//
//	// Для мульти-гильдий тип всегда "mg"
//	newToken := "mg" + guildid + "." + userid + GenerateToken()
//	return newToken
//}
