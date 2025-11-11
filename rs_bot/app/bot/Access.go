package bot

import (
	"encoding/json"
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
			if b.logicScoreboardSetting(in) {
				return
			}
			if b.EventStatistic(in) {
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
func (b *Bot) logicScoreboardSetting(in models.InMessage) bool {
	afterScoreboard, found := strings.CutPrefix(in.Mtext, ".scoreboard")
	if found {
		text := ""
		afterWebhook, foundWebhook := strings.CutPrefix(afterScoreboard, " webhook ")
		if afterWebhook == "name" {
			text = "you can't use the 'name', it's not unique"
			foundWebhook = false
		}
		if foundWebhook {
			if in.Tip != ds {
				return false
			}
			scoreboardReadName := b.storage.Scoreboard.ScoreboardReadName(afterWebhook)
			if scoreboardReadName == nil {
				s := models.ScoreboardParams{
					Name: afterWebhook,
				}
				if in.Tip == ds {
					s.ChannelWebhook = in.Config.DsChannel
				}
				b.storage.Scoreboard.ScoreboardInsertParam(s)
				text = "now the bot will wait here for webhooks from the game, connect another channel to display the leaderboard"
			} else {
				text = "this channel is already listened to by a bot to receive webhooks from the game.Name " + scoreboardReadName.Name
				b.log.Info("found " + afterWebhook + " in scoreboard")
			}
		}
		afterHere, fountHere := strings.CutPrefix(afterScoreboard, " here ")
		if afterHere == "name" {
			text = "you can't use the 'name', it's not unique"
			fountHere = false
		}
		if fountHere {
			scoreboard := b.storage.Scoreboard.ScoreboardReadName(afterHere)
			if scoreboard != nil {
				m, str := scoreboard.GetMapOrString()
				if m == nil {
					m = make(map[string]string)
				}
				if str != "" {
					if strings.HasPrefix(str, "-") {
						m["tg"] = str
					} else {
						m["ds"] = str
					}
				}
				if in.Tip == ds {
					m["ds"] = in.Config.DsChannel
				} else if in.Tip == tg {
					m["tg"] = in.Config.TgChannel
				}
				marshal, _ := json.Marshal(m)
				scoreboard.ChannelScoreboardOrMap = string(marshal)
				b.storage.Scoreboard.ScoreboardUpdateParamScoreChannels(*scoreboard)
				text = "now the leaderboard will be displayed here"
			} else {
				b.log.Info("not found " + afterHere + " in scoreboard")
				text = "it is impossible to connect the leaderboard without having a channel of incoming data from the game via webhook"
			}
		}
		if text == "" && !fountHere && !foundWebhook {
			text = "To set up automatic display of red star event leaders, you need to do several things:\n" +
				"1) set up sending webhook to you in the game in a hidden channel in discord\n" +
				"2) execute the command to connect the bot to listening to this channel, come up with a unique name or use the corporation name in the command '.scoreboard webhook name' where name is a unique name that will link the data in the bot.\n" +
				"3) execute the command in the channel open for viewing your corporation where the leaderboard will be displayed '.scoreboard here name'"
		}
		b.iftipdelete(in)
		b.ifTipSendTextDelSecond(in, text, 1800)
		return true
	}
	return false
}
