package bot2

import (
	"rs/models"
)

func (b *Bot) removeUsersFromOtherQueues(u []models.QueueActive) {
	var userIds []string
	for _, user := range u {
		userIds = append(userIds, user.Data.UserID)
		if b.storage.CountUserIdQueue(user.Data.UserID) > 0 {
			go b.elseTrue(user.Data.UserID)
		}
	}

	//check Artzor bot queue
	b.otherQueue.NeedRemoveOtherQueue(userIds)
}

func (b *Bot) elseTrue(userid string) { //удаляем игрока с очереди
	tt, _ := b.storage.GetActiveQueueByUser(userid)
	for _, t := range tt {
		ok, config := b.CheckCorpNameConfig(t.Data.CorporationUuid)
		if ok {
			//var text string
			//after, drs := strings.CutPrefix(t.Data.LvlRS, "drs")
			//if drs {
			//	text = after
			//}
			//afterRs, rs := strings.CutPrefix(t.Data.LvlRS, "rs")
			//if rs {
			//	text = afterRs
			//}

			in := &models.InMessageV2{
				//Text:        text + "-",
				Tip:         t.Data.Tip,
				Username:    t.Data.Name,
				UserId:      t.Data.UserID,
				NameMention: t.Data.Mention,

				Config:  config,
				Options: []string{models.OptionElseTrue, models.OptionUpdate},
			}
			for _, i := range config.Channels {
				if i.Game != nil && i.Game.Alias == t.Data.Alias || i.GuildName == t.Data.Alias {
					in.Messenger = *i
				}
			}

			rs1 := models.NewRs()
			rs1.SetLevelRsOrDrs(in, t.Data.LvlRS)
			b.RsMinus(in, rs1)
		}
	}
}
