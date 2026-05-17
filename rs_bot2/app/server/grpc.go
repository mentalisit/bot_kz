package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"rs/bot2"
	"rs/models"
	"rs/storage"
	"rs/storage/postgresV2"

	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
)

type Server struct {
	UnimplementedLogicServiceServer
	log *logger.Logger
	b2  *bot2.Bot
	S   *grpc.Server
	st  *postgresV2.Db
}

func GrpcMain(b2 *bot2.Bot, log *logger.Logger, st *storage.Storage) (*Server, error) {
	// Слушаем порт для gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.ErrorErr(err)
		return nil, err
	}

	// Создаём новый gRPC сервер
	s := grpc.NewServer()
	serv := &Server{
		b2:  b2,
		log: log,
		S:   s,
		st:  st.V2,
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

	//// Запускаем сервер профилирования (если необходим)
	//go func() {
	//	servprof.NewServer()
	//}()

	// Функция для graceful shutdown
	go func() {
		<-shutdown
		log.Info("Shutting down gRPC server...")
		s.GracefulStop()
		log.Info("gRPC server stopped.")
	}()

	return serv, nil
}

func (s *Server) OtherQueue(ctx context.Context, i *Other) (*Empty, error) {
	fmt.Println(i.GetUsersIds())
	//todo
	//s.b.ElseChat(i.GetUsersIds())
	return &Empty{}, nil
}
func (s *Server) LogicRs2(ctx context.Context, i *InMessageV2) (*Empty, error) {
	in := models.InMessageV2{
		Text:        i.Text,
		Tip:         i.Tip,
		NameNick:    i.NameNick,
		Username:    i.UserName,
		UserId:      i.UserId,
		NameMention: i.NameMention,
		Messenger:   InfoFromMap(i.Messenger),
		Options:     i.Options,
	}
	if i.Config != nil && i.Config.Name != "" {
		config := s.st.ReadConfigV2Uid(i.Config.Name)
		if config != nil {
			in.Config = *config
		}

	}

	s.b2.Inbox <- in
	return &Empty{}, nil
}

func InfoFromMap(m map[string]string) models.Info {
	if m == nil {
		return models.Info{}
	}

	i := models.Info{
		TypeMessenger:  m[models.MType],
		MessageId:      m[models.MMId],
		ChannelId:      m[models.MChId],
		ChannelName:    m[models.MChName],
		GuildId:        m[models.MGuId],
		GuildName:      m[models.MGuName],
		GuildAvatarUrl: m[models.MGuAvatarUrl],
		UserAvatarUrl:  m[models.MUsAvatarUrl],
		Language:       m[models.MLang],
	}
	return i
}
