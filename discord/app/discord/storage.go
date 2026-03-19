package DiscordClient

import (
	"time"

	"github.com/mentalisit/restapi/models"
)

// BridgeCheckChannelConfigDS bridge
func (d *Discord) BridgeCheckChannelConfigDS(ChatId string) (bool, models.Bridge2Config) {
	//if len(d.bridgeConfig) == 0 || d.bridgeConfigUpdateTime+300 < time.Now().Unix() {
	//	bridgeConfig := d.storage.Db.DBReadBridgeConfig()
	//	d.bridgeConfig = bridgeConfig
	//	d.bridgeConfigUpdateTime = time.Now().Unix()
	//}
	//
	//for _, config := range d.bridgeConfig {
	//	if config.Channel["ds"] != nil {
	//		for _, channelD := range config.Channel["ds"] {
	//			if channelD.ChannelId == ChatId {
	//				return true, config
	//			}
	//		}
	//	}
	//}
	//return false, models.Bridge2Config{}
	return d.storage.Db.CheckBridgeChannel(ChatId)
}

// CheckChannelConfigDS RsConfig
func (d *Discord) CheckChannelConfigDS(chatid string) (channelGood bool, config models.CorporationConfig) {
	//conf, err := d.storage.Db.ReadConfigForDsChannel(chatid)
	//if err != nil {
	//	slog.Error(err.Error())
	//}
	//if conf.DsChannel == chatid {
	//	return true, conf
	//}
	//return false, models.CorporationConfig{}
	return d.storage.Db.CheckKzChannel(chatid)
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

func (d *Discord) checkChannelConfig2(chatid string) (channelGood bool, config models.CorporationConfigV2) {
	//conf := d.storage.Db.ReadConfigV2ByChannelId(chatid)
	//if conf != nil {
	//	return true, *conf
	//}
	//return false, models.CorporationConfigV2{
	//	Channels: make(models.ChannelsMap),
	//}
	return d.storage.Db.CheckRsChannel(chatid)
}
