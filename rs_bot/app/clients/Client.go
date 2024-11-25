package clients

import (
	"github.com/mentalisit/logger"
	ds "rs/clients/DsApi"
	"rs/clients/TgApi"
	"rs/storage"
	"time"
)

type Clients struct {
	Ds      *ds.Client
	Tg      *TgApi.Client
	storage *storage.Storage
}

func NewClients(log *logger.Logger, st *storage.Storage) *Clients {
	c := &Clients{
		storage: st,
		Tg:      TgApi.NewClient(log),
		Ds:      ds.NewClient(log),
	}
	go c.monitorPrimary()

	return c
}

func (c *Clients) monitorPrimary() {
	for {
		time.Sleep(10 * time.Second) // Check interval
		c.Ds.MonitorPrimary()
		c.Tg.MonitorPrimary()
	}
}

func (c *Clients) Shutdown() {
	c.Tg.Close()
	c.Ds.Close()
}
