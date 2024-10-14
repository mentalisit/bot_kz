package main

import (
	"discord/config"
	DiscordClient "discord/discord"
	"discord/server"
	"discord/storage"
	"github.com/mentalisit/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.InitConfig()

	log := logger.LoggerZap(cfg.Logger.Token, cfg.Logger.ChatId, cfg.Logger.Webhook)

	st := storage.NewStorage(log, cfg)

	ds := DiscordClient.NewDiscord(log, st, cfg)

	server.NewServer(ds, log)

	log.Info("Service discord load")

	//ожидаем сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

}
