package DiscordClient

import (
	"github.com/bwmarrin/discordgo"
	"regexp"
	"strings"
)

type replace struct {
	session       *discordgo.Session
	guildsChannel map[string][]*discordgo.Channel
	guildsRoles   map[string][]*discordgo.Role
	guildsMembers map[string][]*discordgo.Member
}

func newReplace(session *discordgo.Session) *replace {
	return &replace{
		session:       session,
		guildsChannel: make(map[string][]*discordgo.Channel),
		guildsRoles:   make(map[string][]*discordgo.Role),
		guildsMembers: make(map[string][]*discordgo.Member),
	}
}

var (
	channelMentionRE = regexp.MustCompile("<#[0-9]+>")
	userMentionRE    = regexp.MustCompile("<@(\\d+)>")
	roleMentionRE    = regexp.MustCompile("<@&(\\d+)>")
)

func (d *Discord) ReplaceTextMessage(text string, guildid string) (newtext string) {
	return d.re.ReplaceTextMessage(text, guildid)
}

func (r *replace) ReplaceTextMessage(text string, guildid string) (newtext string) {
	newtext = r.replaceChannelMentions(text, guildid)
	newtext = r.replaceUserMentions(newtext, guildid)
	newtext = r.replaceRoleMentions(newtext, guildid)
	return newtext
}

func (r *replace) replaceChannelMentions(text string, guildid string) string {
	replaceChannelMentionFunc := func(match string) string {
		channelID := match[2 : len(match)-1]
		if len(r.guildsChannel[guildid]) == 0 || r.guildsChannel[guildid] == nil {
			r.guildsChannel[guildid], _ = r.session.GuildChannels(guildid)
		}
		channelName := "unknownchannel"
		for _, channel := range r.guildsChannel[guildid] {
			if channelID == channel.ID {
				channelName = channel.Name
			}
		}
		return "#" + channelName
	}
	return channelMentionRE.ReplaceAllStringFunc(text, replaceChannelMentionFunc)
}
func (r *replace) replaceRoleMentions(text string, guildid string) string {
	mentionIds := roleMentionRE.FindAllStringSubmatch(text, -1)
	for _, match := range mentionIds {
		mention := match[0]
		roleId := match[1]
		role := r.getRoleById(roleId, guildid)
		if role != nil {
			text = strings.Replace(text, mention, "@&"+role.Name, 1)
		}
	}
	return text
}
func (r *replace) getRoleById(roleId string, guildId string) *discordgo.Role {
	if r.guildsRoles[guildId] == nil || len(r.guildsRoles[guildId]) == 0 {
		r.guildsRoles[guildId], _ = r.session.GuildRoles(guildId)
	}

	for _, role := range r.guildsRoles[guildId] {
		if role.ID == roleId {
			return role
		}
	}
	return nil
}
func (r *replace) replaceUserMentions(text string, guildid string) string {
	mentionIds := userMentionRE.FindAllStringSubmatch(text, -1)
	for _, match := range mentionIds {
		mention := match[0]
		userId := match[1]
		username := r.getUserNameById(userId, guildid)
		text = strings.Replace(text, mention, "@"+username, 1)
	}
	return text
}
func (r *replace) getUserNameById(userId string, guildId string) string {
	if r.guildsMembers[guildId] == nil || len(r.guildsMembers) == 0 {
		r.guildsMembers[guildId], _ = r.session.GuildMembers(guildId, "", 999)
	}
	for _, member := range r.guildsMembers[guildId] {
		if member.User.ID == userId {
			if member.Nick != "" {
				return member.Nick
			} else {
				return member.User.Username
			}
		}
	}
	return "Unknown user"
}
func (r *replace) GetGuildRoles(guildId string) []*discordgo.Role {
	if r.guildsRoles[guildId] == nil || len(r.guildsRoles[guildId]) == 0 {
		r.guildsRoles[guildId], _ = r.session.GuildRoles(guildId)
	}
	return r.guildsRoles[guildId]
}
