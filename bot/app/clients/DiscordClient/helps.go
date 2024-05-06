package DiscordClient

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"kz_bot/models"
)

func (d *Discord) Help(Channel, lang string) {
	mId := d.Hhelp1(Channel, lang)
	d.DeleteMesageSecond(Channel, mId, 184)
}

func (d *Discord) HelpChannelUpdate(c models.CorporationConfig) string {
	if c.MesidDsHelp == "" {
		c.MesidDsHelp = d.Hhelp1(c.DsChannel, c.Country)
		return c.MesidDsHelp
	} else {
		messages, err := d.S.ChannelMessages(c.DsChannel, 10, "", c.MesidDsHelp, "")
		if err != nil {
			go d.DeleteMessage(c.DsChannel, c.MesidDsHelp)
			c.MesidDsHelp = d.Hhelp1(c.DsChannel, c.Country)
			return c.MesidDsHelp
		}
		if len(messages) > 3 {
			go d.DeleteMessage(c.DsChannel, c.MesidDsHelp)
			c.MesidDsHelp = d.Hhelp1(c.DsChannel, c.Country)
		}
	}
	return c.MesidDsHelp
}

func (d *Discord) Hhelp1(chatid, lang string) string {
	Emb := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       16711680,
		Description: fmt.Sprintf("%s \n\n%s", d.getLanguage(lang, "info_bot_delete_msg"), d.getLanguage(lang, "info_help_text")),
		Title:       d.getLanguage(lang, "information"),
	}

	m, err := d.S.ChannelMessageSendComplex(chatid, &discordgo.MessageSend{
		//Content:    "for RS",
		Components: d.AddButtonsStartQueue(chatid),
		Embed:      Emb,
	})
	if err != nil {
		d.log.ErrorErr(err)
	}

	return m.ID
}
