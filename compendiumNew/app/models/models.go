package models

type IncomingMessage struct {
	Text        string
	DmChat      string
	Name        string
	MentionName string
	NameId      string
	NickName    string
	Avatar      string
	AvatarF     string
	ChannelId   string
	//GuildId      string
	//GuildName    string
	//GuildAvatar  string
	//GuildAvatarF string
	Type         string
	Language     string
	MultiAccount *MultiAccount
	MultiGuild   *MultiAccountGuild
}
