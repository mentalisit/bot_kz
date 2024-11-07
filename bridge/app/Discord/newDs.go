package ds

import (
	"bridge/models"
	"context"
	"fmt"
	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
	"sync"
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

func (c *Client) sendBridgeArrayMessages(text, username string, channelID []string, extra []models.FileInfo, Avatar string, reply *models.BridgeMessageReply) (MessageIds []models.MessageIds) {
	req := &SendBridgeArrayMessagesRequest{
		Text:      text,
		Username:  username,
		ChannelID: channelID,
		Avatar:    Avatar,
	}
	if len(extra) > 0 {
		req.Extra = make([]*FileInfo, len(extra))
		for _, i := range extra {
			req.Extra = append(req.Extra, &FileInfo{
				Name:   i.Name,
				Data:   i.Data,
				Url:    i.URL,
				Size:   i.Size,
				FileID: i.FileID,
			})
		}

	}
	if reply != nil {
		req.Reply = &BridgeMessageReply{
			TimeMessage: reply.TimeMessage,
			Text:        reply.Text,
			Avatar:      reply.Avatar,
			UserName:    reply.UserName,
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

func (c *Client) SendBridgeArrayMessage(text, username string, channelID []string, extra []models.FileInfo, Avatar string, reply *models.BridgeMessageReply, resultChannel chan<- models.MessageIds, wg *sync.WaitGroup) {
	defer wg.Done()

	ids := c.sendBridgeArrayMessages(text, username, channelID, extra, Avatar, reply)
	for _, id := range ids {
		resultChannel <- id
	}
}
