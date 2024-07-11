package DiscordClient

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mentalisit/logger"
	"kz_bot/clients/DiscordClient/transmitter"
	"kz_bot/config"
	"kz_bot/models"
	"kz_bot/pkg/clientDiscord"
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
		webhook:       transmitter.New(ds, "", "KzBot", true, log),
		log:           log,
		storage:       st,
		ChanRsMessage: make(chan models.InMessage, 10),
		//ChanBridgeMessage: make(chan models.BridgeMessage, 20),
		bridgeConfig: &st.BridgeConfigs,
		corpConfigRS: st.CorpConfigRS,
	}
	go ds.AddHandler(DS.messageHandler)
	go ds.AddHandler(DS.messageUpdate)
	go ds.AddHandler(DS.messageReactionAdd)
	go ds.AddHandler(DS.onMessageDelete)
	go ds.AddHandler(DS.slash)
	//go DS.ready()
	go DS.loadSlashCommand()

	return DS
}
func (d *Discord) Shutdown() {
	err := d.S.Close()
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
}

func (d *Discord) QueueSend(text, username string) {
	var id string
	var avatarurl string
	chatid := "1232711859690406042"
	if username == "RsBot" {
		avatarurl = d.S.State.User.AvatarURL("128")
	} else if username == "rssoyzbot" {
		avatarurl = "https://www.superherodb.com/pictures2/portraits/10/050/10409.jpg"
	}
	messages, err := d.S.ChannelMessages(chatid, 10, "", "", "")
	if err != nil {
		d.log.ErrorErr(err)
	}
	if len(messages) > 0 {
		for _, message := range messages {
			if message.Author.Username == username {
				id = message.ID
				break
			}
		}
	}

	if text != "" && id == "" {
		send, err := d.webhook.Send(chatid, &discordgo.WebhookParams{Content: text, Username: username, AvatarURL: avatarurl})
		if err != nil {
			d.log.ErrorErr(err)
			return
		}
		id = send.ID
	} else if id != "" {
		message, errs := d.S.ChannelMessage(chatid, id)
		if errs != nil {
			d.log.ErrorErr(errs)
			return
		}

		if message.Content == "нет активных очередей" && (text == "" || text == "нет активных очередей") {
			return
		}

		ts := fmt.Sprintf("\n<t:%d:f>", time.Now().UTC().Unix())
		err = d.webhook.Edit(chatid, id, &discordgo.WebhookParams{Content: text + ts, Username: username, AvatarURL: avatarurl})
		if err != nil {
			d.log.ErrorErr(err)
			return
		}
	}
}
