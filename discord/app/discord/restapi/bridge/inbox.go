package bridge

import (
	"context"
	"discord/models"
	"fmt"
	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client BridgeServiceClient
	log    *logger.Logger
}

func NewClient(log *logger.Logger) *Client {
	conn, err := grpc.Dial("bridge:50051", grpc.WithInsecure())
	if err != nil {
		log.ErrorErr(err)
		return nil
	}
	fmt.Println("connect to bridge grpc ok")
	return &Client{
		conn:   conn,
		client: NewBridgeServiceClient(conn),
		log:    log,
	}
}

// Close закрывает соединение
func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) SendToBridge(i models.ToBridgeMessage) error {
	in := ToBridgeMessage{
		Text:        i.Text,
		Sender:      i.Sender,
		Tip:         i.Tip,
		ChatId:      i.ChatId,
		MesId:       i.MesId,
		GuildId:     i.GuildId,
		TimeMessage: i.TimestampUnix,
		Avatar:      i.Avatar,
	}
	if len(i.Extra) > 0 {
		for _, info := range i.Extra {
			in.Extra = append(in.Extra, &FileInfo{
				Name:   info.Name,
				Data:   info.Data,
				Url:    info.URL,
				Size:   info.Size,
				FileId: info.FileID,
			})
		}
	}
	if i.Reply != nil && i.Reply.UserName != "" {
		in.Reply = &BridgeMessageReply{
			TimeMessage: i.Reply.TimeMessage,
			Text:        i.Reply.Text,
			Avatar:      i.Reply.Avatar,
			UserName:    i.Reply.UserName,
		}
	}

	if i.Config != nil && i.Config.HostRelay != "" {
		conf := &BridgeConfig{
			Id:                int32(i.Config.Id),
			NameRelay:         i.Config.NameRelay,
			HostRelay:         i.Config.HostRelay,
			Role:              i.Config.Role,
			ForbiddenPrefixes: i.Config.ForbiddenPrefixes,
		}
		if len(i.Config.ChannelDs) > 0 {
			for _, d := range i.Config.ChannelDs {
				conf.ChannelDs = append(conf.ChannelDs, &BridgeConfigDs{
					ChannelId:       d.ChannelId,
					GuildId:         d.GuildId,
					CorpChannelName: d.CorpChannelName,
					AliasName:       d.AliasName,
					MappingRoles:    d.MappingRoles,
				})
			}
		}
		if len(i.Config.ChannelTg) > 0 {
			for _, t := range i.Config.ChannelTg {
				conf.ChannelTg = append(conf.ChannelTg, &BridgeConfigTg{
					ChannelId:       t.ChannelId,
					CorpChannelName: t.CorpChannelName,
					AliasName:       t.AliasName,
					MappingRoles:    t.MappingRoles,
				})
			}
		}

		in.Config = conf
	}

	_, err := c.client.InboxBridge(context.Background(), &in)
	return err
}
