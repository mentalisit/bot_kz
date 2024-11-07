package bot

import "rs/models"

func (b *Bot) CheckCorpNameConfig(corpName string) (bool, models.CorporationConfig) {
	//for _, config := range b.storage.ConfigRs.ReadConfigRs() {
	//	if config.CorpName == corpName {
	//		return true, config
	//	}
	//}
	//return false, models.CorporationConfig{}
	conf := b.storage.ConfigRs.ReadConfigForCorpName(corpName)
	if conf.CorpName != "" {
		return true, conf
	}
	return false, models.CorporationConfig{}
}
func (b *Bot) checkConfig(in models.InMessage) (bool, models.CorporationConfig) {
	//for _, config := range b.storage.ConfigRs.ReadConfigRs() {
	//	if config.CorpName == in.Config.CorpName {
	//		return true, config
	//	} else if config.DsChannel != "" && config.DsChannel == in.Config.DsChannel {
	//		return true, config
	//	} else if config.TgChannel != "" && config.TgChannel == in.Config.TgChannel {
	//		return true, config
	//	}
	//}
	if in.Config.CorpName != "" {
		confName := b.storage.ConfigRs.ReadConfigForCorpName(in.Config.CorpName)
		if confName.CorpName == in.Config.CorpName {
			return true, confName
		}
	}

	if in.Config.DsChannel != "" {
		confDs := b.storage.ConfigRs.ReadConfigForDsChannel(in.Config.DsChannel)
		if confDs.DsChannel == in.Config.DsChannel {
			return true, confDs
		}
	}

	if in.Config.TgChannel != "" {
		confTg := b.storage.ConfigRs.ReadConfigForTgChannel(in.Config.TgChannel)
		if confTg.TgChannel == in.Config.TgChannel {
			return true, confTg
		}
	}

	return false, models.CorporationConfig{}
}
