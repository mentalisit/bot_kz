package DiscordClient

type ChatAccess struct {
	GuildID   string `json:"chat_id"`
	GuildName string `json:"chat_name"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
}

//// GetUserChatAccess возвращает подробный доступ пользователя ко всем общим серверам
//func (d *Discord) GetUserChatAccess(userID string) ([]ChatAccess, error) {
//	var accessList []ChatAccess
//
//	d.S.State.RLock()
//	defer d.S.State.RUnlock()
//
//	for _, guild := range d.S.State.Guilds {
//		var status string
//
//		// 1. Проверяем, является ли пользователь владельцем
//		if guild.OwnerID == userID {
//			status = "creator"
//		} else {
//			// 2. Если не владелец, проверяем права
//			permissions, err := d.S.UserPermissions(userID, guild.ID)
//			if err != nil {
//				// Если пользователя нет в кэше сервера, он не может быть участником
//				continue
//			}
//
//			// Проверяем флаги администратора или управления сервером
//			isAdmin := (permissions&discordgo.PermissionAdministrator) != 0 ||
//				(permissions&discordgo.PermissionManageGuild) != 0
//
//			if isAdmin {
//				status = "administrator"
//			} else {
//				status = "member"
//			}
//		}
//
//		// Добавляем данные в результирующий слайс
//		accessList = append(accessList, ChatAccess{
//			GuildID:   guild.ID,
//			GuildName: guild.Name,
//			UserID:    userID,
//			Status:    status,
//		})
//	}
//
//	if len(accessList) == 0 {
//		return nil, fmt.Errorf("user not found in any common guilds")
//	}
//
//	return accessList, nil
//}
