package telegram

import (
	"github.com/mentalisit/restapi/models"
)

func (t *Telegram) bridgeCheckChannelConfigTg(channelId string) (bool, models.Bridge2Config) {
	//if len(t.bridgeConfig) == 0 || t.bridgeConfigUpdateTime+300 < time.Now().Unix() {
	//	t.bridgeConfig = t.Storage.Db.DBReadBridgeConfig()
	//	t.bridgeConfigUpdateTime = time.Now().Unix()
	//}
	//
	//for _, config := range t.bridgeConfig {
	//	if config.Channel["tg"] != nil {
	//		for _, channelD := range config.Channel["tg"] {
	//			if channelD.ChannelId == channelId {
	//				return true, config
	//			}
	//		}
	//	}
	//}
	return t.Storage.Db.CheckBridgeChannel(channelId)

	//return false, models.Bridge2Config{}
}

func (t *Telegram) checkChannelConfigTG(chatid string) (channelGood bool, config models.CorporationConfig) {
	//conf := t.Storage.Db.ReadConfigForTgChannel(chatid)
	//if conf.TgChannel == chatid {
	//	return true, conf
	//}
	//return false, models.CorporationConfig{}
	return t.Storage.Db.CheckKzChannel(chatid)
}

func (t *Telegram) checkChannelConfig2(chatid string) (channelGood bool, config models.CorporationConfigV2) {
	//conf := t.Storage.Db.ReadConfigV2ByChannelId(chatid)
	//if conf != nil {
	//	return true, *conf
	//}
	//return false, models.CorporationConfigV2{}
	return t.Storage.Db.CheckRsChannel(chatid)
}
