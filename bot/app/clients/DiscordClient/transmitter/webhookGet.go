package transmitter

import (
	"github.com/bwmarrin/discordgo"
	"time"
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

	webhooks, err := t.session.ChannelWebhooks(channel)
	if err != nil {
		t.log.ErrorErr(err)
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

	wh, err := t.session.WebhookCreate(channel, t.title+time.Now().Format(" 3:04:05PM"), "")
	if err != nil {
		t.log.ErrorErr(err)
		return nil, err
	}
	t.channelWebhooks[channel] = wh
	return wh, nil
}
