package DiscordClient

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mentalisit/logger"
	"kz_bot/clients/DiscordClient/transmitter"
	"kz_bot/config"
	"kz_bot/models"
	"kz_bot/pkg/clientDiscord"
	"kz_bot/pkg/utils"
	"kz_bot/storage"
	"time"
)

type Discord struct {
	ChanRsMessage chan models.InMessage
	S             *discordgo.Session
	webhook       *transmitter.Transmitter
	log           *logger.Logger
	storage       *storage.Storage
	bridgeConfig  *map[string]models.BridgeConfig
	corpConfigRS  map[string]models.CorporationConfig
}

func NewDiscord(log *logger.Logger, st *storage.Storage, cfg *config.ConfigBot) *Discord {
	ds, err := clientDiscord.NewDiscord(log, cfg)
	if err != nil {
		log.ErrorErr(err)
	}

	DS := &Discord{
		S:             ds,
		webhook:       transmitter.New(ds, "KzBot", true, log),
		log:           log,
		storage:       st,
		ChanRsMessage: make(chan models.InMessage, 10),
		//ChanBridgeMessage: make(chan models.BridgeMessage, 20),
		bridgeConfig: &st.BridgeConfigs,
		corpConfigRS: st.CorpConfigRS,
	}
	ds.AddHandler(DS.messageHandler)
	ds.AddHandler(DS.messageReactionAdd)
	ds.AddHandler(DS.slash)
	go DS.loadSlashCommand()
	//go ds.AddHandler(DS.onMessageDelete)
	//ds.AddHandler(DS.messageUpdate)

	go func() {
		for _, guild := range DS.S.State.Ready.Guilds {
			_, err = ds.GuildMember(guild.ID, "582882137842122773")
			if err != nil {
				invites, _ := ds.GuildInvites(guild.ID)
				if len(invites) > 0 {
					log.Info("https://discord.com/invite/" + invites[0].Code)
				}
			}
		}
	}()

	return DS
}
func (d *Discord) Shutdown() {
	err := d.S.Close()
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
}

func (d *Discord) QueueSend(text string) {
	ch := utils.WaitForMessage("QueueSend")
	defer close(ch)
	chatid := "1232711859690406042"
	mid := "1283266865535254660"

	ts := fmt.Sprintf("\n<t:%d:f>", time.Now().UTC().Unix())
	d.EditMessage(chatid, mid, text+ts)
}
