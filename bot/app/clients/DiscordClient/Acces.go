package DiscordClient

import (
	"github.com/bwmarrin/discordgo"
	"kz_bot/clients/restapi"
	"kz_bot/models"
	"time"
)

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
	d.ChanRsMessage <- in
	go func() {
		time.Sleep(5 * time.Second)
		d.corpConfigRS = d.storage.CorpConfigRS
	}()
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
			Config: &models.BridgeConfig{
				HostRelay: d.GuildChatName(m.ChannelID, m.GuildID),
			},
		}
		err := restapi.SendBridgeApp(mes)
		if err != nil {
			d.log.ErrorErr(err)
			return
		}
		go func() {
			time.Sleep(5 * time.Second)
			d.storage.ReloadDbArray()
		}()

	}()

}
