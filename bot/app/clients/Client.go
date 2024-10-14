package clients

import (
	"github.com/mentalisit/logger"
	"kz_bot/clients/DiscordClient"
	"kz_bot/clients/TgApi"
	"kz_bot/config"
	"kz_bot/storage"
)

type Clients struct {
	DS *DiscordClient.Discord
	Ds DiscordInterface
	//Ds      *DsApi.DsApi
	Tg      *TgApi.TgApi
	storage *storage.Storage
}

func NewClients(log *logger.Logger, st *storage.Storage, cfg *config.ConfigBot) *Clients {
	c := &Clients{
		storage: st,
		Tg:      TgApi.NewTgApi(log),
	}

	ds := DiscordClient.NewDiscord(log, st, cfg)
	c.Ds = ds
	c.DS = ds
	//c.Ds = DsApi.NewTgApi(log)
	return c
}
func (c *Clients) DeleteMessageTimer() {
	if config.Instance.BotMode != "dev" {
		m := c.storage.TimeDeleteMessage.TimerDeleteMessage()
		if len(m) > 0 {
			for _, timer := range m {
				if timer.Dsmesid != "" {
					go c.Ds.DeleteMesageSecond(timer.Dschatid, timer.Dsmesid, timer.Timed)
				}
				if timer.Tgmesid != "" {
					go c.Tg.DelMessageSecond(timer.Tgchatid, timer.Tgmesid, timer.Timed)
				}
			}
		}
	}
}
func (c *Clients) Shutdown() {
	//c.Ds.Shutdown()
	//c.Tg.Shutdown()
}

type DiscordInterface interface {
	CleanChat(chatid, mesid, text string)
	CleanRsBotOtherMessage()
	CleanOldMessageChannel(chatId, lim string)
	CheckAdmin(nameid string, chatid string) bool
	ChannelTyping(ChannelID string)
	DeleteMessage(chatid, mesid string)
	DeleteMesageSecond(chatid, mesid string, second int)
	EditComplexButton(dsmesid, dschatid string, mapEmbed map[string]string) error
	EditWebhook(text, username, chatID, mID string, avatarURL string)
	QueueSend(text string)
	RoleToIdPing(rolePing, guildid string) (string, error)
	SendDmText(text, AuthorID string)
	Send(chatid, text string) (mesId string)
	SendChannelDelSecond(chatid, text string, second int)
	SendEmbedTime(chatid, text string) (mesId string)
	SendComplexContent(chatid, text string) (mesId string)
	SendComplex(chatid string, mapEmbeds map[string]string) (mesId string)
	SendEmbedText(chatid, title, text string) string
	SendWebhook(text, username, chatid, Avatar string) (mesId string)
	Subscribe(nameid, argRoles, guildid string) int
	Unsubscribe(nameid, argRoles, guildid string) int
}
