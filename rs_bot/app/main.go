package main

import (
	"context"
	"fmt"
	"os/signal"
	"rs/bot"
	"rs/clients"
	"rs/server"
	"rs/storage"
	"syscall"
	"time"

	"github.com/mentalisit/conf"
)

func main() {
	fmt.Println("Bot loading ")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := RunNew(ctx)
	if err != nil {
		fmt.Println("Error loading bot", err)
		time.Sleep(10 * time.Second)
		panic(err.Error())
	}
}

func RunNew(ctx context.Context) error {
	_, log, db := conf.InitConf("RS_BOT")

	//storage
	st := storage.NewStorage(log, db)

	//clients Discord, Telegram
	cl := clients.NewClients(log, st)

	b := bot.NewBot(st, cl, log)

	g, _ := server.GrpcMain(b, log)

	//ожидаем сигнала завершения
	<-ctx.Done()
	cl.Shutdown()
	st.Shutdown()
	g.S.GracefulStop()

	//need write code save session and stop all services
	log.Info("shutdown")
	return nil
}

//go:generate protoc --go_out=. --go-grpc_out=. rs.proto
