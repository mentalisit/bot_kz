package models

type SendText struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
}
type CheckRoleStruct struct {
	GuildId  string `json:"guild_id"`
	MemberId string `json:"member_id"`
	RoleId   string `json:"role_id"`
}
type SendPic struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
	Pic     []byte `json:"pic"`
}
