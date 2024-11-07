package grpc_server

import (
	"bridge/models"
	"bridge/server"
	"context"
	"fmt"
	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	UnimplementedBridgeServiceServer
	b   *server.Bridge
	log *logger.Logger
}

func GrpcMain(b *server.Bridge, log *logger.Logger) *Server {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.ErrorErr(err)
	}

	s := grpc.NewServer()
	serv := &Server{
		b:   b,
		log: log,
	}

	RegisterBridgeServiceServer(s, serv)

	fmt.Println("Server is running on port :50051")
	go func() {
		if err = s.Serve(lis); err != nil {
			log.ErrorErr(err)
		}
	}()

	return serv
}
func (s *Server) InboxBridge(ctx context.Context, i *ToBridgeMessage) (*Empty, error) {
	in := models.ToBridgeMessage{
		Text:          i.Text,
		Sender:        i.Sender,
		Tip:           i.Tip,
		ChatId:        i.ChatId,
		MesId:         i.MesId,
		GuildId:       i.GuildId,
		TimestampUnix: i.TimeMessage,
		Avatar:        i.Avatar,
	}
	if len(i.Extra) > 0 {
		for _, info := range i.Extra {
			in.Extra = append(in.Extra, models.FileInfo{
				Name:   info.Name,
				Data:   info.Data,
				URL:    info.Url,
				Size:   info.Size,
				FileID: info.FileId,
			})
		}
	}
	if i.Reply != nil && i.Reply.UserName != "" {
		in.Reply = &models.BridgeMessageReply{
			TimeMessage: i.Reply.TimeMessage,
			Text:        i.Reply.Text,
			Avatar:      i.Reply.Avatar,
			UserName:    i.Reply.UserName,
		}
	}

	if i.Config != nil && i.Config.HostRelay != "" {
		conf := &models.BridgeConfig{
			Id:                int(i.Config.Id),
			NameRelay:         i.Config.NameRelay,
			HostRelay:         i.Config.HostRelay,
			Role:              i.Config.Role,
			ForbiddenPrefixes: i.Config.ForbiddenPrefixes,
		}
		if len(i.Config.ChannelDs) > 0 {
			for _, d := range i.Config.ChannelDs {
				conf.ChannelDs = append(conf.ChannelDs, models.BridgeConfigDs{
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
				conf.ChannelTg = append(conf.ChannelTg, models.BridgeConfigTg{
					ChannelId:       t.ChannelId,
					CorpChannelName: t.CorpChannelName,
					AliasName:       t.AliasName,
					MappingRoles:    t.MappingRoles,
				})
			}
		}

		in.Config = conf
	}

	go s.b.Logic(in)
	return nil, nil
}
