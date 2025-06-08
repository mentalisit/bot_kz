package grpc_server

import (
	DiscordClient "discord/discord"
	"fmt"
	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	UnimplementedBotServiceServer
	log *logger.Logger
	ds  *DiscordClient.Discord
}

func GrpcMain(ds *DiscordClient.Discord, log *logger.Logger) *Server {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.ErrorErr(err)
	}

	s := grpc.NewServer(grpc.MaxRecvMsgSize(10 * 1024 * 1024)) // Устанавливаем максимальный размер принимаемого сообщения в 10MB
	server := &Server{
		log: log,
		ds:  ds,
	}

	RegisterBotServiceServer(s, server)
	fmt.Printf("gRPC server is starting on %s\n", lis.Addr().String())
	//fmt.Println("Server is running on port :50051")
	go func() {
		if err = s.Serve(lis); err != nil {
			log.ErrorErr(err)
		}
	}()

	return server
}
