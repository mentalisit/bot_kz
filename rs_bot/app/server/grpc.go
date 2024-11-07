package server

import (
	"context"
	"fmt"
	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
	"net"
	"rs/bot"
	"rs/models"
)

type Server struct {
	UnimplementedLogicServiceServer
	log *logger.Logger
	b   *bot.Bot
}

func GrpcMain(b *bot.Bot, log *logger.Logger) *Server {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.ErrorErr(err)
	}

	s := grpc.NewServer()
	serv := &Server{
		b:   b,
		log: log,
	}
	RegisterLogicServiceServer(s, serv)

	//fmt.Println("Server is running on port :50051")
	fmt.Printf("gRPC server is starting on %s\n", lis.Addr().String())
	go func() {
		if err = s.Serve(lis); err != nil {
			log.ErrorErr(err)
		}
	}()

	return serv
}
func (s *Server) LogicRs(ctx context.Context, i *InMessage) (*Empty, error) {
	in := models.InMessage{
		Mtext:       i.Mtext,
		Tip:         i.Tip,
		NameNick:    i.NameNick,
		Username:    i.Username,
		UserId:      i.UserId,
		NameMention: i.NameMention,
		Lvlkz:       i.Lvlkz,
		Timekz:      i.Timekz,
		Ds: struct {
			Mesid   string
			Guildid string
			Avatar  string
		}{Mesid: i.Ds.Mesid, Guildid: i.Ds.Guildid, Avatar: i.Ds.Avatar},
		Tg:     struct{ Mesid int }{Mesid: int(i.Tg.Mesid)},
		Config: models.CorporationConfig{},
		Option: models.Option{},
	}
	if i.Config != nil {
		in.Config = models.CorporationConfig{
			Type:           int(i.Config.Type),
			CorpName:       i.Config.CorpName,
			DsChannel:      i.Config.DsChannel,
			TgChannel:      i.Config.TgChannel,
			WaChannel:      i.Config.WaChannel,
			Country:        i.Config.Country,
			DelMesComplite: int(i.Config.DelMesComplite),
			MesidDsHelp:    i.Config.MesidDsHelp,
			MesidTgHelp:    i.Config.MesidTgHelp,
			Forward:        i.Config.Forward,
			Guildid:        i.Config.Guildid,
		}
	}
	if i.Option != nil {
		in.Option = models.Option{
			Reaction: i.Option.Reaction,
			InClient: i.Option.InClient,
			Queue:    i.Option.Queue,
			Pl30:     i.Option.Pl30,
			MinusMin: i.Option.MinusMin,
			Edit:     i.Option.Edit,
			Update:   i.Option.Update,
			Elsetrue: i.Option.Elsetrue,
		}
	}

	s.b.Inbox <- in
	return &Empty{}, nil
}
