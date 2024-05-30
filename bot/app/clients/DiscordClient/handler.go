package DiscordClient

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"kz_bot/clients/DiscordClient/slashCommand"
	"kz_bot/config"
	"time"
)

func (d *Discord) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Message.WebhookID != "" {
		return
	}
	if m.GuildID == "" {
		d.log.Info(m.Author.Username + ": " + m.Content)
		if m.Content == ".паника" {
			d.log.Panic(".паника " + m.Author.Username)
		} else {
			d.Send(m.ChannelID, "эээ я же бот че ты мне пишешь тут, пиши в канале ")
		}
		//DM message
		return
	}
	d.logicMix(m)

}
func (d *Discord) messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if m.Message.WebhookID != "" {
		return
	}

	if m.Message.EditedTimestamp != nil && m.Content != "" {
		//good, config := d.BridgeCheckChannelConfigDS(m.ChannelID)
		//if good {
		//	username := m.Author.Username
		//	if m.Member != nil && m.Member.Nick != "" {
		//		username = m.Member.Nick
		//	}
		//	mes := models.BridgeMessage{
		//		Text:          d.replaceTextMessage(m.Content, m.GuildID),
		//		Sender:        username,
		//		Tip:           "dse",
		//		Avatar:        m.Author.AvatarURL("128"),
		//		ChatId:        m.ChannelID,
		//		MesId:         m.ID,
		//		GuildId:       m.GuildID,
		//		TimestampUnix: m.Timestamp.Unix(),
		//		Config:        &config,
		//	}
		//
		//	if len(m.Attachments) > 0 {
		//		if len(m.Attachments) != 1 {
		//			d.log.Info(fmt.Sprintf("вложение %d", len(m.Attachments)))
		//		}
		//		for _, attachment := range m.Attachments {
		//			mes.FileUrl = append(mes.FileUrl, attachment.URL)
		//		}
		//	}
		//
		//	if m.ReferencedMessage != nil {
		//		usernameR := m.ReferencedMessage.Author.String() //.Username
		//		if m.ReferencedMessage.Member != nil && m.ReferencedMessage.Member.Nick != "" {
		//			usernameR = m.ReferencedMessage.Member.Nick
		//		}
		//		mes.Reply = &models.BridgeMessageReply{
		//			TimeMessage: m.ReferencedMessage.Timestamp.Unix(),
		//			Text:        d.replaceTextMessage(m.ReferencedMessage.Content, m.GuildID),
		//			Avatar:      m.ReferencedMessage.Author.AvatarURL("128"),
		//			UserName:    usernameR,
		//		}
		//	}
		//
		//	d.ChanBridgeMessage <- mes
		//}
	}
}
func (d *Discord) onMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	//good, config := d.BridgeCheckChannelConfigDS(m.ChannelID)
	//if good {
	//	d.ChanBridgeMessage <- models.BridgeMessage{
	//		Tip:    "delDs",
	//		MesId:  m.ID,
	//		Config: &config,
	//	}
	//}
}

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
		d.readReactionQueue(r, message)
	} else {
		d.readReactionTranslate(r, message)
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
		d.handleButtonPressed(i)

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
			commandsModuleWeapon := slashCommand.AddSlashCommandModuleWeaponLocale()
			if len(commandsModuleWeapon) == 0 {
				return
			}
			for _, v := range commandsModuleWeapon {
				_, err := d.S.ApplicationCommandCreate(d.S.State.User.ID, configrs.Guildid, v)
				if err != nil {
					d.log.ErrorErr(err)
				}
			}
		}
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

//func (d *Discord) addSlashHandler() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
//	var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
//		"help": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
//			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//				Type: discordgo.InteractionResponseChannelMessageWithSource,
//				Data: &discordgo.InteractionResponseData{
//
//					Content: models.Help,
//				},
//			})
//			go func() {
//				time.Sleep(1 * time.Minute)
//				s.InteractionResponseDelete(i.Interaction)
//			}()
//		},
//		"helpqueue": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
//			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//				Type: discordgo.InteractionResponseChannelMessageWithSource,
//				Data: &discordgo.InteractionResponseData{
//					Content: models.HelpQueue,
//				},
//			})
//			go func() {
//				time.Sleep(1 * time.Minute)
//				s.InteractionResponseDelete(i.Interaction)
//			}()
//		},
//		"helpnotification": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
//			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//				Type: discordgo.InteractionResponseChannelMessageWithSource,
//				Data: &discordgo.InteractionResponseData{
//					Content: "Уведомления:\n" +
//						"	Подписаться на уведомления о начале очереди: +[4-11]\n" +
//						"+10 -подписаться на уведомления о начале очереди на КЗ 10ур.\n\n" +
//						"	Подписаться на уведомление, если в очереди 3 человека: ++[4-11]\n" +
//						"++10 -подписаться на уведомления о наличии 3х человек в очереди на КЗ 10ур.\n\n" +
//						"	Отключить уведомления о начале сбора: -[5-11]\n" +
//						"-9 -отключить уведомления о начале сборе на КЗ 9ур.\n\n" +
//						"	Отключить уведомления 3/4 в очереди: --[5-11]\n" +
//						"--9 -отключить уведомления о наличии 3х человек в очереди на КЗ 9ур.",
//				},
//			})
//			go func() {
//				time.Sleep(1 * time.Minute)
//				s.InteractionResponseDelete(i.Interaction)
//			}()
//		},
//		"helpevent": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
//			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//				Type: discordgo.InteractionResponseChannelMessageWithSource,
//				Data: &discordgo.InteractionResponseData{
//					Content: models.HelpEvent,
//				},
//			})
//			go func() {
//				time.Sleep(1 * time.Minute)
//				s.InteractionResponseDelete(i.Interaction)
//			}()
//		},
//		"helptop": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
//			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//				Type: discordgo.InteractionResponseChannelMessageWithSource,
//				Data: &discordgo.InteractionResponseData{
//					Content: models.HelpTop,
//				},
//			})
//			go func() {
//				time.Sleep(1 * time.Minute)
//				s.InteractionResponseDelete(i.Interaction)
//			}()
//		},
//		"helpicon": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
//			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//				Type: discordgo.InteractionResponseChannelMessageWithSource,
//				Data: &discordgo.InteractionResponseData{
//					Content: models.HelpIcon,
//				},
//			})
//			go func() {
//				time.Sleep(1 * time.Minute)
//				s.InteractionResponseDelete(i.Interaction)
//			}()
//		},
//	}
//
//	return commandHandlers
//}

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	commands = []*discordgo.ApplicationCommand{
		{
			Name:                     "permission-overview",
			Description:              "Command for demonstration of default command permissions",
			DefaultMemberPermissions: &defaultMemberPermissions,
			DMPermission:             &dmPermission,
		},
		{
			Name:        "options",
			Description: "Command for demonstrating options",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "string-option",
					Description: "String option",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "integer-option",
					Description: "Integer option",
					MinValue:    &integerOptionMinValue,
					MaxValue:    10,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "number-option",
					Description: "Float option",
					MaxValue:    10.1,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "bool-option",
					Description: "Boolean option",
					Required:    true,
				},

				// Required options must be listed first since optional parameters
				// always come after when they're used.
				// The same concept applies to Discord's Slash-commands API

				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel-option",
					Description: "Channel option",
					// Channel type mask
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText,
						discordgo.ChannelTypeGuildVoice,
					},
					Required: false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-option",
					Description: "User option",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role-option",
					Description: "Role option",
					Required:    false,
				},
			},
		},
		{
			Name:        "subcommands",
			Description: "Subcommands and command groups example",
			Options: []*discordgo.ApplicationCommandOption{
				// When a command has subcommands/subcommand groups
				// It must not have top-level options, they aren't accesible in the UI
				// in this case (at least not yet), so if a command has
				// subcommands/subcommand any groups registering top-level options
				// will cause the registration of the command to fail

				{
					Name:        "subcommand-group",
					Description: "Subcommands group",
					Options: []*discordgo.ApplicationCommandOption{
						// Also, subcommand groups aren't capable of
						// containing options, by the name of them, you can see
						// they can only contain subcommands
						{
							Name:        "nested-subcommand",
							Description: "Nested subcommand",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommandGroup,
				},
				// Also, you can create both subcommand groups and subcommands
				// in the command at the same time. But, there's some limits to
				// nesting, count of subcommands (top level and nested) and options.
				// Read the intro of slash-commands docs on Discord dev portal
				// to get more information
				{
					Name:        "subcommand",
					Description: "Top-level subcommand",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		{
			Name:        "followups",
			Description: "Followup messages",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"options": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			// This example stores the provided arguments in an []interface{}
			// which will be used to format the bot's response
			margs := make([]interface{}, 0, len(options))
			msgformat := "You learned how to use command options! " +
				"Take a look at the value(s) you entered:\n"

			// Get the value from the option map.
			// When the option exists, ok = true
			if option, ok := optionMap["string-option"]; ok {
				// Option values must be type asserted from interface{}.
				// Discordgo provides utility functions to make this simple.
				margs = append(margs, option.StringValue())
				msgformat += "> string-option: %s\n"
			}

			if opt, ok := optionMap["integer-option"]; ok {
				margs = append(margs, opt.IntValue())
				msgformat += "> integer-option: %d\n"
			}

			if opt, ok := optionMap["number-option"]; ok {
				margs = append(margs, opt.FloatValue())
				msgformat += "> number-option: %f\n"
			}

			if opt, ok := optionMap["bool-option"]; ok {
				margs = append(margs, opt.BoolValue())
				msgformat += "> bool-option: %v\n"
			}

			if opt, ok := optionMap["channel-option"]; ok {
				margs = append(margs, opt.ChannelValue(nil).ID)
				msgformat += "> channel-option: <#%s>\n"
			}

			if opt, ok := optionMap["user-option"]; ok {
				margs = append(margs, opt.UserValue(nil).ID)
				msgformat += "> user-option: <@%s>\n"
			}

			if opt, ok := optionMap["role-option"]; ok {
				margs = append(margs, opt.RoleValue(nil, "").ID)
				msgformat += "> role-option: <@&%s>\n"
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, they will be discussed in "responses"
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
		},
		"permission-overview": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			perms, err := s.ApplicationCommandPermissions(s.State.User.ID, i.GuildID, i.ApplicationCommandData().ID)

			var restError *discordgo.RESTError
			if errors.As(err, &restError) && restError.Message != nil && restError.Message.Code == discordgo.ErrCodeUnknownApplicationCommandPermissions {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: ":x: No permission overwrites",
					},
				})
				return
			} else if err != nil {
				panic(err)
			}

			if err != nil {
				panic(err)
			}
			format := "- %s %s\n"

			channels := ""
			users := ""
			roles := ""

			for _, o := range perms.Permissions {
				emoji := "❌"
				if o.Permission {
					emoji = "☑"
				}

				switch o.Type {
				case discordgo.ApplicationCommandPermissionTypeUser:
					users += fmt.Sprintf(format, emoji, "<@!"+o.ID+">")
				case discordgo.ApplicationCommandPermissionTypeChannel:
					allChannels, _ := discordgo.GuildAllChannelsID(i.GuildID)

					if o.ID == allChannels {
						channels += fmt.Sprintf(format, emoji, "All channels")
					} else {
						channels += fmt.Sprintf(format, emoji, "<#"+o.ID+">")
					}
				case discordgo.ApplicationCommandPermissionTypeRole:
					if o.ID == i.GuildID {
						roles += fmt.Sprintf(format, emoji, "@everyone")
					} else {
						roles += fmt.Sprintf(format, emoji, "<@&"+o.ID+">")
					}
				}
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Permissions overview",
							Description: "Overview of permissions for this command",
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:  "Users",
									Value: users,
								},
								{
									Name:  "Channels",
									Value: channels,
								},
								{
									Name:  "Roles",
									Value: roles,
								},
							},
						},
					},
					AllowedMentions: &discordgo.MessageAllowedMentions{},
				},
			})
		},
		"responses": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Responses to a command are very important.
			// First of all, because you need to react to the interaction
			// by sending the response in 3 seconds after receiving, otherwise
			// interaction will be considered invalid and you can no longer
			// use the interaction token and ID for responding to the user's request

			content := ""
			// As you can see, the response type names used here are pretty self-explanatory,
			// but for those who want more information see the official documentation
			switch i.ApplicationCommandData().Options[0].IntValue() {
			case int64(discordgo.InteractionResponseChannelMessageWithSource):
				content =
					"You just responded to an interaction, sent a message and showed the original one. " +
						"Congratulations!"
				content +=
					"\nAlso... you can edit your response, wait 5 seconds and this message will be changed"
			default:
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
				})
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
				}
				return
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong",
				})
				return
			}
			time.AfterFunc(time.Second*5, func() {
				content := content + "\n\nWell, now you know how to create and edit responses. " +
					"But you still don't know how to delete them... so... wait 10 seconds and this " +
					"message will be deleted."
				_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &content,
				})
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
					return
				}
				time.Sleep(time.Second * 10)
				s.InteractionResponseDelete(i.Interaction)
			})
		},
		"followups": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Followup messages are basically regular messages (you can create as many of them as you wish)
			// but work as they are created by webhooks and their functionality
			// is for handling additional messages after sending a response.

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					// Note: this isn't documented, but you can use that if you want to.
					// This flag just allows you to create messages visible only for the caller of the command
					// (user who triggered the command)
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "Surprise!",
				},
			})
			msg, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "Followup message has been created, after 5 seconds it will be edited",
			})
			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong",
				})
				return
			}
			time.Sleep(time.Second * 5)

			content := "Now the original message is gone and after 10 seconds this message will ~~self-destruct~~ be deleted."
			s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
				Content: &content,
			})

			time.Sleep(time.Second * 10)

			s.FollowupMessageDelete(i.Interaction, msg.ID)

			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "For those, who didn't skip anything and followed tutorial along fairly, " +
					"take a unicorn :unicorn: as reward!\n" +
					"Also, as bonus... look at the original interaction response :D",
			})
		},
	}
)
