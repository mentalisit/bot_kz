package DiscordClient

import (
	"discord/discord/helpers"
	"discord/models"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	emOK      = "✅"
	emCancel  = "❎"
	emRsStart = "🚀"
	emPl30    = "⌛"
	emPlus    = "➕"
	emMinus   = "➖"
)

func (d *Discord) readReactionQueue(r *discordgo.MessageReactionAdd, message *discordgo.Message) {
	user, err := d.S.User(r.UserID)
	if err != nil {
		d.log.ErrorErr(err)
	}
	if user.ID != message.Author.ID {
		ok, config := d.CheckChannelConfigDS(r.ChannelID)
		if ok {
			in := models.InMessage{
				Tip:         "ds",
				Username:    user.Username,
				UserId:      user.ID,
				NameNick:    "",
				NameMention: user.Mention(),
				Ds: struct {
					Mesid   string
					Guildid string
					Avatar  string
				}{
					Mesid:   r.MessageID,
					Guildid: config.Guildid,
					Avatar:  user.AvatarURL("128"),
				},

				Config: config,
				Option: models.Option{
					Reaction: true},
			}

			if r.Member != nil && r.Member.Nick != "" {
				in.NameNick = r.Member.Nick
			} else if in.NameNick == "" && user.GlobalName != "" {
				in.NameNick = user.GlobalName
			}

			d.reactionUserRemove(r)

			if r.Emoji.Name == emPlus {
				in.Mtext = "+"
			} else if r.Emoji.Name == emMinus {
				in.Mtext = "-"
			} else if r.Emoji.Name == emOK || r.Emoji.Name == emCancel || r.Emoji.Name == emRsStart || r.Emoji.Name == emPl30 {
				in.Lvlkz, err = d.storage.Db.ReadMesIdDS(r.MessageID)
				if err == nil && in.Lvlkz != "" {
					if r.Emoji.Name == emOK {
						in.Timekz = "30"
						in.Mtext = in.Lvlkz + "+"
					} else if r.Emoji.Name == emCancel {
						in.Mtext = in.Lvlkz + "-"
					} else if r.Emoji.Name == emRsStart {
						in.Mtext = in.Lvlkz + "++"
					} else if r.Emoji.Name == emPl30 {
						in.Mtext = in.Lvlkz + "+++"
					}
				}
			}
			d.api.SendRsBotAppRecover(in)
		}
	}
}
func (d *Discord) reactionUserRemove(r *discordgo.MessageReactionAdd) {
	err := d.S.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.Name, r.UserID)
	if err != nil {
		fmt.Println("Ошибка удаления эмоджи", err)
	}
}
func (d *Discord) logicMix(m *discordgo.MessageCreate) {
	if d.ifMentionBot(m) {
		return
	}
	if d.avatar(m) {
		return
	}
	go d.latinOrNot(m) //пытаемся переводить гостевой чат

	if m.Author != nil && m.Author.Locale != "" {
		go d.log.Info(m.Author.Username + " " + m.Author.Locale)
	}
	if m.Member != nil && m.Member.User != nil && m.Member.User.Locale != "" {
		go d.log.Info(m.Member.User.Username + " " + m.Member.User.Locale)
	}

	if strings.HasPrefix(m.Content, ".") {
		d.ifPrefixPoint(m)
		return
	}

	//filter Rs
	ok, config := d.CheckChannelConfigDS(m.ChannelID)
	if ok {
		d.SendToRsFilter(m, config)
		return
	}

	//bridge
	ds, bridgeConfig := d.BridgeCheckChannelConfigDS(m.ChannelID)
	if ds {
		d.SendToBridge(m, bridgeConfig)
	}

	if strings.HasPrefix(m.Content, "%") {
		d.SendToCompendium(m)
		return
	}
}
func (d *Discord) SendToRsFilter(m *discordgo.MessageCreate, config models.CorporationConfig) {
	if len(m.Attachments) > 0 {
		m.Content += m.Attachments[0].URL
	}
	if len(m.Message.Embeds) > 0 {
		m.Content += "\u200B"
	}
	in := models.InMessage{
		Mtext:       d.ReplaceTextMessage(m.Content, m.GuildID),
		Tip:         "ds",
		Username:    m.Author.Username,
		UserId:      m.Author.ID,
		NameNick:    "",
		NameMention: m.Author.Mention(),
		Ds: struct {
			Mesid   string
			Guildid string
			Avatar  string
		}{
			Mesid:   m.ID,
			Guildid: m.GuildID,
			Avatar:  m.Author.AvatarURL("128"),
		},
		Config: config,
		Option: models.Option{InClient: true},
	}
	if m.Member != nil && m.Member.Nick != "" {
		in.NameNick = m.Member.Nick
	}

	d.api.SendRsBotAppRecover(in)
}
func (d *Discord) ifMentionBot(m *discordgo.MessageCreate) bool {
	after, found := strings.CutPrefix(m.Content, d.S.State.User.Mention())
	if found {
		d.DeleteMesageSecond(m.ChannelID, m.ID, 30)
		goodRs, _ := d.CheckChannelConfigDS(m.ChannelID)
		if goodRs {
			d.SendChannelDelSecond(m.ChannelID, fmt.Sprintf("%s че пингуешь? пиши Справка,или пиши создателю бота @Mentalisit#5159 ", m.Author.Mention()), 30)
		} else {
			m.Content = "%" + after
			d.SendToCompendium(m)
		}
	}
	return found
}
func (d *Discord) readReactionTranslate(r *discordgo.MessageReactionAdd, m *discordgo.Message) {
	user, err := d.S.User(r.UserID)
	if err != nil {
		d.log.ErrorErr(err)
	}
	if user.ID != m.Author.ID {

		switch r.Emoji.Name {
		case "🇺🇸":
			d.transtale(m, "en", r)
		case "🇷🇺":
			d.transtale(m, "ru", r)
		case "🇺🇦":
			d.transtale(m, "uk", r)
		case "🇬🇧":
			d.transtale(m, "en", r)
		case "🇧🇾":
			d.transtale(m, "be", r)
		case "🇩🇪":
			d.transtale(m, "de", r)
		case "🇵🇱":
			d.transtale(m, "pl", r)
		}
	}
}
func (d *Discord) SendToCompendium(m *discordgo.MessageCreate) {
	g, err := d.S.Guild(m.GuildID)
	if err != nil {
		d.log.ErrorErr(err)
	}
	channel, _ := d.S.Channel(m.ChannelID)

	user := m.Author
	if m.Member != nil && m.Member.User != nil {
		user = m.Member.User
	}
	nick := ""
	if m.Member != nil && m.Member.Nick != "" {
		nick = m.Member.Nick
	}

	i := models.IncomingMessage{
		Text:         m.Content,
		DmChat:       d.dmChannel(user.ID),
		Name:         user.Username,
		MentionName:  user.Mention(),
		NameId:       user.ID,
		NickName:     nick,
		Avatar:       user.AvatarURL(""),
		AvatarF:      user.Avatar,
		ChannelId:    m.ChannelID,
		GuildId:      m.GuildID,
		GuildName:    g.Name,
		GuildAvatar:  g.IconURL(""),
		GuildAvatarF: g.Icon,
		Type:         "ds",
	}
	if channel != nil {
		i.Language = helpers.DetectLanguage(g.Name + "/" + channel.Name)
	} else {
		i.Language = helpers.DetectLanguage(g.Name)
	}

	d.api.SendCompendiumAppRecover(i)
}
func (d *Discord) ifPrefixPoint(m *discordgo.MessageCreate) {
	good, config := d.CheckChannelConfigDS(m.ID)
	in := models.InMessage{
		Mtext:       m.Content,
		Tip:         "ds",
		Username:    m.Author.Username,
		UserId:      m.Author.ID,
		NameMention: m.Author.Mention(),
		Ds: struct {
			Mesid   string
			Guildid string
			Avatar  string
		}{Mesid: m.ID, Guildid: m.GuildID, Avatar: m.Author.AvatarURL("")},

		Option: models.Option{
			InClient: true,
		},
	}
	if m.Member != nil && m.Member.Nick != "" {
		in.NameNick = m.Member.Nick
	}
	if good {
		in.Config = config
	} else {
		in.Config = models.CorporationConfig{
			CorpName:  d.GuildChatName(m.ChannelID, m.GuildID),
			DsChannel: m.ChannelID,
			Guildid:   m.GuildID,
		}
	}
	d.api.SendRsBotAppRecover(in)
	go func() {
		mes := models.ToBridgeMessage{
			Text:          m.Content,
			Sender:        m.Author.Username,
			Tip:           "ds",
			Avatar:        m.Author.AvatarURL("128"),
			ChatId:        m.ChannelID,
			MesId:         m.ID,
			GuildId:       m.GuildID,
			TimestampUnix: m.Timestamp.Unix(),
		}
		ds, bridgeConfig := d.BridgeCheckChannelConfigDS(m.ChannelID)
		if ds {
			mes.Config = &bridgeConfig
		} else {
			mes.Config = &models.Bridge2Config{
				HostRelay: d.GuildChatName(m.ChannelID, m.GuildID),
			}
		}
		d.api.SendBridgeAppRecover(mes)
	}()

}
func (d *Discord) SendToBridge(m *discordgo.MessageCreate, bridgeConfig models.Bridge2Config) {
	mes := models.ToBridgeMessage{
		ChatId:        m.ChannelID,
		Extra:         []models.FileInfo{},
		Config:        &bridgeConfig,
		Text:          d.ReplaceTextMessage(m.Content, m.GuildID),
		Sender:        d.getAuthorName(m),
		Tip:           "ds",
		MesId:         m.ID,
		GuildId:       m.GuildID,
		TimestampUnix: m.Timestamp.Unix(),
		Avatar:        m.Author.AvatarURL(""),
	}

	d.handleDownloadBridge(&mes, m)

	if m.ReferencedMessage != nil {
		mes.ReplyMap = make(map[string]string)
		mes.ReplyMap[m.ChannelID] = m.ReferencedMessage.ID
		usernameR := m.ReferencedMessage.Author.String()
		if m.ReferencedMessage.Member != nil && m.ReferencedMessage.Member.Nick != "" {
			usernameR = m.ReferencedMessage.Member.Nick
		}
		mes.Reply = &models.BridgeMessageReply{
			TimeMessage: m.ReferencedMessage.Timestamp.Unix(),
			Text:        d.ReplaceTextMessage(m.ReferencedMessage.Content, m.GuildID),
			Avatar:      m.ReferencedMessage.Author.AvatarURL(""),
			UserName:    usernameR,
		}
	}
	if mes.Text != "" || len(mes.Extra) > 0 {
		d.api.SendBridgeAppRecover(mes)
	}
}
