package models

type IncomingMessage struct {
	Text         string
	DmChat       string
	Name         string
	MentionName  string
	NameId       string
	NickName     string
	Avatar       string
	AvatarF      string
	ChannelId    string
	GuildId      string
	GuildName    string
	GuildAvatar  string
	GuildAvatarF string
	Type         string
	Language     string
}
type ActionStruct struct {
	Action  string
	Message interface{}
}
type SendText struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
}
type DeleteMessageStruct struct {
	MessageId string `json:"message_id"`
	Channel   string `json:"channel"`
}
type SendTextDeleteSeconds struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
	Seconds int    `json:"seconds"`
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
