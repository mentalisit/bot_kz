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
	Type         string
	Language     string
	MultiAccount *MultiAccount
	//MultiGuild   *MultiAccountGuild
	MAcc   *MultiAccount
	MGuild *MultiAccountGuildV2
}
