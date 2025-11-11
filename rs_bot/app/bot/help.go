package bot

import (
	"fmt"
	"rs/models"
	"rs/pkg/utils"
	"strings"
	"time"
)

func (b *Bot) AutoHelp() {
	tm := time.Now().UTC()
	mtime := tm.Format("15:04")
	if mtime == "12:00" {
		go b.client.Ds.CleanRsBotOtherMessage()
	}
	if mtime == "00:00" {
		go b.EventAutoStart()
	}
	if tm.Minute() == 0 {
		for _, s := range b.storage.ConfigRs.ReadConfigRs() {
			configTemp := s
			if s.DsChannel != "" {
				configTemp = b.sendHelpDs(configTemp, false)
			}
			if s.TgChannel != "" {
				configTemp = b.sendHelpTg(configTemp, false)
				if configTemp.MesidTgHelp == "Бот отключен" && configTemp.DsChannel != "" {
					b.client.Ds.Send(configTemp.DsChannel, "Бот отключен, для активации бота напишите команду \n.добавить")
				}
			}
			if s != configTemp {
				b.storage.ConfigRs.UpdateConfigRs(configTemp)
				in := models.InMessage{
					Config: configTemp,
					Opt:    []string{models.OptionUpdateAutoHelp},
				}
				b.Inbox <- in
				time.Sleep(1 * time.Second)
			}
		}
	}

	utils.PrintGoroutine(b.log)
	time.Sleep(time.Minute)
}
func (b *Bot) hhelp(in models.InMessage) {
	b.iftipdelete(in)
	go func() {
		ch := utils.WaitForMessage("hhelp")
		defer close(ch)
		var conf models.CorporationConfig
		if in.Tip == ds {
			conf = b.sendHelpDs(in.Config, true)
		} else if in.Tip == tg {
			conf = b.sendHelpTg(in.Config, true)
		}
		b.storage.ConfigRs.UpdateConfigRs(conf)
		//b.storage.ReloadDbArray()
		//b.configCorp = b.storage.CorpConfigRS
	}()

}

//	func (b *Bot) hhelp2(in models.InMessage) {
//		b.iftipdelete(in)
//		if in.Tip == ds {
//			go b.client.Ds.Help(in.Config.DsChannel, in.Config.Country)
//		} else if in.Tip == tg {
//			go func() {
//				text := fmt.Sprintf("%s\n%s ", b.getLanguageText(in.Config.Country, "information"), b.getLanguageText(in.Config.Country, "info_help_text"))
//				mid := b.client.Tg.SendHelp(in.Config.TgChannel, text, []string{})
//				b.client.Tg.DelMessageSecond(in.Config.TgChannel, strconv.Itoa(mid), 180)
//			}()
//			//go b.client.Tg.Help(in.Config.TgChannel, in.Config.Country)
//		}
//	}

func (b *Bot) sendHelpDs(c models.CorporationConfig, ifUser bool) models.CorporationConfig {
	text := b.getLanguageText(c.Country, "info_help_text3")

	mId := b.client.Ds.SendHelp(
		c.DsChannel,
		b.getLanguageText(c.Country, "info_bot_delete_msg"),
		text, c.MesidDsHelp, ifUser)

	if mId == "" {
		b.log.InfoStruct("sendHelpDs", c)
	} else {
		c.MesidDsHelp = mId
	}

	return c
}
func (b *Bot) sendHelpTg(c models.CorporationConfig, ifUser bool) models.CorporationConfig {
	text := b.getLanguageText(c.Country, "info_help_text3")

	if IsThisTopicTG(c.TgChannel) {
		text = fmt.Sprintf("%s\n\n%s",
			b.getLanguageText(c.Country, "info_bot_delete_msg"),
			b.getLanguageText(c.Country, "info_help_text3"))
	}

	//if c.TgChannel == "-1002298028181/4" {
	//	text = fmt.Sprintf("%s\n%s\n%s",
	//		b.getLanguageText(c.Country, "information"),
	//		b.getLanguageText(c.Country, "info_bot_delete_msg"),
	//		b.getLanguageText(c.Country, "info_help_text2"))
	//}

	mId := b.client.Tg.SendHelp(c.TgChannel, text, c.MesidTgHelp, ifUser)

	if mId == "" {
		b.log.InfoStruct("sendHelpTg", c)
	} else {
		c.MesidTgHelp = mId
	}

	return c
}

func IsThisTopicTG(tgchannel string) bool {
	split := strings.Split(tgchannel, "/")
	if split[1] != "0" {
		return true
	}
	return false
}
