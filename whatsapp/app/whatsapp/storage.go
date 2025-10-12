package wa

import (
	"time"
	"whatsapp/models"
)

// BridgeCheckChannelConfigWA bridge
func (b *Whatsapp) BridgeCheckChannelConfigWA(ChatId string) (bool, models.Bridge2Config) {
	if len(b.bridgeConfig) == 0 || b.bridgeConfigUpdateTime+300 < time.Now().Unix() {
		bridgeConfig := b.Storage.Db.DBReadBridgeConfig()
		b.bridgeConfig = bridgeConfig
		b.bridgeConfigUpdateTime = time.Now().Unix()
	}

	for _, config := range b.bridgeConfig {
		if config.Channel["wa"] != nil {
			for _, channelD := range config.Channel["wa"] {
				if channelD.ChannelId == ChatId {
					return true, config
				}
			}
		}
	}
	return false, models.Bridge2Config{}
}

// CheckChannelConfigWA RsConfig
func (b *Whatsapp) CheckChannelConfigWA(chatId string) (channelGood bool, config models.CorporationConfig) {
	conf, err := b.Storage.Db.ReadConfigForDsChannel(chatId)
	if err != nil {
		b.log.ErrorErr(err)
	}
	if conf.WaChannel == chatId {
		return true, conf
	}
	return false, models.CorporationConfig{}
}
