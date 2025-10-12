package wa

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

func (c *Client) sendBridgeArrayMessages(inMessenger models.BridgeSendToMessenger) (MessageIds []models.MessageIds) {
	req := &SendBridgeArrayMessagesRequest{
		Text:      inMessenger.Text,
		ChannelID: inMessenger.ChannelId,
		ReplyMap:  inMessenger.ReplyMap,
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

	arrayMessages, err := c.client.SendBridgeArrayMessages(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
	var mids []models.MessageIds
	for _, id := range arrayMessages.GetMessageIds() {
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
