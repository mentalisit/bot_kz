package bot

import (
	"fmt"
	"kz_bot/models"
	"regexp"
	"strings"
)

func (b *Bot) accessChat() {
	after, res := strings.CutPrefix(b.in.Mtext, ".")
	if res {
		switch after {
		case "add":
			go b.iftipdelete()
			b.accessAddChannel("en")
		case "добавить":
			go b.iftipdelete()
			b.accessAddChannel("ru")
		case "додати":
			go b.iftipdelete()
			b.accessAddChannel("ua")
		case "del":
			go b.iftipdelete()
			b.accessDelChannel("en")
		case "удалить":
			go b.iftipdelete()
			b.accessDelChannel("ru")
		case "видалити":
			go b.iftipdelete()
			b.accessDelChannel("ua")
		case "паника":
			b.log.Panic("перезагрузка по требованию")
		default:
			if b.setLang() {
				return
			}
			if b.cleanOldMessage() {
				return
			}
		}
	}
}
func (b *Bot) accessAddChannel(lang string) {
	ok, _ := b.checkConfig()
	if ok {
		go b.ifTipSendTextDelSecond(b.getLanguageText(lang, "info_activation_not_required"), 20)
	} else {
		c := b.in.Config
		c.Country = lang

		b.storage.ConfigRs.InsertConfigRs(c)
		b.configCorp[c.CorpName] = c
		b.log.Info(c.CorpName + " Добавлена в конфиг корпораций ")
		go b.ifTipSendTextDelSecond(b.getLanguageText(lang, "tranks_for_activation"), 60)
		if c.DsChannel != "" {
			b.client.Ds.Hhelp1(c.DsChannel, lang)
		}
		if c.TgChannel != "" {
			b.client.Tg.Help(c.TgChannel, lang)
		}
	}
}
func (b *Bot) accessDelChannel(lang string) {
	ok, config := b.checkConfig()
	if !ok {
		go b.ifTipSendTextDelSecond(b.getLanguageText(lang, "channel_not_connected"), 60)
	} else {
		b.storage.ConfigRs.DeleteConfigRs(config)
		b.storage.ReloadDbArray()
		b.configCorp = b.storage.CorpConfigRS
		b.log.Info("отключение корпорации " + config.CorpName)
		go b.ifTipSendTextDelSecond(b.getLanguageText(lang, "you_disabled_bot_functions"), 60)
	}
}
func (b *Bot) setLang() bool {
	re := regexp.MustCompile(`^\.set lang (\w{2})$`)
	matches := re.FindStringSubmatch(b.in.Mtext)
	if len(matches) > 0 {
		langUpdate := matches[1]
		ok, config := b.checkConfig()
		if ok {
			if b.storage.Dictionary.CheckTranslateLanguage(langUpdate) {
				go b.updateLanguage(langUpdate, config)
			} else {
				b.ifTipSendTextDelSecond("please wait, I'm trying to translate the bot's language via Google to "+langUpdate, 30)
				err := b.storage.Dictionary.TranslateViaGoogle(langUpdate)
				if err != nil {
					b.log.Info(config.CorpName + " " + langUpdate)
					b.log.ErrorErr(err)
					b.ifTipSendTextDelSecond("failed to translate to "+langUpdate, 30)
				} else {
					go b.updateLanguage(langUpdate, config)
				}
			}
			return true
		}
	}
	return false
}
func (b *Bot) checkConfig() (bool, models.CorporationConfig) {
	for corpName, config := range b.configCorp {
		if corpName != "" && corpName == b.in.Config.CorpName {
			return true, config
		} else if config.DsChannel != "" && config.DsChannel == b.in.Config.DsChannel {
			return true, config
		} else if config.TgChannel != "" && config.TgChannel == b.in.Config.TgChannel {
			return true, config
		}
	}
	return false, models.CorporationConfig{}
}

func (b *Bot) updateLanguage(langUpdate string, config models.CorporationConfig) {
	go b.iftipdelete()
	if config.MesidDsHelp != "" {
		go b.client.Ds.DeleteMessage(config.DsChannel, config.MesidDsHelp)
	}
	config.Country = langUpdate
	config.MesidDsHelp = b.client.Ds.Hhelp1(config.DsChannel, langUpdate)

	b.configCorp[config.CorpName] = config

	b.storage.ConfigRs.AutoHelpUpdateMesid(config)

	go b.ifTipSendTextDelSecond(b.getLanguageText(config.Country, "language_switched_to"), 20)
	b.log.Info(fmt.Sprintf("замена языка в %s на %s", config.CorpName, config.Country))
}
func (b *Bot) cleanOldMessage() bool {
	if b.in.Tip != "ds" {
		return false
	}
	re := regexp.MustCompile(`^\.очистка (\d{1,2}|100)`)
	matches := re.FindStringSubmatch(b.in.Mtext)
	if len(matches) > 0 {
		fmt.Println("limitMessage " + matches[1])
		b.client.Ds.CleanOldMessageChannel(b.in.Config.DsChannel, matches[1])
		return true
	}
	return false
}
