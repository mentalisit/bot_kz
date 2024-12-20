package DiscordClient

import (
	"discord/models"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

//lang ok

//func (d *Discord) AddEnojiRsQueue1(chatid, mesid string) {
//	err := d.S.MessageReactionAdd(chatid, mesid, emOK)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//	err = d.S.MessageReactionAdd(chatid, mesid, emCancel)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//	err = d.S.MessageReactionAdd(chatid, mesid, emRsStart)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//	err = d.S.MessageReactionAdd(chatid, mesid, emPl30)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//}

func (d *Discord) AddButtonsQueue(level string) []discordgo.MessageComponent {
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
	_ = d.S.ChannelMessageDelete(chatid, mesid)
}
func (d *Discord) DeleteMesageSecond(chatid, mesid string, second int) {
	tu := int(time.Now().UTC().Unix())
	d.storage.Db.TimerInsert(models.Timer{
		Dsmesid:  mesid,
		Dschatid: chatid,
		Timed:    tu + second,
	})
}

//	func (d *Discord) EditComplex1(dsmesid, dschatid string, Embeds discordgo.MessageEmbed) error {
//		_, err := d.S.ChannelMessageEditComplex(&discordgo.MessageEdit{
//			Content: &mesContentNil,
//			Embed:   &Embeds,
//			ID:      dsmesid,
//			Channel: dschatid,
//		})
//		if err != nil {
//			return err
//		}
//		return nil
//	}
var mesContentNil string

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
		messages, _ := d.S.ChannelMessages(dschatid, 20, "", "", "")
		for i, message := range messages {
			if message.Author.String() == "Rs_bot#9945" {
				fmt.Printf("1test%d message.Author.Username==\"Rs_bot#9945\"  %+v\n", i, message)
				if len(message.Embeds) > 0 && message.Embeds[0].Title == "Очередь ТКЗ" {
					d.DeleteMessage(message.ChannelID, message.ID)
				} else if len(message.Embeds) == 0 {
					d.DeleteMessage(message.ChannelID, message.ID)
				}
			}
		}
		return err
	}
	return nil
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

func (d *Discord) Subscribe(nameid, argRoles, guildid string) int {
	exist, role := d.roleExists(d.re.GetGuildRoles(guildid), argRoles)

	if !exist { //если нет роли
		var err error
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

	exist, role := d.roleExists(d.re.GetGuildRoles(guildid), argRoles)
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

func (d *Discord) EmbedDS(mapa map[string]string, numkz int, count int) discordgo.MessageEmbed {
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
	title := mapa["title"]
	return discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  16711680,
		Description: fmt.Sprintf("👇 %s <:rs:918545444425072671> %s (%d) ",
			mapa["wishing_to"], mapa["lvlkz"], numkz) +
			textcount,

		Fields: []*discordgo.MessageEmbedField{{
			Name: fmt.Sprintf(" %s %s\n%s %s\n%s %s",
				emOK, mapa["to_add_to_queue"],
				emCancel, mapa["to_exit_the_queue"],
				emRsStart, mapa["forced_start"]),
			Value:  mapa["data_updated"] + ": ",
			Inline: true,
		}},
		Timestamp: time.Now().Format(time.RFC3339), // ТЕКУЩЕЕ ВРЕМЯ ДИСКОРДА
		Title:     title,
	}
}

func (d *Discord) ChannelTyping(ChannelID string) {
	err := d.S.ChannelTyping(ChannelID)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
}

func (d *Discord) SendDmText(text, AuthorID string) {
	dm := d.dmChannel(AuthorID)
	mes, err := d.S.ChannelMessageSend(dm, text)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
	d.DeleteMesageSecond(dm, mes.ID, 600)
}
