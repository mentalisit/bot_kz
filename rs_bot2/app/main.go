package main

import (
	"context"
	"fmt"
	"os/signal"
	"rs/bot2"
	"rs/clients"
	"rs/config"
	"rs/server"
	"rs/storage"
	"rs/webServer"
	"syscall"
	"time"

	"github.com/mentalisit/logger"
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
	//читаем конфигурацию с ENV
	cfg := config.InitConfig()

	//создаем логгер
	log := logger.LoggerZap(cfg.Logger.Token, cfg.Logger.ChatId, cfg.Logger.Webhook, "RS2")

	//storage
	st := storage.NewStorage(log, cfg)

	//clients Discord, Telegram
	cl := clients.NewClients(log, st)

	b := bot2.NewBot(st, cl, log)

	ws := webServer.NewServer(log, st, cl, b)

	b.AddLinkCodeFunc = ws.AddLinkCode

	g, _ := server.GrpcMain(b, log, st)

	//ожидаем сигнала завершения
	<-ctx.Done()
	cl.Shutdown()
	st.Shutdown()
	g.S.GracefulStop()

	//need write code save session and stop all services
	log.Info("shutdown")
	return nil
}

//go:generate protoc --go_out=. --go-grpc_out=. rs2.proto
