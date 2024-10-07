package DiscordClient

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"kz_bot/clients/helper"
	"kz_bot/config"
	"kz_bot/models"
	"path/filepath"
)

func (d *Discord) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Message.WebhookID != "" {
		return
	}
	if m.GuildID == "" {
		if m.Content == ".паника" {
			d.log.Panic(".паника " + m.Author.Username)
		} else {
			in := models.InMessage{
				Mtext:       m.Content,
				Tip:         "dsDM",
				Username:    m.Author.Username,
				UserId:      m.Author.ID,
				NameMention: m.Author.Mention(),
				Ds: struct {
					Mesid   string
					Guildid string
					Avatar  string
				}{
					Mesid:   m.ID,
					Guildid: "",
					Avatar:  m.Author.AvatarURL("128"),
				},
				Config: models.CorporationConfig{
					DsChannel: m.ChannelID},
			}

			d.ChanRsMessage <- in
		}
		//DM message
		return
	}
	go d.logicMix(m)

}

//func (d *Discord) messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
//	if m.Message.WebhookID != "" {
//		return
//	}
//
//	if m.Message.EditedTimestamp != nil && m.Content != "" {
//		//good, config := d.BridgeCheckChannelConfigDS(m.ChannelID)
//		//if good {
//		//	username := m.Author.Username
//		//	if m.Member != nil && m.Member.Nick != "" {
//		//		username = m.Member.Nick
//		//	}
//		//	mes := models.BridgeMessage{
//		//		Text:          d.replaceTextMessage(m.Content, m.GuildID),
//		//		Sender:        username,
//		//		Tip:           "dse",
//		//		Avatar:        m.Author.AvatarURL("128"),
//		//		ChatId:        m.ChannelID,
//		//		MesId:         m.ID,
//		//		GuildId:       m.GuildID,
//		//		TimestampUnix: m.Timestamp.Unix(),
//		//		Config:        &config,
//		//	}
//		//
//		//	if len(m.Attachments) > 0 {
//		//		if len(m.Attachments) != 1 {
//		//			d.log.Info(fmt.Sprintf("вложение %d", len(m.Attachments)))
//		//		}
//		//		for _, attachment := range m.Attachments {
//		//			mes.FileUrl = append(mes.FileUrl, attachment.URL)
//		//		}
//		//	}
//		//
//		//	if m.ReferencedMessage != nil {
//		//		usernameR := m.ReferencedMessage.Author.String() //.Username
//		//		if m.ReferencedMessage.Member != nil && m.ReferencedMessage.Member.Nick != "" {
//		//			usernameR = m.ReferencedMessage.Member.Nick
//		//		}
//		//		mes.Reply = &models.BridgeMessageReply{
//		//			TimeMessage: m.ReferencedMessage.Timestamp.Unix(),
//		//			Text:        d.replaceTextMessage(m.ReferencedMessage.Content, m.GuildID),
//		//			Avatar:      m.ReferencedMessage.Author.AvatarURL("128"),
//		//			UserName:    usernameR,
//		//		}
//		//	}
//		//
//		//	d.ChanBridgeMessage <- mes
//		//}
//	}
//}

//func (d *Discord) onMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
//good, config := d.BridgeCheckChannelConfigDS(m.ChannelID)
//if good {
//	d.ChanBridgeMessage <- models.BridgeMessage{
//		Tip:    "delDs",
//		MesId:  m.ID,
//		Config: &config,
//	}
//}
//}

func (d *Discord) messageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		channel, err1 := s.Channel(r.ChannelID)
		if err1 != nil {
			d.log.Error(err1.Error())
			return
		}
		user, err2 := s.User(r.UserID)
		if err2 != nil {
			d.log.Error(err2.Error())
			return
		}
		d.log.Info(fmt.Sprintln(channel.Name, r.Emoji.Name, user.Username, err.Error()))
		return
	}

	if message.Author.ID == s.State.User.ID {
		go d.readReactionQueue(r, message)
	} else {
		go d.readReactionTranslate(r, message)
	}
}

func (d *Discord) slash(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//if config.Instance.BotMode == "dev" {
	//	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
	//			h(s, i)
	//		}
	//	})
	//}
	if i.Interaction != nil && i.Interaction.Member != nil && i.Interaction.Member.Nick != "" {
		i.Interaction.Member.User.Username = i.Interaction.Member.Nick
	}

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		{
			locale := ""
			switch i.Locale.String() {
			case "Russian":
				locale = "ru"
			case "Ukrainian":
				locale = "ua"
			default:
				locale = "en"
			}
			switch i.ApplicationCommandData().Name {
			case "module":
				d.handleModuleCommand(i, locale)
			case "weapon":
				d.handleWeaponCommand(i, locale)
			default:
				d.log.Info(fmt.Sprintf("slash InteractionApplicationCommand  %+v\n", i.ApplicationCommandData()))
			}
		}
	case discordgo.InteractionMessageComponent:
		go d.handleButtonPressed(i)

	default:
		d.log.Info(fmt.Sprintf("slash %+v\n", i.Type))
	}

}

func (d *Discord) ready() {
	if config.Instance.BotMode == "dev" {
		//d.removeCommand("") //700238199070523412
		//commandsTest := commands
		//if len(commandsTest) == 0 {
		//	return
		//}
		//for _, v := range commandsTest {
		//	_, err := d.S.ApplicationCommandCreate(d.S.State.User.ID, "", v)
		//	if err != nil {
		//		d.log.ErrorErr(err)
		//	}
		//}
		return
	}
	for _, configrs := range d.corpConfigRS {
		if configrs.DsChannel != "" && configrs.Guildid != "" {
			//d.removeCommand(configrs.Guildid)
			commandsModuleWeapon := AddSlashCommandModuleWeaponLocale()
			if len(commandsModuleWeapon) == 0 {
				return
			}
			for _, v := range commandsModuleWeapon {
				_, err := d.S.ApplicationCommandCreate(d.S.State.User.ID, configrs.Guildid, v)
				if err != nil {
					d.log.ErrorErr(err)
					break
				}
			}
		}
	}
}
func AddSlashCommandModuleWeaponLocale() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "module",
			Description: "Select the desired module and level",
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian:   "модули",
				discordgo.Ukrainian: "модулі",
				discordgo.EnglishUS: "module",
			},
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian:   "Выберите нужный модуль и уровень",
				discordgo.Ukrainian: "Виберіть потрібний модуль та рівень",
				discordgo.EnglishUS: "Select the desired module and level",
			},
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "module",
					Description: "Select module",
					NameLocalizations: map[discordgo.Locale]string{
						discordgo.Russian:   "модули",
						discordgo.Ukrainian: "модулі",
						discordgo.EnglishUS: "module",
					},
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.Russian:   "Выберите модуль",
						discordgo.Ukrainian: "Виберіть модуль",
						discordgo.EnglishUS: "Select module",
					},
					Required: true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name: "RSE",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Ингибитор КЗ",
								discordgo.Ukrainian: "Інгібітор ЧЗ",
								discordgo.EnglishUS: "RSE",
							},
							Value: "RSE",
						},
						{
							Name: "Genesis",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Генезис",
								discordgo.Ukrainian: "Генезис",
								discordgo.EnglishUS: "Genesis",
							},
							Value: "GENESIS",
						},
						{
							Name: "Enrich",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Обогатить",
								discordgo.Ukrainian: "Збагатити",
								discordgo.EnglishUS: "Enrich",
							},
							Value: "ENRICH",
						},
						// Добавьте другие модули по мере необходимости
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "level",
					Description: "Select level",
					NameLocalizations: map[discordgo.Locale]string{
						discordgo.Russian:   "уровень",
						discordgo.Ukrainian: "рівень",
						discordgo.EnglishUS: "level",
					},
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.Russian:   "Выберите уровень",
						discordgo.Ukrainian: "Виберіть рівень",
						discordgo.EnglishUS: "Select level",
					},
					Required: true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name: "Level 0",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 0",
								discordgo.Ukrainian: "Рівень 0",
								discordgo.EnglishUS: "Level 0",
							},
							Value: 0,
						},
						{
							Name: "Level 1",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 1",
								discordgo.Ukrainian: "Рівень 1",
								discordgo.EnglishUS: "Level 1",
							},
							Value: 1,
						}, {
							Name: "Level 2",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 2",
								discordgo.Ukrainian: "Рівень 2",
								discordgo.EnglishUS: "Level 2",
							},
							Value: 2,
						}, {
							Name: "Level 3",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 3",
								discordgo.Ukrainian: "Рівень 3",
								discordgo.EnglishUS: "Level 3",
							},
							Value: 3,
						}, {
							Name: "Level 4",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 4",
								discordgo.Ukrainian: "Рівень 4",
								discordgo.EnglishUS: "Level 4",
							},
							Value: 4,
						}, {
							Name: "Level 5",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 5",
								discordgo.Ukrainian: "Рівень 5",
								discordgo.EnglishUS: "Level 5",
							},
							Value: 5,
						}, {
							Name: "Level 6",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 6",
								discordgo.Ukrainian: "Рівень 6",
								discordgo.EnglishUS: "Level 6",
							},
							Value: 6,
						}, {
							Name: "Level 7",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 7",
								discordgo.Ukrainian: "Рівень 7",
								discordgo.EnglishUS: "Level 7",
							},
							Value: 7,
						}, {
							Name: "Level 8",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 8",
								discordgo.Ukrainian: "Рівень 8",
								discordgo.EnglishUS: "Level 8",
							},
							Value: 8,
						}, {
							Name: "Level 9",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 9",
								discordgo.Ukrainian: "Рівень 9",
								discordgo.EnglishUS: "Level 9",
							},
							Value: 9,
						}, {
							Name: "Level 10",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 10",
								discordgo.Ukrainian: "Рівень 10",
								discordgo.EnglishUS: "Level 10",
							},
							Value: 10,
						}, {
							Name: "Level 11",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 11",
								discordgo.Ukrainian: "Рівень 11",
								discordgo.EnglishUS: "Level 11",
							},
							Value: 11,
						}, {
							Name: "Level 12",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 12",
								discordgo.Ukrainian: "Рівень 12",
								discordgo.EnglishUS: "Level 12",
							},
							Value: 12,
						}, {
							Name: "Level 13",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 13",
								discordgo.Ukrainian: "Рівень 13",
								discordgo.EnglishUS: "Level 13",
							},
							Value: 13,
						}, {
							Name: "Level 14",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 14",
								discordgo.Ukrainian: "Рівень 14",
								discordgo.EnglishUS: "Level 14",
							},
							Value: 14,
						}, {
							Name: "Level 15",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 15",
								discordgo.Ukrainian: "Рівень 15",
								discordgo.EnglishUS: "Level 15",
							},
							Value: 15,
						},
						// Добавьте другие уровни по мере необходимости
					},
				},
			},
		},
		{
			Name:        "weapon",
			Description: "Select your main weapon",
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian:   "оружие",
				discordgo.Ukrainian: "зброя",
				discordgo.EnglishUS: "weapon",
			},
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian:   "Выберите оружие",
				discordgo.Ukrainian: "Виберіть основну зброю",
				discordgo.EnglishUS: "Select your main weapon",
			},
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "weapon",
					Description: "Select weapon",
					NameLocalizations: map[discordgo.Locale]string{
						discordgo.Russian:   "оружие",
						discordgo.Ukrainian: "зброя",
						discordgo.EnglishUS: "weapon",
					},
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.Russian:   "Выберите оружие",
						discordgo.Ukrainian: "Виберіть зброю",
						discordgo.EnglishUS: "Select weapon",
					},
					Required: true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name: "Barrage",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Артобстрел",
								discordgo.Ukrainian: "Артилерія",
								discordgo.EnglishUS: "Barrage",
							},
							Value: "barrage",
						},
						{
							Name: "Laser",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Лазер",
								discordgo.Ukrainian: "Лазер",
								discordgo.EnglishUS: "Laser",
							},
							Value: "laser",
						},
						{
							Name: "Chain ray",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Цепной луч",
								discordgo.Ukrainian: "Ланцюговий промінь",
								discordgo.EnglishUS: "Chain ray",
							},
							Value: "chainray",
						},
						{
							Name: "Battery",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Батарея",
								discordgo.Ukrainian: "Батарея",
								discordgo.EnglishUS: "Battery",
							},
							Value: "battery",
						},
						{
							Name: "Mass battery",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Залповая батарея",
								discordgo.Ukrainian: "Залпова батарея",
								discordgo.EnglishUS: "Mass battery",
							},
							Value: "massbattery",
						},
						{
							Name: "Dart launcher",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Пусковая установка",
								discordgo.Ukrainian: "Пускова установка",
								discordgo.EnglishUS: "Dart launcher",
							},
							Value: "dartlauncher",
						},
						{
							Name: "Rocket launcher",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Ракетная установка",
								discordgo.Ukrainian: "Ракетна установка",
								discordgo.EnglishUS: "Rocket launcher",
							},
							Value: "rocketlauncher",
						},
						{
							Name: "Remove weapon",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Удалить оружие",
								discordgo.Ukrainian: "Видалити зброю",
								discordgo.EnglishUS: "Remove weapon",
							},
							Value: "Remove",
						},
					},
				},
			},
		},
	}
}

func (d *Discord) removeCommand(guildId string) {
	registeredCommands, err := d.S.ApplicationCommands(d.S.State.User.ID, guildId)
	if err != nil {
		d.log.Fatal(err.Error())
	}

	fmt.Println(len(registeredCommands))
	for _, v := range registeredCommands {
		fmt.Printf("%+v\n", v)
	}

	for _, v := range registeredCommands {
		err = d.S.ApplicationCommandDelete(d.S.State.User.ID, guildId, v.ID)
		if err != nil {
			d.log.ErrorErr(err)
		}
	}
	fmt.Println("удалены")
}

func (d *Discord) GetAvatarUrl(userId string) string {
	user, err := d.S.User(userId)
	if err != nil {
		return ""
	}
	return user.AvatarURL("")
}

func (d *Discord) handleDownloadBridge(mes *models.ToBridgeMessage, m *discordgo.MessageCreate) {
	if len(m.StickerItems) > 0 {
		mes.Text = fmt.Sprintf("https://cdn.discordapp.com/stickers/%s.png", m.Message.StickerItems[0].ID)
	}
	if len(m.Attachments) > 0 {
		for _, a := range m.Attachments {
			f := models.FileInfo{
				Name: a.Filename,
				Data: nil,
				URL:  a.URL,
				Size: int64(a.Size),
			}
			if filepath.Ext(a.Filename) == ".apk" {
				f.URL = ""
				data, err := helper.DownloadFile(a.URL)
				if err != nil {
					d.log.ErrorErr(err)
				}
				f.Data = data
			}

			mes.Extra = append(mes.Extra, f)
		}
	}
}
