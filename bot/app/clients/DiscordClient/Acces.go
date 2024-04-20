package DiscordClient

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"kz_bot/models"
	"os"
	"regexp"
	"strings"
	"time"
)

//nujno sdelat lang

func (d *Discord) AccesChatDS(m *discordgo.MessageCreate) {
	after, res := strings.CutPrefix(m.Content, ".")
	if res {
		switch after {
		case "add":
			go d.DeleteMesageSecond(m.ChannelID, m.ID, 10)
			d.accessAddChannelDs(m.ChannelID, m.GuildID, "en")
		case "добавить":
			go d.DeleteMesageSecond(m.ChannelID, m.ID, 10)
			d.accessAddChannelDs(m.ChannelID, m.GuildID, "ru")
		case "додати":
			go d.DeleteMesageSecond(m.ChannelID, m.ID, 10)
			d.accessAddChannelDs(m.ChannelID, m.GuildID, "ua")
		case "del", "удалить", "видалити":
			go d.DeleteMesageSecond(m.ChannelID, m.ID, 10)
			d.accessDelChannelDs(m.ChannelID, m.GuildID)
		case "паника":
			d.log.Panic("перезагрузка по требованию")
			os.Exit(1)
		case "removeCommand":
			d.removeCommand("")
			go d.ready()
		case "мес":
			d.DeleteMessage(m.ChannelID, m.ID)
			d.mes()

		default:
			if d.CleanOldMessage(m) {
				return
			}
			if d.setLang(m) {
				return
			}
		}
	}
}
func (d *Discord) mes() {

}
func (d *Discord) accessAddChannelDs(chatid, guildid, lang string) { // внесение в дб и добавление в масив
	ok, _ := d.CheckChannelConfigDS(chatid)
	if ok {
		go d.SendChannelDelSecond(chatid, d.getLanguage(lang, "info_activation_not_required"), 30)
	} else {
		chatName := d.GuildChatName(chatid, guildid)
		d.log.Info("новая активация корпорации " + chatName)
		d.AddDsCorpConfig(chatName, chatid, guildid, lang)
		go d.SendChannelDelSecond(chatid, d.getLanguage(lang, "tranks_for_activation"), 10)

	}
}
func (d *Discord) accessDelChannelDs(chatid, guildid string) { //удаление с бд и масива для блокировки
	ok, config := d.CheckChannelConfigDS(chatid)
	d.DeleteMessage(chatid, config.MesidDsHelp)
	if !ok {
		go d.SendChannelDelSecond(chatid, d.getLanguage("ru", "channel_not_connected"), 60)
	} else {
		d.SendChannelDelSecond(chatid, d.getLang(chatid, "you_disabled_bot_functions"), 60)
		d.storage.ConfigRs.DeleteConfigRs(config)
		d.storage.ReloadDbArray()
		d.corpConfigRS = d.storage.CorpConfigRS
		d.log.Info("отключение корпорации " + d.GuildChatName(chatid, guildid))
	}
}

func (d *Discord) CleanOldMessage(m *discordgo.MessageCreate) bool {
	re := regexp.MustCompile(`^\.очистка (\d{1,2}|100)`)
	matches := re.FindStringSubmatch(m.Content)
	if len(matches) > 0 {
		fmt.Println("limitMessage " + matches[1])
		d.CleanOldMessageChannel(m.ChannelID, matches[1])
		return true
	}
	return false
}
func (d *Discord) setLang(m *discordgo.MessageCreate) bool {
	re := regexp.MustCompile(`^\.set lang (\w{2})$`)
	matches := re.FindStringSubmatch(m.Content)
	if len(matches) > 0 {
		langUpdate := matches[1]
		ok, config := d.CheckChannelConfigDS(m.ChannelID)
		if ok {
			if d.storage.Dictionary.CheckTranslateLanguage(langUpdate) {
				go d.updateLanguage(langUpdate, config, m)
			} else {
				d.SendChannelDelSecond(m.ChannelID, "please wait, I'm trying to translate the bot's language via Google to "+langUpdate, 30)
				err := d.storage.Dictionary.TranslateViaGoogle(langUpdate)
				if err != nil {
					d.log.Info(config.CorpName + " " + langUpdate)
					d.log.ErrorErr(err)
					d.SendChannelDelSecond(m.ChannelID, "failed to translate to "+langUpdate, 30)
				} else {
					go d.updateLanguage(langUpdate, config, m)
				}
			}
			return true
		}
	}
	return false
}
func (d *Discord) CleanRsBotOtherMessage() {
	defer func() {
		if r := recover(); r != nil {
			d.log.Info(fmt.Sprintf("recover() %+v", r))
		}
	}()
	for _, config := range d.corpConfigRS {
		if config.DsChannel != "" {
			channelMessages, err := d.S.ChannelMessages(config.DsChannel, 100, "", "", "")
			if err != nil {
				restErr, _ := err.(*discordgo.RESTError)
				if restErr.Message != nil && restErr.Message.Code == discordgo.ErrCodeUnknownChannel {
					d.log.Info("нужно сделать удаление этого канала : " + config.CorpName)
				} else {
					d.log.ErrorErr(err)
				}
				continue
			}
			if len(channelMessages) > 0 {
				t := time.Now().Unix()
				for _, message := range channelMessages {
					if message.Author.String() != "Rs_bot#9945" && message.Author.String() != "КзБот#0000" {
						if t-message.Timestamp.Unix() < 1209600 && t-message.Timestamp.Unix() > 180 {
							if message.Content == "" || !strings.HasPrefix(message.Content, ".") {
								errd := d.S.ChannelMessageDelete(config.DsChannel, message.ID)
								if errd != nil {
									restErr, _ := err.(*discordgo.RESTError)
									if restErr.Message != nil && restErr.Message.Code == discordgo.ErrCodeMissingPermissions {
										d.log.Info("ошибка удаления ErrCodeMissingPermissions : " + config.CorpName)
									} else {
										d.log.ErrorErr(err)
									}
									d.log.InfoStruct("не удалено", message)
								}

							}
							//fmt.Printf(" message Author: %s Content: %+v\n", message.Author.String(), message.Content)
						}
						if t-message.Timestamp.Unix() < 1209600 && t-message.Timestamp.Unix() > 260000 {
							if strings.HasPrefix(message.Content, ".") {
								_ = d.S.ChannelMessageDelete(config.DsChannel, message.ID)

							}
						}
					}
				}
				fmt.Println("clean OK " + config.CorpName)
			}
		}
	}
	fmt.Println("clean OK")
}
func (d *Discord) updateLanguage(langUpdate string, config models.CorporationConfig, m *discordgo.MessageCreate) {
	go d.DeleteMesageSecond(m.ChannelID, m.ID, 30)
	if config.MesidDsHelp != "" {
		go d.DeleteMessage(config.DsChannel, config.MesidDsHelp)
	}
	config.Country = langUpdate
	d.corpConfigRS[config.CorpName] = config
	config.MesidDsHelp = d.hhelp1(config.DsChannel, langUpdate)

	d.corpConfigRS[config.CorpName] = config
	d.storage.ConfigRs.AutoHelpUpdateMesid(config)
	go d.SendChannelDelSecond(m.ChannelID, d.getLanguage(config.Country, "language_switched_to"), 20)
	d.log.Info(fmt.Sprintf("замена языка в %s на %s", config.CorpName, config.Country))
}
