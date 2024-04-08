package transmitter

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/mentalisit/logger"
	"sync"
)

type Transmitter struct {
	session         *discordgo.Session
	guild           string
	title           string
	autoCreate      bool
	channelWebhooks map[string]*discordgo.Webhook
	mutex           sync.RWMutex
	log             *logger.Logger
}

var ErrWebhookNotFound = errors.New("webhook for this channel and message does not exist")

func New(session *discordgo.Session, guild string, title string, autoCreate bool, log *logger.Logger) *Transmitter {
	return &Transmitter{
		session:    session,
		guild:      guild,
		title:      title,
		autoCreate: autoCreate,

		channelWebhooks: make(map[string]*discordgo.Webhook),

		log: log,
	}
}
