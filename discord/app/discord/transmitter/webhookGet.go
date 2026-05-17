package transmitter

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (t *Transmitter) getOrCreateWebhook(channelID string) (*discordgo.Webhook, error) {
	wh := t.getWebhook(channelID)
	if wh != nil {
		return wh, nil
	}

	t.log.Info("Creating a webhook for " + channelID)
	wh, err := t.createWebhook(channelID)
	if err != nil {
		t.log.ErrorErr(err)
		return nil, err
	}

	return wh, nil
}

func (t *Transmitter) getWebhook(channel string) *discordgo.Webhook {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if t.HasWebhook(channel) {
		return t.channelWebhooks[channel]
	}

	// Add retry mechanism with exponential backoff
	var webhooks []*discordgo.Webhook
	var err error

	for attempt := 0; attempt < 3; attempt++ {
		webhooks, err = t.session.ChannelWebhooks(channel)
		if err == nil {
			break
		}

		if attempt < 2 {
			// Exponential backoff with jitter: 1s, 2s, 4s
			backoffDuration := time.Duration(1<<uint(attempt)) * time.Second
			jitter := time.Duration(rand.Float64() * float64(backoffDuration) * 0.1) // 10% jitter
			time.Sleep(backoffDuration + jitter)

			// Log the error with attempt info
			t.log.ErrorErr(fmt.Errorf("Attempt %d/3 failed to get webhooks for channel %s: %w", attempt+1, channel, err))

		}
	}

	if err != nil {
		t.log.ErrorErr(fmt.Errorf("All attempts failed to get webhooks for channel %s: %w", channel, err))
		return nil
	}

	var webhook *discordgo.Webhook
	for _, i := range webhooks {
		if i.User.Bot && i.User.Username == t.session.State.User.Username {
			webhook = i
			t.channelWebhooks[channel] = webhook
			return webhook
		}
	}

	if webhook == nil {
		webhookCreate, err1 := t.session.WebhookCreate(channel, t.title, "")
		if err1 != nil {
			if len(webhooks) > 0 {
				return webhooks[0]
			}
			t.log.ErrorErr(err1)
			return nil
		}

		return webhookCreate
	}

	return nil
}
func (t *Transmitter) HasWebhook(id string) bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	for _, wh := range t.channelWebhooks {
		if wh.ID == id {
			return true
		}
	}

	return false
}
func (t *Transmitter) createWebhook(channel string) (*discordgo.Webhook, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	var wh *discordgo.Webhook
	var err error

	// Add retry mechanism with exponential backoff
	for attempt := 0; attempt < 3; attempt++ {
		wh, err = t.session.WebhookCreate(channel, t.title+time.Now().Format(" 3:04:05PM"), "")
		if err == nil {
			break
		}

		// Log the error with attempt info
		t.log.ErrorErr(fmt.Errorf("Attempt %d/3 failed to create webhook for channel %s: %w", attempt+1, channel, err))

		if attempt < 2 {
			// Exponential backoff with jitter: 1s, 2s, 4s
			backoffDuration := time.Duration(1<<uint(attempt)) * time.Second
			jitter := time.Duration(rand.Float64() * float64(backoffDuration) * 0.1) // 10% jitter
			time.Sleep(backoffDuration + jitter)
		}
	}

	if err != nil {
		t.log.ErrorErr(fmt.Errorf("All attempts failed to create webhook for channel %s: %w", channel, err))
		return nil, err
	}

	t.channelWebhooks[channel] = wh
	return wh, nil
}
