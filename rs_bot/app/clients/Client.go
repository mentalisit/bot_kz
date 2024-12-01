package clients

import (
	"github.com/mentalisit/logger"
	ds "rs/clients/DsApi"
	"rs/clients/TgApi"
	"rs/storage"
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

	return c
}

func (c *Clients) Shutdown() {
	c.Tg.Close()
	c.Ds.Close()
}
