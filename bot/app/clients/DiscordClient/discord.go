package DiscordClient

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mentalisit/logger"
	"kz_bot/clients/DiscordClient/transmitter"
	"kz_bot/config"
	"kz_bot/models"
	"kz_bot/pkg/clientDiscord"
	"kz_bot/storage"
)

type Discord struct {
	ChanRsMessage chan models.InMessage
	//ChanBridgeMessage chan models.BridgeMessage
	S            *discordgo.Session
	webhook      *transmitter.Transmitter
	log          *logger.Logger
	storage      *storage.Storage
	bridgeConfig map[string]models.BridgeConfig
	corpConfigRS map[string]models.CorporationConfig
}

func NewDiscord(log *logger.Logger, st *storage.Storage, cfg *config.ConfigBot) *Discord {
	ds, err := clientDiscord.NewDiscord(log, cfg)
	if err != nil {
		log.ErrorErr(err)
	}

	DS := &Discord{
		S:             ds,
		webhook:       transmitter.New(ds, "", "KzBot", true, log),
		log:           log,
		storage:       st,
		ChanRsMessage: make(chan models.InMessage, 10),
		//ChanBridgeMessage: make(chan models.BridgeMessage, 20),
		bridgeConfig: st.BridgeConfigs,
		corpConfigRS: st.CorpConfigRS,
	}
	go ds.AddHandler(DS.messageHandler)
	go ds.AddHandler(DS.messageUpdate)
	go ds.AddHandler(DS.messageReactionAdd)
	go ds.AddHandler(DS.onMessageDelete)
	go ds.AddHandler(DS.slash)
	go DS.ready()

	return DS
}
