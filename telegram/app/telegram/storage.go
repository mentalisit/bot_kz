package telegram

import (
	"telegram/models"
	"time"
)

func (t *Telegram) bridgeCheckChannelConfigTg(channelId string) (bool, models.Bridge2Config) {
	if len(t.bridgeConfig) == 0 || t.bridgeConfigUpdateTime+300 < time.Now().Unix() {
		t.bridgeConfig = t.Storage.Db.DBReadBridgeConfig()
		t.bridgeConfigUpdateTime = time.Now().Unix()
	}

	for _, config := range t.bridgeConfig {
		if config.Channel["tg"] != nil {
			for _, channelD := range config.Channel["tg"] {
				if channelD.ChannelId == channelId {
					return true, config
				}
			}
		}
	}

	return false, models.Bridge2Config{}
}

func (t *Telegram) checkChannelConfigTG(chatid string) (channelGood bool, config models.CorporationConfig) {
	//for _, corpporationConfig := range t.corpConfigRS {
	//	if corpporationConfig.TgChannel == chatid {
	//		return true, corpporationConfig
	//	}
	//}
	conf := t.Storage.Db.ReadConfigForTgChannel(chatid)
	if conf.TgChannel == chatid {
		return true, conf
	}
	return false, models.CorporationConfig{}
}
