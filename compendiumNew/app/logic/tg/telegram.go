package tg

import (
	"context"
	"errors"
	"fmt"
	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client TelegramServiceClient
	log    *logger.Logger
}

func NewClient(log *logger.Logger) *Client {
	conn, err := grpc.Dial("telegram:50051", grpc.WithInsecure())
	if err != nil {
		log.ErrorErr(err)
		return nil
	}
	fmt.Println("connect to grpc discord ok")
	return &Client{
		conn:   conn,
		client: NewTelegramServiceClient(conn),
		log:    log,
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Send(chatId string, text string) (string, error) {
	textResponse, err := c.client.Send(context.Background(), &SendMessageRequest{
		Text:   text,
		ChatID: chatId,
	})
	if err != nil {
		return "", err
	}
	if textResponse.Text == "Forbidden" {
		return "", errors.New("forbidden")
	}
	if textResponse.GetText() == "" {
		return "", errors.New("empty")
	}
	return textResponse.GetText(), nil
}
func (c *Client) SendPic(channelId string, text string, pic []byte) error {
	er, err := c.client.SendPic(context.Background(), &SendPicRequest{
		Chatid:     channelId,
		Text:       text,
		ImageBytes: pic,
	})
	if err != nil {
		c.log.ErrorErr(err)
		return err
	}
	if er.GetErrorMessage() != "" {
		return errors.New(er.GetErrorMessage())
	}
	return nil

}
func (c *Client) DeleteMessage(ChatId string, messageID string) error {
	er, err := c.client.DeleteMessage(context.Background(), &DeleteMessageRequest{
		Chatid: ChatId,
		Mesid:  messageID,
	})
	if err != nil {
		c.log.ErrorErr(err)
		return err
	}
	if er.GetErrorMessage() != "" {
		return errors.New(er.GetErrorMessage())
	}
	return nil
}
func (c *Client) EditMessage(channel, messageID string, text, parse string) error {
	er, err := c.client.EditMessage(context.Background(), &EditMessageRequest{
		TextEdit:  text,
		ChatID:    channel,
		MID:       messageID,
		ParseMode: parse,
	})
	if err != nil {
		c.log.ErrorErr(err)
		return err
	}
	if er.GetErrorMessage() != "" {
		return errors.New(er.GetErrorMessage())
	}
	return nil
}
func (c *Client) GetAvatarUrl(userid string) string {
	tr, err := c.client.GetAvatarUrl(context.Background(), &GetAvatarUrlRequest{
		Userid: userid,
	})
	if err != nil {
		c.log.ErrorErr(err)
		return ""
	}
	return tr.GetText()
}
