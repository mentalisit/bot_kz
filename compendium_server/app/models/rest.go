package models

type CheckRoleStruct struct {
	GuildId  string `json:"guild_id"`
	MemberId string `json:"member_id"`
	RoleId   string `json:"role_id"`
}

type DsMembersRoles struct {
	Userid  string   `json:"userid"`
	RolesId []string `json:"rolesId"`
}
