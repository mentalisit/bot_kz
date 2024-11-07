package clients

import (
	"github.com/mentalisit/logger"
	ds "rs/clients/DsApi"
	"rs/clients/TgApi"
	"rs/config"
	"rs/storage"
)

type Clients struct {
	Ds      *ds.Client
	Tg      *TgApi.Client
	storage *storage.Storage
}

func NewClients(log *logger.Logger, st *storage.Storage, cfg *config.ConfigBot) *Clients {
	c := &Clients{
		storage: st,
		Tg:      TgApi.NewClient(log),
	}
	c.Ds = ds.NewClient(log)
	return c
}

//func (c *Clients) DeleteMessageTimer() {
//	if config.Instance.BotMode != "dev" {
//		m := c.storage.TimeDeleteMessage.TimerDeleteMessage()
//		if len(m) > 0 {
//			for _, timer := range m {
//				if timer.Dsmesid != "" {
//					go c.Ds.DeleteMessageSecond(timer.Dschatid, timer.Dsmesid, timer.Timed)
//				}
//				if timer.Tgmesid != "" {
//					go c.Tg.DelMessageSecond(timer.Tgchatid, timer.Tgmesid, timer.Timed)
//				}
//			}
//		}
//	}
//}

func (c *Clients) Shutdown() {
	c.Tg.Close()
	c.Ds.Close()
}
