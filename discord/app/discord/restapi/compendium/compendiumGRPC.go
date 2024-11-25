package compendium

import (
	"context"
	"discord/models"
	"fmt"
	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client LogicServiceClient
	log    *logger.Logger
}

func NewClient(log *logger.Logger) *Client {
	conn, err := grpc.Dial("compendiumnew:50051", grpc.WithInsecure())
	if err != nil {
		log.ErrorErr(err)
		return nil
	}
	fmt.Println("connect to compendiumNew grpc ok")
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

func (c *Client) SendToCompendium(inbox models.IncomingMessage) error {
	in := &IncomingMessage{
		Text:        inbox.Text,
		DmChat:      inbox.DmChat,
		Name:        inbox.Name,
		MentionName: inbox.MentionName,
		NameId:      inbox.NameId,
		NickName:    inbox.NickName,
		Avatar:      inbox.Avatar,
		ChannelId:   inbox.ChannelId,
		GuildId:     inbox.GuildId,
		GuildName:   inbox.GuildName,
		GuildAvatar: inbox.GuildAvatar,
		Type:        inbox.Type,
		Language:    inbox.Language,
	}
	_, err := c.client.InboxMessage(context.Background(), in)
	if err != nil {
		return err
	}
	return nil
}
