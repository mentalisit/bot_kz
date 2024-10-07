package DiscordClient

import (
	"github.com/bwmarrin/discordgo"
	"strconv"
)

//func (d *Discord) Help(Channel, lang string) {
//	mId := d.HelpEmbed(Channel,"","", []string{})
//	d.DeleteMesageSecond(Channel, mId, 184)
//}
//
//func (d *Discord) HelpChannelUpdate(c models.CorporationConfig) string {
//	if c.MesidDsHelp == "" {
//		c.MesidDsHelp = d.Hhelp1(c.DsChannel, c.Country)
//		return c.MesidDsHelp
//	} else {
//		messages, err := d.S.ChannelMessages(c.DsChannel, 10, "", c.MesidDsHelp, "")
//		if err != nil {
//			go d.DeleteMessage(c.DsChannel, c.MesidDsHelp)
//			c.MesidDsHelp = d.Hhelp1(c.DsChannel, c.Country)
//			return c.MesidDsHelp
//		}
//		if len(messages) > 3 {
//			go d.DeleteMessage(c.DsChannel, c.MesidDsHelp)
//			c.MesidDsHelp = d.Hhelp1(c.DsChannel, c.Country)
//		}
//	}
//	return c.MesidDsHelp
//}
//func (d *Discord) HelpEmbed(chatid, Title, Description string, levels []string) string {
//	Emb := &discordgo.MessageEmbed{
//		Author:      &discordgo.MessageEmbedAuthor{},
//		Color:       16711680,
//		Description: fmt.Sprintf("%s \n\n%s", d.getLanguage(lang, "info_bot_delete_msg"), d.getLanguage(lang, "info_help_text")),
//		Title:       d.getLanguage(lang, "information"),
//	}
//
//	m, err := d.S.ChannelMessageSendComplex(chatid, &discordgo.MessageSend{
//		//Content:    "for RS",
//		Components: d.AddButtonsStartQueue(chatid, levels),
//		Embed:      Emb,
//	})
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//
//	return m.ID
//}

func (d *Discord) SendHelp(chatid, Title, Description string, levels []string) string {
	Emb := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       16711680,
		Description: Description,
		Title:       Title,
	}

	m, err := d.S.ChannelMessageSendComplex(chatid, &discordgo.MessageSend{
		Components: d.AddButtonsStartQueue(chatid, levels),
		Embed:      Emb,
	})
	if err != nil {
		d.log.ErrorErr(err)
	}

	return m.ID
}
func (d *Discord) AddButtonsStartQueue(chatid string, levels []string) []discordgo.MessageComponent {
	var mc []discordgo.MessageComponent
	var components []discordgo.MessageComponent
	//_, config := d.CheckChannelConfigDS(chatid)
	//levels := d.storage.Count.ReadTop5Level(config.CorpName)
	if len(levels) > 0 {
		for _, level := range levels {
			button := discordgo.Button{}

			if level[:1] == "d" {
				button.Style = discordgo.DangerButton
				button.Label = level[1:] + "*"
				button.CustomID = level[1:] + "*"
			} else {
				button.Style = discordgo.SecondaryButton
				button.Label = level + "+"
				button.CustomID = level + "+"
			}
			components = append(components, button)
		}
	}

	if len(components) == 0 {
		for i := 7; i < 12; i++ {
			l := strconv.Itoa(i)

			button := discordgo.Button{
				Label:    l + "+",
				Style:    discordgo.SecondaryButton,
				CustomID: l + "+",
			}
			components = append(components, button)

		}
	}
	mc = append(mc, discordgo.ActionsRow{Components: components})

	good, CC := d.CheckChannelConfigDS(chatid)
	if good {
		event := d.storage.Db.NumActiveEvent(CC.CorpName)
		if event > 0 {
			var componentsEvent []discordgo.MessageComponent
			for i := 7; i < 12; i++ {
				l := "s" + strconv.Itoa(i)
				button := discordgo.Button{
					Label:    l + "+",
					Style:    discordgo.DangerButton,
					CustomID: l + "+",
				}
				componentsEvent = append(componentsEvent, button)
			}
			mc = append(mc, discordgo.ActionsRow{Components: componentsEvent})
		}
	}

	return mc
}
