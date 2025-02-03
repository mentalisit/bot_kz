package bot

import (
	"fmt"
	"rs/models"
)

//lang ok

func (b *Bot) RsMinus(in models.InMessage) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.iftipdelete(in)

	CountNames, err := b.storage.Count.СountName(in.UserId, in.RsTypeLevel, in.Config.CorpName) //проверяем есть ли игрок в очереди
	if err != nil {
		b.log.ErrorErr(err)
		return
	}
	if CountNames == 0 {
		b.ifTipSendMentionText(in, b.getText(in, "you_out_of_queue"))
	} else if CountNames > 0 {
		//чтение айди очечреди
		u := b.storage.DbFunc.ReadAll(in.RsTypeLevel, in.Config.CorpName)
		//удаление с БД
		b.storage.DbFunc.DeleteQueue(in.UserId, in.RsTypeLevel, in.Config.CorpName)
		//проверяем очередь
		countQueue, err2 := b.storage.Count.CountQueue(in.RsTypeLevel, in.Config.CorpName)
		if err2 != nil {
			b.log.Error(err2.Error())
			return
		}

		darkStar, level := in.TypeRedStar()
		var text string
		if darkStar {
			text = fmt.Sprintf("%s %s.", b.getText(in, "queue_drs")+level, b.getText(in, "was_deleted"))
		} else {
			text = fmt.Sprintf("%s %s.", b.getText(in, "rs_queue")+level, b.getText(in, "was_deleted"))
		}
		userLeft := fmt.Sprintf("%s %s", in.Username, b.getText(in, "left_queue"))

		if in.Config.DsChannel != "" {
			go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, userLeft, 10)
			if countQueue == 0 {
				go b.client.Ds.SendChannelDelSecond(in.Config.DsChannel, text, 10)
				go b.client.Ds.DeleteMessage(in.Config.DsChannel, u.User1.Dsmesid)
			}
		}
		if in.Config.TgChannel != "" {
			go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, userLeft, 10)
			if countQueue == 0 {
				go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text, 10)
				go b.client.Tg.DelMessage(in.Config.TgChannel, u.User1.Tgmesid)
			}
		}
		if countQueue > 0 {
			b.QueueLevel(in)
		}
	}
}
