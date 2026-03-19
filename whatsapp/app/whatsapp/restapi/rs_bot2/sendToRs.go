package rs_bot2

import (
	"context"
	"fmt"
	"whatsapp/models"

	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client LogicServiceClient
	log    *logger.Logger
}

func NewClient(log *logger.Logger) *Client {
	conn, err := grpc.Dial("rs_bot:50051", grpc.WithInsecure())
	if err != nil {
		log.ErrorErr(err)
		return nil
	}
	fmt.Println("connect to rs grpc ok")
	return &Client{
		conn:   conn,
		client: NewLogicServiceClient(conn),
		log:    log,
	}
}

// Close закрывает соединение
func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) SendToRs2(i models.InMessageV2) error {
	in := InMessageV2{
		Text:        i.Text,
		Tip:         i.Tip,
		NameNick:    i.NameNick,
		UserName:    i.Username,
		UserId:      i.UserId,
		NameMention: i.NameMention,
		Messenger:   i.Messenger.ToMap(),
		Config: &CorporationConfigV2{
			//Id:       int32(i.Config.Id),
			Name: i.Config.Uid,
			//Language: i.Config.Language,
		},
		Options: i.Options,
	}
	for ch, inf := range i.Config.Channels {
		if in.Config.Channels == nil {
			in.Config.Channels = make(map[string]*ChannelInfo)
		}
		if in.Config.Channels[ch] == nil {
			in.Config.Channels[ch] = &ChannelInfo{}
		}
		in.Config.Channels[ch].Data = inf.ToMap()
	}
	for ch, inf := range i.Config.HelpMessage {
		if in.Config.HelpMessage == nil {
			in.Config.HelpMessage = make(map[string]*HelpMessageInfo)
		}
		if in.Config.HelpMessage[ch] == nil {
			in.Config.HelpMessage[ch] = &HelpMessageInfo{}
		}
		in.Config.HelpMessage[ch].Data = inf.ToMap()
	}

	_, err := c.client.LogicRs2(context.Background(), &in)
	if err != nil {
		return err
	}
	return nil
}
