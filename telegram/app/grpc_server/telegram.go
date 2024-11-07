package grpc_server

import (
	"fmt"
	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
	"net"
	"telegram/telegram"
)

type Server struct {
	UnimplementedTelegramServiceServer
	log *logger.Logger
	tg  *telegram.Telegram
}

func GrpcMain(tg *telegram.Telegram, log *logger.Logger) *Server {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.ErrorErr(err)
	}

	s := grpc.NewServer()
	server := &Server{
		log: log,
		tg:  tg,
	}

	RegisterTelegramServiceServer(s, server)
	fmt.Printf("gRPC server is starting on %s\n", lis.Addr().String())
	//fmt.Println("Server is running on port :50051")
	go func() {
		if err = s.Serve(lis); err != nil {
			log.ErrorErr(err)
		}
	}()

	return server
}
