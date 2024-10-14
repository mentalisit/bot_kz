package bot

import "kz_bot/models"

func (b *Bot) CheckCorpNameConfig(corpName string) (bool, models.CorporationConfig) {
	for _, config := range b.configCorp {
		if config.CorpName == corpName {
			return true, config
		}
	}
	return false, models.CorporationConfig{}
}
func (b *Bot) checkConfig(in models.InMessage) (bool, models.CorporationConfig) {
	for corpName, config := range b.configCorp {
		if corpName != "" && corpName == in.Config.CorpName {
			return true, config
		} else if config.DsChannel != "" && config.DsChannel == in.Config.DsChannel {
			return true, config
		} else if config.TgChannel != "" && config.TgChannel == in.Config.TgChannel {
			return true, config
		}
	}
	return false, models.CorporationConfig{}
}
