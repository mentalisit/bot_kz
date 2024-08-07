package DiscordClient

import (
	"kz_bot/models"
	"time"
)

// BridgeCheckChannelConfigDS bridge
func (d *Discord) BridgeCheckChannelConfigDS(ChatId string) (bool, models.BridgeConfig) {
	for _, config := range *d.bridgeConfig {
		for _, channelD := range config.ChannelDs {
			if channelD.ChannelId == ChatId {
				return true, config
			}
		}
	}
	return false, models.BridgeConfig{}
}

// CheckChannelConfigDS RsConfig
func (d *Discord) CheckChannelConfigDS(chatid string) (channelGood bool, config models.CorporationConfig) {
	for _, corpporationConfig := range d.corpConfigRS {
		if corpporationConfig.DsChannel == chatid {
			return true, corpporationConfig
		}
	}
	return false, models.CorporationConfig{}
}

//// AddDsCorpConfig add RsConfig
//func (d *Discord) AddDsCorpConfig(chatName, chatid, guildid, lang string) {
//	c := models.CorporationConfig{
//		CorpName:  chatName,
//		DsChannel: chatid,
//		Country:   lang,
//		Guildid:   guildid,
//	}
//	c.MesidDsHelp = d.HelpChannelUpdate(c)
//	d.storage.ConfigRs.InsertConfigRs(c)
//	d.corpConfigRS[c.CorpName] = c
//	d.log.Info(chatName + " Добавлена в конфиг корпораций ")
//}

func (d *Discord) loadSlashCommand() {
	for {
		if len(d.corpConfigRS) > 0 {
			d.ready()
			break
		}
		time.Sleep(3 * time.Second)
	}
}
