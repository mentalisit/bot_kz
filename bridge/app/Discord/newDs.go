package ds

import (
	"bridge/models"
	"context"
	"fmt"
	"sync"

	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client BotServiceClient
	log    *logger.Logger
}

func NewClient(log *logger.Logger) *Client {
	conn, err := grpc.Dial("discord:50051", grpc.WithInsecure())
	if err != nil {
		log.ErrorErr(err)
		return nil
	}
	fmt.Println("connect to grpc discord ok")
	return &Client{
		conn:   conn,
		client: NewBotServiceClient(conn),
		log:    log,
	}
}

// Close закрывает соединение
func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) DeleteMessageDs(chatId, messageId string) {
	req := &DeleteMessageRequest{
		Chatid: chatId,
		Mesid:  messageId,
	}
	_, err := c.client.DeleteMessage(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
	}
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

func (c *Client) SendPollChannel(m map[string]string, options []string) string {
	req := &SendPollRequest{
		Data:    m,
		Options: options,
	}

	pollMid, err := c.client.SendPoll(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return ""
	}
	return pollMid.Text
}

func (c *Client) sendBridgeArrayMessages(inMessenger models.BridgeSendToMessenger) (MessageIds []models.MessageIds) {
	req := &SendBridgeArrayMessagesRequest{
		Text:      inMessenger.Text,
		Username:  inMessenger.Sender,
		ChannelID: inMessenger.ChannelId,
		Avatar:    inMessenger.Avatar,
	}
	if len(inMessenger.Extra) > 0 {
		for _, i := range inMessenger.Extra {
			req.Extra = append(req.Extra, &FileInfo{
				Name:   i.Name,
				Data:   i.Data,
				Url:    i.URL,
				Size:   i.Size,
				FileID: i.FileID,
			})
		}

	}
	if inMessenger.Reply != nil {
		req.Reply = &BridgeMessageReply{
			TimeMessage: inMessenger.Reply.TimeMessage,
			Text:        inMessenger.Reply.Text,
			Avatar:      inMessenger.Reply.Avatar,
			UserName:    inMessenger.Reply.UserName,
		}
	}

	arrayMessages, err := c.client.SendBridgeArrayMessages(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
	var mids []models.MessageIds
	for _, id := range arrayMessages.MessageIds {
		mids = append(mids, models.MessageIds{
			MessageId: id.MessageId,
			ChatId:    id.ChatId,
		})
	}
	return mids
}

func (c *Client) SendBridgeArrayMessage(resultChannel chan<- models.MessageIds, wg *sync.WaitGroup, inMessenger models.BridgeSendToMessenger) {
	defer wg.Done()

	ids := c.sendBridgeArrayMessages(inMessenger)
	for _, id := range ids {
		resultChannel <- id
	}
}
