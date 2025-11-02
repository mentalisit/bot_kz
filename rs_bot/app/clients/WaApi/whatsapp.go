package WaApi

import (
	"context"
	"errors"
	"fmt"

	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client WhatsappServiceClient
	log    *logger.Logger
}

func NewClient(log *logger.Logger) *Client {
	conn, err := grpc.Dial("whatsapp:50051", grpc.WithInsecure())
	if err != nil {
		log.ErrorErr(err)
		return nil
	}
	fmt.Println("connect to grpc whatsapp ok")
	return &Client{
		conn:   conn,
		client: NewWhatsappServiceClient(conn),
		log:    log,
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) DeleteMessage(ChatId string, messageID string) {
	er, err := c.client.DeleteMessage(context.Background(), &DeleteMessageRequest{
		Chatid: ChatId,
		Mesid:  messageID,
	})
	if err != nil {
		c.log.ErrorErr(err)
	} else if er.GetErrorMessage() != "" {
		c.log.Error(er.GetErrorMessage())
	}
}

func (c *Client) SendChannelDelSecond(chatId, text string, second int) {
	_, err := c.client.SendChannelDelSecond(context.Background(), &SendMessageRequest{
		Text:   text,
		ChatID: chatId,
		Second: int32(second),
	})
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
}

func (c *Client) SendPicScoreboard(chatId, text, filename string) (mid string, err error) {
	scoreboardResponse, err := c.client.SendPicScoreboard(context.Background(), &ScoreboardRequest{
		ChaatId:            chatId,
		Text:               text,
		FileNameScoreboard: filename,
	})
	if err != nil {
		return "", err
	}
	if scoreboardResponse.ErrorMessage != "" {
		return "", errors.New(scoreboardResponse.ErrorMessage)
	}
	return scoreboardResponse.Mid, nil
}
