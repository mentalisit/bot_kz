package bot2

import (
	"fmt"
	"rs/models"
	"sync"
	"time"
)

func (b *Bot) AutoHelp() {
	configV2 := b.storage.ReadConfigV2()
	for _, v2 := range configV2 {
		conf := b.autoSendHelp(v2, false)
		if !b.helpMessageEqual(v2.HelpMessage, conf.HelpMessage) {
			b.storage.UpdateConfigV2HelpMessage(conf)
			time.Sleep(1 * time.Second)
		}
	}
}

// helpMessageEqual сравнивает два HelpMessage (map[string]*Info)
func (b *Bot) helpMessageEqual(a, c models.HelpMessage) bool {
	if len(a) != len(c) {
		return false
	}
	for key, valA := range a {
		valC, exists := c[key]
		if !exists || valA == nil || valC == nil {
			return false
		}
		// Сравниваем поля Info
		if valA.TypeMessenger != valC.TypeMessenger ||
			valA.MessageId != valC.MessageId {
			return false
		}
	}
	return true
}

// autoSendHelp рассылает справку по всем каналам корпорации параллельно
func (b *Bot) autoSendHelp(c models.CorporationConfigV2, ifUser bool) models.CorporationConfigV2 {
	var wg sync.WaitGroup
	var mu sync.Mutex

	for channelID, info := range c.Channels {
		wg.Add(1)
		go func(chID string, chInfo *models.Info) {
			defer wg.Done()
			if chInfo != nil && chInfo.Corp != nil && chInfo.Corp.AutoHelp {
				mId := b.sendHelpGeneric(chID, chInfo, &c, ifUser)
				if mId != "" {
					mu.Lock()
					if c.HelpMessage == nil {
						c.HelpMessage = make(models.HelpMessage)
					}
					if c.HelpMessage[chID] == nil {
						c.HelpMessage[chID] = &models.Info{}
					}
					c.HelpMessage[chID].TypeMessenger = chInfo.TypeMessenger
					c.HelpMessage[chID].MessageId = mId
					mu.Unlock()
				}
			}
		}(channelID, info)
	}
	wg.Wait()
	return c
}

// SendHelpInMessenger отправляет справку в конкретный канал из входящего сообщения
func (b *Bot) SendHelpInMessenger(in *models.InMessageV2) models.CorporationConfigV2 {
	info := in.Config.Channels[in.Messenger.ChannelId]
	mId := b.sendHelpGeneric(in.Messenger.ChannelId, info, &in.Config, true)

	if mId != "" {
		if in.Config.HelpMessage == nil {
			in.Config.HelpMessage = make(models.HelpMessage)
		}
		if in.Config.HelpMessage[in.Messenger.ChannelId] == nil {
			in.Config.HelpMessage[in.Messenger.ChannelId] = &models.Info{}
		}
		in.Config.HelpMessage[in.Messenger.ChannelId].TypeMessenger = in.Messenger.TypeMessenger
		in.Config.HelpMessage[in.Messenger.ChannelId].MessageId = mId
	}
	return in.Config
}

// sendHelpGeneric — внутренняя унифицированная функция для отправки справки
func (b *Bot) sendHelpGeneric(channelID string, info *models.Info, c *models.CorporationConfigV2, ifUser bool) string {
	text := b.getLanguageText(info.Language, "info_help_text3")
	if info.Corp != nil && info.Corp.CustomText && info.Corp.HelpText != "" {
		text = info.Corp.HelpText
	}
	oldMID := ""
	if help, exist := c.HelpMessage[channelID]; exist {
		oldMID = help.MessageId
	}

	var mId string
	var msgDelete string
	if info.Corp != nil && info.Corp.DeleteMessages && info.Corp.DeleteMessagesDelay != 0 {
		msgDelete = fmt.Sprintf(b.getLanguageText(info.Language, "info_bot_delete_msg_new"), info.Corp.DeleteMessagesDelay)
	}

	switch info.TypeMessenger {
	case "ds":
		mId = b.client.Ds.SendHelp(channelID, msgDelete, text, oldMID, ifUser)
	case "tg":
		if info.Corp.DeleteMessages {
			text = fmt.Sprintf("%s\n\n%s", msgDelete, text)
		}
		mId = b.client.Tg.SendHelp(channelID, text, oldMID, ifUser)
	}

	if mId == "" {
		b.log.InfoStruct("sendHelpGeneric_fail", info)
	}
	return mId
}

//func IsThisTopicTG(tgchannel string) bool {
//	split := strings.Split(tgchannel, "/")
//	if len(split) < 2 {
//		return false
//	}
//	return split[1] != "0"
//}
