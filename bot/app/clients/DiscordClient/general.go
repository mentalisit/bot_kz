package DiscordClient

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"kz_bot/models"
	"kz_bot/pkg/utils"
	"time"
)

// lang ok
const (
	emOne  = "1️⃣"
	emTwo  = "2️⃣"
	emTree = "3️⃣"
	emFour = "4️⃣"
	emFive = "5️⃣"
)

func (d *Discord) AddButtonPoll(createTime string, option []string) []discordgo.MessageComponent {
	var components []discordgo.MessageComponent
	if len(option) > 0 {
		if len(option) > 0 && option[0] != "" {
			button := discordgo.Button{
				Style:    discordgo.SecondaryButton,
				Label:    "",
				CustomID: createTime + ".1",
				Emoji: &discordgo.ComponentEmoji{
					Name: emOne,
				},
			}
			components = append(components, button)
		}
		if len(option) > 1 && option[1] != "" {
			button := discordgo.Button{
				Style:    discordgo.SecondaryButton,
				Label:    "",
				CustomID: createTime + ".2",
				Emoji: &discordgo.ComponentEmoji{
					Name: emTwo,
				},
			}
			components = append(components, button)
		}
		if len(option) > 2 && option[2] != "" {
			button := discordgo.Button{
				Style:    discordgo.SecondaryButton,
				Label:    "",
				CustomID: createTime + ".3",
				Emoji: &discordgo.ComponentEmoji{
					Name: emTree,
				},
			}
			components = append(components, button)
		}

		if len(option) > 3 && option[3] != "" {
			button := discordgo.Button{
				Style:    discordgo.SecondaryButton,
				Label:    "",
				CustomID: createTime + ".4",
				Emoji: &discordgo.ComponentEmoji{
					Name: emFour,
				},
			}
			components = append(components, button)
		}
		if len(option) > 4 && option[4] != "" {
			button := discordgo.Button{
				Style:    discordgo.SecondaryButton,
				Label:    "",
				CustomID: createTime + ".5",
				Emoji: &discordgo.ComponentEmoji{
					Name: emFive,
				},
			}
			components = append(components, button)
		}
	}
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: components,
		},
	}
}

func (d *Discord) addButtonsQueue(level string) []discordgo.MessageComponent {
	// Создание кнопки
	buttonOk := discordgo.Button{
		Style:    discordgo.SecondaryButton,
		Label:    level + "+",
		CustomID: level + "+",
		Emoji: &discordgo.ComponentEmoji{
			Name: emOK,
		},
	}
	buttonCancel := discordgo.Button{
		Style:    discordgo.SecondaryButton,
		Label:    level + "-",
		CustomID: level + "-",
		Emoji: &discordgo.ComponentEmoji{
			Name: emCancel,
		},
	}
	buttonRsStart := discordgo.Button{
		Style:    discordgo.SecondaryButton,
		Label:    level + "++",
		CustomID: level + "++",
		Emoji: &discordgo.ComponentEmoji{
			Name: emRsStart,
		},
	}
	buttonPl30 := discordgo.Button{
		Style:    discordgo.SecondaryButton,
		Label:    "+30",
		CustomID: level + "+++",
		Emoji: &discordgo.ComponentEmoji{
			Name: emPl30,
		},
	}

	// Создание компонентов с кнопкой
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				buttonOk,
				buttonCancel,
				buttonRsStart,
				buttonPl30,
			},
		},
	}
}

func (d *Discord) DeleteMessage(chatid, mesid string) {
	ch := utils.WaitForMessage("DelMessage")
	defer close(ch)
	_ = d.S.ChannelMessageDelete(chatid, mesid)
}
func (d *Discord) DeleteMesageSecond(chatid, mesid string, second int) {
	if second > 60 {
		d.storage.TimeDeleteMessage.TimerInsert(models.Timer{
			Dsmesid:  mesid,
			Dschatid: chatid,
			Timed:    second,
		})
	} else {
		go func() {
			time.Sleep(time.Duration(second) * time.Second)
			err := d.S.ChannelMessageDelete(chatid, mesid)
			if err != nil {
				fmt.Println("Ошибка удаления сообщения дискорда ", chatid, mesid, second)
			}
		}()
	}
}

//func (d *Discord) EditComplex1(dsmesid, dschatid string, Embeds discordgo.MessageEmbed) error {
//	_, err := d.S.ChannelMessageEditComplex(&discordgo.MessageEdit{
//		Content: &mesContentNil,
//		Embed:   &Embeds,
//		ID:      dsmesid,
//		Channel: dschatid,
//	})
//	if err != nil {
//		return err
//	}
//	return nil
//}

func (d *Discord) EditComplexButton(dsmesid, dschatid string, mapEmbed map[string]string) error {
	components := d.addButtonsQueue(mapEmbed["buttonLevel"])
	embed := d.embedDS(mapEmbed)
	_, err := d.S.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Content:    &mesContentNil,
		Embed:      embed,
		ID:         dsmesid,
		Channel:    dschatid,
		Components: &components,
	})
	if err != nil {
		return err
	}
	return nil
}
func (d *Discord) EditComplexButton2(dsmesid, dschatid string, Embeds discordgo.MessageEmbed, component []discordgo.MessageComponent) error {
	_, err := d.S.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Content:    &mesContentNil,
		Embed:      &Embeds,
		ID:         dsmesid,
		Channel:    dschatid,
		Components: &component,
	})
	if err != nil {
		return err
	}
	return nil
}
func (d *Discord) Subscribe(nameid, argRoles, guildid string) int {
	g, err := d.S.State.Guild(guildid)
	if err != nil {
		d.log.ErrorErr(err)
		g, err = d.S.Guild(guildid)
		if err != nil {
			d.log.ErrorErr(err)
		}
	}

	exist, role := d.roleExists(g, argRoles)

	if !exist { //если нет роли
		role, err = d.createRole(argRoles, guildid)
		if err != nil {
			d.log.ErrorErr(err)
		}
	}

	member, err := d.S.GuildMember(guildid, nameid)
	if err != nil {
		d.log.ErrorErr(err)
	}
	var subscribe int = 0
	if exist {
		for _, _role := range member.Roles {
			if _role == role.ID {
				subscribe = 1
			}
		}
	}

	err = d.S.GuildMemberRoleAdd(guildid, nameid, role.ID)
	if err != nil {
		d.log.ErrorErr(err)
		subscribe = 2
	}

	return subscribe
}
func (d *Discord) Unsubscribe(nameid, argRoles, guildid string) int {
	var unsubscribe int = 0
	g, err := d.S.State.Guild(guildid)
	if err != nil {
		d.log.ErrorErr(err)
		g, err = d.S.Guild(guildid)
		if err != nil {
			d.log.ErrorErr(err)
		}
	}

	exist, role := d.roleExists(g, argRoles)
	if !exist { //если нет роли
		unsubscribe = 1
	}

	member, err := d.S.GuildMember(guildid, nameid)
	if err != nil {
		d.log.ErrorErr(err)
	}
	if exist {
		for _, _role := range member.Roles {
			if _role == role.ID {
				unsubscribe = 2
			}
		}
	}
	if unsubscribe == 2 {
		err = d.S.GuildMemberRoleRemove(guildid, nameid, role.ID)
		if err != nil {
			d.log.ErrorErr(err)
			unsubscribe = 3
		}
	}

	return unsubscribe
}

func (d *Discord) EditMessage(chatID, messageID, content string) {
	_, err := d.S.ChannelMessageEdit(chatID, messageID, content)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Discord) EditWebhook(text, username, chatID, mID string, avatarURL string) {
	if text == "" {
		return
	}

	params := &discordgo.WebhookParams{
		Content:   text,
		Username:  username,
		AvatarURL: avatarURL,
	}
	err := d.webhook.Edit(chatID, mID, params)
	if err != nil {
		return
	}
}

func (d *Discord) EmbedDS2(mapa map[string]string, numkz int, count int, dark bool) discordgo.MessageEmbed {
	textcount := ""
	if count == 1 {
		textcount = fmt.Sprintf("\n1️⃣ %s \n\n",
			mapa["name1"])
	} else if count == 2 {
		textcount = fmt.Sprintf("\n1️⃣ %s \n2️⃣ %s \n\n",
			mapa["name1"], mapa["name2"])
	} else if count == 3 {
		textcount = fmt.Sprintf("\n1️⃣ %s \n2️⃣ %s \n3️⃣ %s \n\n",
			mapa["name1"], mapa["name2"], mapa["name3"])
	} else {
		textcount = fmt.Sprintf("\n1️⃣ %s \n2️⃣ %s \n3️⃣ %s \n4️⃣ %s \n",
			mapa["name1"], mapa["name2"], mapa["name3"], mapa["name4"])
	}
	title := d.getLanguage(mapa["lang"], "rs_queue")
	if dark {
		title = d.getLanguage(mapa["lang"], "queue_drs")
	}
	return discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  16711680,
		Description: fmt.Sprintf("👇 %s <:rs:918545444425072671> %s (%d) ",
			d.getLanguage(mapa["lang"], "wishing_to"), mapa["lvlkz"], numkz) +
			textcount,

		Fields: []*discordgo.MessageEmbedField{{
			Name: fmt.Sprintf(" %s %s\n%s %s\n%s %s",
				emOK, d.getLanguage(mapa["lang"], "to_add_to_queue"),
				emCancel, d.getLanguage(mapa["lang"], "to_exit_the_queue"),
				emRsStart, d.getLanguage(mapa["lang"], "forced_start")),
			Value:  d.getLanguage(mapa["lang"], "data_updated") + ": ",
			Inline: true,
		}},
		Timestamp: time.Now().Format(time.RFC3339), // ТЕКУЩЕЕ ВРЕМЯ ДИСКОРДА
		Title:     title,
	}
}
func (d *Discord) embedDS(mapa map[string]string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       16711680,
		Description: mapa["description"] + mapa["textcount"],

		Fields: []*discordgo.MessageEmbedField{{
			Name:   mapa["EmbedFieldName"],
			Value:  mapa["EmbedFieldValue"],
			Inline: true,
		}},
		Timestamp: time.Now().Format(time.RFC3339), // ТЕКУЩЕЕ ВРЕМЯ ДИСКОРДА
		Title:     mapa["title"],
	}
}

func (d *Discord) ChannelTyping(ChannelID string) {
	ch := utils.WaitForMessage("ChannelTyping")
	defer close(ch)
	err := d.S.ChannelTyping(ChannelID)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
}

func (d *Discord) SendDmText(text, AuthorID string) {
	ch := utils.WaitForMessage("SendDmText")
	defer close(ch)
	dm := d.dmChannel(AuthorID)
	mes, err := d.S.ChannelMessageSend(dm, text)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
	d.DeleteMesageSecond(dm, mes.ID, 600)
}
