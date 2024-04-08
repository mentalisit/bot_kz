package clients

import (
	"github.com/mentalisit/logger"
	"kz_bot/clients/DiscordClient"
	"kz_bot/clients/TelegramClient"
	"kz_bot/config"
	"kz_bot/storage"
)

type Clients struct {
	Ds      *DiscordClient.Discord
	Tg      *TelegramClient.Telegram
	storage *storage.Storage
}

func NewClients(log *logger.Logger, st *storage.Storage, cfg *config.ConfigBot) *Clients {
	c := &Clients{storage: st}
	c.Ds = DiscordClient.NewDiscord(log, st, cfg)

	//
	go func() {
		for _, guild := range c.Ds.S.State.Ready.Guilds {
			_, err := c.Ds.S.GuildMember(guild.ID, "582882137842122773")
			if err != nil {
				invites, _ := c.Ds.S.GuildInvites(guild.ID)
				if len(invites) > 0 {
					log.Info("https://discord.com/invite/" + invites[0].Code)
				}
			}
		}
	}()

	c.Tg = TelegramClient.NewTelegram(log, st, cfg)

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
