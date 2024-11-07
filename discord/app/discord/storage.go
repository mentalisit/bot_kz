package DiscordClient

import (
	"discord/models"
	"time"
)

// BridgeCheckChannelConfigDS bridge
func (d *Discord) BridgeCheckChannelConfigDS(ChatId string) (bool, models.BridgeConfig) {
	if len(d.bridgeConfig) == 0 || d.bridgeConfigUpdateTime+300 < time.Now().Unix() {
		d.bridgeConfig = d.storage.Db.DBReadBridgeConfig()
		d.bridgeConfigUpdateTime = time.Now().Unix()
	}

	for _, config := range d.bridgeConfig {
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
	//for _, corpporationConfig := range d.corpConfigRS {
	//	if corpporationConfig.DsChannel == chatid {
	//		return true, corpporationConfig
	//	}
	//}
	conf := d.storage.Db.ReadConfigForDsChannel(chatid)
	if conf.DsChannel == chatid {
		return true, conf
	}
	return false, models.CorporationConfig{}
}

func (d *Discord) loadSlashCommand() {
	for {
		if len(d.storage.Db.ReadConfigRs()) > 0 {
			d.ready()
			break
		}
		time.Sleep(3 * time.Second)
	}
}
