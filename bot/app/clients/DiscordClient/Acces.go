package DiscordClient

import (
	"github.com/bwmarrin/discordgo"
	"kz_bot/clients/restapi"
	"kz_bot/models"
	"time"
)

//nujno sdelat lang
//
//func (d *Discord) AccesChatDS(m *discordgo.MessageCreate) {
//	after, res := strings.CutPrefix(m.Content, ".")
//	if res {
//		switch after {
//		case "add":
//			go d.DeleteMesageSecond(m.ChannelID, m.ID, 10)
//			d.accessAddChannelDs(m.ChannelID, m.GuildID, "en")
//		case "добавить":
//			go d.DeleteMesageSecond(m.ChannelID, m.ID, 10)
//			d.accessAddChannelDs(m.ChannelID, m.GuildID, "ru")
//		case "додати":
//			go d.DeleteMesageSecond(m.ChannelID, m.ID, 10)
//			d.accessAddChannelDs(m.ChannelID, m.GuildID, "ua")
//		case "del", "удалить", "видалити":
//			go d.DeleteMesageSecond(m.ChannelID, m.ID, 10)
//			d.accessDelChannelDs(m.ChannelID, m.GuildID)
//		case "паника":
//			d.log.Panic("перезагрузка по требованию")
//			os.Exit(1)
//		case "removeCommand":
//			d.removeCommand("")
//			go d.ready()
//		default:
//			if d.CleanOldMessage(m) {
//				return
//			}
//			if d.setLang(m) {
//				return
//			}
//		}
//	}
//}
//func (d *Discord) accessAddChannelDs(chatid, guildid, lang string) { // внесение в дб и добавление в масив
//	ok, _ := d.CheckChannelConfigDS(chatid)
//	if ok {
//		go d.SendChannelDelSecond(chatid, d.getLanguage(lang, "info_activation_not_required"), 30)
//	} else {
//		chatName := d.GuildChatName(chatid, guildid)
//		d.log.Info("новая активация корпорации " + chatName)
//		d.AddDsCorpConfig(chatName, chatid, guildid, lang)
//		go d.SendChannelDelSecond(chatid, d.getLanguage(lang, "tranks_for_activation"), 10)
//
//	}
//}
//func (d *Discord) accessDelChannelDs(chatid, guildid string) { //удаление с бд и масива для блокировки
//	ok, config := d.CheckChannelConfigDS(chatid)
//	d.DeleteMessage(chatid, config.MesidDsHelp)
//	if !ok {
//		go d.SendChannelDelSecond(chatid, d.getLanguage("ru", "channel_not_connected"), 60)
//	} else {
//		d.SendChannelDelSecond(chatid, d.getLang(chatid, "you_disabled_bot_functions"), 60)
//		d.storage.ConfigRs.DeleteConfigRs(config)
//		d.storage.ReloadDbArray()
//		d.corpConfigRS = d.storage.CorpConfigRS
//		d.log.Info("отключение корпорации " + d.GuildChatName(chatid, guildid))
//	}
//}
//func (d *Discord) setLang(m *discordgo.MessageCreate) bool {
//	re := regexp.MustCompile(`^\.set lang (\w{2})$`)
//	matches := re.FindStringSubmatch(m.Content)
//	if len(matches) > 0 {
//		langUpdate := matches[1]
//		ok, config := d.CheckChannelConfigDS(m.ChannelID)
//		if ok {
//			if d.storage.Dictionary.CheckTranslateLanguage(langUpdate) {
//				go d.updateLanguage(langUpdate, config, m)
//			} else {
//				d.SendChannelDelSecond(m.ChannelID, "please wait, I'm trying to translate the bot's language via Google to "+langUpdate, 30)
//				err := d.storage.Dictionary.TranslateViaGoogle(langUpdate)
//				if err != nil {
//					d.log.Info(config.CorpName + " " + langUpdate)
//					d.log.ErrorErr(err)
//					d.SendChannelDelSecond(m.ChannelID, "failed to translate to "+langUpdate, 30)
//				} else {
//					go d.updateLanguage(langUpdate, config, m)
//				}
//			}
//			return true
//		}
//	}
//	return false
//}
//func (d *Discord) updateLanguage(langUpdate string, config models.CorporationConfig, m *discordgo.MessageCreate) {
//	go d.DeleteMesageSecond(m.ChannelID, m.ID, 30)
//	if config.MesidDsHelp != "" {
//		go d.DeleteMessage(config.DsChannel, config.MesidDsHelp)
//	}
//	config.Country = langUpdate
//	d.corpConfigRS[config.CorpName] = config
//	config.MesidDsHelp = d.Hhelp1(config.DsChannel, langUpdate)
//
//	d.corpConfigRS[config.CorpName] = config
//	d.storage.ConfigRs.AutoHelpUpdateMesid(config)
//	go d.SendChannelDelSecond(m.ChannelID, d.getLanguage(config.Country, "language_switched_to"), 20)
//	d.log.Info(fmt.Sprintf("замена языка в %s на %s", config.CorpName, config.Country))
//}
//
//func (d *Discord) CleanOldMessage(m *discordgo.MessageCreate) bool {
//	re := regexp.MustCompile(`^\.очистка (\d{1,2}|100)`)
//	matches := re.FindStringSubmatch(m.Content)
//	if len(matches) > 0 {
//		fmt.Println("limitMessage " + matches[1])
//		d.CleanOldMessageChannel(m.ChannelID, matches[1])
//		return true
//	}
//	return false
//}

func (d *Discord) ifPrefixPoint(m *discordgo.MessageCreate) {
	in := models.InMessage{
		Mtext:       m.Content,
		Tip:         "ds",
		Name:        m.Author.String(),
		NameMention: m.Author.Mention(),
		Ds: struct {
			Mesid   string
			Nameid  string
			Guildid string
			Avatar  string
		}{Mesid: m.ID, Nameid: m.Author.ID, Guildid: m.GuildID, Avatar: m.Author.AvatarURL("")},
		Config: models.CorporationConfig{
			CorpName:  d.GuildChatName(m.ChannelID, m.GuildID),
			DsChannel: m.ChannelID,
			Guildid:   m.GuildID,
		},
		Option: models.Option{
			InClient: true,
		},
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
			time.Sleep(3 * time.Second)
			d.storage.ReloadDbArray()
			d.bridgeConfig = d.storage.BridgeConfigs
		}()

	}()

}
