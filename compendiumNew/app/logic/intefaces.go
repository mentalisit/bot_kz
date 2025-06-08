package logic

import "compendium/models"

type CorpMember interface {
	CorpMemberInsert(cm models.CorpMember) error
	CorpMembersRead(guildid string) ([]models.CorpMember, error)
	CorpMemberTZUpdate(userid, guildid, timeZone string, offset int) error
	CorpMemberAvatarUpdate(userid, guildid, avatarurl string) error
	CorpMemberByUserId(userId string) (*models.CorpMember, error)
	CorpMemberDelete(guildid string, nameId string) error
	CorpMemberDeleteAlt(guildid string, nameId string, name string) error
	CorpMembersReadByUserId(UserId string) ([]models.CorpMember, error)
}

type Tech interface {
	TechInsert(username, userid, guildid string, tech []byte) error
	TechGet(username, userid, guildid string) ([]byte, error)
	TechUpdate(username, userid, guildid string, tech []byte) error
	TechDelete(username, userid, guildid string) error
	TechGetName(username, guildid string) ([]byte, string, error)
	TechGetAllUserId(userid string) ([]models.TechTable, error)
}

//type Guilds interface {
//	GuildInsert(u models.Guild) error
//	GuildGet(guildid string) (*models.Guild, error)
//	GuildGetAll() ([]models.Guild, error)
//	GuildGetCountByGuildId(guildid string) (int, error)
//	GuildUpdate(u models.Guild) error
//}

//type ListUser interface {
//	ListUserInsert(token, userid, guildid string) error
//	ListUserGetCountByGuildIdByUserId(guildid, userid string) (int, error)
//	ListUserUpdate(token, userid, guildid string) error
//	ListUserDelete(token string) error
//	ListUserDeleteByUserIdByGuildId(userid, guildid string) error
//	ListUserGetToken(userid, guildid string) (string, error)
//	ListUserUpdateToken(tokenOld, tokenNew string) error
//	ListUserGetUserIdAndGuildId(token string) (userid string, guildid string, err error)
//}

type Users interface {
	UsersInsert(u models.User) error
	UsersGetByUserId(userid string) (*models.User, error)
	UsersGetByUserName(username string) (*models.User, error)
	UsersFindByGameName(gameName string) ([]models.User, error)
	UserGetCountByUserId(userid string) (int, error)
	UsersUpdate(u models.User) error
	UsersGetAll() ([]models.User, error)
}

type GuildRoles interface {
	GuildRoleCreate(guildid string, RoleName string) error
	GuildRoleExist(guildid string, RoleName string) bool
	GuildRoleDelete(guildid string, RoleName string) error
	GuildRolesRead(guildid string) ([]models.CorpRole, error)
	GuildRolesSubscribe(guildid, RoleName, userName, userid string) error
	GuildRolesExistSubscribe(guildid, RoleName, userid string) bool
	GuildRolesDeleteSubscribe(guildid, RoleName, userid string) error
	GuildRolesDeleteSubscribeUser(guildid, RoleName, userName, userid string) error
	//GuildRolesReadSubscribeUsers(guildid, RoleName string) ([]string, error)
}
