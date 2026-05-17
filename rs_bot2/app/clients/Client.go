package clients

import (
	ds "rs/clients/DsApi"
	"rs/clients/TgApi"
	"rs/clients/WaApi"
	"rs/storage"

	"github.com/mentalisit/logger"
)

type Clients struct {
	Ds      *ds.Client
	Tg      *TgApi.Client
	Wa      *WaApi.Client
	storage *storage.Storage
}

func NewClients(log *logger.Logger, st *storage.Storage) *Clients {
	c := &Clients{
		storage: st,
		Tg:      TgApi.NewClient(log),
		Ds:      ds.NewClient(log),
		Wa:      WaApi.NewClient(log),
	}

	return c
}

func (c *Clients) Shutdown() {
	c.Tg.Close()
	c.Ds.Close()
	c.Wa.Close()
}
