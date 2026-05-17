package bot2

import (
	"fmt"
	"rs/models"
	"rs/pkg/utils"
	"strconv"
	"time"

	gt "github.com/bas24/googletranslatefree"
)

const (
	emOK      = "✅"
	emCancel  = "❎"
	emRsStart = "🚀"
)

func (b *Bot) getDsPing(rs *models.Rs) string {
	if rs.Info.TypeMessenger != ds || rs.Info.GuildId == "" {
		return rs.NameRsLevel
	}

	roleMention, err := b.client.Ds.RoleToIdPing(rs.NameRsLevel, rs.Info.GuildId)
	if err != nil {
		fmt.Printf("RoleToIdPing nameRsLevel %s GuildId %s err: %+v\n", rs.NameRsLevel, rs.Info.GuildId, err)
		return rs.NameRsLevel
	}
	return roleMention
}

func (b *Bot) deleteInMessage(in *models.InMessageV2) {
	if !in.Options.Contains(models.OptionReaction) &&
		!in.Options.Contains(models.OptionUpdate) &&
		!in.Options.Contains(models.OptionEdit) {

		if in.Tip == ds {
			go b.client.Ds.DeleteMessage(in.Messenger.ChannelId, in.Messenger.MessageId)
			go b.client.Ds.ChannelTyping(in.Messenger.ChannelId)
		} else if in.Tip == tg {
			go b.client.Tg.ChatTyping(in.Messenger.ChannelId)
			go b.client.Tg.DeleteMessage(in.Messenger.ChannelId, in.Messenger.MessageId)
		} else {
			b.log.Info(fmt.Sprintf("need make for %s \n", in.Tip))
		}
	}
}
func (b *Bot) sendTextAfterDeleteSecond(in *models.InMessageV2, text string, time int) {
	if in.Tip == ds {
		go b.client.Ds.SendChannelDelSecond(in.Messenger.ChannelId, text, time)
	} else if in.Tip == tg {
		go b.client.Tg.SendChannelDelSecond(in.Messenger.ChannelId, text, time)
	} else {
		b.log.Info(fmt.Sprintf("need make for %s \n", in.Tip))
	}
}
func (b *Bot) getLanguageText(lang, key string) string {
	if lang == "uk" {
		lang = "ua"
	}
	return b.Dictionary.GetText(lang, key)
}
func (b *Bot) getTextForInfo(i *models.Info, key string) string {
	return b.Dictionary.GetText(i.Language, key)
}

func (b *Bot) getText(in *models.InMessageV2, key string) string {
	channelInfo := in.Config.Channels[in.Messenger.ChannelId]
	if channelInfo == nil {
		return "KeyNotFound"
	}

	return b.getTextForInfo(channelInfo, key)
}

func (b *Bot) deleteMessage(in *models.InMessageV2, channelId, messageId string) {
	if in.Tip == ds {
		go b.client.Ds.DeleteMessage(channelId, messageId)
	} else if in.Tip == tg {
		go b.client.Tg.DeleteMessage(channelId, messageId)
	} else if in.Tip == wa {
		go b.client.Wa.DeleteMessage(channelId, messageId)
	} else {
		b.log.Info(fmt.Sprintf("need make for %s \n", in.Tip))
	}
}
func (b *Bot) SendMentionTextDel(in *models.InMessageV2, text string) {
	text = fmt.Sprintf("%s %s", in.GetNameMention(), text)
	b.sendTextAfterDeleteSecond(in, text, 10)
}

func (b *Bot) SendText(in *models.InMessageV2, text string) {
	if in.Tip == ds {
		go b.client.Ds.Send(in.Messenger.ChannelId, text)
	} else if in.Tip == tg {
		go b.client.Tg.SendChannel(in.Messenger.ChannelId, text)
	} else {
		b.log.Info(fmt.Sprintf("need make for %s \n", in.Tip))
	}
}

func (b *Bot) SendTextReturnId(in *models.InMessageV2, text string) (id string) {
	if in.Tip == ds {
		return b.client.Ds.Send(in.Messenger.ChannelId, text)
	} else if in.Tip == tg {
		return strconv.Itoa(b.client.Tg.SendChannel(in.Messenger.ChannelId, text))
	} else {
		b.log.Info(fmt.Sprintf("need make for %s \n", in.Tip))
	}
	return id
}

func (b *Bot) SendInfoChannelsMessage(info *models.Info, text string, seconds int) {
	if info.TypeMessenger == ds {
		go b.client.Ds.SendChannelDelSecond(info.ChannelId, text, seconds)
	} else if info.TypeMessenger == tg {
		go b.client.Tg.SendChannelDelSecond(info.ChannelId, text, seconds)
	} else if info.TypeMessenger == wa {
		go b.client.Wa.SendChannelDelSecond(info.ChannelId, text, seconds)
	}
}
func (b *Bot) deleteQueueActiveMessage(u map[string]models.QueueMessages) {
	if len(u) == 0 {
		return
	}
	for ch, message := range u {
		if message.TypeMessenger == ds {
			go b.client.Ds.DeleteMessage(ch, message.MessageID)
		} else if message.TypeMessenger == tg {
			go b.client.Tg.DeleteMessage(ch, message.MessageID)
		} else if message.TypeMessenger == wa {
			go b.client.Wa.DeleteMessage(ch, message.MessageID)
		} else {
			b.log.Info(fmt.Sprintf("need make for %s \n", message.TypeMessenger))
		}
	}
}

// UpdateQueueMessages безопасно добавляет новые сообщения в текущий список сообщений очереди
// и возвращает true, если были внесены какие-либо изменения (новые ключи или измененные значения).
func (b *Bot) UpdateQueueMessages(rs *models.Rs, newMessages map[string]models.QueueMessages) (changed bool) {
	if rs.QueueMessages == nil {
		rs.QueueMessages = make(map[string]models.QueueMessages)
	}

	for ch, newMsg := range newMessages {
		oldMsg, exists := rs.QueueMessages[ch]
		// Если ключа не было ИЛИ старое сообщение отличается от нового
		if !exists || oldMsg.MessageID != newMsg.MessageID || oldMsg.TypeMessenger != newMsg.TypeMessenger {
			rs.QueueMessages[ch] = newMsg
			changed = true
		}
	}
	return changed
}

// createQueueUser создает запись пользователя для очереди userIn
func (b *Bot) createQueueUser(in *models.InMessageV2, rs *models.Rs) *models.QueueActive {
	// Определяем алиас канала для сохранения
	channelAlias := ""
	if in.Config.Channels[in.Messenger.ChannelId].Game != nil && in.Config.Channels[in.Messenger.ChannelId].Game.Alias != "" {
		channelAlias = in.Config.Channels[in.Messenger.ChannelId].Game.Alias
	} else {
		channelAlias = in.Config.Channels[in.Messenger.ChannelId].GuildName
	}

	q := &models.QueueActive{
		Data: models.QueueData{
			CorporationUuid: in.Config.Uid,
			Alias:           channelAlias,
			Name:            in.Username,
			UserID:          in.UserId,
			Mention:         in.NameMention,
			Alt:             rs.AltName,
			Tip:             in.Messenger.TypeMessenger,
			Time:            time.Now().UTC().Format("15:04"),
			Date:            time.Now().UTC().Format("2006-01-02"),
			LvlRS:           rs.RsTypeLevel,
			NumRSName:       rs.NumberName,
			NumRSLevel:      rs.NumberLevel,
			MAcc:            in.MAcc,
		},
		RemainingTime: rs.GetTimeRs(),
	}
	if in.Config.Channels[in.Messenger.ChannelId].Game != nil && in.Config.Channels[in.Messenger.ChannelId].Game.GameCorporation != "" {
		q.Data.GameCorporation = in.Config.Channels[in.Messenger.ChannelId].Game.GameCorporation
		q.Data.GamePercent = b.GetTextPercent(in.Config.Channels[in.Messenger.ChannelId].Game, rs.GetTypeRs())
	}
	return q
}

func (b *Bot) getStringForDsPing(rs models.Rs) string {
	if rs.Info.TypeMessenger != ds {
		return ""
	}

	darkOrRed, level := rs.TypeRedStar()

	// Use dictionary text like in old version
	title := rs.GetTitle(b.Dictionary.GetText)
	rs.NameRsLevel = rs.GetRsNameLevel(b.Dictionary.GetText)
	pingLevel := b.getDsPing(&rs)
	description := fmt.Sprintf("👇 %s <:rs:918545444425072671> %s (%d) ",
		b.getTextForInfo(rs.Info, "wishing_to"), pingLevel, rs.NumberLevel)
	embedFieldName := fmt.Sprintf(" %s %s\n%s %s\n%s %s",
		emOK, b.getTextForInfo(rs.Info, "to_add_to_queue"),
		emCancel, b.getTextForInfo(rs.Info, "to_exit_the_queue"),
		emRsStart, b.getTextForInfo(rs.Info, "forced_start"))
	embedFieldValue := b.getTextForInfo(rs.Info, "data_updated") + ": "

	tgText1 := fmt.Sprintf("%s%s (%d)\n", b.getTextForInfo(rs.Info, "queue_drs"), level, rs.NumberLevel)
	if !darkOrRed {
		tgText1 = fmt.Sprintf("%s%s (%d)\n", b.getTextForInfo(rs.Info, "rs_queue"), level, rs.NumberLevel)
	}

	tgText2 := fmt.Sprintf("\n%s++ - %s", level, b.getTextForInfo(rs.Info, "forced_start"))

	// Return formatted string for DS ping
	return title + "|" + description + "|" + embedFieldName + "|" + embedFieldValue + "|" + tgText1 + "|" + tgText2
}

func (b *Bot) getMapPre(rs *models.Rs) map[string]string {
	n := make(map[string]string)

	darkOrRed, level := rs.TypeRedStar()

	n["lang"] = rs.Info.Language
	n["title"] = rs.GetTitle(b.Dictionary.GetText)
	n["levelRs"] = level
	n["description"] = fmt.Sprintf("👇 %s <:rs:918545444425072671> %s (%d) ",
		b.getTextForInfo(rs.Info, "wishing_to"), n["levelRs"], rs.NumberLevel)
	n["EmbedFieldName"] = fmt.Sprintf(" %s %s\n%s %s\n%s %s",
		emOK, b.getTextForInfo(rs.Info, "to_add_to_queue"),
		emCancel, b.getTextForInfo(rs.Info, "to_exit_the_queue"),
		emRsStart, b.getTextForInfo(rs.Info, "forced_start"))
	n["EmbedFieldValue"] = b.getTextForInfo(rs.Info, "data_updated") + ": "
	n["tgText1"] = fmt.Sprintf("%s%s (%d)\n", b.getTextForInfo(rs.Info, "queue_drs"), level, rs.NumberLevel)
	if !darkOrRed {
		n["tgText1"] = fmt.Sprintf("%s%s (%d)\n", b.getTextForInfo(rs.Info, "rs_queue"), level, rs.NumberLevel)
	}
	n["tgText2"] = fmt.Sprintf("\n%s++ - %s", level, b.getTextForInfo(rs.Info, "forced_start"))

	n["buttonLevel"] = level

	return n
}

func (b *Bot) getMapForDsChannel(rs models.Rs) map[string]string {
	n := make(map[string]string)
	n["version"] = "2"

	// Use channel info language
	n["title"] = rs.GetTitle(b.getLanguageText)

	n["description"] = fmt.Sprintf("👇 %s <:rs:918545444425072671> %s (%d) ",
		b.getTextForInfo(rs.Info, "wishing_to"), rs.LevelRsPingDs, rs.NumberLevel)
	n["EmbedFieldName"] = fmt.Sprintf(" %s %s\n%s %s\n%s %s",
		emOK, b.getTextForInfo(rs.Info, "to_add_to_queue"),
		emCancel, b.getTextForInfo(rs.Info, "to_exit_the_queue"),
		emRsStart, b.getTextForInfo(rs.Info, "forced_start"))
	n["EmbedFieldValue"] = b.getTextForInfo(rs.Info, "data_updated") + ": "
	n["buttonLevel"] = rs.GetLevelRs()

	return n
}

func (b *Bot) checkAdmin(in *models.InMessageV2) bool {
	admin := false
	var err error
	if in.Messenger.TypeMessenger == ds {
		admin = b.client.Ds.CheckAdmin(in.UserId, in.Messenger.ChannelId)
	} else if in.Messenger.TypeMessenger == tg {
		admin, err = b.client.Tg.CheckAdminTg(in.Messenger.ChannelId, in.Username)
		if err != nil {
			b.log.ErrorErr(err)
		}
	} else if in.Username == "Mentalisit" || in.Username == "mentalisit" {
		admin = true
	}
	return admin
}

func (b *Bot) Translate(in models.InMessage) {
	text2, err := gt.Translate(in.Mtext, "auto", in.Config.Country)
	if err == nil {
		if in.Mtext != text2 {
			if in.Tip == ds {
				go func(textCopy string, usernameCopy string, dsChannelCopy string, avatarCopy string) {
					ch := utils.WaitForMessage("Translate")
					m := b.client.Ds.SendWebhook(textCopy, usernameCopy, dsChannelCopy, avatarCopy)
					b.client.Ds.DeleteMessageSecond(dsChannelCopy, m, 90)
					close(ch)
				}(text2, in.Username, in.Config.DsChannel, in.Ds.Avatar)
			} else if in.Tip == tg {
				go b.client.Tg.SendChannelDelSecond(in.Config.TgChannel, text2, 90)
			}
		}
	}

}

func (b *Bot) getListUsers(u *[]models.QueueActive, inf *models.Info, mention bool) string {
	forType := inf.TypeMessenger
	forAlias := inf.GuildName
	if inf.Game != nil && inf.Game.Alias != "" {
		forAlias = inf.Game.Alias
	}
	var text string
	if len(*u) != 0 {
		for i, q := range *u {
			var name string
			if q.Data.Tip == ds && q.Data.Tip == forType && q.Data.Alias == forAlias {
				name = b.helpers.EmReadName(&q, forType, true)
			} else if mention && q.Data.Tip == forType && q.Data.Alias == forAlias {
				name = b.helpers.EmReadName(&q, forType, mention)
			} else {
				name = b.helpers.EmReadName(&q, forType)
			}

			switch i {
			case 0:
				text = fmt.Sprintf("1️⃣ %s", name)
			case 1:
				text += fmt.Sprintf("\n2️⃣ %s", name)
			case 2:
				text += fmt.Sprintf("\n3️⃣ %s", name)
			case 3:
				text += fmt.Sprintf("\n4️⃣ %s", name)
			}
			if !mention {
				text += fmt.Sprintf("  🕒  %d (%d)", q.RemainingTime, q.Data.NumRSName)
			}
			if q.Data.Alias != forAlias {
				text = fmt.Sprintf("%s\n        %s", text, q.Data.Alias)
			}
		}
	}

	return text
}
