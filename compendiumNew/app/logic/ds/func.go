package ds

import (
	"context"
	"errors"
)

func (c *Client) DeleteMessage(ChatId, MesId string) error {
	req := &DeleteMessageRequest{
		Chatid: ChatId,
		Mesid:  MesId,
	}
	_, err := c.client.DeleteMessage(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}
func (c *Client) SendChannelDelSecondDs(ChatId, text string, second int) {
	req := &SendChannelDelSecondRequest{
		Chatid: ChatId,
		Text:   text,
		Second: int32(second),
	}

	_, err := c.client.SendChannelDelSecond(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
	}
}
func (c *Client) SendChannelPic(chatId, text string, pic []byte) error {
	req := &SendPicRequest{
		Chatid:     chatId,
		Text:       text,
		ImageBytes: pic,
	}

	errResponse, err := c.client.SendPic(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return err
	}
	if errResponse.ErrorMessage != "" {
		return errors.New(errResponse.ErrorMessage)
	}
	return nil
}
func (c *Client) EditMessage(chatId, mid, text string) error {
	req := &EditMessageRequest{
		Content: text,
		ChatID:  chatId,
		MID:     mid,
	}

	_, err := c.client.EditMessage(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return err
	}
	return nil
}
func (c *Client) SendChannel(chatId, text string) (string, error) {
	req := &SendRequest{
		Chatid: chatId,
		Text:   text,
	}

	textR, err := c.client.Send(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return "", err
	}
	return textR.Text, nil
}
func (c *Client) GetAvatarUrl(userid string) string {
	req := &GetAvatarUrlRequest{
		Userid: userid,
	}

	textR, err := c.client.GetAvatarUrl(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return ""
	}
	return textR.Text
}
