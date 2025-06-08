package main

import (
	"discord/config"
	DiscordClient "discord/discord"
	"discord/grpc_server"
	"discord/storage"
	"github.com/mentalisit/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.InitConfig("DS")

	log := logger.LoggerZap(cfg.Logger.Token, cfg.Logger.ChatId, cfg.Logger.Webhook, "DS")

	st := storage.NewStorage(log, cfg)

	ds := DiscordClient.NewDiscord(log, st, cfg)
	grpc_server.GrpcMain(ds, log)

	//server.NewServer(ds, log)

	log.Info("Service discord load")

	//ожидаем сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ds.Shutdown()

}
