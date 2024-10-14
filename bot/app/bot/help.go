package bot

import (
	"fmt"
	"kz_bot/models"
	"kz_bot/pkg/utils"
	"strconv"
	"strings"
	"time"
)

func (b *Bot) Autohelp() {
	tm := time.Now()
	mtime := tm.Format("15:04")
	EvenOrOdd, _ := strconv.Atoi((tm.Format("2006-01-02"))[8:])
	if mtime == "12:00" {
		a := b.storage.ConfigRs.ReadConfigRs()
		for _, s := range a {
			if s.DsChannel != "" {
				s = b.sendHelpDs(s)
			}
			if s.Forward && s.TgChannel != "" && EvenOrOdd%2 == 0 {
				text := fmt.Sprintf("%s \n%s", b.getLanguageText(s.Country, "info_bot_delete_msg"), b.getLanguageText(s.Country, "info_help_text"))
				if s.MesidTgHelp != "" {
					mID, _ := strconv.Atoi(s.MesidTgHelp)
					if mID != 0 {
						go b.client.Tg.DelMessage(s.TgChannel, mID)
					}
				}
				s.MesidTgHelp = b.client.Tg.SendHelp(s.TgChannel, strings.Replace(text, "3", "10", 1), s.MesidTgHelp)

			}
			b.storage.ConfigRs.UpdateConfigRs(s)
		}
		time.Sleep(time.Minute)
		go b.client.Ds.CleanRsBotOtherMessage()
	} else if tm.Minute() == 0 {
		go func() {
			a := b.storage.ConfigRs.ReadConfigRs()
			for _, s := range a {
				configTemp := s
				if s.DsChannel != "" {
					configTemp = b.sendHelpDs(configTemp)
				}
				if s.TgChannel != "" {
					split := strings.Split(s.TgChannel, "/")
					if split[1] != "0" {
						configTemp = b.sendHelpTg(configTemp)
					}
				}
				if s != configTemp {
					b.storage.ConfigRs.UpdateConfigRs(configTemp)
					in := models.InMessage{Config: configTemp}
					b.QueueAll(in)
				}
			}
		}()
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
			conf = b.sendHelpDs(in.Config)
		} else if in.Tip == tg {
			conf = b.sendHelpTg(in.Config)
		}
		b.storage.ConfigRs.UpdateConfigRs(conf)
	}()

}

//func (b *Bot) hhelp2(in models.InMessage) {
//	b.iftipdelete(in)
//	if in.Tip == ds {
//		go b.client.Ds.Help(in.Config.DsChannel, in.Config.Country)
//	} else if in.Tip == tg {
//		go func() {
//			text := fmt.Sprintf("%s\n%s ", b.getLanguageText(in.Config.Country, "information"), b.getLanguageText(in.Config.Country, "info_help_text"))
//			mid := b.client.Tg.SendHelp(in.Config.TgChannel, text, []string{})
//			b.client.Tg.DelMessageSecond(in.Config.TgChannel, strconv.Itoa(mid), 180)
//		}()
//		//go b.client.Tg.Help(in.Config.TgChannel, in.Config.Country)
//	}
//}

func (b *Bot) sendHelpDs(c models.CorporationConfig) models.CorporationConfig {
	text := fmt.Sprintf("%s \n\n%s",
		b.getLanguageText(c.Country, "info_bot_delete_msg"),
		b.getLanguageText(c.Country, "info_help_text"))

	c.MesidDsHelp = b.client.DS.SendHelp(
		c.DsChannel,
		b.getLanguageText(c.Country, "information"),
		text, c.MesidDsHelp)

	if c.MesidDsHelp == "" {
		b.log.InfoStruct("sendHelpDs", c)
	}

	return c
}
func (b *Bot) sendHelpTg(c models.CorporationConfig) models.CorporationConfig {
	text := fmt.Sprintf("%s\n%s ",
		b.getLanguageText(c.Country, "information"),
		b.getLanguageText(c.Country, "info_help_text"))
	c.MesidTgHelp = b.client.Tg.SendHelp(c.TgChannel, text, c.MesidTgHelp)
	return c
}
