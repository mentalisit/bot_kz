package DiscordClient

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"kz_bot/clients/helper"
	"kz_bot/clients/restapi"
	"kz_bot/models"
	"strings"
	"time"
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
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				in.Lvlkz, err = d.storage.DbFunc.ReadMesIdDS(ctx, r.MessageID)
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
			d.ChanRsMessage <- in
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
	//d.AccesChatDS(m)
	if m.Author != nil && m.Author.Locale != "" {
		go d.log.Info(m.Author.Username + " " + m.Author.Locale)
	}
	if m.Member != nil && m.Member.User != nil && m.Member.User.Locale != "" {
		go d.log.Info(m.Member.User.Username + " " + m.Member.User.Locale)
	}

	if strings.HasPrefix(m.Content, "%") {
		d.SendToCompendium(m)
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
		mes := models.ToBridgeMessage{
			ChatId: m.ChannelID,
			Extra:  []models.FileInfo{},
			Config: &bridgeConfig,
		}
		d.filterNewBridge(m, mes)
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
		Mtext:       m.Content,
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

	d.ChanRsMessage <- in

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
		i.Language = helper.DetectLanguage(g.Name + "/" + channel.Name)
	} else {
		i.Language = helper.DetectLanguage(g.Name)
	}

	err = restapi.SendCompendiumApp(i)
	if err != nil {
		d.log.InfoStruct("SendCompendiumApp", i)
		d.log.ErrorErr(err)
		return
	}

}
