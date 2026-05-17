package bot2

import "rs/models"

func (b *Bot) checkConfig(in *models.InMessageV2) (bool, models.CorporationConfigV2) {
	if in.Messenger.ChannelId != "" {
		for _, v2 := range b.storage.ReadConfigV2() {
			for channel, _ := range v2.Channels {
				if channel == in.Messenger.ChannelId {
					return true, v2

				}
			}
		}

	}

	return false, models.CorporationConfigV2{}
}

func (b *Bot) CheckCorpNameConfig(corpname string) (bool, models.CorporationConfigV2) {
	for _, v2 := range b.storage.ReadConfigV2() {
		if v2.Uid == corpname {
			return true, v2
		}
	}

	return false, models.CorporationConfigV2{}
}
