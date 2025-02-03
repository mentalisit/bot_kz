package bot

import (
	"fmt"
	"rs/models"
)

func (b *Bot) CheckCorpNameConfig(corpName string) (bool, models.CorporationConfig) {
	conf := b.storage.ConfigRs.ReadConfigForCorpName(corpName)
	if conf.CorpName != "" {
		return true, conf
	}
	return false, models.CorporationConfig{}
}
func (b *Bot) checkConfig(in models.InMessage) (bool, models.CorporationConfig) {
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

func (b *Bot) insertUserAccount(in models.InMessage) {
	if in.Mtext == "Получить имена без привязки" {
		getAll, err := b.storage.UserAccount.UserAccountGetAll()
		if err != nil {
			b.log.ErrorErr(err)
		}
		if len(getAll) == 0 {
			b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, "err", 60)
			return
		}
		var tt []models.UserAccount
		for _, account := range getAll {
			if len(account.GameId) != 0 {
				if account.DsId == "" && account.TgId == "" {
					tt = append(tt, account)
				}
			}
		}
		text := "Список из игры:"
		for _, account := range tt {
			text = fmt.Sprintf("%s\n%s", text, account.GeneralName)
		}
		b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 600)

		var ttt []models.UserAccount
		for _, account := range getAll {
			if len(account.GameId) == 0 && account.TgId != "" {
				ttt = append(ttt, account)
			}
		}
		text = "Список из телеграм:"
		for _, account := range ttt {
			text = fmt.Sprintf("%s\n%s", text, account.GeneralName)
		}
		b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 600)

		var ttd []models.UserAccount
		for _, account := range getAll {
			if len(account.GameId) == 0 && account.DsId != "" {
				ttd = append(ttd, account)
			}
		}
		text = "Список из дискорд:"
		for _, account := range ttd {
			text = fmt.Sprintf("%s\n%s", text, account.GeneralName)
		}
		b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 600)

		return
	}
	if in.Username != "" && in.UserId != "" && (in.Opt.Contains(models.OptionInClient) || in.Opt.Contains(models.OptionReaction)) {

		u := models.UserAccount{
			GeneralName: in.Username,
		}

		exist := false

		if in.Tip == ds {
			u.DsId = in.UserId
			_, err := b.storage.UserAccount.UserAccountGetByDsUserId(in.UserId)
			if err == nil {
				exist = true
			}
		} else if in.Tip == tg {
			u.TgId = in.UserId
			_, err := b.storage.UserAccount.UserAccountGetByTgUserId(in.UserId)
			if err == nil {
				exist = true
			}

		}
		if !exist {
			err := b.storage.UserAccount.UserAccountInsert(u)
			if err != nil {
				b.log.ErrorErr(err)
			}
		}
	}
}
