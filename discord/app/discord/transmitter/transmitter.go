package transmitter

import (
	"errors"
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/mentalisit/logger"
)

type Transmitter struct {
	session         *discordgo.Session
	title           string
	autoCreate      bool
	channelWebhooks map[string]*discordgo.Webhook
	mutex           sync.RWMutex
	log             *logger.Logger
}

var ErrWebhookNotFound = errors.New("webhook for this channel and message does not exist")

func New(session *discordgo.Session, title string, autoCreate bool, log *logger.Logger) *Transmitter {
	return &Transmitter{
		session:    session,
		title:      title,
		autoCreate: autoCreate,

		channelWebhooks: make(map[string]*discordgo.Webhook),

		log: log,
	}
}

func (t *Transmitter) CreateChannelAndWebhook(guildID, channelName, categoryID string) (webhookUrl, chatId string, err error) {
	// Проверяем, есть ли уже вебхук для этого канала в мапе
	t.mutex.RLock()
	if webhook, exists := t.channelWebhooks[channelName]; exists {
		t.mutex.RUnlock()
		return fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", webhook.ID, webhook.Token), webhook.ChannelID, nil
	}
	t.mutex.RUnlock()

	// Проверяем, существует ли канал с таким именем
	channels, err := t.session.GuildChannels(guildID)
	if err != nil {
		t.log.ErrorErr(err)
		return "", "", err
	}

	var existingChannel *discordgo.Channel
	for _, ch := range channels {
		if ch.Name == channelName && ch.Type == discordgo.ChannelTypeGuildText {
			existingChannel = ch
			break
		}
	}

	var channelID string
	if existingChannel != nil {
		// Канал уже существует
		channelID = existingChannel.ID

		// Проверяем, есть ли уже вебхуки в этом канале
		webhooks, err := t.session.ChannelWebhooks(channelID)
		if err != nil {
			t.log.ErrorErr(err)
			return "", "", err
		}

		// Ищем вебхук с нашим именем
		for _, wh := range webhooks {
			if wh.Name == t.title {
				t.mutex.Lock()
				t.channelWebhooks[channelName] = wh
				t.mutex.Unlock()
				return fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", wh.ID, wh.Token), wh.ChannelID, nil
			}
		}
	} else {
		// Создаем новый канал с указанием категории
		channelCreateData := discordgo.GuildChannelCreateData{
			Name:     channelName,
			Type:     discordgo.ChannelTypeGuildText,
			ParentID: categoryID,
		}
		channel, err := t.session.GuildChannelCreateComplex(guildID, channelCreateData)
		if err != nil {
			t.log.ErrorErr(err)
			return "", "", err
		}
		channelID = channel.ID
	}

	// Создаем вебхук для канала
	webhook, err := t.session.WebhookCreate(channelID, t.title, "")
	if err != nil {
		t.log.ErrorErr(err)
		return "", "", err
	}

	// Сохраняем вебхук в мапе
	t.mutex.Lock()
	t.channelWebhooks[channelName] = webhook
	t.mutex.Unlock()

	return fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", webhook.ID, webhook.Token), webhook.ChannelID, nil
}
