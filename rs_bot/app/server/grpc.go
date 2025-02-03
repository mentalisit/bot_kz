package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
	"net"
	"rs/bot"
	"rs/models"
	servprof "rs/server/serv"
	"strconv"
)

type Server struct {
	UnimplementedLogicServiceServer
	log *logger.Logger
	b   *bot.Bot
	S   *grpc.Server
}

func GrpcMain(b *bot.Bot, log *logger.Logger) (*Server, error) {
	// Слушаем порт для gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.ErrorErr(err)
		return nil, err
	}

	// Создаём новый gRPC сервер
	s := grpc.NewServer()
	serv := &Server{
		b:   b,
		log: log,
		S:   s,
	}
	RegisterLogicServiceServer(s, serv)

	// Сообщение о старте сервера
	fmt.Printf("gRPC server is starting on %s\n", lis.Addr().String())

	// Создаём канал для завершения работы
	shutdown := make(chan struct{})

	// Запускаем gRPC сервер в горутине
	go func() {
		if err = s.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.ErrorErr(err)
			close(shutdown) // Закрываем канал при ошибке
		}
	}()

	// Запускаем сервер профилирования (если необходим)
	go func() {
		servprof.NewServer()
	}()

	// Функция для graceful shutdown
	go func() {
		<-shutdown
		log.Info("Shutting down gRPC server...")
		s.GracefulStop()
		log.Info("gRPC server stopped.")
	}()

	return serv, nil
}

func (s *Server) LogicRs(ctx context.Context, i *InMessage) (*Empty, error) {
	atoi, _ := strconv.Atoi(i.Timekz)

	in := models.InMessage{
		Mtext:       i.Mtext,
		Tip:         i.Tip,
		NameNick:    i.NameNick,
		Username:    i.Username,
		UserId:      i.UserId,
		NameMention: i.NameMention,
		//Lvlkz:       i.Lvlkz,
		//Timekz:      i.Timekz,
		RsTypeLevel: i.Lvlkz,
		TimeRs:      atoi,
		Ds: struct {
			Mesid   string
			Guildid string
			Avatar  string
		}{Mesid: i.Ds.Mesid, Guildid: i.Ds.Guildid, Avatar: i.Ds.Avatar},
		Tg:     struct{ Mesid int }{Mesid: int(i.Tg.Mesid)},
		Config: models.CorporationConfig{},
		//Option: models.Option{},
		Opt: models.Options{},
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
		o := i.Option
		if o.Edit {
			in.Opt.Add(models.OptionEdit)
		}
		if o.Queue {
			in.Opt.Add(models.OptionQueue)
		}
		if o.MinusMin {
			in.Opt.Add(models.OptionMinusMin)
		}
		if o.Pl30 {
			in.Opt.Add(models.OptionPl30)
		}
		if o.Update {
			in.Opt.Add(models.OptionUpdate)
		}
		if o.InClient {
			in.Opt.Add(models.OptionInClient)
		}
		if o.Elsetrue {
			in.Opt.Add(models.OptionElseTrue)
		}
		if o.Reaction {
			in.Opt.Add(models.OptionReaction)
		}
		//in.Option = models.Option{
		//	Reaction: i.Option.Reaction,
		//	InClient: i.Option.InClient,
		//	Queue:    i.Option.Queue,
		//	Pl30:     i.Option.Pl30,
		//	MinusMin: i.Option.MinusMin,
		//	Edit:     i.Option.Edit,
		//	Update:   i.Option.Update,
		//	Elsetrue: i.Option.Elsetrue,
		//}
	}

	s.b.Inbox <- in
	return &Empty{}, nil
}
