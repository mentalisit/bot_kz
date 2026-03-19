package wa

import (
	"whatsapp/models"
)

// BridgeCheckChannelConfigWA bridge
func (b *Whatsapp) BridgeCheckChannelConfigWA(ChatId string) (bool, models.Bridge2Config) {
	return b.Storage.Db.CheckBridgeChannel(ChatId)
}

func (b *Whatsapp) checkChannelConfig2(chatId string) (channelGood bool, config models.CorporationConfigV2) {
	return b.Storage.Db.CheckRsChannel(chatId)
}
