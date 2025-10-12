package grpc_server

import (
	"bridge/logic"
	"bridge/models"
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
)

type Server struct {
	UnimplementedBridgeServiceServer
	b   *logic.Bridge
	log *logger.Logger
}

func GrpcMain(b *logic.Bridge, log *logger.Logger) *Server {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.ErrorErr(err)
	}
	s := grpc.NewServer(grpc.MaxRecvMsgSize(10 * 1024 * 1024)) // Устанавливаем максимальный размер принимаемого сообщения в 10MB
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
		ReplyMap:      i.ReplyMap,
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
		conf := &models.Bridge2Config{
			Id:                int(i.Config.Id),
			NameRelay:         i.Config.NameRelay,
			HostRelay:         i.Config.HostRelay,
			Role:              i.Config.Role,
			ForbiddenPrefixes: i.Config.ForbiddenPrefixes,
			Channel:           make(map[string][]models.Bridge2Configs),
		}
		if len(i.Config.Channel) > 0 {
			for x, dd := range i.Config.Channel {
				if conf.Channel[x] == nil {
					conf.Channel[x] = []models.Bridge2Configs{}
				}
				for _, d := range dd.Configs {
					conf.Channel[x] = append(conf.Channel[x], models.Bridge2Configs{
						ChannelId:       d.ChannelId,
						GuildId:         d.GuildId,
						CorpChannelName: d.CorpChannelName,
						AliasName:       d.AliasName,
						MappingRoles:    d.MappingRoles,
					})
				}
			}
		}
		in.Config = conf
	}
	if in.Tip == "wa" {
		if strings.Contains(in.ChatId, "/") {
			split := strings.Split(in.ChatId, "/")
			in.ChatId = split[1]
		}
	}
	go s.b.Logic(in)
	return nil, nil
}
