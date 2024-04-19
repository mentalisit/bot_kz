package main

import (
	"context"
	"fmt"
	"github.com/mentalisit/logger"
	"kz_bot/bot"
	"kz_bot/clients"
	"kz_bot/config"
	"kz_bot/server"
	"kz_bot/storage"
	"os"
	"os/signal"
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
	log := logger.LoggerZap(cfg.Logger.Token, cfg.Logger.ChatId, cfg.Logger.Webhook)

	if cfg.BotMode == "dev" {
		log = logger.LoggerZapDEV()

		go func() {
			time.Sleep(5 * time.Minute)
			os.Exit(1)
		}()
		//os.Exit(1)
	}

	log.Info("🚀  загрузка  🚀 " + cfg.BotMode)

	//storage
	st := storage.NewStorage(log, cfg)

	//clients Discord, Telegram
	cl := clients.NewClients(log, st, cfg)

	server.NewServer(cl, log)

	go bot.NewBot(st, cl, log, cfg)
	//go BridgeChat.NewBridge(log, cl, st)

	//ожидаем сигнала завершения
	<-ctx.Done()
	cl.Shutdown()
	st.Shutdown()

	//need write code save session and stop all services
	log.Info("shutdown")
	return nil
}
