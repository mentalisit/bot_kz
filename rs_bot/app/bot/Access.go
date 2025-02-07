package bot

import (
	"fmt"
	"regexp"
	"rs/models"
	"strconv"
	"strings"
)

func (b *Bot) accessChat(in models.InMessage) {
	after, res := strings.CutPrefix(in.Mtext, ".")
	if res {
		switch after {
		case "add":
			b.accessAddChannel(in, "en")
		case "добавить":
			b.accessAddChannel(in, "ru")
		case "додати":
			b.accessAddChannel(in, "ua")
		case "del":
			b.accessDelChannel(in, "en")
		case "удалить":
			b.accessDelChannel(in, "ru")
		case "видалити":
			b.accessDelChannel(in, "ua")
		case "паника":
			b.log.Panic("перезагрузка по требованию")
		default:
			if b.SendALLChannel(in) {
				return
			}
			if b.setLang(in) {
				return
			}
			if b.cleanOldMessage(in) {
				return
			}
		}
	}
}
func (b *Bot) accessAddChannel(in models.InMessage, lang string) {
	go b.iftipdelete(in)
	ok, _ := b.checkConfig(in)
	if ok {
		go b.ifTipSendTextDelSecond(in, b.getLanguageText(lang, "info_activation_not_required"), 20)
	} else {
		c := in.Config
		c.Country = lang

		b.log.Info(c.CorpName + " Добавлена в конфиг корпораций ")
		go b.ifTipSendTextDelSecond(in, b.getLanguageText(lang, "tranks_for_activation"), 60)

		if c.DsChannel != "" {
			c = b.sendHelpDs(c, true)
		}
		if c.TgChannel != "" {
			c = b.sendHelpTg(c, true)
		}
		fmt.Printf("accessAddChannel in %+v \n", in)
		fmt.Printf("accessAddChannel conf %+v \n", c)

		b.storage.ConfigRs.InsertConfigRs(c)

	}
}
func (b *Bot) accessDelChannel(in models.InMessage, lang string) {
	go b.iftipdelete(in)
	ok, config := b.checkConfig(in)
	if !ok {
		go b.ifTipSendTextDelSecond(in, b.getLanguageText(lang, "channel_not_connected"), 60)
	} else {
		b.storage.ConfigRs.DeleteConfigRs(config)
		b.log.Info("отключение корпорации " + config.CorpName)
		go b.ifTipSendTextDelSecond(in, b.getLanguageText(lang, "you_disabled_bot_functions"), 60)
		if config.MesidDsHelp != "" {
			go b.client.Ds.DeleteMessage(config.DsChannel, config.MesidDsHelp)
		}
		if config.MesidTgHelp != "" {
			mid, _ := strconv.Atoi(config.MesidTgHelp)
			if mid != 0 {
				go b.client.Tg.DelMessage(config.TgChannel, mid)
			}

		}
	}
}
func (b *Bot) setLang(in models.InMessage) bool {
	re := regexp.MustCompile(`^\.set lang (\w{2})$`)
	matches := re.FindStringSubmatch(in.Mtext)
	if len(matches) > 0 {
		langUpdate := matches[1]
		ok, config := b.checkConfig(in)
		if ok {
			if b.storage.Dictionary.CheckTranslateLanguage(langUpdate) {
				go b.updateLanguage(in, langUpdate, config)
			} else {
				b.ifTipSendTextDelSecond(in, "please wait, I'm trying to translate the bot's language via Google to "+langUpdate, 30)
				err := b.storage.Dictionary.TranslateViaGoogle(langUpdate)
				if err != nil {
					b.log.Info(config.CorpName + " " + langUpdate)
					b.log.ErrorErr(err)
					b.ifTipSendTextDelSecond(in, "failed to translate to "+langUpdate, 30)
				} else {
					go b.updateLanguage(in, langUpdate, config)
				}
			}
			return true
		}
	}

	return false
}

func (b *Bot) updateLanguage(in models.InMessage, langUpdate string, config models.CorporationConfig) {
	go b.iftipdelete(in)
	if config.MesidDsHelp != "" {
		go b.client.Ds.DeleteMessage(config.DsChannel, config.MesidDsHelp)
	}
	config.Country = langUpdate

	text := fmt.Sprintf("%s \n\n%s", b.getLanguageText(langUpdate, "info_bot_delete_msg"), b.getLanguageText(langUpdate, "info_help_text"))
	config.MesidDsHelp = b.client.Ds.SendHelp(config.DsChannel, b.getLanguageText(langUpdate, "information"), text, "", false)

	b.storage.ConfigRs.UpdateConfigRs(config)

	go b.ifTipSendTextDelSecond(in, b.getLanguageText(config.Country, "language_switched_to"), 20)
	b.log.Info(fmt.Sprintf("замена языка в %s на %s", config.CorpName, config.Country))
}
func (b *Bot) cleanOldMessage(in models.InMessage) bool {
	if in.Tip != "ds" {
		return false
	}
	re := regexp.MustCompile(`^\.очистка (\d{1,2}|100)`)
	matches := re.FindStringSubmatch(in.Mtext)
	if len(matches) > 0 {
		fmt.Println("limitMessage " + matches[1])
		b.client.Ds.CleanOldMessageChannel(in.Config.DsChannel, matches[1])
		return true
	}
	return false
}
