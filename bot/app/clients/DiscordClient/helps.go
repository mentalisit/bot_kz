package DiscordClient

import (
	"fmt"
	"kz_bot/models"
)

func (d *Discord) Help(Channel, lang string) {
	mId := d.hhelp1(Channel, lang)
	d.DeleteMesageSecond(Channel, mId, 184)
}

func (d *Discord) HelpChannelUpdate(c models.CorporationConfig) string {
	if c.MesidDsHelp == "" {
		c.MesidDsHelp = d.hhelp1(c.DsChannel, c.Country)
		return c.MesidDsHelp
	} else {
		messages, err := d.S.ChannelMessages(c.DsChannel, 10, "", c.MesidDsHelp, "")
		if err != nil {
			go d.DeleteMessage(c.DsChannel, c.MesidDsHelp)
			c.MesidDsHelp = d.hhelp1(c.DsChannel, c.Country)
			return c.MesidDsHelp
		}
		if len(messages) > 2 {
			go d.DeleteMessage(c.DsChannel, c.MesidDsHelp)
			c.MesidDsHelp = d.hhelp1(c.DsChannel, c.Country)
		}
	}
	return c.MesidDsHelp
}

func (d *Discord) hhelp1(chatid, lang string) string {
	m := d.SendEmbedText(chatid, d.getLanguage(lang, "spravka"),
		fmt.Sprintf("%s \n\n%s", d.getLanguage(lang, "botUdalyaet"), d.getLanguage(lang, "hhelpText")))
	return m.ID
}
