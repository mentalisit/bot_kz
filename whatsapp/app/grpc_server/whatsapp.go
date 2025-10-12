package grpc_server

import (
	"fmt"
	"net"
	wa "whatsapp/whatsapp"

	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
)

type Server struct {
	UnimplementedWhatsappServiceServer
	log *logger.Logger
	wa  *wa.Whatsapp
}

func GrpcMain(wa *wa.Whatsapp, log *logger.Logger) *Server {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.ErrorErr(err)
	}

	s := grpc.NewServer()
	server := &Server{
		log: log,
		wa:  wa,
	}

	RegisterWhatsappServiceServer(s, server)
	fmt.Printf("gRPC server is starting on %s\n", lis.Addr().String())
	go func() {
		if err = s.Serve(lis); err != nil {
			log.ErrorErr(err)
		}
	}()

	return server
}
