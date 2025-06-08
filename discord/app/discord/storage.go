package DiscordClient

import (
	"discord/models"
	"log/slog"
	"time"
)

// BridgeCheckChannelConfigDS bridge
func (d *Discord) BridgeCheckChannelConfigDS(ChatId string) (bool, models.BridgeConfig) {
	if len(d.bridgeConfig) == 0 || d.bridgeConfigUpdateTime+300 < time.Now().Unix() {
		bridgeConfig, err := d.storage.Db.DBReadBridgeConfig()
		if err != nil {
			slog.Error(err.Error())
			return false, models.BridgeConfig{}
		}
		d.bridgeConfig = bridgeConfig
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
	conf, err := d.storage.Db.ReadConfigForDsChannel(chatid)
	if err != nil {
		slog.Error(err.Error())
	}
	if conf.DsChannel == chatid {
		return true, conf
	}
	return false, models.CorporationConfig{}
}

func (d *Discord) loadSlashCommand() {
	for {
		configRs, _ := d.storage.Db.ReadConfigRs()
		if configRs != nil && len(configRs) > 0 {
			d.ready()
			break
		}
		time.Sleep(3 * time.Second)
	}
}
