package rs_bot

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

func (c *Client) SendToRs(i models.InMessage) error {
	in := InMessage{
		Mtext:       i.Mtext,
		Tip:         i.Tip,
		NameNick:    i.NameNick,
		Username:    i.Username,
		UserId:      i.UserId,
		NameMention: i.NameMention,
		Lvlkz:       i.Lvlkz,
		Timekz:      i.Timekz,
		Ds: &Ds{
			Mesid:   i.Ds.Mesid,
			Guildid: i.Ds.Guildid,
			Avatar:  i.Ds.Avatar,
		},
		Tg: &Tg{Mesid: int32(i.Tg.Mesid)},
		Config: &CorporationConfig{
			Type:           int32(i.Config.Type),
			CorpName:       i.Config.CorpName,
			DsChannel:      i.Config.DsChannel,
			TgChannel:      i.Config.TgChannel,
			WaChannel:      i.Config.WaChannel,
			Country:        i.Config.Country,
			DelMesComplite: int32(i.Config.DelMesComplite),
			MesidDsHelp:    i.Config.MesidDsHelp,
			MesidTgHelp:    i.Config.MesidTgHelp,
			Forward:        i.Config.Forward,
			Guildid:        i.Config.Guildid,
		},
		Option: &Option{
			Reaction: i.Option.Reaction,
			InClient: i.Option.InClient,
			Queue:    i.Option.Queue,
			Pl30:     i.Option.Pl30,
			MinusMin: i.Option.MinusMin,
			Edit:     i.Option.Edit,
			Update:   i.Option.Update,
			Elsetrue: i.Option.Elsetrue,
		},
	}

	_, err := c.client.LogicRs(context.Background(), &in)
	if err != nil {
		return err
	}
	return nil
}
