package bot2

import (
	"fmt"
	"rs/models"
)

func (b *Bot) RsMinus(in *models.InMessageV2, rs *models.Rs) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.deleteInMessage(in)

	active, err := b.storage.CheckUserInActiveQueue(in.UserId, rs.RsTypeLevel, in.Config.Uid)
	if err != nil {
		b.log.ErrorErr(err)
		return
	}

	channelInfo := in.Config.Channels[in.Messenger.ChannelId]
	if channelInfo == nil {
		return
	}

	if !active {
		b.SendMentionTextDel(in, b.getTextForInfo(channelInfo, "you_out_of_queue"))
	} else {
		//чтение айди очечреди
		u, _ := b.storage.GetActiveQueueByCorpAndLevel(in.Config.Uid, rs.RsTypeLevel)

		messages, errId := b.storage.ReadQueueMessages(in.Config.Uid, rs.RsTypeLevel)
		if errId == nil {
			rs.QueueMessages = messages
		}

		//удаление с БД
		err = b.storage.DeleteActiveQueueByUserAndLevel(in.UserId, in.Config.Uid, rs.RsTypeLevel)
		if err != nil {
			b.log.ErrorErr(err)
			return
		}

		countQueue := len(u) - 1

		for _, info := range in.Config.Channels {
			userLeft := fmt.Sprintf("%s %s", in.Username, b.getTextForInfo(info, "left_queue"))
			b.SendInfoChannelsMessage(info, userLeft, 10)
			if countQueue == 0 {
				rs.Info = info
				text := fmt.Sprintf("%s %s.", rs.GetTitle(b.getLanguageText), b.getTextForInfo(rs.Info, "was_deleted"))
				b.SendInfoChannelsMessage(info, text, 10)
				// Delete all queue messages when queue becomes empty
				if len(rs.QueueMessages) > 0 {
					b.deleteQueueActiveMessage(rs.QueueMessages)
					b.storage.DeleteQueueMessages(in.Config.Uid, rs.RsTypeLevel)
				}
			}
		}
		if countQueue > 0 {
			// Delete old queue messages before updating
			//if len(rs.QueueMessages) > 0 {
			//	b.deleteQueueActiveMessage(rs.QueueMessages)
			//}
			// Update queue with new user list
			b.QueueLevel(in, rs)
		}
	}
}
