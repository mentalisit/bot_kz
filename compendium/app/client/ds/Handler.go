package ds

//func (d *Discord) GetRoles(guildId string) []models.CorpRole {
//	roles, err := d.s.GuildRoles(guildId)
//	if err != nil {
//		d.log.ErrorErr(err)
//		return nil
//	}
//	var guildRole []models.CorpRole
//	for _, role := range roles {
//		r := models.CorpRole{
//			Name: role.Name,
//			Id:   role.ID,
//		}
//		if r.Name == "@everyone" {
//			r.Id = ""
//		}
//
//		guildRole = append(guildRole, r)
//	}
//	return guildRole
//}
//
//func (d *Discord) CheckRole(guildId, memderId, roleid string) bool {
//	if roleid == "" {
//		return true
//	}
//	member, err := d.s.GuildMember(guildId, memderId)
//	if err != nil {
//		d.log.ErrorErr(err)
//		return false
//	}
//	for _, role := range member.Roles {
//		if roleid == role {
//			return true
//		}
//	}
//	return false
//}
