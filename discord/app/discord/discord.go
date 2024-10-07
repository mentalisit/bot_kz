package DiscordClient

import (
	"discord/config"
	"discord/discord/restapi"
	"discord/discord/transmitter"
	"discord/models"
	"discord/storage"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mentalisit/logger"
	"time"
)

type Discord struct {
	S            *discordgo.Session
	webhook      *transmitter.Transmitter
	log          *logger.Logger
	storage      *storage.Storage
	bridgeConfig map[string]models.BridgeConfig
	corpConfigRS map[string]models.CorporationConfig
	api          *restapi.Recover
}

func NewDiscord(log *logger.Logger, st *storage.Storage, cfg *config.ConfigBot) *Discord {
	ds, err := discordgo.New("Bot " + cfg.TokenDiscord)
	if err != nil {
		log.Panic("Ошибка запуска дискорда" + err.Error())
		return nil
	}

	err = ds.Open()
	if err != nil {
		log.Panic("Ошибка открытия ДС " + err.Error())
	}
	fmt.Println("Бот Дискорд загружен ")
	DS := &Discord{
		S:            ds,
		webhook:      transmitter.New(ds, "KzBot", true, log),
		log:          log,
		storage:      st,
		bridgeConfig: make(map[string]models.BridgeConfig),
		corpConfigRS: make(map[string]models.CorporationConfig),
		api:          restapi.NewRecover(log),
	}
	ds.AddHandler(DS.messageHandler)
	ds.AddHandler(DS.messageReactionAdd)
	ds.AddHandler(DS.slash)
	go DS.loadSlashCommand()

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
	chatid := "1232711859690406042"
	mid := "1283266865535254660"

	ts := fmt.Sprintf("\n<t:%d:f>", time.Now().UTC().Unix())
	d.EditMessage(chatid, mid, text+ts)
}

func (d *Discord) loadConfig() {
	bc := restapi.GetBridgeConfig()
	if len(bc) > 0 {
		d.bridgeConfig = bc
	}

	rs := d.storage.Db.ReadConfigRs()
	if len(rs) > 0 {
		fmt.Printf("rsLoad %d\n", len(rs))
		for _, configRs := range rs {
			d.corpConfigRS[configRs.CorpName] = configRs
		}
	}

}
