package telegram

import "telegram/models"

func (t *Telegram) bridgeCheckChannelConfigTg(mId string) (bool, models.BridgeConfig) {
	for _, config := range t.bridgeConfig {
		for _, channelD := range config.ChannelTg {
			if channelD.ChannelId == mId {
				return true, config
			}
		}
	}
	return false, models.BridgeConfig{}
}

func (t *Telegram) checkChannelConfigTG(chatid string) (channelGood bool, config models.CorporationConfig) {
	for _, corpporationConfig := range t.corpConfigRS {
		if corpporationConfig.TgChannel == chatid {
			return true, corpporationConfig
		}
	}
	return false, models.CorporationConfig{}
}
