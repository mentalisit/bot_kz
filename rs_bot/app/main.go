package main

import (
	"context"
	"fmt"
	"github.com/mentalisit/logger"
	"os/signal"
	"rs/bot"
	"rs/clients"
	"rs/config"
	"rs/server"
	"rs/storage"
	"syscall"
	"time"
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
	log := logger.LoggerZap(cfg.Logger.Token, cfg.Logger.ChatId, cfg.Logger.Webhook, "RS")

	//storage
	st := storage.NewStorage(log, cfg)

	//clients Discord, Telegram
	cl := clients.NewClients(log, st)

	b := bot.NewBot(st, cl, log)

	server.GrpcMain(b, log)

	//ожидаем сигнала завершения
	<-ctx.Done()
	cl.Shutdown()
	st.Shutdown()

	//need write code save session and stop all services
	log.Info("shutdown")
	return nil
}
