package DiscordClient

import (
	"discord/config"
	"discord/discord/helpers"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/mentalisit/restapi/models"
)

const (
	emOK      = "✅"
	emCancel  = "❎"
	emRsStart = "🚀"
	emPl30    = "⌛"
	emPlus    = "➕"
	emMinus   = "➖"
)

func (d *Discord) logicMix(m *discordgo.MessageCreate) {
	// Получаем и кэшируем имена канала и гильдии
	channelName := d.getChannelName(m.ChannelID)
	guildName := d.getGuildName(m.GuildID)

	if d.ifMentionBot(m) {
		return
	}
	if d.avatar(m) {
		return
	}
	go d.latinOrNot(m) //пытаемся переводить гостевой чат

	// Сохраняем сообщение в БД
	go d.saveMessageToStorage(m, channelName, guildName)

	if m.Author != nil && m.Author.Locale != "" {
		go d.log.Info(m.Author.Username + " " + m.Author.Locale)
	}
	if m.Member != nil && m.Member.User != nil && m.Member.User.Locale != "" {
		go d.log.Info(m.Member.User.Username + " " + m.Member.User.Locale)
	}

	if strings.HasPrefix(m.Content, ".") {
		d.ifPrefixPoint(m)
		return
	}

	//filter Rs
	ok, config := d.CheckChannelConfigDS(m.ChannelID)
	if ok {
		d.SendToRsFilter(m, config)
		return
	}
	good2, config2 := d.checkChannelConfig2(m.ChannelID)
	if good2 {
		d.SendToRs2Filter(m, config2)
		return
	}

	//bridge
	ds, bridgeConfig := d.BridgeCheckChannelConfigDS(m.ChannelID)
	if ds {
		d.SendToBridge(m, bridgeConfig)
	}

	if strings.HasPrefix(m.Content, "%") {
		d.SendToCompendium(m)
		return
	}
}

// saveMessageToStorage сохраняет сообщение Discord в БД
func (d *Discord) saveMessageToStorage(m *discordgo.MessageCreate, channelName, guildName string) {
	// Получаем communityID (UUID) для гильдии
	communityID, err := d.storage.Db.GetGildUUIDMyCompendium(m.GuildID)
	if err != nil {
		mg := config.MultiAccountGuildV2{
			GId:       uuid.New(),
			GuildName: d.getGuildName(m.GuildID),
			Channels:  make(config.GuildChannels),
		}
		mg.Channels["ds"] = append(mg.Channels["ds"], m.GuildID)
		save, err := d.storage.Db.GuildSave(mg)
		if err != nil {
			d.log.ErrorErr(err)
			return
		}
		if save.GId != uuid.Nil {
			communityID = &save.GId
		} else {
			// Если не нашли guild, используем нулевой UUID или логируем
			d.log.Info(fmt.Sprintf("Не удалось получить communityID для чата %s: %v", m.GuildID, err))
			return
		}

	}

	// Сохраняем сообщение с переданными именами канала и гильдии
	if err := d.storage.Db.SaveDiscordMessageWithNames(*communityID, m, channelName, guildName); err != nil {
		d.log.Error(fmt.Sprintf("Ошибка сохранения сообщения %s: %v", m.ID, err))
	}
}

// getChannelName получает имя канала из кэша или API
func (d *Discord) getChannelName(channelID string) string {
	if name, exists := d.channelNameCache[channelID]; exists {
		return name
	}

	channel, err := d.S.Channel(channelID)
	if err != nil {
		d.log.ErrorErr(err)
		return channelID // Возвращаем ID если не удалось получить имя
	}

	d.channelNameCache[channelID] = channel.Name
	return channel.Name
}

// getGuildName получает имя гильдии из кэша или API
func (d *Discord) getGuildName(guildID string) string {
	if name, exists := d.guildNameCache[guildID]; exists {
		return name
	}

	guild, err := d.S.Guild(guildID)
	if err != nil {
		d.log.ErrorErr(err)
		return guildID // Возвращаем ID если не удалось получить имя
	}

	d.guildNameCache[guildID] = guild.Name
	return guild.Name
}

func (d *Discord) SendToRs2Filter(m *discordgo.MessageCreate, config2 models.CorporationConfigV2) {
	if len(m.Attachments) > 0 {
		m.Content += m.Attachments[0].URL
	}
	if len(m.Message.Embeds) > 0 {
		m.Content += "\u200B"
	}
	in2 := models.InMessageV2{
		Text:        d.ReplaceTextMessage(m.Content, m.GuildID),
		Tip:         "ds",
		NameNick:    "",
		Username:    m.Author.Username,
		UserId:      m.Author.ID,
		NameMention: m.Author.Mention(),
		Messenger: models.Info{
			TypeMessenger:  "ds",
			MessageId:      m.ID,
			ChannelId:      m.ChannelID,
			GuildId:        m.GuildID,
			GuildName:      config2.Channels[m.ChannelID].GuildName,
			GuildAvatarUrl: config2.Channels[m.ChannelID].GuildAvatarUrl,
			UserAvatarUrl:  m.Author.AvatarURL("128"),
		},
		Config:  config2,
		Options: models.Options{models.OptionInClient},
	}
	d.api.SendRsBotV2AppRecover(in2)
}

func (d *Discord) SendToRsFilter(m *discordgo.MessageCreate, config models.CorporationConfig) {
	if len(m.Attachments) > 0 {
		m.Content += m.Attachments[0].URL
	}
	if len(m.Message.Embeds) > 0 {
		m.Content += "\u200B"
	}
	in := models.InMessage{
		Mtext:       d.ReplaceTextMessage(m.Content, m.GuildID),
		Tip:         "ds",
		Username:    m.Author.Username,
		UserId:      m.Author.ID,
		NameNick:    "",
		NameMention: m.Author.Mention(),
		Ds: struct {
			Mesid   string
			Guildid string
			Avatar  string
		}{
			Mesid:   m.ID,
			Guildid: m.GuildID,
			Avatar:  m.Author.AvatarURL("128"),
		},
		Config: config,
		Option: models.Option{InClient: true},
	}
	if m.Member != nil && m.Member.Nick != "" {
		in.NameNick = m.Member.Nick
	}

	d.api.SendRsBotAppRecover(in)
}
func (d *Discord) ifMentionBot(m *discordgo.MessageCreate) bool {
	after, found := strings.CutPrefix(m.Content, d.S.State.User.Mention())
	if found {
		d.DeleteMesageSecond(m.ChannelID, m.ID, 30)
		goodRs, _ := d.CheckChannelConfigDS(m.ChannelID)
		if goodRs {
			d.SendChannelDelSecond(m.ChannelID, fmt.Sprintf("%s че пингуешь? пиши Справка,или пиши создателю бота @Mentalisit#5159 ", m.Author.Mention()), 30)
		} else {
			m.Content = "%" + after
			d.SendToCompendium(m)
		}
	}
	return found
}
func (d *Discord) readReactionTranslate(r *discordgo.MessageReactionAdd, m *discordgo.Message) {
	user, err := d.S.User(r.UserID)
	if err != nil {
		d.log.ErrorErr(err)
	}
	if user.ID != m.Author.ID {

		switch r.Emoji.Name {
		case "🇺🇸":
			d.transtale(m, "en", r)
		case "🇷🇺":
			d.transtale(m, "ru", r)
		case "🇺🇦":
			d.transtale(m, "uk", r)
		case "🇬🇧":
			d.transtale(m, "en", r)
		case "🇧🇾":
			d.transtale(m, "be", r)
		case "🇩🇪":
			d.transtale(m, "de", r)
		case "🇵🇱":
			d.transtale(m, "pl", r)
		}
	}
}
func (d *Discord) SendToCompendium(m *discordgo.MessageCreate) {
	guildName := d.getGuildName(m.GuildID)
	channelName := d.getChannelName(m.ChannelID)

	user := m.Author
	if m.Member != nil && m.Member.User != nil {
		user = m.Member.User
	}
	nick := ""
	if m.Member != nil && m.Member.Nick != "" {
		nick = m.Member.Nick
	}

	i := models.IncomingMessage{
		Text:         m.Content,
		DmChat:       d.dmChannel(user.ID),
		Name:         user.Username,
		MentionName:  user.Mention(),
		NameId:       user.ID,
		NickName:     nick,
		Avatar:       user.AvatarURL(""),
		AvatarF:      user.Avatar,
		ChannelId:    m.ChannelID,
		GuildId:      m.GuildID,
		GuildName:    guildName,
		GuildAvatar:  "", // TODO: добавить кэширование аватаров гильдий
		GuildAvatarF: "", // TODO: добавить кэширование аватаров гильдий
		Type:         "ds",
	}
	if channelName != "" {
		i.Language = helpers.DetectLanguage(guildName + "/" + channelName)
	} else {
		i.Language = helpers.DetectLanguage(guildName)
	}

	d.api.SendCompendiumAppRecover(i)
}
func (d *Discord) ifPrefixPoint(m *discordgo.MessageCreate) {
	good, config := d.CheckChannelConfigDS(m.ID)
	in := models.InMessage{
		Mtext:       m.Content,
		Tip:         "ds",
		Username:    m.Author.Username,
		UserId:      m.Author.ID,
		NameMention: m.Author.Mention(),
		Ds: struct {
			Mesid   string
			Guildid string
			Avatar  string
		}{
			Mesid:   m.ID,
			Guildid: m.GuildID,
			Avatar:  m.Author.AvatarURL("")},

		Option: models.Option{
			InClient: true,
		},
	}
	if m.Member != nil && m.Member.Nick != "" {
		in.NameNick = m.Member.Nick
	}
	if good {
		in.Config = config
	} else {
		in.Config = models.CorporationConfig{
			CorpName:  d.GuildChatName(m.ChannelID, m.GuildID),
			DsChannel: m.ChannelID,
			Guildid:   m.GuildID,
		}
	}
	d.api.SendRsBotAppRecover(in)

	good2, config2 := d.checkChannelConfig2(m.ChannelID)
	in2 := models.InMessageV2{
		Text:        m.Content,
		Tip:         "ds",
		Username:    m.Author.Username,
		UserId:      m.Author.ID,
		NameMention: m.Author.Mention(),
		Messenger: models.Info{
			TypeMessenger: "ds",
			MessageId:     m.ID,
			ChannelId:     m.ChannelID,
			GuildId:       m.GuildID,
			UserAvatarUrl: m.Author.AvatarURL("128"),
			Language:      "ru",
			CreatedAt:     time.Now(),
		},
		Options: models.Options{models.OptionInClient},
	}

	if good2 {
		in2.Config = config2
	}

	if m.Content == ".setup" || strings.HasPrefix(m.Content, ".invite ") || strings.HasPrefix(m.Content, ".setting") {
		in2.Config = models.CorporationConfigV2{
			Channels:    make(models.ChannelsMap),
			HelpMessage: make(models.HelpMessage),
		}
		guildName := d.getGuildName(m.GuildID)
		channelName := d.getChannelName(m.ChannelID)
		in2.Messenger.GuildName = guildName
		in2.Messenger.GuildAvatarUrl = "" // TODO: добавить кэширование аватаров гильдий
		in2.Messenger.ChannelName = channelName
		if in2.Config.Channels[m.ChannelID] == nil {
			in2.Config.Channels[m.ChannelID] = &models.Info{}
		}
		in2.Config.Channels[m.ChannelID] = &in2.Messenger
	}

	d.api.SendRsBotV2AppRecover(in2)

	go func() {
		mes := models.ToBridgeMessage{
			Text:          m.Content,
			Sender:        m.Author.Username,
			SenderId:      m.Author.ID,
			Tip:           "ds",
			Avatar:        m.Author.AvatarURL("128"),
			ChatId:        m.ChannelID,
			MesId:         m.ID,
			GuildId:       m.GuildID,
			TimestampUnix: m.Timestamp.Unix(),
		}
		ds, bridgeConfig := d.BridgeCheckChannelConfigDS(m.ChannelID)
		if ds {
			mes.Config = &bridgeConfig
		} else {
			mes.Config = &models.Bridge2Config{
				HostRelay: d.GuildChatName(m.ChannelID, m.GuildID),
			}
		}
		d.api.SendBridgeAppRecover(mes)
	}()

}
func (d *Discord) SendToBridge(m *discordgo.MessageCreate, bridgeConfig models.Bridge2Config) {
	mes := models.ToBridgeMessage{
		ChatId:        m.ChannelID,
		Extra:         []models.FileInfo{},
		Config:        &bridgeConfig,
		Text:          d.ReplaceTextMessage(m.Content, m.GuildID),
		Sender:        d.getAuthorName(m),
		SenderId:      m.Author.ID,
		Tip:           "ds",
		MesId:         m.ID,
		GuildId:       m.GuildID,
		TimestampUnix: m.Timestamp.Unix(),
		Avatar:        m.Author.AvatarURL(""),
	}

	d.handleDownloadBridge(&mes, m)

	if m.ReferencedMessage != nil {
		mes.ReplyMap = make(map[string]string)
		mes.ReplyMap[m.ChannelID] = m.ReferencedMessage.ID
		usernameR := m.ReferencedMessage.Author.String()
		if m.ReferencedMessage.Member != nil && m.ReferencedMessage.Member.Nick != "" {
			usernameR = m.ReferencedMessage.Member.Nick
		}
		mes.Reply = &models.BridgeMessageReply{
			TimeMessage: m.ReferencedMessage.Timestamp.Unix(),
			Text:        d.ReplaceTextMessage(m.ReferencedMessage.Content, m.GuildID),
			Avatar:      m.ReferencedMessage.Author.AvatarURL(""),
			UserName:    usernameR,
		}
	}
	if mes.Text != "" || len(mes.Extra) > 0 {
		d.api.SendBridgeAppRecover(mes)
	}
}
