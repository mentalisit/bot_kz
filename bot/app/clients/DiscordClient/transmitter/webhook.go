package transmitter

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func (t *Transmitter) Send(channelID string, params *discordgo.WebhookParams) (*discordgo.Message, error) {
	wh, err := t.getOrCreateWebhook(channelID)
	if err != nil {
		t.log.ErrorErr(err)
		return nil, err
	}

	msg, err := t.session.WebhookExecute(wh.ID, wh.Token, true, params)
	if err != nil {
		return nil, fmt.Errorf("execute failed: %w", err)
	}

	return msg, nil
}

func (t *Transmitter) Edit(channelID string, messageID string, params *discordgo.WebhookParams) error {
	wh := t.getWebhook(channelID)

	if wh == nil {
		return ErrWebhookNotFound
	}

	uri := discordgo.EndpointWebhookToken(wh.ID, wh.Token) + "/messages/" + messageID
	_, err := t.session.RequestWithBucketID("PATCH", uri, params, discordgo.EndpointWebhookToken("", ""))
	if err != nil {
		return err
	}

	return nil
}
